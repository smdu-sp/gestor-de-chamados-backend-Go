package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
)

// AbrirConexao recebe as configurações e abre uma conexão com o banco de dados MySQL.
func AbrirConexao(config *cfg.Config) (*sql.DB, error) {
	dsn := montarDSN(config)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir conexão: %w", err)
	}

	configurarPool(db, config)

	if err := testarConexao(db, config); err != nil {
		FecharConexao(slog.Default(), db)
		return nil, fmt.Errorf("abrir conexão: %w", err)
	}

	slog.Info("Conexão com MySQL estabelecida com sucesso")
	return db, nil
}

// FecharConexao recebe o logger e a conexão com o banco de dados e a fecha.
func FecharConexao(logger *slog.Logger, db *sql.DB) {
	logger.Info("Fechando conexão com banco de dados")
	if err := db.Close(); err != nil {
		logger.Error("Erro ao fechar conexão com banco",
			"error", err,
			"database", "mysql",
		)
	} else {
		logger.Info("Conexão com banco fechada com sucesso")
	}
}

// configurarPool recebe a conexão com o banco de dados e as configurações e configura o pool de conexões.
func configurarPool(db *sql.DB, config *cfg.Config) {
	db.SetMaxOpenConns(config.DBMaxOpenConns())       // numero máximo de conexões abertas
	db.SetMaxIdleConns(config.DBMaxIdleConns())       // numero máximo de conexões ociosas
	db.SetConnMaxLifetime(config.DBConnMaxLifetime()) // tempo máximo de vida da conexão
	db.SetConnMaxIdleTime(config.DBConnMaxIdleTime()) // tempo máximo de inatividade da conexão
}

// montarDSN recebe as configurações e monta o DSN para conexão com o MySQL.
func montarDSN(config *cfg.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true"+ // parseTime
			"&multiStatements=true"+ // permite múltiplas statements
			"&charset=utf8mb4"+ // charset
			"&collation=utf8mb4_unicode_ci&loc=Local"+ // localização
			"&timeout=10s&readTimeout=30s&writeTimeout=30s"+ // timeouts
			"&interpolateParams=true"+ // melhora performance
			"&maxAllowedPacket=67108864"+ // 64MB
			"&tls=preferred", // segurança
		config.DBUser(), config.DBPass(),
		config.DBHost(), config.DBPort(),
		config.DBName(),
	)
}

// testarConexao recebe a conexão com o banco de dados e as configurações e testa a conexão com o banco.
func testarConexao(db *sql.DB, config *cfg.Config) error {
	maxRetries := config.DBMaxRetryAttempts()
	var ultimoErro error

	for i := range maxRetries {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err := db.PingContext(ctx)
		cancel()

		if err == nil {
			return nil
		}

		ultimoErro = fmt.Errorf("testar conexão: %w", err)

		if i < maxRetries-1 {
			espera := time.Duration(1<<i) * time.Second // backoff exponencial: 1s, 2s, 4s, 8s, ...

			slog.Warn("Falha ao conectar, tentando novamente...",
				slog.Int("tentativa", i+1), slog.Int("max_tentativas", maxRetries))

			time.Sleep(espera)
		}
	}

	return fmt.Errorf("falha ao conectar após %d tentativas: %w", maxRetries, ultimoErro)
}

// ExecutarSeed recebe o contexto, a conexão com o banco de dados e as configurações,
// e executa o script de seed no banco de dados.
func ExecutarSeed(ctx context.Context, db *sql.DB, config *cfg.Config) error {
	if !config.EmDesenvolvimento() {
		return fmt.Errorf("seed permitido apenas em ambiente de desenvolvimento")
	}

	seedPath := "internal/infra/migrations/seed.sql"

	seedSQL, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("ler seed.sql: %w", err)
	}

	if _, err := db.ExecContext(ctx, string(seedSQL)); err != nil {
		return fmt.Errorf("executar seed.sql: %w", err)
	}

	return nil
}