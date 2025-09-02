package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
)

const (
	testeSecret    = "testesecret"
	refreshSecret  = "refreshsecret"
	expectedStatus = "expected status %d, got %d"
	shouldNotCallNextHandler = "não deveria chamar next handler"
)

// mockUserService implementa user.UserServiceInterface
type mockUserService struct {
	called    bool
	gotUserID string
	err       error
}

func (m *mockUserService) BuscarPorID(ctx context.Context, id string) (*user.Usuario, error) {
	return nil, nil
}
func (m *mockUserService) BuscarPorLogin(ctx context.Context, login string) (*user.Usuario, error) {
	return nil, nil
}
func (m *mockUserService) Criar(ctx context.Context, u *user.Usuario) error {
	return nil
}
func (m *mockUserService) AtualizarUltimoLogin(ctx context.Context, userID string) error {
	m.called = true
	m.gotUserID = userID
	return m.err
}

// --- Testes ---

func TestAuthenticateUserMissingAuthorization(t *testing.T) {
	manager := &jwt.Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}

	mockSvc := &mockUserService{}

	handler := AuthenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error(shouldNotCallNextHandler)
	}), manager, mockSvc)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf(expectedStatus, http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticateUserInvalidAuthorizationFormat(t *testing.T) {
	manager := &jwt.Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}

	mockSvc := &mockUserService{}

	handler := AuthenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error(shouldNotCallNextHandler)
	}), manager, mockSvc)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf(expectedStatus, http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticateUserInvalidToken(t *testing.T) {
	manager := &jwt.Manager{
		AccessSecret:  []byte(errors.New("invalid token").Error()),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}

	mockSvc := &mockUserService{}

	handler := AuthenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error(shouldNotCallNextHandler)
	}), manager, mockSvc)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf(expectedStatus, http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticateUserValidTokenUpdateLastLoginCalled(t *testing.T) {
	manager := &jwt.Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}

	mockSvc := &mockUserService{}

	claims := jwt.Claims{ID: "user123"}
	token, _ := manager.SignAccess(claims)

	handler := AuthenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := UserFromCtx(r)
		if c == nil || c.ID != "user123" {
			t.Error("claims não encontrados no contexto")
		}
		w.WriteHeader(http.StatusOK)
	}), manager, mockSvc)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(expectedStatus, http.StatusOK, rr.Code)
	}
	if !mockSvc.called || mockSvc.gotUserID != "user123" {
		t.Error("AtualizarUltimoLogin não chamado corretamente")
	}
}

func TestAuthenticateUserValidTokenUpdateLastLoginError(t *testing.T) {
	manager := &jwt.Manager{
		AccessSecret:  []byte(testeSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}

	mockSvc := &mockUserService{err: errors.New("erro db")}

	claims := jwt.Claims{ID: "user123"}
	token, _ := manager.SignAccess(claims)

	handler := AuthenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), manager, mockSvc)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(expectedStatus, http.StatusOK, rr.Code)
	}
}

func TestUserFromCtxReturnsClaimsSuccess(t *testing.T) {
	claims := &jwt.Claims{ID: "user456"}
	ctx := context.WithValue(context.Background(), UserKey, claims)
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

	got := UserFromCtx(req)
	if got == nil || got.ID != "user456" {
		t.Error("UserFromCtx não retornou claims corretos")
	}
}

func TestUserFromCtxReturnsNil(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	got := UserFromCtx(req)
	if got != nil {
		t.Error("UserFromCtx deveria retornar nil quando não há claims")
	}
}
