package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
)

type ctxKey string

const UserKey ctxKey = "user"
const bearerPrefix string = "Bearer "

// UpdateLastLoginFn é uma função que atualiza o campo último login do usuário
type UpdateLastLoginFn func(ctx context.Context, userID string) error

// WithUser é um middleware que autentica o usuário a partir do token JWT presente no cabeçalho Authorization
// e adiciona as claims do usuário no contexto da requisição.
func WithUser(next http.Handler, jm *jwt.Manager, updateLastLogin UpdateLastLoginFn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		// verifica se o cabeçalho Authorization está presente
		if !strings.HasPrefix(auth, bearerPrefix) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// extrai o token do cabeçalho Authorization
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token := strings.TrimSpace(parts[1])

		// valida o token e extrai as claims
		claims, err := jm.ParseAccess(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// atualiza ultimoLogin
		if updateLastLogin != nil && claims.ID != "" {
			if err := updateLastLogin(r.Context(), claims.ID); err != nil {
				log.Printf("[DEBUG] erro ao atualizar último login: %v", err)
			}
		}

		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromCtx retorna claims do usuário presente no contexto
func UserFromCtx(r *http.Request) *jwt.Claims {
	if v := r.Context().Value(UserKey); v != nil {
		if claims, ok := v.(*jwt.Claims); ok {
			return claims
		}
	}
	return nil
}
