package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"

	goLdap "github.com/go-ldap/ldap/v3"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
)

// LDAPConnection é a interface que abstrai uma conexão LDAP
type LDAPConnection interface {
	Close() error
	Bind(username, password string) error
	Search(sr *goLdap.SearchRequest) (*goLdap.SearchResult, error)
	StartTLS(*tls.Config) error
}

// realLDAPConn adapta goLdap.Conn para a interface LDAPConnection
type realLDAPConn struct {
	*goLdap.Conn
}

func (c *realLDAPConn) Search(sr *goLdap.SearchRequest) (*goLdap.SearchResult, error) {
	return c.Conn.Search(sr)
}

// Client implementa o Authenticator usando LDAP
type Client struct {
	Server    string
	Domain    string
	Base      string
	User      string
	Pass      string
	LoginAttr string

	// ConnectFunc permite injeção de mock em testes
	ConnectFunc func(user, pass string) (LDAPConnection, error)
}

// Garantia de que Client implementa usecase.Authenticator
var _ usecase.Authenticator = (*Client)(nil)

// Bind autentica usuário no LDAP/AD
func (c *Client) Bind(login, senha string) error {
	bindUser := login
	if c.Domain != "" && login != "" && !strings.Contains(login, "@") {
		bindUser += c.Domain
	}

	ldapConn, err := c.connect(bindUser, senha)
	if err != nil {
		return fmt.Errorf("falha no bind LDAP: %w", err)
	}
	defer ldapConn.Close()

	return nil
}

// SearchByLogin busca usuário pelo atributo LoginAttr
func (c *Client) SearchByLogin(login string) (nome, email, outLogin string, err error) {
	ldapConn, err := c.connect(c.UserWithDomain(), c.Pass)
	if err != nil {
		return "", "", "", fmt.Errorf("erro de conexão do LDAP: %w", err)
	}
	defer ldapConn.Close()

	filter := fmt.Sprintf("(%s=%s)", c.LoginAttr, goLdap.EscapeFilter(login))

	req := goLdap.NewSearchRequest(
		c.Base,
		goLdap.ScopeWholeSubtree,
		goLdap.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{"cn", "mail", c.LoginAttr},
		nil,
	)

	res, err := ldapConn.Search(req)
	if err != nil {
		return "", "", "", fmt.Errorf("erro na pesquisa LDAP: %w", err)
	}

	if len(res.Entries) == 0 {
		return "", "", "", fmt.Errorf("usuário não encontrado")
	}

	entry := res.Entries[0]
	return entry.GetAttributeValue("cn"),
		entry.GetAttributeValue("mail"),
		entry.GetAttributeValue(c.LoginAttr),
		nil
}

// connect cria conexão LDAP e faz bind com usuário e senha
func (c *Client) connect(user, pass string) (LDAPConnection, error) {
	if c.ConnectFunc != nil {
		return c.ConnectFunc(user, pass)
	}

	ldapConn, err := goLdap.DialURL(c.Server)
	if err != nil {
		return nil, err
	}

	conn := &realLDAPConn{ldapConn}

	if strings.HasPrefix(c.Server, "ldaps") {
		if err := conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			conn.Close()
			return nil, err
		}
	}

	if err := conn.Bind(user, pass); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// UserWithDomain retorna usuário completo para bind AD/OpenLDAP
func (c *Client) UserWithDomain() string {
	if c.Domain != "" {
		return c.User + c.Domain
	}
	return fmt.Sprintf("uid=%s,ou=users,%s", c.User, c.Base)
}
