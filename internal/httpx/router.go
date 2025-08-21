package httpx

import "net/http"

// Middleware que garante que a requisição HTTP use apenas o método especificado.
func Method(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != method {
			w.Header().Set("Allow", method)
			http.Error(w, http.StatusText(
				http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed,
			)
			return
		}
		// Chama o handler original se o método estiver correto
		h(w, r)
	}
}

// Retorna um novo http.ServeMux
func NewMux() *http.ServeMux {
	return http.NewServeMux()
}
