package middleware

import (
	"encoding/json"
	"net/http"
)

// RequirePermissions cria middleware que verifica se o usuário possui pelo menos uma das permissões
func RequirePermissions(perms ...string) func(http.Handler) http.Handler {
	set := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		set[p] = struct{}{}
	}

	// Retorna middleware
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := UserFromCtx(r)
			if claims == nil {
				writeJSONError(w, http.StatusUnauthorized, "usuário não autenticado")
				return
			}

			if len(set) == 0 || hasPermission(claims.Permissao, set) {
				next.ServeHTTP(w, r)
				return
			}

			writeJSONError(w, http.StatusForbidden, "permissão insuficiente")
		})
	}
}

// hasPermission verifica se perm existe no conjunto
func hasPermission(perm string, set map[string]struct{}) bool {
	_, ok := set[perm]
	return ok
}

// writeJSONError envia resposta JSON com mensagem de erro
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
