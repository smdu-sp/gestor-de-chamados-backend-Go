package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

// RegistrarRotasSwagger registra as rotas do Swagger
func RegistrarRotasSwagger(mux *http.ServeMux) {
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
}

// RegistrarRotasHealthCheck registra a rota de health check
func RegistrarRotasHealthCheck(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		// HealthResponse estrutura da resposta JSON do health check
		type HealthResp struct {
			Status    string `json:"status"`              // "ok" ou "fail"
			DB        string `json:"database"`            // "ok" ou "fail"
			Timestamp string `json:"timestamp,omitempty"` // horário do check
		}

		resp := HealthResp{
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
