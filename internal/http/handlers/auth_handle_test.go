package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

// ====================== MOCKS ======================

// Mock UserService
type mockUserService struct {
	BuscarPorLoginFn       func(ctx context.Context, login string) (*user.Usuario, error)
	BuscarPorIDFn          func(ctx context.Context, id string) (*user.Usuario, error)
	CriarFn                func(ctx context.Context, u *user.Usuario) error
	AtualizarUltimoLoginFn func(ctx context.Context, id string) error
}

func (m *mockUserService) BuscarPorLogin(ctx context.Context, login string) (*user.Usuario, error) {
	return m.BuscarPorLoginFn(ctx, login)
}
func (m *mockUserService) BuscarPorID(ctx context.Context, id string) (*user.Usuario, error) {
	return m.BuscarPorIDFn(ctx, id)
}
func (m *mockUserService) Criar(ctx context.Context, u *user.Usuario) error {
	return m.CriarFn(ctx, u)
}
func (m *mockUserService) AtualizarUltimoLogin(ctx context.Context, id string) error {
	return m.AtualizarUltimoLoginFn(ctx, id)
}

// Mock LDAP
type mockLDAP struct {
	BindFn          func(login, senha string) error
	SearchByLoginFn func(login string) (string, string, string, error)
}

func (m *mockLDAP) Bind(login, senha string) error {
	return m.BindFn(login, senha)
}
func (m *mockLDAP) SearchByLogin(login string) (string, string, string, error) {
	return m.SearchByLoginFn(login)
}

// Mock JWT
type mockJWT struct {
	SignAccessFn   func(c jwt.Claims) (string, error)
	SignRefreshFn  func(c jwt.Claims) (string, error)
	ParseRefreshFn func(token string) (*jwt.Claims, error)
}

func (m *mockJWT) SignAccess(c jwt.Claims) (string, error) {
	return m.SignAccessFn(c)
}
func (m *mockJWT) SignRefresh(c jwt.Claims) (string, error) {
	return m.SignRefreshFn(c)
}
func (m *mockJWT) ParseRefresh(token string) (*jwt.Claims, error) {
	return m.ParseRefreshFn(token)
}

const (
	AccessToken  = "access-token"
	RefreshToken = "refresh-token"
	Domain       = "@rede.sp"
	Login        = "/login"
	Refresh      = "/refresh"
	Me           = "/me"
)

// ====================== TESTS ======================

// ----- LOGIN -----

