package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type LDAP struct {
	URL                string
	BaseDN             string
	BindDN             string
	BindPassword       string
	UserFilter         string `env:"LDAP_USER_FILTER"`
	UseTLS             bool
	InsecureSkipVerify bool   // somente para ambiente dev
	AttrLogin          string // ex: uid ou sAMAccountName
	AttrName           string // ex: cn ou displayName
	AttrEmail          string // ex: mail
	AttrAvatar         string // opcional
	AttrPerm           string
}

type Config struct {
	Addr      string        // ":8080"
	JWTSecret string        // segredo HMAC
	JWTIssuer string        // emissor
	JWTTTL    time.Duration // tempo de vida
	LDAP      LDAP
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustDuration(envKey, def string) time.Duration {
	v := getenv(envKey, def)
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("config: %s inválido: %v", envKey, err)
	}
	return d
}

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
	return Config{
		Addr:      getenv("HTTP_ADDR", ":8080"),
		JWTSecret: getenv("JWT_SECRET", "dev-secret-change-me"),
		JWTIssuer: getenv("JWT_ISSUER", "authapi"),
		JWTTTL:    mustDuration("JWT_TTL", "2h"),
		LDAP: LDAP{
			URL:                getenv("LDAP_URL", "ldap://localhost:389"),
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
