package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type Claims struct {
	UserID    string `json:"uid"`
	Nome      string `json:"name"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Permissao string `json:"perm"`
	jwt.RegisteredClaims
}

func NewTokenManager(secret, issuer string, ttl time.Duration) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (t *TokenManager) Generate(c Claims) (string, error) {
	now := time.Now()

	if c.RegisteredClaims.Issuer == "" {
		c.RegisteredClaims.Issuer = t.issuer
	}

	c.RegisteredClaims.IssuedAt = jwt.NewNumericDate(now)
	c.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	c.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(t.ttl))

	if c.RegisteredClaims.ID == "" {
		c.RegisteredClaims.ID = newJTI()
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return tok.SignedString(t.secret)
}

func (t *TokenManager) Parse(tokenStr string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("alg não permitido")
			}
			return t.secret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

func newJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Em caso de erro, retorne uma string vazia ou trate conforme necessário
		return ""
	}
	return hex.EncodeToString(b)
}
