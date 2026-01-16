package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	goJwt "github.com/golang-jwt/jwt/v5"
)

// Erros específicos do pacote jwt
var (
	ErrParseWithClaims = errors.New("erro ao fazer parse com claims")
	ErrSignedToken     = errors.New("erro ao assinar token")
)

// TokenJWT gerencia a criação e validação de tokens JWT
type TokenJWTService struct {
	tokenSegredo        []byte
	refreshTokenSegredo []byte
	tokenTTL            time.Duration
	refreshTokenTTL     time.Duration
}

// NovoTokenJWT cria uma nova instância de TokenJWT
func NovoTokenJWTService(
	tokenSegredo, refreshTokenSegredo []byte,
	tokenTTL, refreshTokenTTL time.Duration,
) *TokenJWTService {
	return &TokenJWTService{
		tokenSegredo:        tokenSegredo,
		refreshTokenSegredo: refreshTokenSegredo,
		tokenTTL:            tokenTTL,
		refreshTokenTTL:     refreshTokenTTL,
	}
}

// Asserção de interface para garantir que JWTService implementa auth.JWTService
var _ JWTService = (*TokenJWTService)(nil)

// GerarToken gera um token de acesso
func (g *TokenJWTService) GerarToken(claims Claims) (string, error) {
	tokenGerado, err := g.gerarJWT(claims, g.tokenSegredo, g.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("gerar token: %w", err)
	}

	return tokenGerado, nil
}

// GerarRefreshToken gera um token de refresh
func (g *TokenJWTService) GerarRefreshToken(claims Claims) (string, error) {
	refreshTokenGerado, err := g.gerarJWT(claims, g.refreshTokenSegredo, g.refreshTokenTTL)
	if err != nil {
		return "", fmt.Errorf("gerar refresh token: %w", err)
	}

	return refreshTokenGerado, nil
}

// ValidarToken valida um token de acesso
func (g *TokenJWTService) ValidarToken(token string) (*Claims, error) {
	claimsValidadas, err := g.validarJWT(token, g.tokenSegredo)
	if err != nil {
		return nil, err
	}

	return claimsValidadas, nil
}

// ValidarRefreshToken valida um token de refresh
func (g *TokenJWTService) ValidarRefreshToken(refreshToken string) (*Claims, error) {
	claimsValidadas, err := g.validarJWT(refreshToken, g.refreshTokenSegredo)
	if err != nil {
		slog.Error("Erro ao validar token de refresh", "token", refreshToken, "erro", err)
		return nil, fmt.Errorf("validar refresh token: %w", err)
	}

	return claimsValidadas, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// gerarJWT é uma função helper interna para gerar token
func (g *TokenJWTService) gerarJWT(claims Claims, segredo []byte, ttl time.Duration) (string, error) {
	claims.ExpiresAt = goJwt.NewNumericDate(time.Now().Add(ttl))
	token := goJwt.NewWithClaims(goJwt.SigningMethodHS256, claims)

	tokenJWT, err := token.SignedString(segredo)
	if err != nil {
		slog.Error("Erro ao assinar JWT", "usuario_id", claims.ID, "erro", err)
		return "", fmt.Errorf("gerar JWT: %w", ErrSignedToken)
	}

	return tokenJWT, nil
}

// validar é uma função helper interna para validar token
func (g *TokenJWTService) validarJWT(token string, segredo []byte) (*Claims, error) {
	// Faz o parse do token com as claims
	parsed, err := goJwt.ParseWithClaims(
		token,
		&Claims{},
		func(token *goJwt.Token) (any, error) {
			return segredo, nil
		},
	)

	// Verifica se houve erro na validação
	if err != nil {
		switch {
		case errors.Is(err, goJwt.ErrTokenExpired):
			slog.Error("Token JWT expirado", "token", token)
			return nil, fmt.Errorf("erro ao validar token JWT: %w", err)

		case errors.Is(err, goJwt.ErrTokenSignatureInvalid):
			slog.Error("Assinatura inválida no token JWT", "token", token)
			return nil, fmt.Errorf("erro ao validar token JWT: %w", err)

		default:
			slog.Error("Erro ao fazer parse do token JWT", "token", token, "erro", err)
			return nil, fmt.Errorf("erro ao fazer parse do token JWT: %w", ErrParseWithClaims)
		}
	}

	// Verifica se as claims são válidas
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		slog.Error("Claims inválidas no token JWT", "token", token)
		return nil, fmt.Errorf("erro ao fazer parse do token JWT: %w", ErrParseWithClaims)
	}

	return claims, nil
}
