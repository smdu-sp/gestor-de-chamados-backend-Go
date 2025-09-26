package response

import (
	"encoding/json"
	"log"
	"net/http"
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
	err := encoder.Encode(v)
	if err != nil {
		log.Printf("error encoding JSON response: %v", err)
		http.Error(w, "erro interno", http.StatusInternalServerError)
	}
}

// ErrorJSON escreve uma resposta JSON de erro com o status HTTP fornecido e registra a mensagem de erro.
func ErrorJSON(w http.ResponseWriter, status int, msg string, details any) {
	log.Printf("error %d: %s: %v", status, msg, details)
	JSON(w, status, Error{Message: msg, Details: details})
}

// PageResp representa a resposta de paginação para listagens
type PageResponse[T any] struct {
	Total  int `json:"total"`
	Pagina int `json:"pagina"`
	Limite int `json:"limite"`
	Items   []T `json:"items"`
}

// HealthResponse estrutura da resposta JSON do health check
type HealthResponse struct {
	Status    string `json:"status"`              // "ok" ou "fail"
	DB        string `json:"database"`            // "ok" ou "fail"
	Timestamp string `json:"timestamp,omitempty"` // horário do check
}

// TokenPair representa um par de tokens JWT (access e refresh)
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// MethodErrorResponse representa a estrutura de resposta para erros de método HTTP
type MethodErrorResponse struct {
	MetodoUsado     string `json:"metodo_usado"`
	MetodoPermitido string `json:"metodo_permitido"`
}
