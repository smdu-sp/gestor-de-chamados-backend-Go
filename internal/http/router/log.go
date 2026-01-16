package router

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RegistrarRotasLog registra as rotas de log
func RegistrarRotasLog(
	mux *http.ServeMux,
	handler *handler.LogHandler,
	aplicarPermissoes PermissoesHandler,
) {
	mux.Handle("GET /logs/buscar-por-id/{id}",
		aplicarPermissoes(handler.BuscarPorID, "DEV"))

	mux.Handle("GET /logs/listar-paginado",
		aplicarPermissoes(handler.ListarPaginado, "DEV"))
}
