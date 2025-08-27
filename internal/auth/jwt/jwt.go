package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
	jwt.RegisteredClaims
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
	c.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

// parse é uma função helper interna para validar token
func (m *Manager) parse(tokenStr string, secret []byte) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if cl, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return cl, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
