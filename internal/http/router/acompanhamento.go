package router

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RegistrarRotasAcompanhamento registra as rotas de acompanhamento
func RegistrarRotasAcompanhamento(
	mux *http.ServeMux,
	handler *handler.AcompanhamentoHandler,
	aplicarPermissoes PermissoesHandler,
) {
	mux.Handle(
		"POST /acompanhamentos",
		aplicarPermissoes(handler.Criar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /acompanhamentos/buscar-por-id/{id}",
		aplicarPermissoes(handler.BuscarPorID, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /acompanhamentos/listar-paginado",
		aplicarPermissoes(handler.ListarPaginado, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"PATCH /acompanhamentos/atualizar/{id}",
		aplicarPermissoes(handler.Atualizar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"DELETE /acompanhamentos/deletar/{id}",
		aplicarPermissoes(handler.Deletar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /acompanhamentos/buscar-por-chamado-id/{id}",
		aplicarPermissoes(handler.BuscarPorChamadoID, "ADM", "TEC", "USR", "DEV"),
	)
}
