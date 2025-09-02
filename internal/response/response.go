package response

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
)

type Error struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// JSON escreve uma resposta JSON com o status HTTP fornecido.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(v)
}

// ErrorJSON escreve uma resposta JSON de erro com o status HTTP fornecido e registra a mensagem de erro.
func ErrorJSON(w http.ResponseWriter, status int, msg string, details any) {
	log.Printf("error %d: %s", status, msg)
	JSON(w, status, Error{Message: msg, Details: details})
}

// PageResp representa a resposta de paginação para listagens
type PageResponse struct {
	Total  int            `json:"total"`
	Pagina int            `json:"pagina"`
	Limite int            `json:"limite"`
	Data   []user.Usuario `json:"data"`
}

// ErrorResponse representa o padrão de erro da API
type ErrorResponse struct {
    Message string      `json:"message"`
    Details any `json:"details,omitempty"`
}

// HealthResponse estrutura da resposta JSON do health check
type HealthResponse struct {
	Status    string `json:"status"`              // "ok" ou "fail"
	DB        string `json:"database"`            // "ok" ou "fail"
	Timestamp string `json:"timestamp,omitempty"` // horário do check
}