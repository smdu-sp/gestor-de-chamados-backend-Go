package router

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	mid "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/provider/ldap"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/handler"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/middleware"
	auth "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/usecase"
	uc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/usecase"
)

var prefixosPublicos = []string{
	"/health",
	"/login",
	"/refresh",
	"/swagger",
}

// InicializarRoteadorHTTP configura e retorna o roteador HTTP da aplicação
func InicializarRoteadorHTTP(cfg config.Config, db *sql.DB) http.Handler {
	// Injeção de dependências:

	// Repositório e casos de uso de usuários
	usuarioRepository := repository.NewMySQLUsuarioRepository(db)
	usuarioUsecase := uc.NewUsuarioUsecase(usuarioRepository)

	// Repositório e caso de uso de chamados
	chamadoRepository := repository.NewMySQLChamadoRepository(db)
	chamadoUsecase := uc.NewChamadoUsecase(chamadoRepository)

	// Repositório e caso de uso de categorias
	categoriaRepository := repository.NewMySQLCategoriaRepository(db)
	categoriaUsecase := uc.NewCategoriaUsecase(categoriaRepository)

	// Repositório e caso de uso de subcategorias
	subcategoriaRepository := repository.NewMySQLSubcategoriaRepository(db)
	subcategoriaUsecase := uc.NewSubcategoriaUsecase(subcategoriaRepository)

	// Gerenciador JWT
	gerenteJWT := &jwt.GerenteJWT{
		ChaveAcesso:  []byte(cfg.JWTSecret),
		ChaveRefresh: []byte(cfg.RTSecret),
		TLLAcesso:    converterDuracao(cfg.AccessTTL),
		TLLRefresh:   converterDuracao(cfg.RefreshTTL),
	}

	// Cliente LDAP
	clienteLDAP := &ldap.Client{
		Server:    cfg.LDAPServer,
		Domain:    cfg.LDAPDomain,
		Base:      cfg.LDAPBase,
		User:      cfg.LDAPUser,
		Pass:      cfg.LDAPPass,
		LoginAttr: cfg.LDAPLoginAttr,
	}

	// Caso de uso de autenticação
	authUsecase := auth.NewAuthInternoUsecase(
		usuarioUsecase,
		gerenteJWT,
		clienteLDAP,
		cfg,
	)

	// Handlers
	AuthHandler := &handler.AuthHandler{Usecase: authUsecase}
	usuarioHandler := &handler.UsuarioHandler{UsecaseUsr: usuarioUsecase, UsecaseAuth: authUsecase, UsecaseLDAP: clienteLDAP,}
	chamadoHandler := &handler.ChamadoHandler{Usecase: chamadoUsecase}
	categoriaHandler := &handler.CategoriaHandler{Usecase: categoriaUsecase}
	subcategoriaHandler := &handler.SubcategoriaHandler{Usecase: subcategoriaUsecase}

	// Rotas públicas
	publico := http.NewServeMux()
	SwaggerRegistrarRotas(publico)
	HealthCheckRegistrarRotas(publico, db)
	AuthRegistrarRotas(publico, AuthHandler)

	// Rotas protegidas
	muxProtegido := http.NewServeMux()
	muxProtegido.HandleFunc("/eu", AuthHandler.Me)
	UsuarioRegistrarRotas(muxProtegido, usuarioHandler, gerenteJWT, usuarioUsecase)
	ChamadoRegistrarRotas(muxProtegido, chamadoHandler, gerenteJWT, usuarioUsecase)
	CategoriaRegistrarRotas(muxProtegido, categoriaHandler, gerenteJWT, usuarioUsecase)
	SubcategoriaRegistrarRotas(muxProtegido, subcategoriaHandler, gerenteJWT, usuarioUsecase)

	// Roteador principal com CORS
	rotas := CriarRoteadorAutenticacao(publico, muxProtegido, gerenteJWT, usuarioUsecase)
	rotas = middleware.CORS(cfg.CORSOrigin)(rotas)

	log.Println("CORS liberado para:", cfg.CORSOrigin)
	return rotas
}

// CriarRoteadorAutenticacao cria um roteador que diferencia rotas públicas de protegidas com autenticação
func CriarRoteadorAutenticacao(publico, protegido http.Handler, gerenteJWT jwt.JWTUsecase, usrUsecase *uc.UsuarioUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Rotas públicas
		for _, prefixo := range prefixosPublicos {
			if strings.HasPrefix(r.URL.Path, prefixo) {
				publico.ServeHTTP(w, r)
				return
			}
		}

		// Rotas protegidas com middleware de autenticação
		protegidoAuth := mid.AutenticarUsuario(protegido, gerenteJWT, usrUsecase)
		protegidoAuth.ServeHTTP(w, r)
	})
}

// converterDuracao converte string em time.Duration
func converterDuracao(d string) time.Duration {
	t, _ := time.ParseDuration(d)
	return t
}
