package config

// LDAPConfig contém configurações do LDAP
type LDAPConfig struct {
	server    string // Servidor LDAP
	domain    string // Domínio LDAP
	base      string // Base DN do LDAP
	user      string // Usuário de ligação LDAP
	pass      string // Senha do usuário de ligação LDAP
	loginAttr string // Atributo de login LDAP
}

// carregarLDAPConfig carrega as configurações do LDAP
func carregarLDAPConfig() LDAPConfig {
	return LDAPConfig{
		server:    getEnv("LDAP_SERVER", ""),
		domain:    getEnv("LDAP_DOMAIN", ""),
		base:      getEnv("LDAP_BASE", ""),
		user:      getEnv("LDAP_USER", ""),
		pass:      getEnv("LDAP_PASS", ""),
		loginAttr: getEnv("LDAP_LOGIN_ATTR", "uid"),
	}
}

// Validar valida as configurações do LDAP.
//
// Em caso de erros de validação, retorna um erro do tipo ErrosConfig.
func (l *LDAPConfig) Validar() error {
	erros := NovoErrosConfig()

	// Se a base estiver definida, os outros campos são obrigatórios
	if l.base != "" {
		mensagemErro := "é obrigatório quando LDAP_BASE está definido"

		campos := map[string]string{
			"LDAP_SERVER":     l.server,
			"LDAP_DOMAIN":     l.domain,
			"LDAP_USER":       l.user,
			"LDAP_PASS":       l.pass,
			"LDAP_LOGIN_ATTR": l.loginAttr,
		}

		for nome, valor := range campos {
			if valor == "" {
				erros.Add(nome, mensagemErro)
			}
		}
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// Métodos de acesso

func (l LDAPConfig) Server() string    { return l.server }
func (l LDAPConfig) Domain() string    { return l.domain }
func (l LDAPConfig) Base() string      { return l.base }
func (l LDAPConfig) User() string      { return l.user }
func (l LDAPConfig) Pass() string      { return l.pass }
func (l LDAPConfig) LoginAttr() string { return l.loginAttr }
