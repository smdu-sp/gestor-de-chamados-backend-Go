package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/handler"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	goSwagger "github.com/swaggo/http-swagger"
)

// RegisterSwaggerRoutes registra as rotas do Swagger
func RegisterSwaggerRoutes(mux *http.ServeMux) {
	// Swagger UI
	mux.Handle("/swagger/",
		goSwagger.Handler(goSwagger.URL("/swagger/swagger.json")))

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
func RegisterAuthRoutes(mux *http.ServeMux, authH *handler.AuthHandler) {
	mux.HandleFunc("/login", authH.Login)
	mux.HandleFunc("/refresh", authH.Refresh)
}

// RegisterUsuarioRoutes registra as rotas de usuário
func RegisterUsuarioRoutes(mux *http.ServeMux, usrH *handler.UsersHandler, jwtManager *jwt.Manager, svc usecase.UserUsecase) {
    // helper para aplicar autenticação + permissões
    secure := func(handler http.HandlerFunc, perms ...string) http.Handler {
        return middleware.AuthenticateUser(
            middleware.RequirePermissions(perms...)(handler),
            jwtManager, svc,
        )
    }

    mux.Handle("/usuarios/criar", secure(usrH.Criar, "ADM"))
    mux.Handle("/usuarios/buscar-tudo", secure(usrH.BuscarTudo, "ADM"))
    mux.Handle("/usuarios/buscar-por-id/", secure(usrH.BuscarPorID, "ADM"))
    mux.Handle("/usuarios/atualizar/", secure(usrH.Atualizar, "ADM"))
    mux.Handle("/usuarios/lista-completa", secure(usrH.ListaCompleta, "ADM"))
    mux.Handle("/usuarios/buscar-tecnicos", secure(usrH.BuscarTecnicos, "ADM"))
    mux.Handle("/usuarios/desativar/", secure(usrH.Desativar, "ADM"))
    mux.Handle("/usuarios/autorizar/", secure(usrH.Autorizar, "ADM"))
    mux.Handle("/usuarios/buscar-novo/", secure(usrH.BuscarNovo, "ADM"))

    mux.HandleFunc("/usuarios/valida-usuario", usrH.ValidaUsuario) // não precisa de ADM
}
