package auth

import (
	"net/http"
	"slices"
	"strings"
)

// RequerPermissoes libera se o usuário tiver QUALQUER uma das permissões
func RequerPermissoes(permissoes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Chama o próximo handler se autorizado
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtém as claims do usuário
			claims, err := ClaimsdoContexto(r.Context())
			if err != nil {
				ResponderJSONUnauthorized(w, r, ErrUsuarioNaoAutenticado.Error())
				return
			}

			// Normaliza a permissão do usuário
			permissaoUsuario := strings.ToLower(strings.TrimSpace(claims.Permissao))
			ok := slices.ContainsFunc(permissoes, func(p string) bool {
				return permissaoUsuario == strings.ToLower(strings.TrimSpace(p))
			})

			if !ok {
				ResponderJSONForbidden(w, r, ErrUsuarioNaoAutorizado.Error())
				return
			}

			// Chama o próximo handler
			next.ServeHTTP(w, r)
		})
	}
}
