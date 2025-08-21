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

func New(cfg config.LDAP) *Client {
	return &Client{cfg: cfg}
}

// Realiza a autenticação do usuário no LDAP
// Passos:
// 1) Bind de serviço (opcional) para localizar DN do usuário via search.
// 2) Bind com DN do usuário + senha fornecida para validar credenciais.
// 3) Retorna um model.User mapeado (não persiste no banco local).
func (c *Client) Authenticate(login, password string) (*model.User, error) {

	if password == "" {
		return nil, fmt.Errorf("senha vazia")
	}

	// Conecta ao servidor LDAP
	conn, err := ldap.DialURL(c.cfg.URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Inicia TLS se configurado
	if c.cfg.UseTLS {
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: c.cfg.InsecureSkipVerify})
		if err != nil {
			return nil, err
		}
	}

	// Bind de serviço (opcional) para poder procurar DN do usuário
	if c.cfg.BindDN != "" {
		if err = conn.Bind(c.cfg.BindDN, c.cfg.BindPassword); err != nil {
			return nil, fmt.Errorf("bind de serviço falhou: %w", err)
		}
	}

	// Monta filtro LDAP para localizar usuário (uid, sAMAccountName ou email)
	filter := fmt.Sprintf(
		"(|(uid=%s)(sAMAccountName=%s)(mail=%s))", 
		login, 
		login, 
		login,
	)

	// Cria a requisição de busca
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

	// Executa a busca
	res, err := conn.Search(req)
	if err != nil {
		return nil, err
	}

	// Verifica se algum usuário foi encontrado
	if len(res.Entries) == 0 {
		return nil, fmt.Errorf("usuário não encontrado")
	}

	// Pega a primeira entrada encontrada
	entry := res.Entries[0]
	userDN := entry.DN

	// Valida credenciais do usuário pelo bind do próprio usuário
	if err = conn.Bind(userDN, password); err != nil {
		return nil, fmt.Errorf("credenciais inválidas: %w", err)
	}

	// Mapeia atributos LDAP para model.User
	loginAttr := entry.GetAttributeValue(c.cfg.AttrLogin)
	nameAttr := entry.GetAttributeValue(c.cfg.AttrName)
	emailAttr := entry.GetAttributeValue(c.cfg.AttrEmail)
	avatarAttr := entry.GetAttributeValue(c.cfg.AttrAvatar)
	permAttr := entry.GetAttributeValue(c.cfg.AttrPerm)

	var avatarPtr *string
	if avatarAttr != "" {
		avatarPtr = &avatarAttr
	}

	// Se permissão não estiver definida, assume padrão USR
	if permAttr == "" {
		permAttr = string(model.USR)
	}

	// Retorna o usuário mapeado (sem persistência)
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
