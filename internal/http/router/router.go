package router

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	ldapauth "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/ldap"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handlers"
)

var publicPrefixes = []string{
	"/health",
	"/login",
	"/refresh",
	"/swagger",
}

// Build faz a construção do handler principal
func Build(cfg config.Config, db *sql.DB) http.Handler {
	repo := user.NewRepository(db)
	svc := user.NewService(repo)

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

	authH := &handlers.AuthHandler{
		Users:  svc,
		JWT:    jwtManager,
		LDAP:   ldapClient,
		Config: cfg,
	}

	usrH := &handlers.UsersHandler{
		Svc:  svc,
		LDAP: ldapClient,
	}

	// Rotas públicas
	public := http.NewServeMux()
	RegisterSwaggerRoutes(public)
	RegisterHealthRoute(public, db)
	RegisterAuthRoutes(public, authH)

	// Rotas protegidas
	protected := http.NewServeMux()
	protected.HandleFunc("/eu", authH.Me)
	RegisterUsuarioRoutes(protected, usrH)

	// Handler principal
	routes := BuildHandler(public, protected, jwtManager, svc)
	routes = middleware.CORS(cfg.CORSOrigin)(routes)

	log.Println("CORS liberado para:", cfg.CORSOrigin)
	return routes
}

// BuildHandler gerencia as rotas publicas e protegidas
func BuildHandler(public, protected http.Handler, jwtManager *jwt.Manager, svc user.UserServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Rotas públicas
		for _, prefix := range publicPrefixes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				public.ServeHTTP(w, r)
				return
			}
		}

		// Rotas protegidas
		protectedWithAuth := middleware.AuthenticateUser(protected, jwtManager, svc)
		protectedWithAuth.ServeHTTP(w, r)
	})
}

// parseDuration converte uma string em time.Duration
func parseDuration(d string) time.Duration {
	t, _ := time.ParseDuration(d)
	return t
}
