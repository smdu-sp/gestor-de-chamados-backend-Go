package ldapauth

import (
	"crypto/tls"
	"fmt"
	"strings"

	goLdap "github.com/go-ldap/ldap/v3"
)

// LDAPConnection é a interface que abstrai uma conexão LDAP
type LDAPConnection interface {
	Close() error
	Bind(username, password string) error
	Search(sr *goLdap.SearchRequest) (*goLdap.SearchResult, error)
	StartTLS(*tls.Config) error
}

// LDAPInterface é a interface que define os métodos que a implementação LDAP deve fornecer
type LDAPInterface interface {
	Bind(login, senha string) error
	SearchByLogin(login string) (nome, email, outLogin string, err error)
}

// realLDAPConn adapta goLdap.Conn para a interface LDAPConnection
type realLDAPConn struct {
	*goLdap.Conn
}

// Search implementa a busca de usuários no LDAP
func (c *realLDAPConn) Search(sr *goLdap.SearchRequest) (*goLdap.SearchResult, error) {
	return c.Conn.Search(sr)
}

// Client representa cliente LDAP configurado
type Client struct {
	Server    string // ex: ldap://localhost:389
	Domain    string // ex: @rede.sp (AD)
	Base      string // ex: DC=rede,DC=sp
	User      string // bind user
	Pass      string // bind password
	LoginAttr string // ex: "uid" para OpenLDAP, "sAMAccountName" para AD

	// ConnectFunc permite injeção de mock em testes
	ConnectFunc func(user, pass string) (LDAPConnection, error)
}

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
	fmt.Println("Tentando conectar ao LDAP com user:", c.UserWithDomain())
	ldapConn, err := c.connect(c.UserWithDomain(), c.Pass)
	if err != nil {
		return "", "", "", fmt.Errorf("erro de conexão do LDAP: %w", err)
	}
	defer ldapConn.Close()

	// Filtro de pesquisa
	filter := fmt.Sprintf("(%s=%s)", c.LoginAttr, goLdap.EscapeFilter(login))
	fmt.Println("Filtro LDAP usado:", filter)
	fmt.Println("Base DN:", c.Base)

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

	fmt.Println("Número de entradas encontradas:", len(res.Entries))

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
	// Se for teste, usa ConnectFunc mockada
	if c.ConnectFunc != nil {
		return c.ConnectFunc(user, pass)
	}

	ldapConn, err := goLdap.DialURL(c.Server)
	if err != nil {
		return nil, err
	}

	conn := &realLDAPConn{ldapConn}

	// Usar StartTLS se for ldaps
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
