package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
)

type authUsecase struct {
	UsecaseUsuario usecase.UsuarioUsecase
	UsecaseJWT     jwt.JWTUsecase
	UsecaseLDAP    usecase.AuthExternoUsecase
	UsecaseLog     usecase.LogUsecase
	Config         config.Config
}

func NewAuthInternoUsecase(
	usecaseUsuario usecase.UsuarioUsecase,
	usecaseJWT jwt.JWTUsecase,
	usecaseLDAP usecase.AuthExternoUsecase,
	usecaseLog usecase.LogUsecase,
	config config.Config,
) usecase.AuthInternoUsecase {

	return &authUsecase{usecaseUsuario, usecaseJWT, usecaseLDAP, usecaseLog, config}
}

// --- Auxiliares internos ---

// getBindString retorna a string de bind para o LDAP, dependendo da configuração.
func (a *authUsecase) getBindString(login string) string {
	if a.Config.LDAPDomain != "" {
		return login + a.Config.LDAPDomain
	}

	ou := "users"
	if login == "admin1" {
		ou = "admins"
	}
	return "uid=" + login + ",ou=" + ou + "," + a.Config.LDAPBase
}

// criarUsuarioSeNecessario cria um novo usuário no banco de dados se ele não existir.
func (a *authUsecase) criarUsuarioSeNecessario(ctx context.Context, login string, u *model.Usuario) (*model.Usuario, error) {
	const metodo = "[usecase.auth.criarUsuarioSeNecessario]: %w"
	if u != nil {
		return u, nil
	}

	name, mail, sLogin, err := a.UsecaseLDAP.PesquisarPorLogin(login)
	if err != nil {
		return nil, fmt.Errorf(metodo, err)
	}

	usuario := &model.Usuario{
		Nome:      name,
		Login:     sLogin,
		Email:     mail,
		Permissao: model.PermUSR,
		Status:    true,
	}

	if err := a.UsecaseUsuario.CriarUsuario(ctx, usuario); err != nil {
		return nil, fmt.Errorf(metodo, err)
	}

	usuarioSalvo, err := a.UsecaseUsuario.BuscarUsuarioPorLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf(metodo, err)
	}
	return usuarioSalvo, nil
}

// createClaims cria as claims do JWT a partir do usuário.
func createClaims(u *model.Usuario) jwt.Claims {
	return jwt.Claims{
		ID:        u.ID,
		Login:     u.Login,
		Nome:      u.Nome,
		Email:     u.Email,
		Permissao: string(u.Permissao),
	}
}

// --- Implementações da interface ---

// Login autentica o usuário e retorna um par de tokens (access e refresh).
func (a *authUsecase) Login(ctx context.Context, login, senha string) (*response.TokenPair, error) {
	const metodo = "[usecase.auth.Login]: %w"

	usuario, err := a.UsecaseUsuario.BuscarUsuarioPorLogin(ctx, login)
	if err != nil && !errors.Is(err, repository.ErrUsuarioNaoEncontrado) {
		return nil, fmt.Errorf(metodo, err)
	}
	if errors.Is(err, repository.ErrUsuarioNaoEncontrado) {
		usuario = nil
	}

	bind := a.getBindString(login)
	if err := a.UsecaseLDAP.Bind(bind, senha); err != nil {
		return nil, fmt.Errorf(metodo, err)
	}

	usuario, err = a.criarUsuarioSeNecessario(ctx, login, usuario)
	if err != nil {
		return nil, fmt.Errorf(metodo, err)
	}

	_ = a.UsecaseUsuario.AtualizarUltimoLoginUsuario(ctx, usuario.ID)

	claims := createClaims(usuario)
	access, err := a.UsecaseJWT.GerarToken(claims)
	if err != nil {
		return nil, fmt.Errorf(metodo, err)
	}
	refresh, err := a.UsecaseJWT.GerarRefreshToken(claims)
	if err != nil {
		return nil, fmt.Errorf(metodo, err)
	}

	return &response.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

// Refresh valida o refresh token e retorna um novo par de tokens (access e refresh).
func (a *authUsecase) Refresh(ctx context.Context, refreshToken string) (*response.TokenPair, error) {
	const metodo = "[usecase.auth.Refresh]: %w"
	claims, err := a.UsecaseJWT.ValidarRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf(metodo, err)
	}

	usuario, err := a.UsecaseUsuario.BuscarUsuarioPorID(ctx, claims.ID)
	if err != nil || usuario == nil {
		return nil, fmt.Errorf(metodo, err)
	}

	_ = a.UsecaseUsuario.AtualizarUltimoLoginUsuario(ctx, usuario.ID)

	claims.RegisteredClaims.IssuedAt = gojwt.NewNumericDate(time.Now())

	access, err := a.UsecaseJWT.GerarToken(*claims)
	if err != nil {
		return nil, fmt.Errorf(metodo, err)
	}
	refresh, err := a.UsecaseJWT.GerarRefreshToken(*claims)
	if err != nil {
		return nil, fmt.Errorf(metodo, err)
	}

	return &response.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

// Me retorna os dados do usuário autenticado.
func (a *authUsecase) Me(ctx context.Context, userID string) (*response.UsuarioResponse, error) {
	usuario, err := a.UsecaseUsuario.BuscarUsuarioPorID(ctx, userID)
	if err != nil || usuario == nil {
		return nil, fmt.Errorf("[usecase.auth.Me]: %w", err)
	}
	return response.ToUsuarioResponse(usuario), nil
}
