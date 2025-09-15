package router

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	auth "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/repository"
	handlers "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/handler"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/middleware"
	ldapauth "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/provider/ldap"

	authusecase "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/usecase/auth"
	userusecase "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/usecase/user"
)

var publicPrefixes = []string{
	"/health",
	"/login",
	"/refresh",
	"/swagger",
}

// Build constrói o handler principal da aplicação
func Build(cfg config.Config, db *sql.DB) http.Handler {
	// Repositório e casos de uso de usuários
	repo := repository.NewMySQLUserRepository(db)
	userSvc := userusecase.NewUserUsecase(repo)

	// JWT Manager
	jwtManager := &jwt.Manager{
		AccessSecret:  []byte(cfg.JWTSecret),
		RefreshSecret: []byte(cfg.RTSecret),
		AccessTTL:     parseDuration(cfg.AccessTTL),
		RefreshTTL:    parseDuration(cfg.RefreshTTL),
	}

ldapClient := &ldapauth.Client{
	Server:    cfg.LDAPServer,
	Domain:    cfg.LDAPDomain,
	Base:      cfg.LDAPBase,
	User:      cfg.LDAPUser,
	Pass:      cfg.LDAPPass,
	LoginAttr: cfg.LDAPLoginAttr,
}

// Caso de uso de autenticação
authSvc := authusecase.NewAuthUsecase(userSvc, jwtManager, ldapClient, cfg)

// Handlers
authH := &handlers.AuthHandler{Auth: authSvc}
	usrH := &handlers.UsersHandler{
		Svc:     userSvc,
		AuthSvc: authSvc,
		LDAP:    ldapClient,
	}

	// Rotas públicas
	public := http.NewServeMux()
	RegisterSwaggerRoutes(public)
	RegisterHealthRoute(public, db)
	RegisterAuthRoutes(public, authH)

	// Rotas protegidas
	protected := http.NewServeMux()
	protected.HandleFunc("/eu", authH.Me)
	RegisterUsuarioRoutes(protected, usrH, jwtManager, userSvc)

	// Handler principal com CORS
	routes := BuildHandler(public, protected, jwtManager, userSvc)
	routes = middleware.CORS(cfg.CORSOrigin)(routes)

	log.Println("CORS liberado para:", cfg.CORSOrigin)
	return routes
}

// BuildHandler gerencia rotas públicas e protegidas com autenticação
func BuildHandler(public, protected http.Handler, jwtManager *jwt.Manager, userSvc *userusecase.Usercase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Rotas públicas
		for _, prefix := range publicPrefixes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				public.ServeHTTP(w, r)
				return
			}
		}

		// Rotas protegidas com middleware de autenticação
		protectedWithAuth := auth.AuthenticateUser(protected, jwtManager, userSvc)
		protectedWithAuth.ServeHTTP(w, r)
	})
}

// parseDuration converte string em time.Duration
func parseDuration(d string) time.Duration {
	t, _ := time.ParseDuration(d)
	return t
}
