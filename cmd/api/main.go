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

	// Carrega variáveis de ambiente
	err := godotenv.Load()
	if err != nil {
		log.Println("⛔ Aviso: .env não foi carregado automaticamente")
	}

	// Carrega a configuração da aplicação (HTTP, JWT, LDAP, etc.)
	cfg := config.Load()

	// Cria repositório de usuários em memória (poderá ser trocado por banco no futuro).
	users := memory.NewUserRepo()

	// Inicializa cliente LDAP para autenticação contra o diretório.
	ldapClient := ldapx.New(cfg.LDAP)

	// Cria gerenciador de tokens JWT com base nas configs carregadas.
	tm := auth.NewTokenManager(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		cfg.JWTTTL,
	)

	// Instancia handler de autenticação (login, /me), injetando dependências.
	authH := &handler.AuthHandler{
		LDAP:  ldapClient,
		Users: users,
		TM:    tm,
	}

	// Instancia handler de usuário (rotas protegidas de exemplo).
	userH := &handler.UserHandler{}

	// Cria roteador HTTP nativo (ServeMux).
	mux := httpx.NewMux()

	// ------------------------
	// Rotas públicas
	// ------------------------

	// Endpoint público raiz "/"
	mux.HandleFunc("/", httpx.Method(userH.Publico, http.MethodGet))

	// Endpoint de login "/login" (gera token JWT)
	mux.HandleFunc("/login", httpx.Method(authH.Login, http.MethodPost))

	// ------------------------
	// Rotas protegidas
	// ------------------------

	// Cria uma cadeia de middlewares: autenticação JWT + timeout + recover
	authChain := func(h http.Handler) http.Handler {
		return httpx.Chain(h, 
			httpx.AuthMiddleware(tm),			// valida token
			httpx.Timeout(5*time.Second), // timeout por request
			httpx.Recover,								// captura panics
		)
	}

	// Rota autenticada: retorna dados do usuário logado (/me)
	mux.Handle("/me", authChain(http.HandlerFunc(authH.Me)))

	// ------------------------
	// RBAC (controle de acesso por papel)
	// ------------------------

	// Apenas usuários com papel "ADM"
	mux.Handle("/admin/ping", authChain(
		httpx.RequireRoles("ADM")(http.HandlerFunc(userH.AdminOnly)),
	))

	// Apenas usuários com papel "SUP" ou "DEV"
	mux.Handle("/suporte/ping", authChain(
		httpx.RequireRoles("SUP", "DEV")(http.HandlerFunc(userH.SuporteOuDev)),
	))

	// ------------------------
	// Inicialização do servidor
	// ------------------------

	// Exibe no log a porta em que a API está escutando
	log.Printf("listening on %s", cfg.Addr)

	// Inicia servidor HTTP com o roteador configurado
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal(err) // encerra se não conseguir subir
	}
}
