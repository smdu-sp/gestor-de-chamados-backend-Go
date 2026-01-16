package router

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
)

// RegistrarRotasAuth registra as rotas de autenticação
func RegistrarRotasAuth(mux *http.ServeMux, handler *auth.AuthHandler) {
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/refresh", handler.Refresh)
}

// RegistrarRotaEu registra a rota para consulta dos dados do usuário autenticado
func RegistrarRotaEu(mux *http.ServeMux, handler *auth.AuthHandler) {
	mux.HandleFunc("/eu", handler.Me)
}