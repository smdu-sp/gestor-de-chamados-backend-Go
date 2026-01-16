package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// AuthHandler gerencia as requisições HTTP relacionadas à autenticação.
type AuthHandler struct {
	Svc AuthInternoService
}

// NovoAuthHandler cria uma nova instância de AuthHandler.
func NovoAuthHandler(authInSvc AuthInternoService) *AuthHandler {
	return &AuthHandler{Svc: authInSvc}
}

// Login godoc
// @Summary Login de usuário
// @Description Autentica um usuário e retorna tokens JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginReq true "Credenciais de login"
// @Success 200 {object} TokenResp "Tokens JWT retornados com sucesso"
// @Failure 400 {object} ErroResposta "Payload inválido ou JSON malformado"
// @Failure 401 {object} ErroResposta "Credenciais inválidas"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req LoginReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	tokens, err := h.Svc.Login(r.Context(), req)
	if err != nil {
		partes := strings.Split(err.Error(), ":")
		erroRaiz := strings.TrimSpace(partes[len(partes)-1])
		slog.Error("Erro ao logar usuário", "error", err)
		ResponderJSONUnauthorized(w, r, erroRaiz)
		return
	}

	ResponderJSONStatusOK(w, tokens)
}

// Refresh godoc
// @Summary Refresh de tokens
// @Description Atualiza os tokens JWT usando um refresh token válido.
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh body RefreshReq true "Refresh token"
// @Success 200 {object} TokenResp "Tokens JWT atualizados com sucesso"
// @Failure 400 {object} ErroResposta "Payload inválido ou JSON malformado"
// @Failure 401 {object} ErroResposta "Refresh token inválido ou expirado"
// @Router /refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	tokens, err := h.Svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		slog.Error("Erro ao atualizar tokens", "error", err)
		ResponderJSONUnauthorized(w, r, "Erro ao atualizar tokens")
		return
	}

	ResponderJSONStatusOK(w, tokens)
}

// Me godoc
// @Summary Dados do usuário autenticado
// @Description Retorna os dados do usuário autenticado com base no token JWT.
// @Tags Auth
// @Produce json
// @Success 200 {object} UsuarioResp "Dados do usuário retornados com sucesso"
// @Failure 401 {object} ErroResposta "Token inválido ou ausente"
// @Failure 404 {object} ErroResposta "Usuário não encontrado"
// @Failure 500 {object} ErroResposta "Erro interno ao buscar dados do usuário"
// @Router /me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, err := ClaimsDoRequest(r)
	if err == ErrClaimsInvalidasOuAusentes  {
		ResponderJSONUnauthorized(w, r, "token inválido ou ausente")
		return
	}

	usuario, err := h.Svc.Me(r.Context(), claims.ID)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrUsuarioNaoEncontrado):
			ResponderJSONNotFound(w, r, "usuário não encontrado")
			return

		default:
			slog.Error("Erro interno ao buscar dados do usuário", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar dados do usuário")
			return
		}
	}

	ResponderJSONStatusOK(w, ParaUsuarioResp(*usuario))
}
