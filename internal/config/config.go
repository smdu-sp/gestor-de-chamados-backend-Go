package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Estrutura para configuração de conexão com LDAP.
type LDAP struct {
	URL                string // URL do servidor LDAP (ex.: ldap://localhost:389)
	BaseDN             string // DN base da árvore LDAP (ex.: dc=example,dc=org)
	BindDN             string // DN usado para autenticar e realizar buscas (ex.: cn=admin,dc=example,dc=org)
	BindPassword       string // senha do usuário BindDN
	UserFilter         string `env:"LDAP_USER_FILTER"` // filtro para localizar usuários (ex.: "(|(uid=%s)(mail=%s))")
	UseTLS             bool   // habilita StartTLS/LDAPS
	InsecureSkipVerify bool   // ignora verificação de certificado (apenas em ambiente de dev)
	AttrLogin          string // atributo que representa o login (uid ou sAMAccountName)
	AttrName           string // atributo que representa o nome do usuário (cn, displayName)
	AttrEmail          string // atributo que representa o e-mail
	AttrAvatar         string // atributo opcional para foto/avatar
	AttrPerm           string // atributo usado para determinar a permissão/papel (ex.: department, memberOf)
}

type Config struct {
	Addr      string        // ":8080"
	JWTSecret string        // segredo HMAC
	JWTIssuer string        // emissor
	JWTTTL    time.Duration // tempo de vida
	LDAP      LDAP
}

// getenv retorna o valor de uma variável de ambiente
// ou um valor padrão caso ela não esteja definida.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// mustDuration valida valores de duração.
func mustDuration(envKey, def string) time.Duration {
	v := getenv(envKey, def)
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("config: %s inválido: %v", envKey, err)
	}
	return d
}

// mustBool valida valores booleanos.
func mustBool(envKey string, def bool) bool {
	v := os.Getenv(envKey)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatalf("config: %s inválido: %v", envKey, err)
	}
	return b
}

func Load() Config {
	host := getenv("LDAP_HOST", "localhost")
	port := getenv("LDAP_PORT", "389")
	url := func() string {
		// Se LDAP_URL existir, usa; senão monta a partir de host:port
		if url := getenv("LDAP_URL", ""); url != "" {
			return url
		}
		return fmt.Sprintf("ldap://%s:%s", host, port)
	}()

	return Config{
		Addr:      getenv("HTTP_ADDR", ":8080"),
		JWTSecret: getenv("JWT_SECRET", "dev-secret-change-me"),
		JWTIssuer: getenv("JWT_ISSUER", "authapi"),
		JWTTTL:    mustDuration("JWT_TTL", "2h"),
		LDAP: LDAP{
			URL:                url,
			BaseDN:             getenv("LDAP_BASE_DN", "dc=example,dc=org"),
			BindDN:             getenv("LDAP_BIND_DN", ""),
			BindPassword:       getenv("LDAP_BIND_PASS", ""),
			UserFilter:         getenv("LDAP_USER_FILTER", ""),
			UseTLS:             mustBool("LDAP_TLS", false),
			InsecureSkipVerify: mustBool("LDAP_INSECURE", true),
			AttrLogin:          getenv("LDAP_ATTR_LOGIN", "uid"),
			AttrName:           getenv("LDAP_ATTR_NAME", "cn"),
			AttrEmail:          getenv("LDAP_ATTR_EMAIL", "mail"),
			AttrAvatar:         getenv("LDAP_ATTR_AVATAR", ""),
			AttrPerm:           getenv("LDAP_ATTR_PERM", "department"),
		},
	}
}
