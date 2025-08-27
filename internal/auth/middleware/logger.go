package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger é middleware que loga cada requisição HTTP
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(lw, r)

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lw.status, time.Since(start))
	})
}

// responseWriter wrapa http.ResponseWriter para capturar status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
