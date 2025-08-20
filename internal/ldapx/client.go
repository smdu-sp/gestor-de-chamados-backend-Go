package ldapx

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"
)

type Client struct {
	cfg config.LDAP
}

func New(cfg config.LDAP) *Client { return &Client{cfg: cfg} }

// Authenticate:
// 1) Faz bind de serviço (opcional) para localizar DN do usuário via search.
// 2) Faz bind com o DN do usuário + senha do login fornecido.
// 3) Retorna dados mapeados para o domínio User (sem persistir).
func (c *Client) Authenticate(login, password string) (*model.User, error) {

	fmt.Println("URL:", c.cfg.URL)
	fmt.Println("BaseDN:", c.cfg.BaseDN)
	fmt.Println("BindDN:", c.cfg.BindDN)
	fmt.Println("UserFilter:", c.cfg.UserFilter)

	if password == "" {
		return nil, fmt.Errorf("senha vazia")
	}

	conn, err := ldap.DialURL(c.cfg.URL)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if c.cfg.UseTLS {
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: c.cfg.InsecureSkipVerify})

		if err != nil {
			return nil, err
		}
	}

	// Bind de serviço (se informado) para procurar DN do usuário
	if c.cfg.BindDN != "" {
		if err = conn.Bind(c.cfg.BindDN, c.cfg.BindPassword); err != nil {
			return nil, fmt.Errorf("bind de serviço falhou: %w", err)
		}
	}

	filter := fmt.Sprintf("(|(uid=%s)(sAMAccountName=%s)(mail=%s))", login, login, login)
	fmt.Println("LDAP filter gerado:", filter)

	req := ldap.NewSearchRequest(
		c.cfg.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{
			c.cfg.AttrLogin,
			c.cfg.AttrName,
			c.cfg.AttrEmail,
			c.cfg.AttrAvatar,
			c.cfg.AttrPerm,
		},
		nil,
	)

	res, err := conn.Search(req)

	if err != nil {
		return nil, err
	}

	if len(res.Entries) == 0 {
		return nil, fmt.Errorf("usuário não encontrado")
	}

	entry := res.Entries[0]
	userDN := entry.DN
	// Valida credenciais do usuário pelo bind do próprio usuário
	if err = conn.Bind(userDN, password); err != nil {
		return nil, fmt.Errorf("credenciais inválidas: %w", err)
	}

	loginAttr := entry.GetAttributeValue(c.cfg.AttrLogin)
	nameAttr := entry.GetAttributeValue(c.cfg.AttrName)
	emailAttr := entry.GetAttributeValue(c.cfg.AttrEmail)
	avatarAttr := entry.GetAttributeValue(c.cfg.AttrAvatar)
	permAttr := entry.GetAttributeValue(c.cfg.AttrPerm)

	var avatarPtr *string

	if avatarAttr != "" {
		avatarPtr = &avatarAttr
	}

	if permAttr == "" {
		permAttr = string(model.USR)
	}

	return &model.User{
		ID:           "", // preenchido no repositório
		Nome:         nameAttr,
		Login:        loginAttr,
		Email:        emailAttr,
		Status:       true,
		Avatar:       avatarPtr,
		UltimoLogin:  nil,
		CriadoEm:     (func() (t time.Time) { return t })(),
		AtualizadoEm: (func() (t time.Time) { return t })(),
		Permissao:    model.Permissao(permAttr),
	}, nil
}
