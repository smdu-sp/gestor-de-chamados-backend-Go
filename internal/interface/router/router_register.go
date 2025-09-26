package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/handler"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
	goSwagger "github.com/swaggo/http-swagger"
)

// SwaggerRegistrarRotas registra as rotas do Swagger
func SwaggerRegistrarRotas(mux *http.ServeMux) {
	mux.HandleFunc("/swagger/", goSwagger.WrapHandler)
}

// HealthCheckRegistrarRotas registra a rota de health check
func HealthCheckRegistrarRotas(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		resp := response.HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
		}

		// Verifica conexão com o banco de dados
		if err := db.Ping(); err != nil {
			resp.Status = "fail"
			resp.DB = "fail"
		} else {
			resp.DB = "ok"
		}

		w.Header().Set("Content-Type", "application/json")

		if resp.Status == "fail" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}

// AuthRegistrarRotas registra as rotas de autenticação
func AuthRegistrarRotas(mux *http.ServeMux, authH *handler.AuthHandler) {
	mux.HandleFunc("/login", authH.Login)
	mux.HandleFunc("/refresh", authH.Refresh)
}

// UsuarioRegistrarRotas registra as rotas de usuário
func UsuarioRegistrarRotas(mux *http.ServeMux, usrH *handler.UsuarioHandler, jwtManager *jwt.GerenteJWT, svc usecase.UsuarioUsecase) {
	// helper para aplicar autenticação + permissões
	aplicarPermissoes := func(handler http.HandlerFunc, perms ...string) http.Handler {
		return middleware.AutenticarUsuario(
			middleware.RequerPermissoes(perms...)(handler),
			jwtManager, svc,
		)
	}

	mux.Handle("/usuarios/criar", aplicarPermissoes(usrH.Criar, "ADM"))
	mux.Handle("/usuarios/buscar-tudo", aplicarPermissoes(usrH.BuscarTudo, "ADM"))
	mux.Handle("/usuarios/buscar-por-id/", aplicarPermissoes(usrH.BuscarPorID, "ADM"))
	mux.Handle("/usuarios/atualizar/", aplicarPermissoes(usrH.Atualizar, "ADM"))
	mux.Handle("/usuarios/atualizar-permissao/", aplicarPermissoes(usrH.AtualizarPermissao, "ADM"))
	mux.Handle("/usuarios/lista-completa", aplicarPermissoes(usrH.ListaCompleta, "ADM"))
	mux.Handle("/usuarios/buscar-tecnicos", aplicarPermissoes(usrH.BuscarTecnicos, "ADM"))
	mux.Handle("/usuarios/desativar/", aplicarPermissoes(usrH.Desativar, "ADM"))
	mux.Handle("/usuarios/autorizar/", aplicarPermissoes(usrH.Autorizar, "ADM"))
	mux.Handle("/usuarios/buscar-novo/", aplicarPermissoes(usrH.BuscarNovo, "ADM"))
	mux.Handle("/usuarios/ativar/", aplicarPermissoes(usrH.Ativar, "ADM"))
	mux.HandleFunc("/usuarios/valida-usuario", usrH.ValidaUsuario) // não precisa de permissão ADM
}

// ChamadoRegistrarRotas registra as rotas de chamado
func ChamadoRegistrarRotas(mux *http.ServeMux, chmH *handler.ChamadoHandler, jwtManager *jwt.GerenteJWT, svc usecase.UsuarioUsecase) {
	// helper para aplicar autenticação + permissões
	aplicarPermissoes := func(handler http.HandlerFunc, perms ...string) http.Handler {
		return middleware.AutenticarUsuario(
			middleware.RequerPermissoes(perms...)(handler),
			jwtManager, svc,
		)
	}

	mux.Handle("/chamados/criar", aplicarPermissoes(chmH.Criar, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/atualizar/", aplicarPermissoes(chmH.Atualizar, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/buscar-por-id/", aplicarPermissoes(chmH.BuscarPorID, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/buscar-tudo", aplicarPermissoes(chmH.BuscarTudo, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/lista-completa", aplicarPermissoes(chmH.ListaCompleta, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/atualizar-status/", aplicarPermissoes(chmH.AtualizarStatus, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/arquivar/", aplicarPermissoes(chmH.Arquivar, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/desarquivar/", aplicarPermissoes(chmH.Desarquivar, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/atribuir-tecnico/", aplicarPermissoes(chmH.AtribuirTecnico, "ADM", "TEC", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/chamados/remover-tecnico/", aplicarPermissoes(chmH.RemoverTecnicoChamado, "ADM", "TEC", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
}

// CategoriaRegistrarRotas registra as rotas de categoria
func CategoriaRegistrarRotas(mux *http.ServeMux, catH *handler.CategoriaHandler, jwtManager *jwt.GerenteJWT, svc usecase.UsuarioUsecase) {
	// helper para aplicar autenticação + permissões
	aplicarPermissoes := func(handler http.HandlerFunc, perms ...string) http.Handler {
		return middleware.AutenticarUsuario(
			middleware.RequerPermissoes(perms...)(handler),
			jwtManager, svc,
		)
	}

	mux.Handle("/categorias/criar", aplicarPermissoes(catH.Criar, "ADM"))
	mux.Handle("/categorias/atualizar/", aplicarPermissoes(catH.Atualizar, "ADM"))
	mux.Handle("/categorias/buscar-por-id/", aplicarPermissoes(catH.BuscarPorID, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/categorias/buscar-tudo", aplicarPermissoes(catH.BuscarTudo, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/categorias/lista-completa", aplicarPermissoes(catH.ListaCompleta, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/categorias/desativar/", aplicarPermissoes(catH.Desativar, "ADM"))
	mux.Handle("/categorias/ativar/", aplicarPermissoes(catH.Ativar, "ADM"))
}

// SubcategoriaRegistrarRotas registra as rotas de subcategoria
func SubcategoriaRegistrarRotas(mux *http.ServeMux, subcatH *handler.SubcategoriaHandler, jwtManager *jwt.GerenteJWT, svc usecase.UsuarioUsecase) {
	// helper para aplicar autenticação + permissões
	aplicarPermissoes := func(handler http.HandlerFunc, perms ...string) http.Handler {
		return middleware.AutenticarUsuario(
			middleware.RequerPermissoes(perms...)(handler),
			jwtManager, svc,
		)
	}

	mux.Handle("/subcategorias/criar", aplicarPermissoes(subcatH.Criar, "ADM"))
	mux.Handle("/subcategorias/atualizar/", aplicarPermissoes(subcatH.Atualizar, "ADM"))
	mux.Handle("/subcategorias/buscar-por-id/", aplicarPermissoes(subcatH.BuscarPorID, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/subcategorias/buscar-tudo", aplicarPermissoes(subcatH.BuscarTudo, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/subcategorias/lista-completa", aplicarPermissoes(subcatH.ListaCompleta, "ADM", "TEC", "USR", "CAD", "DEV", "SUP", "INF", "VOIP", "IMP"))
	mux.Handle("/subcategorias/desativar/", aplicarPermissoes(subcatH.Desativar, "ADM"))
	mux.Handle("/subcategorias/ativar/", aplicarPermissoes(subcatH.Ativar, "ADM"))
}