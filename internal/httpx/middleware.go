package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"
)

type ctxKey string

const ctxClaims ctxKey = "claims"

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

// Autenticação por JWT (Bearer <token>)
func AuthMiddleware(tm *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authz := r.Header.Get("Authorization")

			if authz == "" {
				Error(w, http.StatusUnauthorized, "token ausente")
				return
			}

			parts := strings.SplitN(authz, " ", 2)

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				Error(w, http.StatusUnauthorized, "formato de autorização inválido")
				return
			}

			claims, err := tm.Parse(parts[1])

			if err != nil {
				Error(w, http.StatusUnauthorized, "token inválido")
				return
			}

			ctx := context.WithValue(r.Context(), ctxClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RBAC: exige um dos papéis
func RequireRoles(allowed ...model.Permissao) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := r.Context().Value(ctxClaims)

			if v == nil {
				Error(w, http.StatusUnauthorized, "não autenticado")
				return
			}

			c := v.(*auth.Claims)

			if !auth.IsAllowed(c.Permissao, allowed...) {
				Error(w, http.StatusForbidden, "acesso negado")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Recover + Timeout
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {

			if rec := recover(); rec != nil {
				Error(w,
					http.StatusInternalServerError, "erro interno")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, "request timeout")
	}
}

// Chain simples de middlewares (no framework!)
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Helper para obter Claims no handler
func ClaimsFrom(r *http.Request) (*auth.Claims, error) {
	v := r.Context().Value(ctxClaims)

	if v == nil {
		return nil, errors.New("sem claims")
	}

	c, ok := v.(*auth.Claims)
	
	if !ok {
		return nil, errors.New("claims inválidas")
	}
	return c, nil
}
