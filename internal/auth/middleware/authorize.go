package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
)

// RequirePermissions libera se o usuário tiver QUALQUER uma das permissões
func RequirePermissions(perms ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtém as claims do usuário
			claims := UserFromCtx(r)
			if claims == nil {
				response.ErrorJSON(w, http.StatusUnauthorized, "unauthorized", "Você não está autorizado a acessar este recurso")
				return
			}

			// Normaliza a permissão do usuário
			userPerm := strings.ToLower(strings.TrimSpace(claims.Permissao))
			ok := slices.ContainsFunc(perms, func(p string) bool {
				return userPerm == strings.ToLower(strings.TrimSpace(p))
			})

			if !ok {
				response.ErrorJSON(w, http.StatusForbidden, "forbidden", "Você não tem permissão para acessar este recurso")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
