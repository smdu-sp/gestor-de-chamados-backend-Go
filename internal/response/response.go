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

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

// ErrorJSON escreve uma resposta JSON de erro com o status HTTP fornecido e registra a mensagem de erro.
func ErrorJSON(w http.ResponseWriter, status int, msg string, details any) {
	log.Printf("error %d: %s", status, msg)
	JSON(w, status, Error{Message: msg, Details: details})
}
