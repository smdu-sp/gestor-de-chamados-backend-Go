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

// NewTokenManager cria uma nova instância de TokenManager
// recebendo segredo, emissor e tempo de expiração.
func NewTokenManager(secret, issuer string, ttl time.Duration) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

// Generate cria um novo token JWT com base nos Claims recebidos.
func (t *TokenManager) Generate(c Claims) (string, error) {
	now := time.Now()

	// Preenche emissor, timestamps e expiração, caso não estejam definidos.
	if c.RegisteredClaims.Issuer == "" {
		c.RegisteredClaims.Issuer = t.issuer
	}
	c.RegisteredClaims.IssuedAt = jwt.NewNumericDate(now)
	c.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	c.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(t.ttl))

	// Gera um identificador único para o token (jti), se não existir.
	if c.RegisteredClaims.ID == "" {
		c.RegisteredClaims.ID = newJTI()
	}

	// Cria token usando algoritmo HS256 e retorna string assinada.
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return tok.SignedString(t.secret)
}

// Valida o algoritmo e retorna a chave secreta
func (t *TokenManager) keyFunc(token *jwt.Token) (any, error) {
	if token.Method != jwt.SigningMethodHS256 {
		return nil, errors.New("alg não permitido")
	}
	return t.secret, nil
}

// Parse valida e interpreta um token JWT.
// Retorna os claims se o token for válido, senão retorna erro.
func (t *TokenManager) Parse(tokenStr string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		t.keyFunc,
	)

	if err != nil {
		return nil, err
	}

	// Retorna os claims se o token for válido.
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

// newJTI gera um identificador único para cada token (jti).
// Usa 16 bytes aleatórios convertidos em hexadecimal.
func newJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Em caso de erro, retorne uma string vazia ou trate conforme necessário
		return ""
	}
	return hex.EncodeToString(b)
}
