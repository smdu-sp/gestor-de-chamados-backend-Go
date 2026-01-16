package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

type AuthInterno struct {
	UsuarioSvc usr.Service
	JWTSvc     JWTService
	LDAPSvc    AuthExternoService
	LogSvc     log.Service
	Config     *cfg.Config
}

func NovoAuthInterno(
	usuarioSvc usr.Service,
	jwtSvc JWTService,
	ldapSvc AuthExternoService,
	logSvc log.Service,
	config *cfg.Config,
) *AuthInterno {

	return &AuthInterno{
		UsuarioSvc: usuarioSvc,
		JWTSvc:     jwtSvc,
		LDAPSvc:    ldapSvc,
		LogSvc:     logSvc,
		Config:     config,
	}
}

// Asserção de interface para garantir que AuthInterno implementa AuthInternoService
var _ AuthInternoService = (*AuthInterno)(nil)

// Login recebe as credenciais do usuário, autentica via LDAP, cria o usuário se necessário,
// atualiza o último login e retorna um par de tokens JWT (access e refresh).
//
// Erros sentinelas possíveis:  
func (a *AuthInterno) Login(ctx context.Context, req LoginReq) (*TokenResp, error) {
	// 1- Verificar se o usuário existe no sistema
	usuario, err := a.UsuarioSvc.BuscarPorLogin(ctx, req.Login)
	if err != nil && !errors.Is(err, mysql.ErrUsuarioNaoEncontrado) {
		return nil, fmt.Errorf("Login: %w", err)
	}

	// Se o usuário não for encontrado, definir como nil para criar depois
	if errors.Is(err, mysql.ErrUsuarioNaoEncontrado) {
		usuario = nil
	}

	// 2- Autenticar o usuário via LDAP
	bind := a.getBindString(req.Login)
	if err := a.LDAPSvc.Bind(bind, req.Senha); err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	// 3- Criar o usuário no sistema se não existir
	usuarioCriado, err := a.criarUsuarioSeNecessario(ctx, req.Login, usuario)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	// 4- Atualizar o último login do usuário
	_, err = a.UsuarioSvc.AtualizarUltimoLogin(ctx, usuarioCriado.ID())
	if err != nil {
		slog.Error("Login", "error", err)
	}

	// 5- Gerar tokens JWT
	claims := novoClaims(usuarioCriado)
	access, err := a.JWTSvc.GerarToken(claims)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}
	refresh, err := a.JWTSvc.GerarRefreshToken(claims)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	return ParaTokenResp(access, refresh), nil
}

// Refresh valida o refresh token e retorna um novo par de tokens (access e refresh).
func (a *AuthInterno) Refresh(ctx context.Context, refreshToken string) (*TokenResp, error) {
	// 1- Validar o refresh token
	claims, err := a.JWTSvc.ValidarRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}

	// 2- Buscar o usuário associado ao token
	usuario, err := a.UsuarioSvc.BuscarPorID(ctx, claims.ID)
	if err != nil || usuario == nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}

	// 3- Atualizar o último login do usuário
	_, err = a.UsuarioSvc.AtualizarUltimoLogin(ctx, usuario.ID())
	if err != nil {
		slog.Error("Refresh", "error", err)
	}

	// 4- Atualizar o campo IssuedAt para o tempo atual
	claims.RegisteredClaims.IssuedAt = gojwt.NewNumericDate(time.Now())

	// 5- Gerar novos tokens
	access, err := a.JWTSvc.GerarToken(*claims)
	if err != nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}
	refresh, err := a.JWTSvc.GerarRefreshToken(*claims)
	if err != nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}

	return ParaTokenResp(access, refresh), nil
}

// Me retorna os dados do usuário autenticado.
func (a *AuthInterno) Me(ctx context.Context, id string) (*usr.Usuario, error) {
	usuario, err := a.UsuarioSvc.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Me: %w", err)
	}
	return usuario, nil
}

//==============================================================================
// Métodos auxiliares
//=============================================================================

// getBindString retorna a string de bind para o LDAP, dependendo da configuração.
// Se LDAPDomain estiver configurado, retorna "login + domain".
// Caso contrário, retorna o DN completo baseado no login.
// Util para autenticação LDAP com diferentes esquemas, como Active Directory ou OpenLDAP.
func (a *AuthInterno) getBindString(login string) string {
	if a.Config.LDAPDomain() != "" {
		return login + a.Config.LDAPDomain()
	}

	ou := "users"
	if login == "admin1" {
		ou = "admins"
	}
	return "uid=" + login + ",ou=" + ou + "," + a.Config.LDAPBase()
}

// criarUsuarioSeNecessario cria um novo usuário no banco de dados se ele não existir.
func (a *AuthInterno) criarUsuarioSeNecessario(ctx context.Context, login string, usuario *usr.Usuario) (*usr.Usuario, error) {
	// 1- Se o usuário já existe, retornar
	if usuario != nil {
		return usuario, nil
	}

	// 2- Buscar informações do usuário no LDAP
	usuarioLDAP, err := a.LDAPSvc.PesquisarPorLogin(login)
	if err != nil {
		return nil, fmt.Errorf("criarUsuarioSeNecessario: %w", err)
	}

	// 3- Criar o usuário no sistema
	usuarioNovo := &usr.CriarParams{
		Nome:      usuarioLDAP.Nome,
		Login:     usuarioLDAP.Login,
		Email:     usr.NovoEmail(usuarioLDAP.Email),
		Permissao: usr.PermUSR,
	}

	// 4- Salvar o usuário no banco de dados
	usuarioSalvo, err := a.UsuarioSvc.Criar(ctx, *usuarioNovo)
	if err != nil {
		return nil, fmt.Errorf("criarUsuarioSeNecessario: %w", err)
	}

	return usuarioSalvo, nil
}
