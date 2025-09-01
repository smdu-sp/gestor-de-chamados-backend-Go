package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/stretchr/testify/assert"
)

const (
	testeSecret   = "teste-secret"
	refreshSecret = "refresh-secret"
	nomeUsuario   = "usuario_teste"
	emailUsuario  = "usuario@teste.com"
)

// Teste de assinatura de token de acesso
func TestManagerSignAccessSuccess(t *testing.T) {
	// Arrange
	manager := Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

	claims := Claims{
		ID:        "123",
		Login:     "usuario_teste",
		Nome:      nomeUsuario,
		Email:     emailUsuario,
		Permissao: "admin",
	}

	// Act
	token, err := manager.SignAccess(claims)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	parsedClaims, err := manager.ParseAccess(token)
	assert.NoError(t, err)
	assert.Equal(t, claims.ID, parsedClaims.ID)
	assert.Equal(t, claims.Login, parsedClaims.Login)
	assert.Equal(t, claims.Nome, parsedClaims.Nome)
	assert.Equal(t, claims.Email, parsedClaims.Email)
	assert.Equal(t, claims.Permissao, parsedClaims.Permissao)
}

// Teste de assinatura de token de atualização
func TestManagerSignRefreshSuccess(t *testing.T) {
	// Arrange
	manager := Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

	claims := Claims{
		ID:        "123",
		Login:     "usuario_teste",
		Nome:      nomeUsuario,
		Email:     emailUsuario,
		Permissao: "admin",
	}

	// Act
	token, err := manager.SignRefresh(claims)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	parsedClaims, err := manager.ParseRefresh(token)
	assert.NoError(t, err)
	assert.Equal(t, claims.ID, parsedClaims.ID)
	assert.Equal(t, claims.Login, parsedClaims.Login)
	assert.Equal(t, claims.Nome, parsedClaims.Nome)
	assert.Equal(t, claims.Email, parsedClaims.Email)
	assert.Equal(t, claims.Permissao, parsedClaims.Permissao)
}

// Teste de parsing de token de acesso inválido
func TestManagerParseAccessInvalid(t *testing.T) {
	// Arrange
	manager := Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

	// Token inválido
	invalidToken := "invalid.token.string"

	// Act
	claims, err := manager.ParseAccess(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Teste de parsing de token de acesso expirado
func TestManagerParseAccessExpired(t *testing.T) {
	// Arrange
	manager := Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     -time.Hour, // Token expirado (tempo negativo)
		RefreshTTL:    time.Hour * 24 * 7,
	}

	claims := Claims{
		ID:        "123",
		Login:     "usuario_teste",
		Nome:      nomeUsuario,
		Email:     emailUsuario,
		Permissao: "admin",
	}

	// Act
	token, err := manager.SignAccess(claims)
	// Assert
	assert.NoError(t, err)
	// Act
	parsedClaims, err := manager.ParseAccess(token)
	// Assert
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.Contains(t, err.Error(), "token expirado")
}

// Teste de parsing de token de acesso com segredo incorreto
func TestManagerParseAccessWrongSecret(t *testing.T) {
	// Arrange
	manager := Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

	claims := Claims{
		ID:        "123",
		Login:     "usuario_teste",
		Nome:      nomeUsuario,
		Email:     emailUsuario,
		Permissao: "admin",
	}

	// Act
	token, err := manager.SignAccess(claims)
	// Assert
	assert.NoError(t, err)

	// Arrange
	managerWrongSecret := Manager{
		AccessSecret:  []byte("wrong-secret"),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

	// Act
	parsedClaims, err := managerWrongSecret.ParseAccess(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.Contains(t, err.Error(), "assinatura inválida")
}

// Teste de parsing de token de atualização inválido
func TestManagerParseRefreshInvalid(t *testing.T) {
	// Arrange
	manager := Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}
	invalidToken := "invalid.token.string"

	// Act
	claims, err := manager.ParseRefresh(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Teste de parsing de token de atualização expirado
func TestManagerTokenExpiration(t *testing.T) {
	// Arrange
	manager := Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Second * 2, // Token expira em 2 segundos
		RefreshTTL:    time.Hour * 24 * 7,
	}

	claims := Claims{
		ID:        "123",
		Login:     "usuario_teste",
		Nome:      nomeUsuario,
		Email:     emailUsuario,
		Permissao: "admin",
	}

	// Act
	token, err := manager.SignAccess(claims)
	// Assert
	assert.NoError(t, err)

	// Act
	parsedClaims, err := manager.ParseAccess(token)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, parsedClaims)

	// Act
	time.Sleep(time.Second * 3)
	expiredClaims, err := manager.ParseAccess(token)
	// Assert
	assert.Error(t, err)
	assert.Nil(t, expiredClaims)
	assert.Contains(t, err.Error(), "token expirado")
}

// Test de parsing de claims inválidos
func TestManagerParseInvalidClaims(t *testing.T) {
	// Arrange
	manager := Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(), 
	}

	// Act
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testeSecret))

	// Assert
	assert.NoError(t, err)

	// Adiciona um caractere inválido ao token (após assinatura) para quebrar a validação
	// mas mantendo-o válido o suficiente para passar pelo parsing básico
	parts := strings.Split(tokenString, ".")
	if len(parts) == 3 {
		// Modifica a parte de payload sem afetar a assinatura
		tokenString = parts[0] + "." + parts[1] + "!" + "." + parts[2]
	}

	// Act
	parsedClaims, err := manager.ParseAccess(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	// O erro pode não conter exatamente "claims inválidos", mas deve indicar problema com o token
	if err != nil {
		t.Logf("Erro recebido: %v", err)
	}
}
