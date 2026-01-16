package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RecuperarDePanico é um middleware que recupera de pânicos em handlers HTTP.
func RecuperarDePanico(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error(
					"PANIC recuperado",
					"erro", err,
					"stack", string(debug.Stack()),
					"path", r.URL.Path,
				)
				handler.ResponderJSONRecoveryPanic(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
