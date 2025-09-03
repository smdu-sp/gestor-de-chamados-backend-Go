package jwt

import (
	"errors"
	"fmt"
	"time"

	goJwt "github.com/golang-jwt/jwt/v5"
)

// JWTInterface define os métodos que a implementação JWT deve fornecer
type JWTInterface interface {
	SignAccess(c Claims) (string, error)
	SignRefresh(c Claims) (string, error)
	ParseRefresh(token string) (*Claims, error)
}

type Manager struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Claims struct {
	ID        string `json:"sub"`
	Login     string `json:"login"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Permissao string `json:"permissao"`
	goJwt.RegisteredClaims
}

// SignAccess gera um token de acesso
func (m *Manager) SignAccess(c Claims) (string, error) {
	return m.sign(c, m.AccessSecret, m.AccessTTL)
}

// SignRefresh gera um token de refresh
func (m *Manager) SignRefresh(c Claims) (string, error) {
	return m.sign(c, m.RefreshSecret, m.RefreshTTL)
}

// ParseAccess valida um token de acesso
func (m *Manager) ParseAccess(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, m.AccessSecret)
}

// ParseRefresh valida um token de refresh
func (m *Manager) ParseRefresh(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, m.RefreshSecret)
}

// sign é uma função helper interna para gerar token
func (m *Manager) sign(c Claims, secret []byte, ttl time.Duration) (string, error) {
	c.RegisteredClaims.ExpiresAt = goJwt.NewNumericDate(time.Now().Add(ttl))
	token := goJwt.NewWithClaims(goJwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

// parse é uma função helper interna para validar token
func (m *Manager) parse(tokenStr string, secret []byte) (*Claims, error) {
	// Parse o token e valide a assinatura
	parsed, err := goJwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *goJwt.Token) (any, error) {
			return secret, nil
		})

	// Verifica se houve erro na validação
	if err != nil {
		if errors.Is(err, goJwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expirado: %w", err)
		}
		if errors.Is(err, goJwt.ErrTokenSignatureInvalid) {
			return nil, fmt.Errorf("assinatura inválida: %w", err)
		}
		return nil, fmt.Errorf("erro ao validar token: %w", err)
	}

	// Verifica se as claims são válidas
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("claims inválidos: %w", goJwt.ErrTokenInvalidClaims)
}
