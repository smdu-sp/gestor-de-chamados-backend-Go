package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	lg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
)

// Erros sentinela do repositório de logs
var ErrLogNaoEncontrado = errors.New("log não encontrado no banco de dados")

// LogRepository implementa a interface LogRepository para MySQL.
type LogRepository struct {
	db *sql.DB
}

// NovoLogRepository cria uma nova instância de LogRepository.
func NovoLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{db: db}
}

// Asserção de interface para garantir que LogRepository implementa LogRepository
var _ lg.Repository = (*LogRepository)(nil)

// BuscarPorID recebe um ID e retorna o log correspondente.
//
// Erros sentinela possíveis: ErrLogNaoEncontrado.
func (r *LogRepository) BuscarPorID(ctx context.Context, id string) (*lg.Log, error) {
	query :=
		`SELECT id, usuario_id, acao, entidade, detalhes, criado_em
		 FROM logs 
		 WHERE id = ?`

	usr, err := r.Buscar(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("buscar por ID no repo: %w", err)
	}

	if usr == nil {
		return nil, ErrLogNaoEncontrado
	}

	return usr, nil
}

// Criar recebe um log e o insere no repositório.
func (r *LogRepository) Criar(ctx context.Context, log lg.Log) error {
	query :=
		`INSERT INTO logs (
		 id, usuario_id, acao, entidade, detalhes, criado_em
		 ) VALUES (?, ?, ?, ?, ?, ?)`

	resultado, err := r.db.ExecContext(ctx, query,
		log.ID(),
		log.UsuarioID(),
		log.Acao(),
		log.Entidade(),
		log.Detalhes(),
		log.CriadoEm(),
	)
	if err != nil {
		return fmt.Errorf("criar no repo: %w", err)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("criar no repo: %w", err)
	}
	if linhasAfetadas == 0 {
		return fmt.Errorf("criar no repo: nenhuma linha afetada")
	}

	return nil
}

// Listar recebe um filtro e retorna uma lista de logs que correspondem aos critérios do filtro,
// juntamente com o total de registros encontrados.
func (r *LogRepository) Listar(ctx context.Context, filtro lg.LogFiltro) ([]lg.Log, int, error) {
	query, args := r.construirQueryListar(filtro)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}
	defer rows.Close()

	// Scaneia os resultados
	logSlice, err := scanRows(rows, r.scanLog)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	// Obtém o total de registros
	total, err := obterTotalRegistros(ctx, r.db)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	return logSlice, total, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// Buscar recebe uma query e argumentos, executa a busca e retorna o log correspondente.
func (r *LogRepository) Buscar(ctx context.Context, query string, args ...any) (*lg.Log, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	log, err := r.scanLog(row)
	if err != nil {
		return nil, fmt.Errorf("buscar: %w", err)
	}
	return log, nil
}

// scanLog recebe um scanner (row ou rows) e retorna um log.
func (r *LogRepository) scanLog(scanner scanner) (*lg.Log, error) {
	var logDB lg.LogDB
	err := scanner.Scan(
		&logDB.ID,
		&logDB.UsuarioID,
		&logDB.Acao,
		&logDB.Entidade,
		&logDB.Detalhes,
		&logDB.CriadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	return lg.CarregarDoBD(logDB), nil
}

// construirQueryListar recebe um filtro e constrói a query SQL correspondente.
func (r *LogRepository) construirQueryListar(filtro lg.LogFiltro) (string, []any) {
	var query strings.Builder
	args := []any{}

	query.WriteString(
		`SELECT SQL_CALC_FOUND_ROWS
		id, usuario_id, acao, entidade, detalhes, criado_em
		FROM logs 
		WHERE 1=1`,
	)

	if filtro.Busca() != nil && *filtro.Busca() != "" {
		padrao := "%" + *filtro.Busca() + "%"
		query.WriteString(" AND (detalhes LIKE ?)")
		args = append(args, padrao)
	}

	if filtro.UsuarioID() != nil && *filtro.UsuarioID() != "" {
		query.WriteString(" AND usuario_id = ?")
		args = append(args, *filtro.UsuarioID())
	}

	if filtro.Acao() != nil && *filtro.Acao() != "" {
		query.WriteString(" AND acao = ?")
		args = append(args, *filtro.Acao())
	}

	if filtro.Entidade() != nil && *filtro.Entidade() != "" {
		query.WriteString(" AND entidade = ?")
		args = append(args, *filtro.Entidade())
	}

	if filtro.DataInicio() != nil {
		query.WriteString(" AND criado_em >= ?")
		args = append(args, *filtro.DataInicio())
	}

	if filtro.DataFim() != nil {
		query.WriteString(" AND criado_em <= ?")
		args = append(args, *filtro.DataFim())
	}

	query.WriteString(" ORDER BY criado_em DESC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite(), filtro.Offset())
	return query.String(), args
}
