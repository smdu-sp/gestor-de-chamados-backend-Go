package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	GoVersion       = "1.25.0"                // GoVersion define a versão mínima do Go necessária
	TimeoutPadrao   = 5 * time.Second         // TimeoutPadrao define um tempo padrão de timeout para operações
	asterisco       = "*"                     // asterisco representa qualquer origem em CORS
	stringVazia     = ""                      // stringVazia representa uma string vazia
	localhost8080   = "http://localhost:8080" // localhost8080 representa o localhost na porta 8080
	ipLocalhost8080 = "http://127.0.0.1:8080" // ipLocalhost8080 representa o IP localhost na porta 8080
	ipLocalhost     = "127.0.0.1"             // ipLocalhost representa o IP localhost
	localhost       = "localhost"             // localhost representa o nome localhost
	production      = "production"            // production representa o ambiente de produção
	development     = "development"           // development representa o ambiente de desenvolvimento
	APP             = "APP"                   // APP representa a seção de configuração da aplicação
	SERVER          = "SERVER"                // SERVER representa a seção de configuração do servidor
	DATABASE        = "DATABASE"              // DATABASE representa a seção de configuração do banco de dados
	AUTH            = "AUTH"                  // AUTH representa a seção de configuração de autenticação
	LDAP            = "LDAP"                  // LDAP representa a seção de configuração do LDAP
	PRODUCAO        = "EM PRODUÇÃO"           // PRODUCAO representa a seção de validação para ambiente de produção
)

// Config representa todas as configurações carregadas da aplicação
type Config struct {
	app      AppConfig      // Configuração da aplicação
	server   ServerConfig   // Configuração do servidor
	database DatabaseConfig // Configuração do banco de dados
	auth     AuthConfig     // Configuração de autenticação
	ldap     LDAPConfig     // Configuração do LDAP
}

// Carregar carrega as configurações do ambiente ou usa valores padrão.
//
// Em caso de erros de validação, retorna um erro do tipo ErrosConfig.
func Carregar() (*Config, error) {
	// Tenta carregar .env (não é crítico se não existir)
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("falha ao carregar .env: %w", err)
	}

	config := Config{
		app:      carregarAppConfig(),
		server:   carregarServerConfig(),
		database: carregarDatabaseConfig(),
		auth:     carregarAuthConfig(),
		ldap:     carregarLDAPConfig(),
	}

	if err := config.validar(); err != nil {
		return nil, err
	}

	return &config, nil
}

