package handlers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-ldap/ldap/v3"
	ldapauth "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/ldap"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

// --- MOCK REPOSITORY ---

type mockUserRepo struct {
	InsertFunc        func(ctx context.Context, u *user.Usuario) error
	UpdateFunc        func(ctx context.Context, id string, u *user.Usuario) error
	FindByIDFunc      func(ctx context.Context, id string) (*user.Usuario, error)
	FindByLoginFunc   func(ctx context.Context, login string) (*user.Usuario, error)
	ListFunc          func(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]user.Usuario, int, error)
	UpdateLastLoginFn func(ctx context.Context, id string) error
}

func (m *mockUserRepo) Insert(ctx context.Context, u *user.Usuario) error {
	return m.InsertFunc(ctx, u)
}

func (m *mockUserRepo) Update(ctx context.Context, id string, u *user.Usuario) error {
	return m.UpdateFunc(ctx, id, u)
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*user.Usuario, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *mockUserRepo) FindByLogin(ctx context.Context, login string) (*user.Usuario, error) {
	return m.FindByLoginFunc(ctx, login)
}

func (m *mockUserRepo) List(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]user.Usuario, int, error) {
	return m.ListFunc(ctx, pagina, limite, busca, status, permissao)
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, id string) error {
	if m.UpdateLastLoginFn != nil {
		return m.UpdateLastLoginFn(ctx, id)
	}
	return nil
}

// --- MOCK LDAP ---
type mockLDAPConnection struct{}

func (m *mockLDAPConnection) Close() error                         { return nil }
func (m *mockLDAPConnection) Bind(username, password string) error { return nil }
func (m *mockLDAPConnection) Search(sr *ldap.SearchRequest) (*ldap.SearchResult, error) {
	entry := &ldap.Entry{
		Attributes: []*ldap.EntryAttribute{
			{Name: "cn", Values: []string{"Nome Teste"}},
			{Name: "mail", Values: []string{"teste@teste.com"}},
			{Name: "uid", Values: []string{"teste"}},
		},
	}
	return &ldap.SearchResult{Entries: []*ldap.Entry{entry}}, nil
}

func (m *mockLDAPConnection) StartTLS(cfg *tls.Config) error { return nil }

// --- TESTES ---