// TestLoginSuccessUserExists verifica se o login é bem-sucedido quando o usuário existe.
func TestLoginSuccessUserExistsSuccess(t *testing.T) {
	// Arrange
	mockUser := &user.Usuario{
		ID:        "123",
		Login:     "teste",
		Nome:      "TestUser",
		Email:     "teste@dominio.com",
		Permissao: user.PermUSR,
		Status:    true,
	}

	users := &mockUserService{
		BuscarPorLoginFn:       func(ctx context.Context, login string) (*user.Usuario, error) { return mockUser, nil },
		AtualizarUltimoLoginFn: func(ctx context.Context, id string) error { return nil },
	}

	ldap := &mockLDAP{BindFn: func(login, senha string) error { return nil }}
	jwtM := &mockJWT{
		SignAccessFn:  func(c jwt.Claims) (string, error) { return AccessToken, nil },
		SignRefreshFn: func(c jwt.Claims) (string, error) { return RefreshToken, nil },
	}

	handler := &AuthHandler{
		Users:  users,
		LDAP:   ldap,
		JWT:    jwtM,
		Config: config.Config{LDAPDomain: Domain},
	}

	reqBody, _ := json.Marshal(user.LoginDTO{Login: "teste", Senha: "1234"})
	req := httptest.NewRequest(http.MethodPost, Login, bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	// Act
	handler.Login(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp user.TokenPair
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, AccessToken, resp.AccessToken)
	assert.Equal(t, RefreshToken, resp.RefreshToken)
}

// TestLoginCreateUserFromLDAP verifica se um novo usuário é criado a partir das informações do LDAP.
func TestLoginCreateUserFromLDAPSuccess(t *testing.T) {
	// Arrange
	// Simula um "banco de usuários" em memória
	usersDB := map[string]*user.Usuario{}

	users := &mockUserService{
		BuscarPorLoginFn: func(ctx context.Context, login string) (*user.Usuario, error) {
			u, ok := usersDB[login]
			if !ok {
				return nil, nil
			}
			return u, nil
		},
		CriarFn: func(ctx context.Context, u *user.Usuario) error {
			u.ID = "456"
			usersDB[u.Login] = u
			return nil
		},
		AtualizarUltimoLoginFn: func(ctx context.Context, id string) error { return nil },
	}

	ldap := &mockLDAP{
		BindFn: func(login, senha string) error { return nil },
		SearchByLoginFn: func(login string) (string, string, string, error) {
			return "LDAP User", "ldap@dominio.com", login, nil
		},
	}

	jwtM := &mockJWT{
		SignAccessFn:  func(c jwt.Claims) (string, error) { return AccessToken, nil },
		SignRefreshFn: func(c jwt.Claims) (string, error) { return RefreshToken, nil },
	}

	handler := &AuthHandler{
		Users:  users,
		LDAP:   ldap,
		JWT:    jwtM,
		Config: config.Config{LDAPDomain: "@rede.sp"},
	}

	reqBody, _ := json.Marshal(user.LoginDTO{Login: "teste", Senha: "1234"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	// Act
	handler.Login(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp user.TokenPair
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "access-token", resp.AccessToken)
	assert.Equal(t, "refresh-token", resp.RefreshToken)
}

// TestLoginFailureLDAPBind verifica se o login falha quando as credenciais do LDAP estão incorretas.
func TestLoginFailureLDAPBind(t *testing.T) {
	// Arrange
	users := &mockUserService{BuscarPorLoginFn: func(ctx context.Context, login string) (*user.Usuario, error) {
		return &user.Usuario{Login: login}, nil
	}}
	ldap := &mockLDAP{BindFn: func(login, senha string) error { return errors.New("credenciais incorretas") }}
	jwtM := &mockJWT{}

	handler := &AuthHandler{
		Users:  users,
		LDAP:   ldap,
		JWT:    jwtM,
		Config: config.Config{LDAPDomain: Domain},
	}

	reqBody, _ := json.Marshal(user.LoginDTO{Login: "teste", Senha: "1234"})
	req := httptest.NewRequest(http.MethodPost, Login, bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	// Act
	handler.Login(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ----- REFRESH -----

// TestRefreshSuccess verifica se o refresh do token é bem-sucedido.
func TestRefreshSuccess(t *testing.T) {
	// Arrange
	users := &mockUserService{
		BuscarPorIDFn: func(ctx context.Context, id string) (*user.Usuario, error) {
			return &user.Usuario{ID: "123", Login: "teste"}, nil
		},
		AtualizarUltimoLoginFn: func(ctx context.Context, id string) error { return nil },
	}

	jwtM := &mockJWT{
		ParseRefreshFn: func(token string) (*jwt.Claims, error) { return &jwt.Claims{ID: "123", Login: "teste"}, nil },
		SignAccessFn:   func(c jwt.Claims) (string, error) { return AccessToken, nil },
		SignRefreshFn:  func(c jwt.Claims) (string, error) { return RefreshToken, nil },
	}

	handler := AuthHandler{
		Users:  users,
		JWT:    jwtM,
		Config: config.Config{},
	}

	reqBody, _ := json.Marshal(RefreshRequest{RefreshToken: RefreshToken})
	req := httptest.NewRequest(http.MethodPost, Refresh, bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	// Act
	handler.Refresh(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp user.TokenPair
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, AccessToken, resp.AccessToken)
	assert.Equal(t, RefreshToken, resp.RefreshToken)
}

// TestRefreshInvalidToken verifica se o refresh do token falha com um token inválido.
func TestRefreshInvalidToken(t *testing.T) {
	// Arrange
	users := &mockUserService{}
	jwtM := &mockJWT{
		ParseRefreshFn: func(token string) (*jwt.Claims, error) { return nil, errors.New("refresh inválido") },
	}

	handler := &AuthHandler{
		Users:  users,
		JWT:    jwtM,
		Config: config.Config{},
	}

	reqBody, _ := json.Marshal(RefreshRequest{RefreshToken: "invalid"})
	req := httptest.NewRequest(http.MethodPost, Refresh, bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	// Act
	handler.Refresh(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ----- ME -----

// TestMeSuccess verifica se a requisição "Me" retorna os dados do usuário corretamente.
func TestMeSuccess(t *testing.T) {
	// Arrange
	mockUser := &user.Usuario{ID: "123", Login: "teste", Nome: "Test User", Email: "teste@dominio.com", Status: true}
	users := &mockUserService{BuscarPorIDFn: func(ctx context.Context, id string) (*user.Usuario, error) { return mockUser, nil }}

	handler := &AuthHandler{Users: users}

	req := httptest.NewRequest(http.MethodGet, Me, nil)
	ctx := context.WithValue(req.Context(), middleware.UserKey, &jwt.Claims{ID: "123"})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Act
	handler.Me(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp user.UsuarioResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "123", resp.ID)
	assert.Equal(t, "Test User", resp.Nome)
}

// TestMeFailureNoClaims verifica se a requisição "Me" falha quando não há claims no contexto.
func TestMeFailureNoClaims(t *testing.T) {
	// Arrange
	users := &mockUserService{}
	handler := &AuthHandler{Users: users}

	req := httptest.NewRequest(http.MethodGet, Me, nil)
	w := httptest.NewRecorder()

	// Act
	handler.Me(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ----- HELPERS -----

// TestParseLoginRequestValid verifica se a requisição de login é analisada corretamente.
func TestParseLoginRequestValidSuccess(t *testing.T) {
	// Arrange
	body := `{"login":"teste","senha":"1234"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	// Act
	loginDTO, err := ParseLoginRequest(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "teste", loginDTO.Login)
	assert.Equal(t, "1234", loginDTO.Senha)
}

// TestParseLoginRequestInvalid verifica se a requisição de login inválida é tratada corretamente.
func TestParseLoginRequestInvalid(t *testing.T) {
	// Arrange
	body := `{"login":"", "senha":""}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	// Act
	_, err := ParseLoginRequest(req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "login/senha obrigatórios", err.Error())
}
