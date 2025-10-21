package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
)

// AuthInternoUsecase é a interface para casos de uso de autenticação
type AuthInternoUsecase interface {
	// Login realiza a autenticação de um usuário e gera um token.
	Login(ctx context.Context, login, senha string) (*response.TokenPair, error)

	// Refresh renova um token de acesso usando um token de atualização.
	Refresh(ctx context.Context, refreshToken string) (*response.TokenPair, error)

	// Me retorna os dados do usuário autenticado.
	Me(ctx context.Context, userID string) (*response.UsuarioResponse, error)
}

// AuthExterno é a interface para sistemas externos de autenticação (LDAP, OAuth, etc.)
type AuthExternoUsecase interface {
	// Bind tenta autenticar um usuário com o sistema externo.
	Bind(login, senha string) error

	// PesquisarPorLogin busca informações de um usuário no sistema externo pelo login.
	PesquisarPorLogin(login string) (nome, email, outLogin string, err error)
}
