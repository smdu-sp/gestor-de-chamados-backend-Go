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
	RM    *auth.RefreshManager
	RRepo repository.RefreshRepository
}

type loginReq struct {
	Login, Password string
}

type loginResp struct {
	AccessToken  string     `json:"accessToken"`
	RefreshToken string     `json:"refreshToken"`
	User         model.User `json:"user"`
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

type refreshResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Login autentica o usuário via LDAP, atualiza o shadow account
// no repositório local e retorna um JWT
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

	// Atualiza a última vez que o usuário fez login
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

	rt := h.RM.Generate()
	if err := h.RRepo.Save(
		r.Context(),
		repository.RefreshToken{
			Token: rt, UserID: u.ID,
		}); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "falha ao salvar refresh token")
		return
	}

	// Retorna JSON com o token, refresh token e dados do usuário
	httpx.JSON(w, http.StatusOK, loginResp{
		AccessToken:  tok,
		RefreshToken: rt,
		User:         *u,
	})
}

// Refresh endpoint
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "json inválido")
		return
	}

	exists, _ := h.RRepo.Exists(r.Context(), req.RefreshToken)
	if !exists {
		httpx.Error(w, http.StatusUnauthorized, "refresh inválido")
		return
	}

	// Em cenário real, associe Refresh -> UserID e recarregue o usuário do repositório
	// Aqui simplificado: invalida o refresh e gera novo Access
	if err := h.RRepo.Delete(r.Context(), req.RefreshToken); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "falha ao invalidar refresh token")
		return
	}

	// Em produção, recarregue dados do usuário
	claims := auth.Claims{
		UserID:    "user-id",
		Nome:      "",
		Login:     "",
		Email:     "",
		Permissao: string(model.USR),
	}
	at, err := h.TM.Generate(claims)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "falha ao gerar token")
		return
	}

	// Novo refresh
	rt := h.RM.Generate()
	if err := h.RRepo.Save(
		r.Context(),
		repository.RefreshToken{Token: rt, UserID: claims.UserID},
	); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "falha ao salvar refresh token")
		return
	}

	// Retorna JSON com o token e refresh token
	httpx.JSON(w, http.StatusOK, refreshResp{
		AccessToken:  at,
		RefreshToken: rt,
	})
}

// Logout: revoga todos refresh do usuário
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := httpx.ClaimsFrom(r)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "não autenticado")
		return
	}
	if err := h.RRepo.DeleteByUser(r.Context(), c.UserID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "falha ao invalidar refresh tokens")
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"ok": "logout"})
}

// Retorna informações do usuário autenticado a partir do token JWT
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