// TestCriarUsuario verifica se o usuário é criado com sucesso
func TestCriarUsuarioSucesso(t *testing.T) {
	// Arrange
	repo := &mockUserRepo{
		InsertFunc: func(ctx context.Context, u *user.Usuario) error {
			return nil
		},
	}
	svc := user.NewService(repo)
	handler := &UsersHandler{Svc: svc}

	usuario := user.Usuario{ID: "1", Nome: "Test", Login: "t", Email: "a@b.com"}
	body, _ := json.Marshal(usuario)
	req := httptest.NewRequest(http.MethodPost, "/usuarios/criar", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Act
	handler.Criar(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp user.Usuario
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, usuario.ID, resp.ID)
}

// TestCriarUsuarioBadPayload verifica se o handler retorna erro ao receber payload inválido
func TestCriarUsuarioBadPayload(t *testing.T) {
	// Arrange
	svc := user.NewService(nil)
	handler := &UsersHandler{Svc: svc}

	req := httptest.NewRequest(http.MethodPost, "/usuarios/criar", bytes.NewReader([]byte("{invalid")))
	w := httptest.NewRecorder()

	// Act
	handler.Criar(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestBuscarPorID verifica se o usuário é buscado pelo ID
func TestBuscarPorIDSucesso(t *testing.T) {
	// Arrange
	repo := &mockUserRepo{
		FindByIDFunc: func(ctx context.Context, id string) (*user.Usuario, error) {
			if id == "123" {
				return &user.Usuario{ID: "123", Nome: "User"}, nil
			}
			return nil, nil
		},
	}
	svc := user.NewService(repo)
	handler := &UsersHandler{Svc: svc}

	// usuário existe
	req := httptest.NewRequest(http.MethodGet, "/usuarios/buscar-por-id/123", nil)
	w := httptest.NewRecorder()

	// Act
	handler.BuscarPorID(w, req)
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// usuário não existe
	req = httptest.NewRequest(http.MethodGet, "/usuarios/buscar-por-id/999", nil)
	w = httptest.NewRecorder()

	// Act
	handler.BuscarPorID(w, req)
	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestAtualizarUsuarioSucesso verifica se o usuário é atualizado com sucesso
func TestAtualizarUsuarioSucesso(t *testing.T) {
	// Arrange
	repo := &mockUserRepo{
		UpdateFunc: func(ctx context.Context, id string, u *user.Usuario) error { return nil },
	}
	svc := user.NewService(repo)
	handler := &UsersHandler{Svc: svc}

	update := user.Usuario{Nome: "Updated"}
	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPatch, "/usuarios/atualizar/123", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Act
	handler.Atualizar(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestBuscarTudo verifica se todos os usuários são buscados com sucesso
func TestBuscarTudoSucesso(t *testing.T) {
	// Arrange
	repo := &mockUserRepo{
		ListFunc: func(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]user.Usuario, int, error) {
			return []user.Usuario{{ID: "1"}, {ID: "2"}}, 2, nil
		},
	}
	svc := user.NewService(repo)
	handler := &UsersHandler{Svc: svc}

	req := httptest.NewRequest(http.MethodGet, "/usuarios/buscar-tudo", nil)
	w := httptest.NewRecorder()

	// Act
	handler.BuscarTudo(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestListaCompleta verifica se a lista completa de usuários é retornada com sucesso
func TestListaCompletaSucesso(t *testing.T) {
	// Arrange
	repo := &mockUserRepo{
		ListFunc: func(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]user.Usuario, int, error) {
			return []user.Usuario{{ID: "1"}, {ID: "2"}}, 2, nil
		},
	}
	svc := user.NewService(repo)
	handler := &UsersHandler{Svc: svc}

	req := httptest.NewRequest(http.MethodGet, "/usuarios/lista-completa", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ListaCompleta(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestBuscarTecnicos verifica se os técnicos são buscados com sucesso
func TestBuscarTecnicos(t *testing.T) {
	// Arrange
	repo := &mockUserRepo{
		ListFunc: func(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]user.Usuario, int, error) {
			return []user.Usuario{{ID: "1", Nome: "A"}, {ID: "2", Nome: "B"}}, 2, nil
		},
	}
	svc := user.NewService(repo)
	handler := &UsersHandler{Svc: svc}

	req := httptest.NewRequest(http.MethodGet, "/usuarios/buscar-tecnicos", nil)
	w := httptest.NewRecorder()

	// Act
	handler.BuscarTecnicos(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestDesativarAutorizar verifica se as ações de desativar e autorizar um usuário funcionam corretamente
func TestDesativarAutorizarSucesso(t *testing.T) {
	// Arrange
	repo := &mockUserRepo{
		UpdateFunc: func(ctx context.Context, id string, u *user.Usuario) error { return nil },
	}
	svc := user.NewService(repo)
	handler := &UsersHandler{Svc: svc}

	req := httptest.NewRequest(http.MethodDelete, "/usuarios/desativar/123", nil)
	w := httptest.NewRecorder()

	// Act
	handler.Desativar(w, req)
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest(http.MethodPatch, "/usuarios/autorizar/123", nil)
	w = httptest.NewRecorder()

	// Act
	handler.Autorizar(w, req)
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestValidaUsuario verifica se a validação do usuário funciona corretamente
func TestValidaUsuarioSucesso(t *testing.T) {
	// Arrange
	handler := &UsersHandler{}
	req := httptest.NewRequest(http.MethodGet, "/usuarios/valida-usuario", nil)
	w := httptest.NewRecorder()
	// Act
	handler.ValidaUsuario(w, req)
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestBuscarNovo verifica se um novo usuário pode ser buscado com sucesso
func TestBuscarNovoSucesso(t *testing.T) {
	// Arrange
	repo := &mockUserRepo{
		FindByLoginFunc: func(ctx context.Context, login string) (*user.Usuario, error) {
			return nil, user.ErrUsuarioNaoEncontrado
		},
	}
	svc := user.NewService(repo)
	ldap := &ldapauth.Client{
		ConnectFunc: func(user, pass string) (ldapauth.LDAPConnection, error) {
			return &mockLDAPConnection{}, nil
		},
		LoginAttr: "uid",
		Base:      "dc=teste,dc=com",
	}
	handler := &UsersHandler{Svc: svc, LDAP: ldap}

	req := httptest.NewRequest(http.MethodGet, "/usuarios/buscar-novo/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.BuscarNovo(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}
