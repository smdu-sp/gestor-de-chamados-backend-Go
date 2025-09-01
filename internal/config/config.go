package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config armazena as configurações da aplicação
type Config struct {
	Port          string // Porta onde o servidor irá escutar
	Env           string // Ambiente: local, development, production
	CORSOrigin    string // Origem permitida para CORS
	DBHost        string // Host do banco de dados
	DBPort        string // Porta do banco de dados
	DBUser        string // Usuário do banco de dados
	DBPass        string // Senha do banco de dados
	DBName        string // Nome do banco de dados
	JWTSecret     string // Segredo para assinar JWTs
	RTSecret      string // Segredo para assinar Refresh Tokens
	AccessTTL     string // Tempo de vida do Access Token
	RefreshTTL    string // Tempo de vida do Refresh Token
	LDAPServer    string // Endereço do servidor LDAP
	LDAPDomain    string // Domínio LDAP
	LDAPBase      string // Base DN para buscas LDAP
	LDAPUser      string // Usuário para bind no LDAP
	LDAPPass      string // Senha para bind no LDAP
	LDAPLoginAttr string // Atributo usado para login (ex: uid, cn, mail)
}

// Load carrega as configurações do ambiente ou usa valores padrão
func Load() Config {
	// Tenta carregar o .env
	if err := godotenv.Load(); err != nil {
		log.Println("Não foi possível carregar o .env, usando variáveis de ambiente do sistema")
	}

	// Carrega as variáveis de ambiente com valores padrão
	cfg := Config{
		Port:          getenv("PORT", "8080"),
		Env:           getenv("ENVIRONMENT", "local"),
		CORSOrigin:    getenv("CORS_ORIGIN", ""),
		DBHost:        getenv("DB_HOST", "127.0.0.1"),
		DBPort:        getenv("DB_PORT", "3306"),
		DBUser:        getenv("DB_USER", "user"),
		DBPass:        getenv("DB_PASS", "userpassword"),
		DBName:        getenv("DB_NAME", "mydatabase"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		RTSecret:      os.Getenv("RT_SECRET"),
		AccessTTL:     getenv("ACCESS_TTL", "24h"),
		RefreshTTL:    getenv("REFRESH_TTL", "168h"),
		LDAPServer:    getenv("LDAP_SERVER", ""),
		LDAPDomain:    getenv("LDAP_DOMAIN", ""),
		LDAPBase:      getenv("LDAP_BASE", ""),
		LDAPUser:      getenv("LDAP_USER", ""),
		LDAPPass:      getenv("LDAP_PASS", ""),
		LDAPLoginAttr: getenv("LDAP_LOGIN_ATTR", "uid"),
	}

	if cfg.JWTSecret == "" || cfg.RTSecret == "" {
		log.Println("[aviso] defina JWT_SECRET e RT_SECRET no .env")
	}

	return cfg
}

// getenv retorna o valor da variável de ambiente ou o valor padrão se não estiver definida
func getenv(k, def string) string {
	if valor := os.Getenv(k); valor != "" {
		return valor
	}
	return def
}
