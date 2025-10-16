package jwt

import (
	"errors"
	"fmt"
	"time"

	goJwt "github.com/golang-jwt/jwt/v5"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

var (
	ErrParseWithClaims = errors.New("erro ao fazer parse com claims")
	ErrSignedString    = errors.New("erro ao assinar token")
)

// JWTUsecase define os métodos que a implementação JWT deve fornecer
type JWTUsecase interface {
	// GerarToken gera um token de acesso
	GerarToken(c Claims) (string, error)

	// GerarRefreshToken gera um token de refresh
	GerarRefreshToken(c Claims) (string, error)

	// ValidarRefreshToken valida um token de refresh
	ValidarRefreshToken(token string) (*Claims, error)

	// ValidarToken valida um token de acesso
	ValidarToken(token string) (*Claims, error)
}

// GerenteJWT gerencia a criação e validação de tokens JWT
type GerenteJWT struct {
	ChaveAcesso  []byte
	ChaveRefresh []byte
	TLLAcesso    time.Duration
	TLLRefresh   time.Duration
}

// NewGerenteJWT cria uma nova instância de GerenteJWT
func NewGerenteJWT(chaveAcesso, chaveRefresh []byte, ttlAcesso, ttlRefresh time.Duration) *GerenteJWT {
	return &GerenteJWT{
		ChaveAcesso:  chaveAcesso,
		ChaveRefresh: chaveRefresh,
		TLLAcesso:    ttlAcesso,
		TLLRefresh:   ttlRefresh,
	}
}

// Claims define as claims personalizadas para o token JWT
type Claims struct {
	ID        string `json:"sub"`
	Login     string `json:"login"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Permissao string `json:"permissao"`
	goJwt.RegisteredClaims
}

// GerarToken gera um token de acesso
func (g *GerenteJWT) GerarToken(c Claims) (string, error) {
	tokenGerado, err := g.gerarJWT(c, g.ChaveAcesso, g.TLLAcesso)
	if err != nil {
		return "", fmt.Errorf("[jwt.GerarToken]: %w", err)
	}
	return tokenGerado, nil
}

// GerarRefreshToken gera um token de refresh
func (g *GerenteJWT) GerarRefreshToken(c Claims) (string, error) {
	refreshTokenGerado, err := g.gerarJWT(c, g.ChaveRefresh, g.TLLRefresh)
	if err != nil {
		return "", fmt.Errorf("[jwt.GerarRefreshToken]: %w", err)
	}
	return refreshTokenGerado, nil
}

// ValidarToken valida um token de acesso
func (g *GerenteJWT) ValidarToken(token string) (*Claims, error) {
	claimsValidadas, err := g.validarJWT(token, g.ChaveAcesso)
	if err != nil {
		return nil, fmt.Errorf("[jwt.ValidarToken]: %w", err)
	}
	return claimsValidadas, nil
}

// ValidarRefreshToken valida um token de refresh
func (g *GerenteJWT) ValidarRefreshToken(token string) (*Claims, error) {
	claimsValidadas, err := g.validarJWT(token, g.ChaveRefresh)
	if err != nil {
		return nil, fmt.Errorf("[jwt.ValidarRefreshToken]: %w", err)
	}
	return claimsValidadas, nil
}

// gerarJWT é uma função helper interna para gerar token
func (g *GerenteJWT) gerarJWT(c Claims, secret []byte, ttl time.Duration) (string, error) {
	c.RegisteredClaims.ExpiresAt = goJwt.NewNumericDate(time.Now().Add(ttl))
	tokenJWT := goJwt.NewWithClaims(goJwt.SigningMethodHS256, c)

	jwtString, err := tokenJWT.SignedString(secret)
	if err != nil {
		return "", utils.NewAppError(
			"[jwt.gerarJWT]",
			utils.LevelError,
			"erro ao tentar assinar token",
			fmt.Errorf(utils.FmtErroWrap, ErrSignedString, err),
		)
	}
	return jwtString, nil
}

// validar é uma função helper interna para validar token
func (g *GerenteJWT) validarJWT(tokenStr string, secret []byte) (*Claims, error) {
	const metodo = "[jwt.validar]"
	parsed, err := goJwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *goJwt.Token) (any, error) {
			return secret, nil
		})

	// Verifica se houve erro na validação
	if err != nil {
		if errors.Is(err, goJwt.ErrTokenExpired) {
			return nil, utils.NewAppError(
				metodo,
				utils.LevelError,
				"erro ao tentar validar token expirado",
				err,
			)
		}
		if errors.Is(err, goJwt.ErrTokenSignatureInvalid) {
			return nil, utils.NewAppError(
				metodo,
				utils.LevelError,
				"erro ao tentar validar token com assinatura inválida",
				err,
			)
		}
		return nil, utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao tentar validar token",
			fmt.Errorf(utils.FmtErroWrap, ErrParseWithClaims, err),
		)
	}

	// Verifica se as claims são válidas
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, utils.NewAppError(
		metodo,
		utils.LevelError,
		"erro ao tentar validar token: claims inválidos",
		goJwt.ErrTokenInvalidClaims,
	)
}
