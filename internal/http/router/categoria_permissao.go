package router

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RegistrarRotasCategoriaPermissao registra as rotas de categoria-permissão
func RegistrarRotasCategoriaPermissao(
	mux *http.ServeMux,
	handler *handler.CategoriaPermissaoHandler,
	aplicarPermissoes PermissoesHandler,
) {
	mux.Handle(
		"POST /categoria-permissoes/criar",
		aplicarPermissoes(handler.Criar, "ADM", "DEV"),
	)

	mux.Handle(
		"GET /categorias-permissoes/{categoriaId}/usuarios/{usuarioId}",
		aplicarPermissoes(handler.BuscarPorID, "ADM", "TEC", "USR", "DEV"),
	)

	mux.Handle(
		"PATCH /categorias-permissoes/atualizar/{categoriaId}/usuarios/{usuarioId}",
		aplicarPermissoes(handler.Atualizar, "ADM", "DEV"),
	)

	mux.Handle(
		"GET /categoria-permissoes/listar-paginado",
		aplicarPermissoes(handler.ListarPaginado, "ADM", "TEC", "DEV"),
	)

	mux.Handle(
		"DELETE /categorias-permissoes/deletar/{categoriaId}/usuarios/{usuarioId}",
		aplicarPermissoes(handler.Deletar, "ADM", "DEV"),
	)
}
