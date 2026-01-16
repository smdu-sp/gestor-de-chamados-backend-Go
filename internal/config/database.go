package config

import (
	"strconv"
	"time"
)

const (
	dbPortDefault             = "3306"         // Porta padrão do banco de dados
	dbUserDefault             = "user"         // Usuário padrão do banco de dados
	dbPassDefault             = "userpassword" // Senha padrão do banco de dados
	dbNameDefault             = "mydatabase"   // Nome padrão do banco de dados
	dbMaxOpenConnsDefault     = 25             // Número máximo padrão de conexões abertas no pool
	minOpenConns              = 1              // Número mínimo de conexões abertas no pool
	minIdleConns              = 0              // Número mínimo de conexões ociosas no pool
	dbMaxIdleConnsDefault     = 10             // Número máximo padrão de conexões ociosas no pool
	dbConnMaxLifetimeDefault  = "30m"          // Tempo máximo padrão de vida de uma conexão
	dbConnMaxIdleTimeDefault  = "5m"           // Tempo máximo padrão que uma conexão pode ficar ociosa
	dbMaxRetryAttemptsDefault = 5              // Número máximo padrão de tentativas de conexão
	dbMigrationTimeoutDefault = "2m"           // Tempo máximo padrão para migrações do banco de dados
	minRetryAttempts          = 0              // Número mínimo de tentativas de conexão
)

// DatabaseConfig contém configurações do banco de dados
type DatabaseConfig struct {
	host             string        // Host do banco de dados
	port             string        // Porta do banco de dados
	user             string        // Usuário do banco de dados
	pass             string        // Senha do banco de dados
	name             string        // Nome do banco de dados
	maxRetryAttempts int           // Número máximo de tentativas de conexão
	maxOpenConns     int           // Número máximo de conexões abertas no pool
	maxIdleConns     int           // Número máximo de conexões ociosas no pool
	connMaxLifetime  time.Duration // Tempo máximo de vida de uma conexão
	connMaxIdleTime  time.Duration // Tempo máximo que uma conexão pode ficar ociosa
	migrationTimeout time.Duration // Tempo máximo para migrações do banco de dados
}

// carregarDatabaseConfig carrega as configurações do banco de dados
func carregarDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		host:             getEnv("DB_HOST", ipLocalhost),
		port:             getEnv("DB_PORT", dbPortDefault),
		user:             getEnv("DB_USER", dbUserDefault),
		pass:             getEnv("DB_PASS", dbPassDefault),
		name:             getEnv("DB_NAME", dbNameDefault),
		maxOpenConns:     getEnvInt("DB_MAX_OPEN_CONNS", dbMaxOpenConnsDefault),
		maxIdleConns:     getEnvInt("DB_MAX_IDLE_CONNS", dbMaxIdleConnsDefault),
		connMaxLifetime:  getEnvDuration("DB_CONN_MAX_LIFETIME", dbConnMaxLifetimeDefault),
		connMaxIdleTime:  getEnvDuration("DB_CONN_MAX_IDLE_TIME", dbConnMaxIdleTimeDefault),
		maxRetryAttempts: getEnvInt("DB_MAX_RETRY_ATTEMPTS", dbMaxRetryAttemptsDefault),
		migrationTimeout: getEnvDuration("DB_MIGRATION_TIMEOUT", dbMigrationTimeoutDefault),
	}
}

// Validar valida as configurações do banco de dados.
//
// Em caso de erros de validação, retorna um erro do tipo ErrosConfig.
func (d *DatabaseConfig) Validar() error {
	erros := NovoErrosConfig()

	// Campos obrigatórios
	obrigatorios := map[string]string{
		"DB_HOST": d.host,
		"DB_USER": d.user,
		"DB_NAME": d.name,
	}

	for campo, valor := range obrigatorios {
		if valor == "" {
			erros.Add(campo, "é obrigatório")
		}
	}

	// Conexões
	if d.maxOpenConns < minOpenConns {
		erros.Add("DB_MAX_OPEN_CONNS",
			"deve ser maior ou igual a "+strconv.Itoa(minOpenConns),
		)
	}

	if d.maxIdleConns < minIdleConns {
		erros.Add("DB_MAX_IDLE_CONNS",
			"deve ser maior ou igual a "+strconv.Itoa(minIdleConns),
		)
	}

	if d.maxOpenConns < d.maxIdleConns {
		erros.Add("DB_MAX_OPEN_CONNS",
			"deve ser maior ou igual a "+strconv.Itoa(d.maxIdleConns),
		)
	}

	if d.maxRetryAttempts < minRetryAttempts {
		erros.Add("DB_MAX_RETRY_ATTEMPTS",
			"deve ser maior ou igual a "+strconv.Itoa(minRetryAttempts),
		)
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// Métodos de acesso

func (c DatabaseConfig) Host() string                    { return c.host }
func (c DatabaseConfig) Port() string                    { return c.port }
func (c DatabaseConfig) User() string                    { return c.user }
func (c DatabaseConfig) Pass() string                    { return c.pass }
func (c DatabaseConfig) Name() string                    { return c.name }
func (c DatabaseConfig) MaxRetryAttempts() int           { return c.maxRetryAttempts }
func (c DatabaseConfig) MaxOpenConns() int               { return c.maxOpenConns }
func (c DatabaseConfig) MaxIdleConns() int               { return c.maxIdleConns }
func (c DatabaseConfig) ConnMaxLifetime() time.Duration  { return c.connMaxLifetime }
func (c DatabaseConfig) ConnMaxIdleTime() time.Duration  { return c.connMaxIdleTime }
func (c DatabaseConfig) MigrationTimeout() time.Duration { return c.migrationTimeout }
