package router

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	hdl "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/middleware"
)

var prefixosPublicos = []string{
	"/login",
	"/refresh",
	"/swagger",
	"/docs",
	"/health",
}

// PermissoesHandler é um tipo para funções que aplicam permissões a handlers HTTP.
type PermissoesHandler func(http.HandlerFunc, ...string) http.HandlerFunc

// Handlers agrupa todos os handlers HTTP da aplicação.
type Handlers struct {
	Log                *hdl.LogHandler
	Usuario            *hdl.UsuarioHandler
	Chamado            *hdl.ChamadoHandler
	Categoria          *hdl.CategoriaHandler
	Atendimento        *hdl.AtendimentoHandler
	Subcategoria       *hdl.SubcategoriaHandler
	Acompanhamento     *hdl.AcompanhamentoHandler
	CategoriaPermissao *hdl.CategoriaPermissaoHandler
	Auth               *auth.AuthHandler
}

// Router encapsula as dependências necessárias para roteamento.
type Router struct {
	handlers   *Handlers
	jwtSvc     auth.JWTService
	usuarioSvc usr.Service
	config     *cfg.Config
}

// NovoRouter cria um novo roteador com as dependências injetadas.
func NovoRouter(
	handlers *Handlers,
	jwtSvc auth.JWTService,
	usuarioSvc usr.Service,
	config *cfg.Config,
) *Router {
	return &Router{
		handlers:   handlers,
		jwtSvc:     jwtSvc,
		usuarioSvc: usuarioSvc,
		config:     config,
	}
}

// Configurar cria e configura o handler HTTP com todas as rotas e middlewares.
func (r *Router) Configurar() http.Handler {
	slog.Info("iniciando configuração do roteador HTTP")

	// 1 - Criar muxes separados para rotas públicas e protegidas
	muxPublico := http.NewServeMux()
	muxProtegido := http.NewServeMux()

	// 2 - Registrar rotas públicas e protegidas
	r.registrarRotasPublicas(muxPublico)
	r.registrarRotasProtegidas(muxProtegido)

	// 3 - Aplicar autenticação ao mux protegido
	muxProtegidoAutenticado := auth.AutenticarUsuario(muxProtegido, r.jwtSvc, r.usuarioSvc)

	// 4 - Criar roteador principal
	var roteador http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rotaEhPublica(r.URL.Path) {
			muxPublico.ServeHTTP(w, r)
			return
		}
		muxProtegidoAutenticado.ServeHTTP(w, r)
	})

	// 5 - Aplicar middlewares globais
	roteador = middleware.CORS(r.config.AppCORSOrigin())(roteador)
	roteador = middleware.RecuperarDePanico(roteador)

	slog.Info("roteador HTTP configurado com sucesso")
	return roteador
}

// registrarRotasPublicas registra todas as rotas públicas.
func (r *Router) registrarRotasPublicas(mux *http.ServeMux) {
	RegistrarRotasAuth(mux, r.handlers.Auth)
	RegistrarRotasSwagger(mux)
}

// registrarRotasProtegidas registra todas as rotas que requerem autenticação.
func (r *Router) registrarRotasProtegidas(mux *http.ServeMux) {

	// Registrar rotas protegidas com permissões
	RegistrarRotaEu(mux, r.handlers.Auth)
	RegistrarRotasLog(mux, r.handlers.Log, aplicarPermissoes)
	RegistrarRotasUsuario(mux, r.handlers.Usuario, aplicarPermissoes)
	RegistrarRotasChamado(mux, r.handlers.Chamado, aplicarPermissoes)
	RegistrarRotasCategoria(mux, r.handlers.Categoria, aplicarPermissoes)
	RegistrarRotasAtendimento(mux, r.handlers.Atendimento, aplicarPermissoes)
	RegistrarRotasSubcategoria(mux, r.handlers.Subcategoria, aplicarPermissoes)
	RegistrarRotasAcompanhamento(mux, r.handlers.Acompanhamento, aplicarPermissoes)
	RegistrarRotasCategoriaPermissao(mux, r.handlers.CategoriaPermissao, aplicarPermissoes)
}

// aplicarPermissoes é uma função auxiliar para aplicar permissões a um handler.
func aplicarPermissoes(handler http.HandlerFunc, perms ...string) http.HandlerFunc {
	middleware := auth.RequerPermissoes(perms...)
	handlerProtegido := middleware(handler)

	return handlerProtegido.ServeHTTP
}

// rotaEhPublica verifica se o caminho da rota é pública.
func rotaEhPublica(path string) bool {
	for _, prefixo := range prefixosPublicos {
		if path == prefixo || strings.HasPrefix(path, prefixo+"/") {
			return true
		}
	}
	return false
}
