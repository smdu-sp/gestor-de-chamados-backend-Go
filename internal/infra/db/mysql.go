package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
)

// OpenMySQL abre uma conexão com o MySQL usando a configuração fornecida.
func OpenMySQL(cfg config.Config) (*sql.DB, error) {
	dsn := buildDSN(cfg)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão com MySQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao conectar no MySQL: %w", err)
	}

	return db, nil
}

// buildDSN monta o Data Source Name para conexão MySQL
func buildDSN(cfg config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
}
