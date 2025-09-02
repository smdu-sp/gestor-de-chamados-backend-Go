package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handlers"
	response "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/response"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterSwaggerRoutes registra as rotas do Swagger
func RegisterSwaggerRoutes(mux *http.ServeMux) {
	// Swagger UI
	mux.Handle("/swagger/",
		httpSwagger.Handler(httpSwagger.URL("/swagger/swagger.json")))

	// Arquivo swagger.json
	mux.Handle(
		"/swagger/swagger.json",
		http.StripPrefix("/swagger/",
			http.FileServer(http.Dir("./docs"))),
	)
}

// RegisterHealthRoute registra a rota de health check
func RegisterHealthRoute(mux *http.ServeMux, db *sql.DB) {
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

// RegisterAuthRoutes registra as rotas de autenticação
func RegisterAuthRoutes(mux *http.ServeMux, authH *handlers.AuthHandler) {
	mux.HandleFunc("/login", authH.Login)
	mux.HandleFunc("/refresh", authH.Refresh)
}

// RegisterUsuarioRoutes registra as rotas de usuário
func RegisterUsuarioRoutes(mux *http.ServeMux, usrH *handlers.UsersHandler) {
	mux.Handle("/usuarios/criar",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.Criar)),
	)

	mux.Handle("/usuarios/buscar-tudo",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.BuscarTudo)),
	)

	mux.Handle("/usuarios/buscar-por-id/",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.BuscarPorID)), // + :id
	)

	mux.Handle("/usuarios/atualizar/",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.Atualizar)), // + :id
	)

	mux.Handle("/usuarios/lista-completa",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.ListaCompleta)),
	)

	mux.Handle("/usuarios/buscar-tecnicos",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.BuscarTecnicos)),
	)

	mux.Handle("/usuarios/desativar/",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.Desativar)), // + :id
	)

	mux.Handle("/usuarios/autorizar/",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.Autorizar)), // + :id
	)

	mux.Handle("/usuarios/buscar-novo/",
		middleware.RequirePermissions("ADM")(http.HandlerFunc(usrH.BuscarNovo)),
	)

	mux.HandleFunc("/usuarios/valida-usuario", usrH.ValidaUsuario)
}
