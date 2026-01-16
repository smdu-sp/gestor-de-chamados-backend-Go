package router

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RegistrarRotasAtendimento registra as rotas de atendimento
func RegistrarRotasAtendimento(
	mux *http.ServeMux,
	handler *handler.AtendimentoHandler,
	aplicarPermissoes PermissoesHandler,
) {
	mux.Handle(
		"POST /atendimentos",
		aplicarPermissoes(handler.Criar, "ADM", "TEC", "DEV"),
	)

	mux.Handle(
		"GET /atendimentos/buscar-por-id/{id}",
		aplicarPermissoes(handler.BuscarPorID, "ADM", "TEC", "DEV"),
	)

	mux.Handle(
		"GET /atendimentos/buscar-por-ids/chamados/{chamadoId}/tecnicos/{tecnicoId}",
		aplicarPermissoes(handler.BuscarPorChamadoEAtribuidoID, "ADM", "TEC", "DEV"),
	)

	mux.Handle(
		"GET /atendimentos/listar-paginado",
		aplicarPermissoes(handler.ListarPaginado, "ADM", "TEC", "DEV"),
	)

	mux.Handle(
		"PATCH /atendimentos/atualizar/{id}",
		aplicarPermissoes(handler.Atualizar, "ADM", "TEC", "DEV"),
	)
}
