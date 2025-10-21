package middleware

import (
		"log"
		"net/http"
		"runtime/debug"
)

// RecuperarDePanico é um middleware que recupera de pânicos em handlers HTTP.
func RecuperarDePanico(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("[PANIC] %v\n%s", err, debug.Stack())
                w.Header().Set("Content-Type", "application/json; charset=utf-8")
                http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
