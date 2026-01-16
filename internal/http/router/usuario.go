package router

import (
	"net/http"

	hdl "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
)

// RegistrarRotasUsuario registra as rotas de usuário com permissões.
func RegistrarRotasUsuario(
	mux *http.ServeMux,
	handler *hdl.UsuarioHandler,
	aplicarPermissoes PermissoesHandler,
) {
	mux.HandleFunc(
		"POST /usuarios",
		aplicarPermissoes(handler.Criar, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"GET /usuarios/listar-paginado",
		aplicarPermissoes(handler.ListarPaginado, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"GET /usuarios/buscar-por-id/{id}",
		aplicarPermissoes(handler.BuscarPorID, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"PATCH /usuarios/atualizar/{id}",
		aplicarPermissoes(handler.Atualizar, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"PATCH /usuarios/atualizar-permissao/{id}",
		aplicarPermissoes(handler.AtualizarPermissao, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"GET /usuarios/listar",
		aplicarPermissoes(handler.Listar, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"GET /usuarios/buscar-tecnicos",
		aplicarPermissoes(handler.BuscarTecnicos, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"DELETE /usuarios/desativar/{id}",
		aplicarPermissoes(handler.Desativar, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"PATCH /usuarios/autorizar/{id}",
		aplicarPermissoes(handler.Autorizar, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"GET /usuarios/buscar-novo/{login}",
		aplicarPermissoes(handler.BuscarNovo, "ADM", "DEV"),
	)

	mux.HandleFunc(
		"PATCH /usuarios/ativar/{id}",
		aplicarPermissoes(handler.Ativar, "ADM", "DEV"),
	)
	
	mux.HandleFunc(
		"GET /usuarios/validar-usuario",
		handler.ValidarUsuario,
	)
}
