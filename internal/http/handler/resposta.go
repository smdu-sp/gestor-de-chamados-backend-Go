package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
)

// -----------------------------------------------------------------------------
// SUCESSO
// -----------------------------------------------------------------------------

// ResponderJSON escreve uma resposta JSON com o status HTTP fornecido.
func ResponderJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(body); err != nil {
		slog.Error("[EscreverJSON]: erro ao serializar", "error", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
}

// ResponderJSONNoContent escreve uma resposta JSON com status 204 No Content.
func ResponderJSONNoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

// ResponderJSONStatusOK escreve uma resposta JSON com status 200 OK.
func ResponderJSONStatusOK(w http.ResponseWriter, body any) {
	ResponderJSON(w, http.StatusOK, body)
}

// ResponderJSONCreated escreve uma resposta JSON com status 201 Created.
func ResponderJSONCreated(w http.ResponseWriter, body any) {
	ResponderJSON(w, http.StatusCreated, body)
}

// ResponderJSONPaginado escreve uma resposta JSON paginada com status 200 OK.
func ResponderJSONPaginado[T any](w http.ResponseWriter, total, pagina, limite int, items []T) {
	resp := dto.NovoRespPaginada(total, pagina, limite, items)
	ResponderJSON(w, http.StatusOK, resp)
}

// -----------------------------------------------------------------------------
// ERRO - RFC 7807
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
func ResponderJSONError(w http.ResponseWriter, status int, e ErroResposta) {
	e.Status = status

	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(e); err != nil {
		slog.Error("[EscreverJSONError]: erro ao serializar", "error", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
}

// ResponderJSONConflict escreve uma resposta JSON com status 409 Conflict.
func ResponderJSONConflict(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusConflict, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Conflito de Recurso",
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

// ResponderJSONTimeout escreve uma resposta JSON com status 408 Request Timeout.
func ResponderJSONTimeout(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusRequestTimeout, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Tempo de Requisição Excedido",
		Detalhes:  erro,
	})
}

// ResponderJSONBadRequest escreve uma resposta JSON com status 400 Bad Request.
func ResponderJSONBadRequest(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusBadRequest, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Requisição Inválida",
		Detalhes:  erro,
	})
}

// ResponderJSONContextCanceled escreve uma resposta JSON com status 400 Bad Request.
func ResponderJSONContextCanceled(w http.ResponseWriter, r *http.Request, erro string) {
	ResponderJSONError(w, http.StatusBadRequest, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Contexto Cancelado",
		Detalhes:  erro,
	})
}

// ResponderJSONValidationError escreve uma resposta JSON com status 422 Unprocessable Entity.
func ResponderJSONValidationError(w http.ResponseWriter, r *http.Request, erro any) {
	ResponderJSONError(w, http.StatusUnprocessableEntity, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Erro de Validação de Dados",
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

// ResponderJSONRecoveryPanic escreve uma resposta JSON de erro para pânicos recuperados em handlers HTTP.
func ResponderJSONRecoveryPanic(w http.ResponseWriter, r *http.Request) {
	ResponderJSONError(w, http.StatusInternalServerError, ErroResposta{
		Instancia: r.URL.Path,
		Titulo:    "Erro interno do servidor",
		Detalhes:  "Ocorreu um erro inesperado ao processar a requisição.",
	})
}

// -----------------------------------------------------------------------------
// DECODIFICAÇÃO
// -----------------------------------------------------------------------------

// DecodificarJSON decodifica o payload JSON da requisição para a estrutura fornecida
func DecodificarJSON(w http.ResponseWriter, r *http.Request, payload any) bool {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(payload); err != nil {

		var ute *json.UnmarshalTypeError
		var se *json.SyntaxError

		switch {
		// Tipo de dado incorreto
		case errors.As(err, &ute):
			erro := fmt.Sprintf(
				"Campo '%s' deveria ser %s, mas recebeu %s",
				ute.Field, ute.Type.String(), ute.Value,
			)
			ResponderJSONInvalidPayload(w, r, erro)
			return false

		// Campo desconhecido
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			field = strings.Trim(field, `"`)
			ResponderJSONInvalidPayload(w, r, "Campo desconhecido: "+field)
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

		default:
			slog.Error("[DecodificarJSON]: erro ao decodificar", "error", err)
			ResponderJSONInternalError(w, r, err.Error())
			return false
		}
	}

	if dec.More() {
		ResponderJSONInvalidPayload(w, r, "JSON possui conteúdo inesperado após o objeto principal")
		return false
	}

	return true
}
