package httpx

import "net/http"

// Método correto? Se não, 405
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
		h(w, r)
	}
}

// ServeMux nativo, sem framework
func NewMux() *http.ServeMux { return http.NewServeMux() }
