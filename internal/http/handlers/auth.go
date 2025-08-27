package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	myjwt "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	ldapauth "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/ldap"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/response"
)

// AuthHandler gerencia autenticação e refresh de tokens
type AuthHandler struct {
	Users  *user.Service
	JWT    *myjwt.Manager
	LDAP   *ldapauth.Client
	Config config.Config
}

// RefreshRequest payload para refresh token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// --- Helpers internos ---

// writeError envia uma resposta de erro padrão
func writeError(w http.ResponseWriter, code int, msg string) {
	response.ErrorJSON(w, code, msg, nil)
}

// parseLoginRequest lê e valida o payload de login
func parseLoginRequest(r *http.Request) (*user.LoginDTO, error) {
	var req user.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("payload inválido: %w", err)
	}
	if req.Login == "" || req.Senha == "" {
		return nil, errors.New("login/senha obrigatórios")
	}
	return &req, nil
}

// getBindString monta a string de bind para o LDAP
func (h *AuthHandler) getBindString(login string) string {
	// Se LDAPDomain estiver configurado, usa o formato user@domain
	if h.Config.LDAPDomain != "" {
		return login + h.Config.LDAPDomain
	}

	// Senão, usa o formato uid=user,ou=users,base
	ou := "users"
	if login == "admin1" {
		ou = "admins"
	}
	return fmt.Sprintf("uid=%s,ou=%s,%s", login, ou, h.Config.LDAPBase)
}

// createClaims cria as claims JWT a partir do usuário
func createClaims(u *user.Usuario) myjwt.Claims {
	return myjwt.Claims{
		ID:        u.ID,
		Login:     u.Login,
		Nome:      u.Nome,
		Email:     u.Email,
		Permissao: string(u.Permissao),
	}
}

// jwtClaimsFromRequest extrai as claims JWT do contexto da requisição
func jwtClaimsFromRequest(r *http.Request) *myjwt.Claims {
	if v := r.Context().Value(middleware.UserKey); v != nil {
		if c, ok := v.(*myjwt.Claims); ok {
			return c
		}
	}
	return nil
}

// criarUsuarioSeNecessario cria um usuário no banco se ele não existir, buscando dados no LDAP
func (h *AuthHandler) criarUsuarioSeNecessario(ctx context.Context, login string, u *user.Usuario) (*user.Usuario, error) {
	if u != nil {
		return u, nil
	}
	// Usuário não existe, buscar no LDAP e criar 
	name, mail, sLogin, err := h.LDAP.SearchByLogin(login)
	if err != nil {
		log.Println("Erro buscando LDAP:", err)
		return nil, err
	}

	novo := &user.Usuario{
		Nome:      name,
		Login:     sLogin,
		Email:     mail,
		Permissao: user.PermUSR,
		Status:    true,
	}

	// Criar usuário no banco
	if err := h.Users.Criar(ctx, novo); err != nil {
		log.Println("Erro criando usuário:", err)
		return nil, err
	}

	log.Println("Usuário criado com sucesso:", novo.Login)
	return h.Users.BuscarPorLogin(ctx, login)
}

// --- Handlers públicos ---

// Login autentica o usuário via LDAP e retorna um par de tokens JWT
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse e validação do request
	req, err := parseLoginRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Buscar usuário no banco
	u, err := h.Users.BuscarPorLogin(r.Context(), req.Login)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro interno")
		return
	}

	// Autenticar no LDAP
	bind := h.getBindString(req.Login)
	if err := h.LDAP.Bind(bind, req.Senha); err != nil {
		writeError(w, http.StatusUnauthorized, "credenciais incorretas")
		return
	}

	// Se usuário não existe no banco, criar
	u, err = h.criarUsuarioSeNecessario(r.Context(), req.Login, u)
	if err != nil || u == nil {
		writeError(w, http.StatusInternalServerError, "erro ao salvar usuário")
		return
	}

	// Atualizar último login
	if err := h.Users.AtualizarUltimoLogin(r.Context(), u.ID); err != nil {
		log.Println("falha ao atualizar último login:", err)
	}

	// Gerar tokens JWT
	claims := createClaims(u)
	acc, err := h.JWT.SignAccess(claims)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro gerando token")
		return
	}

	// Gerar refresh token
	rt, err := h.JWT.SignRefresh(claims)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro gerando refresh")
		return
	}

	response.JSON(w, http.StatusOK, user.TokenPair{
		AccessToken:  acc,
		RefreshToken: rt,
	})
}

// Refresh valida o refresh token e retorna um novo par de tokens JWT
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "payload inválido")
		return
	}

	if body.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh_token obrigatório")
		return
	}

	claims, err := h.JWT.ParseRefresh(body.RefreshToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "refresh inválido")
		return
	}

	u, err := h.Users.BuscarPorID(r.Context(), claims.ID)
	if err != nil || u == nil {
		writeError(w, http.StatusUnauthorized, "usuário inválido")
		return
	}

	if err := h.Users.AtualizarUltimoLogin(r.Context(), u.ID); err != nil {
		log.Println("falha ao atualizar último login:", err)
	}

	// Atualizar issued at para agora
	claims.RegisteredClaims.IssuedAt = gojwt.NewNumericDate(time.Now())

	acc, err := h.JWT.SignAccess(*claims)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro gerando token")
		return
	}
	rt, err := h.JWT.SignRefresh(*claims)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro gerando refresh")
		return
	}

	response.JSON(w, http.StatusOK, user.TokenPair{
		AccessToken:  acc,
		RefreshToken: rt,
	})
}


// Me retorna os dados do usuário autenticado
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	clm := jwtClaimsFromRequest(r)
	if clm == nil {
		writeError(w, http.StatusUnauthorized, "não autenticado")
		return
	}

	u, err := h.Users.BuscarPorID(r.Context(), clm.ID)
	if err != nil || u == nil {
		writeError(w, http.StatusNotFound, "usuário não encontrado")
		return
	}

	resp := user.UsuarioResponse{
		ID:        u.ID,
		Nome:      u.Nome,
		Login:     u.Login,
		Email:     u.Email,
		Permissao: u.Permissao,
		Status:    u.Status,
		Avatar:    u.Avatar,
	}

	response.JSON(w, http.StatusOK, resp)
}
