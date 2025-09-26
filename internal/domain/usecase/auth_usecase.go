package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
)

// AuthInternoUsecase é a interface para casos de uso de autenticação
type AuthInternoUsecase interface {
	Login(ctx context.Context, login, senha string) (*response.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*response.TokenPair, error)
	Me(ctx context.Context, userID string) (*response.UsuarioResponse, error)
}

// AuthExterno é a interface para sistemas externos de autenticação (LDAP, OAuth, etc.)
type AuthExternoUsecase interface {
	Bind(login, senha string) error
	PesquisarPorLogin(login string) (nome, email, outLogin string, err error)
}
