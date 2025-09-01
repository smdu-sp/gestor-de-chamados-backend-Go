package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
)

const (
	testeSecret   = "testesecret"
	refreshSecret = "refreshsecret"
	errorMessage  = "should not call next handler"
	expectedStatus = "expected status %d, got %d"
)

// TestWithUserMissingAuthorization verifica se a ausência do cabeçalho Authorization resulta em um erro 401 Unauthorized.
func TestWithUserMissingAuthorization(t *testing.T) {
	// Arrange
	manager := &jwt.Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

	// Act
	handler := WithUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error(errorMessage)
	}), manager, nil)

	// Assert
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf(expectedStatus, http.StatusUnauthorized, rr.Code)
	}
}

// TestWithUserInvalidAuthorizationFormat verifica se um formato de autorização inválido resulta em um erro 401 Unauthorized.
func TestWithUserInvalidAuthorizationFormat(t *testing.T) {
    // Arrange
	manager := &jwt.Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

    // Act
	handler := WithUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error(errorMessage)
	}), manager, nil)

    // Assert
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf(expectedStatus, http.StatusUnauthorized, rr.Code)
	}
}

// TestWithUserInvalidToken verifica se um token inválido resulta em um erro 401 Unauthorized.
func TestWithUserInvalidToken(t *testing.T) {
    // Arrange
	manager := &jwt.Manager{
		AccessSecret:  []byte(errors.New("invalid token").Error()),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}

    // Act
	handler := WithUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error(errorMessage)
	}), manager, nil)

    // Assert
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf(expectedStatus, http.StatusUnauthorized, rr.Code)
	}
}

// TestWithUserValidTokenUpdateLastLoginCalled verifica se a função de atualização do último login é chamada com o ID correto do usuário.
func TestWithUserValidTokenUpdateLastLoginCalledSuccess(t *testing.T) {
    // Arrange
	called := false
	var gotUserID string
	updateFn := func(ctx context.Context, userID string) error {
		called = true
		gotUserID = userID
		return nil
	}

	manager := &jwt.Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}
	claims := jwt.Claims{ID: "user123"}
	token, _ := manager.SignAccess(claims)

    // Act
	handler := WithUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := UserFromCtx(r)
		if c == nil || c.ID != "user123" {
			t.Error("claims not found in context")
		}
		w.WriteHeader(http.StatusOK)
	}), manager, updateFn)

    // Assert
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if !called || gotUserID != "user123" {
		t.Error("updateLastLogin not called with correct userID")
	}
}

// TestWithUserValidTokenUpdateLastLoginError verifica se um erro na atualização do último login resulta em um erro 500 Internal Server Error.
func TestWithUserValidTokenUpdateLastLoginError(t *testing.T) {
    // Arrange
	updateFn := func(ctx context.Context, userID string) error {
		return errors.New("db error")
	}

	manager := &jwt.Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     time.Minute * 15,
		RefreshTTL:    time.Hour * 24 * 7,
	}
	claims := jwt.Claims{ID: "user123"}
	token, _ := manager.SignAccess(claims)

    // Act
	handler := WithUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), manager, updateFn)

    // Assert
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(expectedStatus, http.StatusOK, rr.Code)
	}
}

// TestUserFromCtxReturnsClaims verifica se a função UserFromCtx retorna os claims corretos.
func TestUserFromCtxReturnsClaimsSuccess(t *testing.T) {
    // Arrange
	claims := &jwt.Claims{ID: "user456"}
	ctx := context.WithValue(context.Background(), UserKey, claims)
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

    // Act
	got := UserFromCtx(req)

    // Assert
	if got == nil || got.ID != "user456" {
		t.Error("UserFromCtx did not return correct claims")
	}
}

// TestUserFromCtxReturnsNil verifica se a função UserFromCtx retorna nil quando não há claims no contexto.
func TestUserFromCtxReturnsNil(t *testing.T) {
    // Arrange
	req := httptest.NewRequest("GET", "/", nil)
    // Act
	got := UserFromCtx(req)
    // Assert
	if got != nil {
		t.Error("UserFromCtx should return nil when no claims in context")
	}
}
