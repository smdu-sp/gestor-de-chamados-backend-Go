package ldapauth

import (
	"crypto/tls"
	"errors"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
)

// mockLDAPConn simula uma conexão LDAP para testes
type mockLDAPConn struct {
	bindCalled    bool
	bindUser      string
	bindPass      string
	bindErr       error
	searchRequest *ldap.SearchRequest
	searchResult  *ldap.SearchResult
	searchErr     error
	closeCalled   bool
}

// Bind simula a operação de bind em uma conexão LDAP
func (m *mockLDAPConn) Bind(username, password string) error {
	m.bindCalled = true
	m.bindUser = username
	m.bindPass = password
	return m.bindErr
}

// Search simula a operação de busca em uma conexão LDAP
func (m *mockLDAPConn) Search(sr *ldap.SearchRequest) (*ldap.SearchResult, error) {
	m.searchRequest = sr
	return m.searchResult, m.searchErr
}

// StartTLS simula a operação de StartTLS em uma conexão LDAP
func (m *mockLDAPConn) StartTLS(_ *tls.Config) error {
	// não usado em teste
	return nil
}

// Close simula a operação de fechamento em uma conexão LDAP
func (m *mockLDAPConn) Close() error {
	m.closeCalled = true
	return nil
}

const (
	ldapServer   = "ldap://localhost:389"
	ldapDomain   = "@rede.sp"
	ldapBase     = "dc=rede,dc=sp"
	ldapUser     = "admin"
	ldapPass     = "adminpass"
	ldapLoginAttr = "uid"
)

// ----------------------
// Testes Bind
// ----------------------

// Teste de Bind com sucesso
func TestClientBindSuccess(t *testing.T) {
	// Arrange
	mock := &mockLDAPConn{}
	client := &Client{
		Server: ldapServer,
		Domain: ldapDomain,
		ConnectFunc: func(user, pass string) (LDAPConnection, error) {
			mock.bindCalled = true
			mock.bindUser = user
			mock.bindPass = pass
			return mock, nil
		},
	}

	// Act
	err := client.Bind("usuario1", "senha123")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "usuario1@rede.sp", mock.bindUser)
	assert.Equal(t, "senha123", mock.bindPass)
	assert.True(t, mock.closeCalled)
}

func TestClientBindFail(t *testing.T) {
	// Arrange
	client := &Client{
		Server: ldapServer,
		Domain: ldapDomain,
		ConnectFunc: func(user, pass string) (LDAPConnection, error) {
			return nil, errors.New("erro no bind")
		},
	}

	// Act
	err := client.Bind("usuario1", "senha123")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "falha no bind LDAP")
}

// Teste de conexão com erro
func TestClientConnectError(t *testing.T) {
	// Arrange
	client := &Client{
		ConnectFunc: func(user, pass string) (LDAPConnection, error) {
			return nil, errors.New("erro na conexão")
		},
	}

	// Act
	err := client.Bind("usuario1", "senha123")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "falha no bind LDAP")
}

// ----------------------
// Testes SearchByLogin
// ----------------------

// Teste de busca por login com sucesso
func TestSearchByLoginSuccess(t *testing.T) {
	// Arrange
	mock := &mockLDAPConn{
		searchResult: &ldap.SearchResult{
			Entries: []*ldap.Entry{
				{
					Attributes: []*ldap.EntryAttribute{
						{Name: "cn", Values: []string{"Usuário Teste"}},
						{Name: "mail", Values: []string{"usuario@teste.com"}},
						{Name: "uid", Values: []string{"usuario1"}},
					},
				},
			},
		},
	}
	client := &Client{
		Base:      ldapBase,
		User:      ldapUser,
		Pass:      ldapPass,
		LoginAttr: ldapLoginAttr,
		ConnectFunc: func(user, pass string) (LDAPConnection, error) {
			return mock, nil
		},
	}

	// Act
	nome, email, login, err := client.SearchByLogin("usuario1")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Usuário Teste", nome)
	assert.Equal(t, "usuario@teste.com", email)
	assert.Equal(t, "usuario1", login)
	assert.True(t, mock.closeCalled)
}

// Teste de busca por login não encontrado
func TestSearchByLoginUserNotFound(t *testing.T) {
	// Arrange
	mock := &mockLDAPConn{
		searchResult: &ldap.SearchResult{Entries: []*ldap.Entry{}},
	}
	client := &Client{
		Base:      ldapBase,
		User:      ldapUser,
		Pass:      ldapPass,
		LoginAttr: ldapLoginAttr,
		ConnectFunc: func(user, pass string) (LDAPConnection, error) {
			return mock, nil
		},
	}

	// Act
	nome, email, login, err := client.SearchByLogin("naoexiste")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usuário não encontrado")
	assert.Empty(t, nome)
	assert.Empty(t, email)
	assert.Empty(t, login)
}

// Teste de busca por login com erro
func TestSearchByLoginSearchError(t *testing.T) {
	// Arrange
	mock := &mockLDAPConn{
		searchErr: errors.New("erro na pesquisa"),
	}
	client := &Client{
		Base:      ldapBase,
		User:      ldapUser,
		Pass:      ldapPass,
		LoginAttr: ldapLoginAttr,
		ConnectFunc: func(user, pass string) (LDAPConnection, error) {
			return mock, nil
		},
	}

	// Act
	nome, email, login, err := client.SearchByLogin("usuario1")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro na pesquisa LDAP")
	assert.Empty(t, nome)
	assert.Empty(t, email)
	assert.Empty(t, login)
}

// Teste de busca por login com erro
func TestSearchByLoginConnectError(t *testing.T) {
	// Arrange
	client := &Client{
		Base:      ldapBase,
		User:      ldapUser,
		Pass:      ldapPass,
		LoginAttr: ldapLoginAttr,
		ConnectFunc: func(user, pass string) (LDAPConnection, error) {
			return nil, errors.New("falha na conexão")
		},
	}

	// Act
	nome, email, login, err := client.SearchByLogin("usuario1")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro de conexão do LDAP")
	assert.Empty(t, nome)
	assert.Empty(t, email)
	assert.Empty(t, login)
}

// ----------------------
// Testes UserWithDomain
// ----------------------

// Teste de busca por login com sucesso
func TestUserWithDomain(t *testing.T) {
	// Teste com domínio
	t.Run("com domínio", func(t *testing.T) {
		// Arrange
		client := &Client{User: "admin", Domain: "@rede.sp"}
		// Act
		assert.Equal(t, "admin@rede.sp", client.UserWithDomain())
	})

	// Teste de busca por login sem domínio
	t.Run("sem domínio", func(t *testing.T) {
		// Arrange
		client := &Client{User: "admin", Base: "dc=rede,dc=sp"}
		// Assert
		assert.Equal(t, "uid=admin,ou=users,dc=rede,dc=sp", client.UserWithDomain())
	})

	// Teste de busca por login com erro
	t.Run("vazio", func(t *testing.T) {
		// Arrange
		client := &Client{}
		// Assert
		assert.Equal(t, "uid=,ou=users,", client.UserWithDomain())
	})
}

// ----------------------
// Testes Close do mock
// ----------------------

// Teste de fechamento do mock
func TestMockCloseCalled(t *testing.T) {
	// Arrange
	mock := &mockLDAPConn{}
	// Act
	_ = mock.Close()
	// Assert
	assert.True(t, mock.closeCalled)
}
