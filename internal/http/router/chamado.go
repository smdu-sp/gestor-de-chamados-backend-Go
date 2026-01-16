package router

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RegistrarRotasChamado registra as rotas de chamado
func RegistrarRotasChamado(
	mux *http.ServeMux,
	handler *handler.ChamadoHandler,
	aplicarPermissoes PermissoesHandler,
) {
	mux.Handle(
		"POST /chamados",
		aplicarPermissoes(handler.Criar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"PATCH /chamados/atualizar/{id}",
		aplicarPermissoes(handler.Atualizar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /chamados/buscar-por-id/{id}",
		aplicarPermissoes(handler.BuscarPorID, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /chamados/listar-paginado",
		aplicarPermissoes(handler.ListarPaginado, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /chamados/listar",
		aplicarPermissoes(handler.Listar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"PATCH /chamados/atualizar-status/{id}",
		aplicarPermissoes(handler.AtualizarStatus, "ADM", "TEC", "DEV"),
	)

	mux.Handle(
		"PATCH /chamados/arquivar/{id}",
		aplicarPermissoes(handler.Arquivar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"PATCH /chamados/desarquivar/{id}",
		aplicarPermissoes(handler.Desarquivar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"PATCH /chamados/atualizar-solucao/{id}",
		aplicarPermissoes(handler.AtualizarSolucao, "ADM", "TEC", "DEV"),
	)
}
