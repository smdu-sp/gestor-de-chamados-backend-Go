package ldapauth

import (
	"crypto/tls"
	"fmt"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
)

type Client struct {
	Server    string // ex: ldap://localhost:389
	Domain    string // ex: @rede.sp (AD)
	Base      string // ex: DC=rede,DC=sp
	User      string // bind user
	Pass      string // bind password
	LoginAttr string // ex: "uid" para OpenLDAP, "sAMAccountName" para AD
}

// Bind autentica usuário no LDAP/AD
func (c *Client) Bind(login, senha string) error {
	bindUser := login
	if c.Domain != "" && login != "" && !strings.Contains(login, "@") {
		bindUser += c.Domain
	}

	l, err := c.connect(bindUser, senha)
	if err != nil {
		return fmt.Errorf("falha no bind LDAP: %w", err)
	}
	defer l.Close()

	return nil
}

// SearchByLogin busca usuário pelo atributo LoginAttr
func (c *Client) SearchByLogin(login string) (nome, email, outLogin string, err error) {
	l, err := c.connect(c.UserWithDomain(), c.Pass)
	if err != nil {
		return "", "", "", fmt.Errorf("erro bind service account: %w", err)
	}
	defer l.Close()

	// Filtro de pesquisa
	filter := fmt.Sprintf("(%s=%s)", c.LoginAttr, ldap.EscapeFilter(login))
	req := ldap.NewSearchRequest(
		c.Base,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{"cn", "mail", c.LoginAttr},
		nil,
	)

	res, err := l.Search(req)
	if err != nil {
		return "", "", "", fmt.Errorf("erro na pesquisa LDAP: %w", err)
	}
	if len(res.Entries) == 0 {
		return "", "", "", fmt.Errorf("usuário não encontrado")
	}

	entry := res.Entries[0]
	return entry.GetAttributeValue("cn"), entry.GetAttributeValue("mail"), entry.GetAttributeValue(c.LoginAttr), nil
}

// connect cria conexão LDAP e faz bind com usuário e senha
func (c *Client) connect(user, pass string) (*ldap.Conn, error) {
	l, err := ldap.DialURL(c.Server)
	if err != nil {
		return nil, err
	}

	// Usar StartTLS se for ldaps
	if strings.HasPrefix(c.Server, "ldaps") {
		if err := l.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			l.Close()
			return nil, err
		}
	}

	if err := l.Bind(user, pass); err != nil {
		l.Close()
		return nil, err
	}

	return l, nil
}

// UserWithDomain retorna usuário completo para bind AD/OpenLDAP
func (c *Client) UserWithDomain() string {
	if c.Domain != "" {
		return c.User + c.Domain
	}
	return fmt.Sprintf("uid=%s,ou=users,%s", c.User, c.Base)
}
