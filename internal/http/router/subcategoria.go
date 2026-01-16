package router

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RegistrarRotasSubcategoria registra as rotas de subcategoria
func RegistrarRotasSubcategoria(
	mux *http.ServeMux,
	handler *handler.SubcategoriaHandler,
	aplicarPermissoes PermissoesHandler,
) {
	mux.Handle(
		"POST /subcategorias",
		aplicarPermissoes(handler.Criar, "ADM", "DEV"),
	)

	mux.Handle(
		"PATCH /subcategorias/atualizar/{id}",
		aplicarPermissoes(handler.Atualizar, "ADM", "DEV"),
	)

	mux.Handle(
		"GET /subcategorias/buscar-por-id/{id}",
		aplicarPermissoes(handler.BuscarPorID, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /subcategorias/listar-paginado",
		aplicarPermissoes(handler.ListarPaginado, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /subcategorias/listar",
		aplicarPermissoes(handler.Listar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"DELETE /subcategorias/desativar/{id}",
		aplicarPermissoes(handler.Desativar, "ADM", "DEV"),
	)

	mux.Handle(
		"PATCH /subcategorias/ativar/{id}",
		aplicarPermissoes(handler.Ativar, "ADM", "DEV"),
	)
}
