package router

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	ldapauth "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/ldap"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handlers"
)

func Build(cfg config.Config, db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	repo := user.NewRepository(db)
	svc := user.NewService(repo)
	accessTTL, _ := time.ParseDuration(cfg.AccessTTL)
	refreshTTL, _ := time.ParseDuration(cfg.RefreshTTL)

	jm := &jwt.Manager{
		AccessSecret:  []byte(cfg.JWTSecret),
		RefreshSecret: []byte(cfg.RTSecret),
		AccessTTL:     accessTTL,
		RefreshTTL:    refreshTTL,
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
		JWT:    jm,
		LDAP:   ldapClient,
		Config: cfg,
	}

	usrH := &handlers.UsersHandler{
		Svc:  svc,
		LDAP: ldapClient,
	}

	// Público
	mux.HandleFunc("/login", authH.Login)
	mux.HandleFunc("/refresh", authH.Refresh) // lê refresh_token do body

	// Protegido (JWT)
	protected := http.NewServeMux()
	// Auth util
	protected.HandleFunc("/eu", authH.Me)

	// Usuários — paths EXATOS do Nest
	protected.HandleFunc("/usuarios/criar", usrH.Criar)
	protected.HandleFunc("/usuarios/buscar-tudo", usrH.BuscarTudo)
	protected.HandleFunc("/usuarios/buscar-por-id/", usrH.BuscarPorID) // + :id
	protected.HandleFunc("/usuarios/atualizar/", usrH.Atualizar)       // + :id
	protected.HandleFunc("/usuarios/lista-completa", usrH.ListaCompleta)
	protected.HandleFunc("/usuarios/buscar-tecnicos", usrH.BuscarTecnicos)
	protected.HandleFunc("/usuarios/desativar/", usrH.Desativar) // + :id (soft delete)
	protected.HandleFunc("/usuarios/autorizar/", usrH.Autorizar) // + :id (status=true)
	protected.HandleFunc("/usuarios/valida-usuario", usrH.ValidaUsuario)
	protected.HandleFunc("/usuarios/buscar-novo/", usrH.BuscarNovo) // + :login

	// Aplica middlewares: CORS, Logger e JWT+update ultimoLogin
	var h http.Handler = mux
	h = chain(
		h,
		mount(middleware.WithUser(protected, jm, func(ctx context.Context, id string) error {
			return svc.AtualizarUltimoLogin(ctx, id)
		})),
		middleware.CORS(cfg.CORSOrigin),
		middleware.Logger,
	)

	log.Println("CORS liberado para:", cfg.CORSOrigin)
	return h
}

// Helpers
func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func mount(h http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" || r.URL.Path == "/login" || r.URL.Path == "/refresh" {
				next.ServeHTTP(w, r) // chama o handler público
				return
			}
			h.ServeHTTP(w, r) // chama o handler protegido
		})
	}
}
