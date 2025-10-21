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

var prefixosPublicos = [4]string{
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

	// Repositório e caso de uso de logs
	logRepository := repository.NewMySQLLogRepository(db)
	logUsecase := uc.NewLogUsecase(logRepository)

	// Repositório e caso de uso de acompanhamentos
	acompanhamentoRepository := repository.NewMySQLAcompanhamentoRepository(db)
	acompanhamentoUsecase := uc.NewAcompanhamentoUsecase(acompanhamentoRepository)

	// Repositório e caso de uso de atendimentos
	atendimentoRepository := repository.NewMySQLAtendimentoRepository(db)
	atendimentoUsecase := uc.NewAtendimentoUsecase(atendimentoRepository)

	// Repositório e caso de uso de categoriaPermissão
	categoriaPermissaoRepository := repository.NewMySQLCategoriaPermissaoRepository(db)
	categoriaPermissaoUsecase := uc.NewCategoriaPermissaoUsecase(categoriaPermissaoRepository)

	// Gerenciador JWT
	gerenteJWT := jwt.NewGerenteJWT(
		[]byte(cfg.JWTSecret),
		[]byte(cfg.RTSecret),
		converterDuracao(cfg.AccessTTL),
		converterDuracao(cfg.RefreshTTL),
	)

	// Cliente LDAP
	clienteLDAP := ldap.NewClienteLDAP(
		cfg.LDAPServer,
		cfg.LDAPDomain,
		cfg.LDAPBase,
		cfg.LDAPUser,
		cfg.LDAPPass,
		cfg.LDAPLoginAttr,
	)

	// Caso de uso de autenticação
	authUsecase := auth.NewAuthInternoUsecase(usuarioUsecase, gerenteJWT, clienteLDAP, logUsecase, cfg)

	// Handlers
	AuthHandler := handler.NewAuthHandler(authUsecase)
	usuarioHandler := handler.NewUsuarioHandler(usuarioUsecase, authUsecase, clienteLDAP, logUsecase)
	chamadoHandler := handler.NewChamadoHandler(chamadoUsecase, logUsecase)
	categoriaHandler := handler.NewCategoriaHandler(categoriaUsecase, logUsecase)
	subcategoriaHandler := handler.NewSubcategoriaHandler(subcategoriaUsecase, logUsecase)
	logHandler := handler.NewLogHandler(logUsecase)
	acompanhamentoHandler := handler.NewAcompanhamentoHandler(acompanhamentoUsecase, logUsecase)
	atendimentoHandler := handler.NewAtendimentoHandler(atendimentoUsecase, logUsecase)
	categoriaPermissaoHandler := handler.NewCategoriaPermissaoHandler(categoriaPermissaoUsecase, logUsecase)

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
	LogRegistrarRotas(muxProtegido, logHandler, gerenteJWT, usuarioUsecase)
	AcompanhamentoRegistrarRotas(muxProtegido, acompanhamentoHandler, gerenteJWT, usuarioUsecase)
	AtendimentoRegistrarRotas(muxProtegido, atendimentoHandler, gerenteJWT, usuarioUsecase)
	CategoriaPermissaoRegistrarRotas(muxProtegido, categoriaPermissaoHandler, gerenteJWT, usuarioUsecase)

	// Roteador principal com CORS
	rotas := CriarRoteadorAutenticacao(publico, muxProtegido, gerenteJWT, usuarioUsecase)
	rotas = middleware.CORS(cfg.CORSOrigin)(rotas)
	rotas = middleware.RecuperarDePanico(rotas)

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
