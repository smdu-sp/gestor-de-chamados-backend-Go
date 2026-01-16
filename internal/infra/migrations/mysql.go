package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

type MySQLMigrador struct {
	db *sql.DB
	fs fs.FS
}

// Asserção de interface para garantir que MySQLMigrador implementa Migrador.
var _ Migrador = (*MySQLMigrador)(nil)

// NovoMySQLMigrador cria um migrador para MySQL.
func NovoMySQLMigrador(db *sql.DB, migrationsFS fs.FS) *MySQLMigrador {
	return &MySQLMigrador{
		db: db,
		fs: migrationsFS,
	}
}

// Migrar executa as migrações pendentes.
func (m *MySQLMigrador) Migrar(ctx context.Context) error {
	if err := m.garantirSchemaMigracoes(ctx); err != nil {
		return err
	}

	aplicadas, err := m.migracoesAplicadas(ctx)
	if err != nil {
		return err
	}

	disponiveis, err := m.migracoesDisponiveis()
	if err != nil {
		return err
	}

	for _, migracao := range disponiveis {
		if aplicadas[migracao] {
			continue
		}

		if err := m.aplicarMigracao(ctx, migracao); err != nil {
			return err
		}
	}

	return nil
}

// garantirSchemaMigracoes garante que a tabela de controle de migrações exista.
func (m *MySQLMigrador) garantirSchemaMigracoes(ctx context.Context) error {
	const query = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version     VARCHAR(255) NOT NULL PRIMARY KEY,
		applied_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := m.db.ExecContext(ctx, query)
	return err
}

// migracoesAplicadas retorna um mapa das migrações já aplicadas.
func (m *MySQLMigrador) migracoesAplicadas(ctx context.Context) (map[string]bool, error) {
	rows, err := m.db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aplicadas := make(map[string]bool)

	for rows.Next() {
		var versao string
		if err := rows.Scan(&versao); err != nil {
			return nil, err
		}
		aplicadas[versao] = true
	}

	return aplicadas, rows.Err()
}

// migracoesDisponiveis retorna a lista de migrações disponíveis no sistema de arquivos.
func (m *MySQLMigrador) migracoesDisponiveis() ([]string, error) {
	entries, err := fs.ReadDir(m.fs, ".")
	if err != nil {
		return nil, err
	}

	var migracoes []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		nome := entry.Name()
		if strings.HasSuffix(nome, ".sql") {
			migracoes = append(migracoes, nome)
		}
	}

	sort.Strings(migracoes)
	return migracoes, nil
}

// aplicarMigracao aplica uma migração específica.
func (m *MySQLMigrador) aplicarMigracao(ctx context.Context, nome string) error {
	sqlBytes, err := fs.ReadFile(m.fs, nome)
	if err != nil {
		return fmt.Errorf("ler migração %s: %w", nome, err)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, string(sqlBytes)); err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf(
				"rollback da migração %s: %v (erro original: %w)",
				nome,
				rbErr,
				err,
			)
		}
		return fmt.Errorf("executar migração %s: %w", nome, err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO schema_migrations (version) VALUES (?)`,
		nome,
	); err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf(
				"rollback da migração %s: %v (erro original: %w)",
				nome,
				rbErr,
				err,
			)
		}
		return err
	}

	return tx.Commit()
}
