package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/httpx"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/ldapx"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/repository"
)

type AuthHandler struct {
	LDAP  *ldapx.Client
	Users repository.UserRepository
	TM    *auth.TokenManager
}

type loginReq struct {
	Login, Password string
}

type loginResp struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "json inválido")
		return
	}

	u, err := h.LDAP.Authenticate(req.Login, req.Password)

	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Upsert no repositório local (shadow account)
	if err := h.Users.Upsert(r.Context(), u); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "falha ao salvar usuário")
		return
	}

	now := time.Now()
	_ = h.Users.TouchLogin(r.Context(), u.Login, now)
	tok, err := h.TM.Generate(auth.Claims{
		UserID:    u.ID,
		Nome:      u.Nome,
		Login:     u.Login,
		Email:     u.Email,
		Permissao: string(u.Permissao),
	})

	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "falha ao gerar token")
		return
	}

	httpx.JSON(w, http.StatusOK, loginResp{Token: tok, User: *u})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	c, err := httpx.ClaimsFrom(r)

	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "não autenticado")
		return
	}

	resp := map[string]any{
		"sub":   c.UserID,
		"name":  c.Nome,
		"login": c.Login,
		"email": c.Email,
		"perm":  c.Permissao,
		"exp":   c.ExpiresAt,
	}
	
	httpx.JSON(w, http.StatusOK, resp)
}