// validar verifica todas as configurações.
//
// Em caso de erros de validação, retorna um erro do tipo ErrosConfig.
func (c *Config) validar() error {
	erros := NovoErrosConfig()

	if err := c.app.Validar(); err != nil {
		erros.Add(APP, err.Error())
	}

	if err := c.server.Validar(); err != nil {
		erros.Add(SERVER, err.Error())
	}

	if err := c.database.Validar(); err != nil {
		erros.Add(DATABASE, err.Error())
	}

	if err := c.auth.Validar(); err != nil {
		erros.Add(AUTH, err.Error())
	}

	if err := c.ldap.Validar(); err != nil {
		erros.Add(LDAP, err.Error())
	}

	if err := c.ValidarEmProducao(); err != nil {
		erros.Add(PRODUCAO, err.Error())
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// ValidarEmProduçao valida as restrições adicionais para ambiente de produção.
//
// Em caso de erros de validação, retorna um erro do tipo ErrosConfig.
func (c *Config) ValidarEmProducao() error {
	if !c.EmProducao() {
		return nil
	}

	erros := NovoErrosConfig()

	origensRestritas := map[string]bool{
		asterisco:       true,
		stringVazia:     true,
		localhost8080:   true,
		ipLocalhost8080: true,
	}

	origin := c.AppCORSOrigin()

	if origensRestritas[origin] {
		erros.Add("CORS_ORIGIN", "não pode ser "+asterisco+", vazio, "+localhost8080+" ou "+ipLocalhost8080+" em produção")
	}

	if c.DBHost() == localhost || c.DBHost() == ipLocalhost {
		erros.Add("DB_HOST", "não pode usar "+localhost+" ou "+ipLocalhost+" em produção")
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// =============================================================================
// Métodos de acesso às configurações
// =============================================================================

//------------------- APP CONFIG --------------------------------------

func (c Config) AppName() string         { return c.app.Name() }
func (c Config) AppVersion() string      { return c.app.Version() }
func (c Config) AppEnvironment() string  { return c.app.Environment() }
func (c Config) AppLogLevel() string     { return c.app.LogLevel() }
func (c Config) AppPort() string         { return c.app.Port() }
func (c Config) AppCORSOrigin() string   { return c.app.CORSOrigin() }
func (c Config) EmProducao() bool        { return c.app.Environment() == production }
func (c Config) EmDesenvolvimento() bool { return c.app.Environment() == development }

//------------------- SERVER CONFIG -----------------------------------

func (c Config) ServerReadTimeout() time.Duration     { return c.server.ReadTimeout() }
func (c Config) ServerWriteTimeout() time.Duration    { return c.server.WriteTimeout() }
func (c Config) ServerIdleTimeout() time.Duration     { return c.server.IdleTimeout() }
func (c Config) ServerShutdownTimeout() time.Duration { return c.server.ShutdownTimeout() }
func (c Config) ServerSignalChannelBufferSize() int   { return c.server.SignalChannelBufferSize() }
func (c Config) ServerErrorChannelBufferSize() int    { return c.server.ErrorChannelBufferSize() }

//------------------- DATABASE CONFIG ---------------------------------

func (c Config) DBHost() string                    { return c.database.Host() }
func (c Config) DBPort() string                    { return c.database.Port() }
func (c Config) DBUser() string                    { return c.database.User() }
func (c Config) DBPass() string                    { return c.database.Pass() }
func (c Config) DBName() string                    { return c.database.Name() }
func (c Config) DBMaxRetryAttempts() int           { return c.database.MaxRetryAttempts() }
func (c Config) DBMaxOpenConns() int               { return c.database.MaxOpenConns() }
func (c Config) DBMaxIdleConns() int               { return c.database.MaxIdleConns() }
func (c Config) DBConnMaxLifetime() time.Duration  { return c.database.ConnMaxLifetime() }
func (c Config) DBConnMaxIdleTime() time.Duration  { return c.database.ConnMaxIdleTime() }
func (c Config) DBMigrationTimeout() time.Duration { return c.database.MigrationTimeout() }

//------------------- AUTH CONFIG -------------------------------------

func (c Config) TokenSecret() string        { return c.auth.TokenSecret() }
func (c Config) RefreshTokenSecret() string { return c.auth.RefreshTokenSecret() }
func (c Config) TokenTTL() string           { return c.auth.TokenTTL() }
func (c Config) RefreshTokenTTL() string    { return c.auth.RefreshTokenTTL() }

//------------------- LDAP CONFIG -------------------------------------

func (c Config) LDAPServer() string    { return c.ldap.Server() }
func (c Config) LDAPDomain() string    { return c.ldap.Domain() }
func (c Config) LDAPBase() string      { return c.ldap.Base() }
func (c Config) LDAPUser() string      { return c.ldap.User() }
func (c Config) LDAPPass() string      { return c.ldap.Pass() }
func (c Config) LDAPLoginAttr() string { return c.ldap.LoginAttr() }

// =============================================================================
// Funções auxiliares
// =============================================================================

// getEnv retorna o valor da variável de ambiente
// ou o valor padrão se não estiver definida
func getEnv(chave, fallback string) string {
	if valor := os.Getenv(chave); valor != "" {
		return valor
	}
	return fallback
}

// getEnvInt retorna o valor inteiro da variável de ambiente
// ou o valor padrão se não estiver definida ou inválida
func getEnvInt(chave string, fallback int) int {
	valor := os.Getenv(chave)
	if valor == "" {
		return fallback
	}

	valorConv, err := strconv.Atoi(valor)
	if err != nil {
		return fallback
	}

	return valorConv
}

// getEnvDuration retorna o valor time.Duration da variável de ambiente
// ou o valor padrão se não estiver definida ou inválida
func getEnvDuration(chave, fallback string) time.Duration {
	valor := getEnv(chave, fallback)

	duracao, err := time.ParseDuration(valor)
	if err != nil {
		slog.Warn(
			"valor inválido para duração, usando fallback",
			"campo", chave,
			"valor", valor,
			"erro", err,
		)

		duracao, _ = time.ParseDuration(fallback)
	}

	return duracao
}
