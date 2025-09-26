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
	users usecase.UsuarioUsecase
	jwt   jwt.JWTUsecase
	ldap  usecase.AuthExternoUsecase
	cfg   config.Config
}

func NewAuthInternoUsecase(users usecase.UsuarioUsecase, jwt jwt.JWTUsecase, ldap usecase.AuthExternoUsecase, cfg config.Config) usecase.AuthInternoUsecase {
	return &authUsecase{users, jwt, ldap, cfg}
}

// --- Auxiliares internos ---

func (a *authUsecase) getBindString(login string) string {
	if a.cfg.LDAPDomain != "" {
		return login + a.cfg.LDAPDomain
	}

	ou := "users"
	if login == "admin1" {
		ou = "admins"
	}
	return "uid=" + login + ",ou=" + ou + "," + a.cfg.LDAPBase
}

func (a *authUsecase) criarUsuarioSeNecessario(ctx context.Context, login string, u *model.Usuario) (*model.Usuario, error) {
	if u != nil {
		return u, nil
	}

	name, mail, sLogin, err := a.ldap.PesquisarPorLogin(login)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário no LDAP: %w", err)
	}

	novo := &model.Usuario{
		Nome:      name,
		Login:     sLogin,
		Email:     mail,
		Permissao: model.PermUSR,
		Status:    true,
	}

	if err := a.users.CriarUsuario(ctx, novo); err != nil {
		return nil, fmt.Errorf("erro ao criar usuário no banco: %w", err)
	}

	return a.users.BuscarUsuarioPorLogin(ctx, login)
}

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

func (a *authUsecase) Login(ctx context.Context, login, senha string) (*response.TokenPair, error) {
	usuario, err := a.users.BuscarUsuarioPorLogin(ctx, login)
	if err != nil && !errors.Is(err, repository.ErrUsuarioNaoEncontrado) {
		return nil, fmt.Errorf("erro buscando usuário: %w", err)
	}
	if errors.Is(err, repository.ErrUsuarioNaoEncontrado) {
		usuario = nil
	}

	bind := a.getBindString(login)
	if err := a.ldap.Bind(bind, senha); err != nil {
		return nil, fmt.Errorf("credenciais inválidas: %w", err)
	}

	usuario, err = a.criarUsuarioSeNecessario(ctx, login, usuario)
	if err != nil {
		return nil, err
	}

	_ = a.users.AtualizarUltimoLoginUsuario(ctx, usuario.ID)

	claims := createClaims(usuario)
	access, err := a.jwt.GerarToken(claims)
	if err != nil {
		return nil, err
	}
	refresh, err := a.jwt.GerarRefreshToken(claims)
	if err != nil {
		return nil, err
	}

	return &response.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (a *authUsecase) Refresh(ctx context.Context, refreshToken string) (*response.TokenPair, error) {
	claims, err := a.jwt.ValidarRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh inválido: %w", err)
	}

	usuario, err := a.users.BuscarUsuarioPorID(ctx, claims.ID)
	if err != nil || usuario == nil {
		return nil, errors.New("usuário inválido")
	}

	_ = a.users.AtualizarUltimoLoginUsuario(ctx, usuario.ID)

	claims.RegisteredClaims.IssuedAt = gojwt.NewNumericDate(time.Now())

	access, err := a.jwt.GerarToken(*claims)
	if err != nil {
		return nil, err
	}
	refresh, err := a.jwt.GerarRefreshToken(*claims)
	if err != nil {
		return nil, err
	}

	return &response.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (a *authUsecase) Me(ctx context.Context, userID string) (*response.UsuarioResponse, error) {
    usuario, err := a.users.BuscarUsuarioPorID(ctx, userID)
    if err != nil || usuario == nil {
        return nil, err
    }
    return response.ToUsuarioResponse(usuario), nil
}
