package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"

	goLdap "github.com/go-ldap/ldap/v3"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
)

var ErrUsuarioNaoEncontradoNoLDAP = fmt.Errorf("usuário não encontrado no LDAP")

// LDAPConnection é a interface que abstrai uma conexão LDAP do go-ldap
type ConexaoLDAP interface {
	Close() error
	Bind(username, password string) error
	Search(sr *goLdap.SearchRequest) (*goLdap.SearchResult, error)
	StartTLS(*tls.Config) error
}

// conexaoLDAPReal adapta goLdap.Conn para a interface LDAPConnection
type conexaoLDAPReal struct {
	*goLdap.Conn
}

func (c *conexaoLDAPReal) Pesquisar(pesquisa *goLdap.SearchRequest) (*goLdap.SearchResult, error) {
	return c.Conn.Search(pesquisa)
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
	ConnectFunc func(user, pass string) (ConexaoLDAP, error)
}

// Garantia de que Client implementa usecase.AuthExternoUsecase
var _ usecase.AuthExternoUsecase = (*Client)(nil)

// Bind autentica usuário no LDAP/AD
func (c *Client) Bind(login, senha string) error {
	bindUsuario := login
	if c.Domain != "" && login != "" && !strings.Contains(login, "@") {
		bindUsuario += c.Domain
	}

	ldapConn, err := c.conectar(bindUsuario, senha)
	if err != nil {
		return fmt.Errorf("[ldap.Bind]: %w", err)
	}
	defer ldapConn.Close()

	return nil
}

// PesquisarPorLogin busca usuário pelo atributo LoginAttr
func (c *Client) PesquisarPorLogin(login string) (nome, email, outLogin string, err error) {
	metodo := "[ldap.PesquisarPorLogin]: %w"
	ldapConn, err := c.conectar(c.UsuarioComDominio(), c.Pass)
	if err != nil {
		return "", "", "", fmt.Errorf(metodo, err)
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
		return "", "", "", fmt.Errorf(metodo, err)
	}

	if len(res.Entries) == 0 {
		return "", "", "", fmt.Errorf(metodo, ErrUsuarioNaoEncontradoNoLDAP)
	}

	entry := res.Entries[0]
	return entry.GetAttributeValue("cn"),
		entry.GetAttributeValue("mail"),
		entry.GetAttributeValue(c.LoginAttr),
		nil
}

// connect cria conexão LDAP e faz bind com usuário e senha
func (c *Client) conectar(user, pass string) (ConexaoLDAP, error) {
	if c.ConnectFunc != nil {
		return c.ConnectFunc(user, pass)
	}

	ldapConn, err := goLdap.DialURL(c.Server)
	if err != nil {
		return nil, err
	}

	conn := &conexaoLDAPReal{ldapConn}

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

// UsuarioComDominio retorna usuário completo para bind AD/OpenLDAP
func (c *Client) UsuarioComDominio() string {
	if c.Domain != "" {
		return c.User + c.Domain
	}
	return fmt.Sprintf("uid=%s,ou=users,%s", c.User, c.Base)
}
