package auth

import (
	"context"

	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// AuthInternoService é a interface para casos de uso de autenticação
type AuthInternoService interface {
	Login(ctx context.Context, req LoginReq) (*TokenResp, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenResp, error)
	Me(ctx context.Context, userID string) (*usr.Usuario, error)
}

// AuthExternoService é a interface para sistemas externos de autenticação (LDAP, OAuth, etc.)
type AuthExternoService interface {
	Bind(login, senha string) error
	PesquisarPorLogin(login string) (*UsuarioLDAP, error)
}

// JWTService define os métodos que a implementação JWT deve fornecer
type JWTService interface {
	GerarToken(claims Claims) (string, error)
	GerarRefreshToken(claims Claims) (string, error)
	ValidarRefreshToken(token string) (*Claims, error)
	ValidarToken(refreshToken string) (*Claims, error)
}
