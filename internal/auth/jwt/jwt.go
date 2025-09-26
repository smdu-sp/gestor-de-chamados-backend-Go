package jwt

import (
	"errors"
	"fmt"
	"time"

	goJwt "github.com/golang-jwt/jwt/v5"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

var ErrParseWithClaims = errors.New("erro ao fazer parse com claims")

// JWTUsecase define os métodos que a implementação JWT deve fornecer
type JWTUsecase interface {
	GerarToken(c Claims) (string, error)
	GerarRefreshToken(c Claims) (string, error)
	ValidarRefreshToken(token string) (*Claims, error)
	ValidarToken(token string) (*Claims, error)
}

// GerenteJWT gerencia a criação e validação de tokens JWT
type GerenteJWT struct {
	ChaveAcesso  []byte
	ChaveRefresh []byte
	TLLAcesso    time.Duration
	TLLRefresh   time.Duration
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
	return g.gerarJWT(c, g.ChaveAcesso, g.TLLAcesso)
}

// GerarRefreshToken gera um token de refresh
func (g *GerenteJWT) GerarRefreshToken(c Claims) (string, error) {
	return g.gerarJWT(c, g.ChaveRefresh, g.TLLRefresh)
}

// ValidarToken valida um token de acesso
func (g *GerenteJWT) ValidarToken(token string) (*Claims, error) {
	return g.validarJWT(token, g.ChaveAcesso)
}

// ValidarRefreshToken valida um token de refresh
func (g *GerenteJWT) ValidarRefreshToken(token string) (*Claims, error) {
	return g.validarJWT(token, g.ChaveRefresh)
}

// gerarJWT é uma função helper interna para gerar token
func (g *GerenteJWT) gerarJWT(c Claims, secret []byte, ttl time.Duration) (string, error) {
	c.RegisteredClaims.ExpiresAt = goJwt.NewNumericDate(time.Now().Add(ttl))
	tokenJWT := goJwt.NewWithClaims(goJwt.SigningMethodHS256, c)
	return tokenJWT.SignedString(secret)
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
