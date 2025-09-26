package handler

import (
	"encoding/json"
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
)

// AuthHandler gerencia as requisições HTTP relacionadas à autenticação.
type AuthHandler struct {
	Usecase usecase.AuthInternoUsecase
}

// RefreshRequest representa o payload para a requisição de refresh de tokens.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LoginDTO representa o payload para a requisição de login.
type LoginDto struct {
	Login string `json:"login"`
	Senha string `json:"senha"`
}

// --- Helpers ---
// parseLoginRequest analisa o corpo da requisição para extrair os dados de login.
func parseLoginRequest(r *http.Request) (*LoginDto, error) {
	var req LoginDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Login == "" || req.Senha == "" {
		return nil, http.ErrNoCookie
	}
	return &req, nil
}

// jwtClaimsFromRequest extrai as claims JWT do contexto da requisição.
func jwtClaimsFromRequest(r *http.Request) *jwt.Claims {
	if valor := r.Context().Value(middleware.ChaveUsuario); valor != nil {
		if claims, ok := valor.(*jwt.Claims); ok {
			return claims
		}
	}
	return nil
}

// --- Handlers ---

// Login godoc
// @Summary      Login
// @Description  Autentica um usuário e retorna tokens JWT.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body     LoginDto true  "Login e senha do usuário"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /login [post]
// Login autentica um usuário e retorna tokens JWT.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := parseLoginRequest(r)
	if err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "payload inválido", err.Error())
		return
	}

	tokens, err := h.Usecase.Login(r.Context(), req.Login, req.Senha)
	if err != nil {
		response.ErrorJSON(w, http.StatusUnauthorized, "falha no login", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, tokens)
}

// Refresh godoc
// @Summary      Refresh
// @Description  Atualiza os tokens JWT usando um token de refresh.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refreshRequest  body      RefreshRequest  true  "Token de refresh"
// @Success      200             {object}  map[string]string
// @Failure      400             {object}  map[string]string
// @Failure      401             {object}  map[string]string
// @Router       /refresh [post]
// Refresh atualiza os tokens JWT usando um token de refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "payload inválido", err.Error())
		return
	}

	tokens, err := h.Usecase.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		response.ErrorJSON(w, http.StatusUnauthorized, "refresh inválido", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := jwtClaimsFromRequest(r)
	if claims == nil {
		response.ErrorJSON(w, http.StatusUnauthorized, "não autenticado", nil)
		return
	}

	usuario, err := h.Usecase.Me(r.Context(), claims.ID)
	if err != nil || usuario == nil {
		response.ErrorJSON(w, http.StatusNotFound, "usuário não encontrado", nil)
		return
	}

	response.JSON(w, http.StatusOK, usuario)
}
