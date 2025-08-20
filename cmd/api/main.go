package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/handler"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/httpx"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/ldapx"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/repository/memory"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("⛔ Aviso: .env não foi carregado automaticamente")
	}

	cfg := config.Load()
	users := memory.NewUserRepo()
	ldapClient := ldapx.New(cfg.LDAP)
	tm := auth.NewTokenManager(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		cfg.JWTTTL,
	)
	authH := &handler.AuthHandler{
		LDAP:  ldapClient,
		Users: users,
		TM:    tm,
	}
	userH := &handler.UserHandler{}
	mux := httpx.NewMux()

	// Rota pública
	mux.HandleFunc("/", httpx.Method(userH.Publico, http.MethodGet))

	// Login (pública)
	mux.HandleFunc("/login", httpx.Method(authH.Login, http.MethodPost))

	// Protegidas
	authChain := func(h http.Handler) http.Handler {
		return httpx.Chain(h, httpx.AuthMiddleware(tm),
			httpx.Timeout(5*time.Second), httpx.Recover)
	}

	mux.Handle("/me", authChain(http.HandlerFunc(authH.Me)))

	// RBAC exemplos
	mux.Handle("/admin/ping", authChain(httpx.RequireRoles("ADM")(http.HandlerFunc(userH.AdminOnly))))
	mux.Handle("/suporte/ping", authChain(httpx.RequireRoles("SUP", "DEV")(http.HandlerFunc(userH.SuporteOuDev))))
	log.Printf("listening on %s", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal(err)
	}
}
