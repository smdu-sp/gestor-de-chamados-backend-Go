package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)


// Erros sentinela comuns de banco de dados
var ErrConexaoMySQL = errors.New("erro de conexão com o banco de dados MySQL")

// ConectarMySQL abre uma conexão com o MySQL usando a configuração fornecida.
func ConectarMySQL(cfg config.Config) (*sql.DB, error) {
	dsn := buildDSN(cfg)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, utils.NewAppError(
			"[db.ConectarMySQL]",
			utils.LevelError,
			"falha ao abrir conexão com MySQL",
			fmt.Errorf(utils.FmtErroWrap, ErrConexaoMySQL, err),
		)
	}

	if err := db.Ping(); err != nil {
		return nil, utils.NewAppError(
			"[db.ConectarMySQL]",
			utils.LevelError,
			"falha ao tentar conectar com MySQL",
			fmt.Errorf(utils.FmtErroWrap, ErrConexaoMySQL, err),
		)
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
