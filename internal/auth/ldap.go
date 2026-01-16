package auth

import (
	"crypto/tls"
	"fmt"
	"strings"

	goLdap "github.com/go-ldap/ldap/v3"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
)

// UsuarioLDAP representa um usuário LDAP
type UsuarioLDAP struct {
	Nome  string
	Email string
	Login string
}

// NovoUsuarioLDAP cria uma nova instância de UsuarioLDAP
func NovoUsuarioLDAP(nome, email, login string) *UsuarioLDAP {
	return &UsuarioLDAP{
		Nome:  nome,
		Email: email,
		Login: login,
	}
}

// LDAPService implementa o Authenticator usando LDAP
type LDAPService struct {
	Server    string
	Domain    string
	Base      string
	User      string
	Pass      string
	LoginAttr string

	// ConnectFunc permite injeção de mock em testes
	ConnectFunc func(user, pass string) (ConexaoLDAP, error)
}

// NovoLDAPService cria uma nova instância de LDAPService
func NovoLDAPService(config *cfg.Config) *LDAPService {
	return &LDAPService{
		Server:    config.LDAPServer(),
		Domain:    config.LDAPDomain(),
		Base:      config.LDAPBase(),
		User:      config.LDAPUser(),
		Pass:      config.LDAPPass(),
		LoginAttr: config.LDAPLoginAttr(),
	}
}

// conexaoLDAPReal adapta goLdap.Conn para a interface ConexaoLDAP
type conexaoLDAPReal struct {
	*goLdap.Conn
}

// ConexaoLDAP é a interface que abstrai uma conexão LDAP do go-ldap
type ConexaoLDAP interface {
	Close() error
	Bind(username, password string) error
	Search(sr *goLdap.SearchRequest) (*goLdap.SearchResult, error)
	StartTLS(*tls.Config) error
}

// Pesquisar executa uma pesquisa LDAP
func (c *conexaoLDAPReal) Pesquisar(pesquisa *goLdap.SearchRequest) (*goLdap.SearchResult, error) {
	return c.Conn.Search(pesquisa)
}

// Asserção de interface para garantir que LDAPService implementa AuthExternoService
var _ AuthExternoService = (*LDAPService)(nil)

// Bind recebe login e senha, faz bind no servidor LDAP.
//
// Em caso de erro, retorna erro do go-ldap.
func (l *LDAPService) Bind(login, senha string) error {
	bindUsuario := login
	if l.Domain != "" && login != "" && !strings.Contains(login, "@") {
		bindUsuario += l.Domain
	}

	ldapConn, err := l.conectar(bindUsuario, senha)
	if err != nil {
		return fmt.Errorf("erro ao fazer bind no LDAP: %w", err)
	}
	defer ldapConn.Close()

	return nil
}

// PesquisarPorLogin busca usuário pelo atributo LoginAttr e retorna nome, email e login.
//
// Em caso de erro, retorna erro do go-ldap.
func (l *LDAPService) PesquisarPorLogin(login string) (usuarioLDAP *UsuarioLDAP, err error) {
	// 1 - Conectar no LDAP com usuário de serviço
	ldapConn, err := l.conectar(l.UsuarioComDominio(), l.Pass)
	if err != nil {
		return nil, fmt.Errorf("PesquisarPorLogin: %w", err)
	}
	defer ldapConn.Close()

	// 2 - Montar filtro de pesquisa
	filtro := fmt.Sprintf("(%s=%s)", l.LoginAttr, goLdap.EscapeFilter(login))

	// 3 - Executar pesquisa
	req := goLdap.NewSearchRequest(
		l.Base,
		goLdap.ScopeWholeSubtree,
		goLdap.NeverDerefAliases,
		0, 0, false,
		filtro,
		[]string{"cn", "mail", l.LoginAttr},
		nil,
	)

	// 4 - Obter resultado
	resultado, err := ldapConn.Search(req)
	if err != nil {
		return nil, fmt.Errorf("PesquisarPorLogin: %w", err)
	}

	// 5 - Verificar se encontrou usuário
	if len(resultado.Entries) == 0 {
		return nil, fmt.Errorf("usuário não encontrado no LDAP")
	}

	// 6 - Retornar dados do usuário
	dados := resultado.Entries[0]
	usuarioLDAP = NovoUsuarioLDAP(
		dados.GetAttributeValue("cn"),
		dados.GetAttributeValue("mail"),
		dados.GetAttributeValue(l.LoginAttr),
	)

	return usuarioLDAP, nil
}

// connect recebe usuário e senha, conecta no servidor LDAP e retorna a conexão autenticada.
//
// Em caso de erro, retorna erro do go-ldap.
func (l *LDAPService) conectar(user, pass string) (ConexaoLDAP, error) {
	// 1 - Usar função mock se definida
	if l.ConnectFunc != nil {
		return l.ConnectFunc(user, pass)
	}

	// 2 - Conectar no servidor LDAP
	ldapConn, err := goLdap.DialURL(l.Server)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no LDAP: %w", err)
	}

	// 3 - Adaptar conexão para interface
	conn := &conexaoLDAPReal{ldapConn}

	// 4 - Iniciar TLS se necessário
	if strings.HasPrefix(l.Server, "ldaps") {
		if err := conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			conn.Close()
			return nil, fmt.Errorf("erro ao iniciar TLS na conexão LDAP: %w", err)
		}
	}

	// 5 - Fazer bind com usuário e senha
	if err := conn.Bind(user, pass); err != nil {
		conn.Close()
		return nil, fmt.Errorf("erro ao fazer bind no LDAP: %w", err)
	}

	// 6 - Retornar conexão autenticada
	return conn, nil
}

// UsuarioComDominio retorna usuário completo para bind AD/OpenLDAP
func (l *LDAPService) UsuarioComDominio() string {
	if l.Domain != "" {
		return l.User + l.Domain
	}

	return fmt.Sprintf("uid=%s,ou=users,%s", l.User, l.Base)
}
