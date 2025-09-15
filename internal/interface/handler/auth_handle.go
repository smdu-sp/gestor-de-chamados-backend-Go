package handler

import (
	"encoding/json"
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
)

type AuthHandler struct {
	Auth usecase.AuthUsecase
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginDTO struct {
	Login string `json:"login"`
	Senha string `json:"senha"`
}

// --- Helpers ---

func parseLoginRequest(r *http.Request) (*LoginDTO, error) {
	var req LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Login == "" || req.Senha == "" {
		return nil, http.ErrNoCookie 
	}
	return &req, nil
}

func jwtClaimsFromRequest(r *http.Request) *jwt.Claims {
	if valor := r.Context().Value(middleware.UserKey); valor != nil {
		if claims, ok := valor.(*jwt.Claims); ok {
			return claims
		}
	}
	return nil
}

// --- Handlers ---

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := parseLoginRequest(r)
	if err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "payload inválido", err.Error())
		return
	}

	tokens, err := h.Auth.Login(r.Context(), req.Login, req.Senha)
	if err != nil {
		response.ErrorJSON(w, http.StatusUnauthorized, "falha no login", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "payload inválido", err.Error())
		return
	}

	tokens, err := h.Auth.Refresh(r.Context(), body.RefreshToken)
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

	usuario, err := h.Auth.Me(r.Context(), claims.ID)
	if err != nil || usuario == nil {
		response.ErrorJSON(w, http.StatusNotFound, "usuário não encontrado", nil)
		return
	}

	response.JSON(w, http.StatusOK, usuario)
}