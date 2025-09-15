package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
)

// AuthUsecase é a interface exposta para a aplicação
type AuthUsecase interface {
	Login(ctx context.Context, login, senha string) (*response.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*response.TokenPair, error)
	Me(ctx context.Context, userID string) (*response.UsuarioResponse, error)
}

// Authenticator é a interface para sistemas externos de autenticação (LDAP, OAuth, etc.)
type Authenticator interface {
	Bind(login, senha string) error
	SearchByLogin(login string) (nome, email, outLogin string, err error)
}
