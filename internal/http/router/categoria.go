package router

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RegistrarRotasCategoria registra as rotas de categoria
func RegistrarRotasCategoria(
	mux *http.ServeMux,
	handler *handler.CategoriaHandler,
	aplicarPermissoes PermissoesHandler,
) {
	mux.Handle(
		"POST /categorias",
		aplicarPermissoes(handler.Criar, "ADM", "DEV"),
	)

	mux.Handle(
		"PATCH /categorias/atualizar/{id}",
		aplicarPermissoes(handler.Atualizar, "ADM", "DEV"),
	)

	mux.Handle(
		"GET /categorias/buscar-por-id/{id}",
		aplicarPermissoes(handler.BuscarPorID, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /categorias/buscar-por-nome/{nome}",
		aplicarPermissoes(handler.BuscarPorNome, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /categorias/listar-paginado",
		aplicarPermissoes(handler.ListarPaginado, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"GET /categorias/listar",
		aplicarPermissoes(handler.Listar, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"DELETE /categorias/desativar/{id}",
		aplicarPermissoes(handler.Desativar, "ADM", "DEV"),
	)

	mux.Handle(
		"PATCH /categorias/ativar/{id}",
		aplicarPermissoes(handler.Ativar, "ADM", "DEV"),
	)
}
