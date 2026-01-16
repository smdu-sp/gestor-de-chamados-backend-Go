package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// -----------------------------------------------------------------------------
// ESCRITA JSON - SUCESSO
// -----------------------------------------------------------------------------

// ResponderJSON escreve uma resposta JSON com o status HTTP fornecido.
func ResponderJSON(w http.ResponseWriter, status int, corpo any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(corpo); err != nil {
		slog.Error("[EscreverJSON]: erro ao serializar", "error", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
}

// ResponderJSONStatusOK escreve uma resposta JSON com status 200 OK.
func ResponderJSONStatusOK(w http.ResponseWriter, corpo any) {
	ResponderJSON(w, http.StatusOK, corpo)
}

// -----------------------------------------------------------------------------
// ESCRITA JSON ERRO - RFC 7807
// -----------------------------------------------------------------------------

// ErroResposta epresenta a estrutura de resposta de erro definida no RFC 7807.
type ErroResposta struct {
	Status    int    `json:"status"`              // Código de status HTTP
	Tipo      string `json:"tipo,omitempty"`      // URI que identifica o tipo de erro
	Instancia string `json:"instancia,omitempty"` // URI que identifica a instância do erro
	Titulo    string `json:"titulo"`              // Título curto do erro
	Detalhes  any    `json:"detalhes,omitempty"`  // Descrição detalhada do erro
}

// ResponderJSONError escreve uma resposta JSON de erro com o status HTTP fornecido e registra a mensagem de erro.
func ResponderJSONError(w http.ResponseWriter, status int, erroResp ErroResposta) {
	erroResp.Status = status

	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(erroResp); err != nil {
		slog.Error("[EscreverJSONError]: erro ao serializar", "error", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
}

// ResponderJSONUnauthorized escreve uma resposta JSON com status 401 Unauthorized.
func ResponderJSONUnauthorized(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusUnauthorized, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Acesso não Autorizado",
		Detalhes:  erro,
	})
}

// ResponderJSONForbidden escreve uma resposta JSON com status 403 Forbidden.
func ResponderJSONForbidden(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusForbidden, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Acesso Negado",
		Detalhes:  erro,
	})
}

// ResponderJSONInvalidPayload escreve uma resposta JSON com status 400 Bad Request para payloads inválidos.
func ResponderJSONInvalidPayload(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusBadRequest, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Payload Inválido",
		Detalhes:  erro,
	})
}

// ResponderJSONInternalError escreve uma resposta JSON com status 500 Internal Server Error.
func ResponderJSONInternalError(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusInternalServerError, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Erro Interno do Servidor",
		Detalhes:  erro,
	})
}

// ResponderJSONNotFound escreve uma resposta JSON com status 404 Not Found.
func ResponderJSONNotFound(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusNotFound, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Recurso não Encontrado",
		Detalhes:  erro,
	})
}

// -----------------------------------------------------------------------------
// DECODIFICAÇÃO DE JSON
// -----------------------------------------------------------------------------

// DecodificarJSON decodifica o payload JSON da requisição para a estrutura fornecida
func DecodificarJSON(w http.ResponseWriter, r *http.Request, payload any) bool {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(payload); err != nil {

		var ute *json.UnmarshalTypeError
		var se *json.SyntaxError

		switch {
		// Tipo incorreto em campo JSON
		case errors.As(err, &ute):
			erro := fmt.Sprintf(
				"Campo '%s' deveria ser %s, mas recebeu %s",
				ute.Field, ute.Type.String(), ute.Value,
			)
			ResponderJSONInvalidPayload(w, r, erro)
			return false

		// Campo desconhecido no JSON
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			field = strings.Trim(field, `"`)
			ResponderJSONInvalidPayload(w, r, "Campo desconhecido: '"+field+"'")
			return false

		// JSON malformado
		case errors.As(err, &se):
			erro := fmt.Sprintf("JSON malformado (posição %d)", se.Offset)
			ResponderJSONInvalidPayload(w, r, erro)
			return false

		// JSON vazio
		case errors.Is(err, io.EOF):
			ResponderJSONInvalidPayload(w, r, "JSON vazio")
			return false

		// Erro desconhecido de decodificação
		default:
			ResponderJSONInvalidPayload(w, r, err.Error())
			return false
		}
	}

	// Verifica se há mais dados no JSON além do objeto principal
	if dec.More() {
		ResponderJSONInvalidPayload(w, r, "JSON possui conteúdo inesperado após o objeto principal")
		return false
	}

	return true
}
