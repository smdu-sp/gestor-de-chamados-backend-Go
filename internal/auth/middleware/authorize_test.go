package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
)

func TestRequirePermissionsUnauthorized(t *testing.T) {
	handler := RequirePermissions("ADM")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Não deveria chamar o handler quando claims são nil")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Esperado status 401, recebeu %d", rr.Code)
	}
}

func TestRequirePermissionsForbidden(t *testing.T) {
	claims := &jwt.Claims{Permissao: "USR"}
	ctx := context.WithValue(context.Background(), UserKey, claims)

	handler := RequirePermissions("ADM")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Não deveria chamar o handler quando permissão não bate")
	}))

	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Esperado status 403, recebeu %d", rr.Code)
	}
}

func TestRequirePermissionsAllowed(t *testing.T) {
	claims := &jwt.Claims{Permissao: "ADM"}
	ctx := context.WithValue(context.Background(), UserKey, claims)

	called := false
	handler := RequirePermissions("ADM", "TEC")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("Handler não foi chamado apesar de permissão correta")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("Esperado status 200, recebeu %d", rr.Code)
	}
}
