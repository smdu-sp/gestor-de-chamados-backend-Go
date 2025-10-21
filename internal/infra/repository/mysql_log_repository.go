package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// Erros sentinela do repositório de logs
var (
	ErrLogNaoEncontrado = errors.New("log não encontrado no banco de dados MySQL")
	ErrScannerLog       = errors.New("erro ao escanear log do banco de dados MySQL")
)

// MySQLLogRepository implementa a interface LogRepository para MySQL.
type MySQLLogRepository struct {
	db *sql.DB
}

// NewMySQLLogRepository cria uma nova instância de MySQLLogRepository.
func NewMySQLLogRepository(db *sql.DB) *MySQLLogRepository {
	return &MySQLLogRepository{db: db}
}

// BuscarPorID busca um log pelo seu ID.
func (r *MySQLLogRepository) BuscarPorID(ctx context.Context, id string) (*model.Log, error) {
	usuario, err := r.Buscar(
		ctx,
		`SELECT id, usuario_id, acao, entidade, detalhes, criado_em
		FROM logs 
		WHERE id = ?`,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("[MySQLLogRepository.BuscarPorID]: %w", err)
	}

	if usuario == nil {
		return nil, utils.NewAppError(
			"[MySQLLogRepository.BuscarPorID]",
			utils.LevelInfo,
			"a busca por ID não retornou resultados",
			ErrLogNaoEncontrado,
		)
	}
	
	return usuario, nil
}

// Salvar insere um novo log no banco de dados.
func (r *MySQLLogRepository) Salvar(ctx context.Context, l *model.Log) error {
	const metodo = "[MySQLLogRepository.Salvar]"

	resultado, err := r.db.ExecContext(
		ctx,
		`INSERT INTO logs (
		id, usuario_id, acao, entidade, detalhes, criado_em
		) VALUES (?, ?, ?, ?, ?, NOW())`,
		l.ID, l.UsuarioID, l.Acao, l.Entidade, l.Detalhes,
	)
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro inesperado ao salvar log no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao obter o número de linhas afetadas ao salvar log no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrRowsAffected, err),
		)
	}
	if linhasAfetadas == 0 {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"nenhuma linha foi afetada ao tentar salvar o log",
			ErrExecContext,
		)
	}

	return nil
}

// Listar retorna uma lista paginada de logs, com filtros opcionais.
func (r *MySQLLogRepository) Listar(ctx context.Context, filtro model.LogFiltro) ([]model.Log, int, error) {
	var query strings.Builder
	args := []any{}

	query.WriteString(
		`SELECT SQL_CALC_FOUND_ROWS
		id, usuario_id, acao, entidade, detalhes, criado_em
		FROM logs 
		WHERE 1=1`,
	)

	if filtro.Busca != nil && *filtro.Busca != "" {
		padrao := "%" + *filtro.Busca + "%"
		query.WriteString(" AND (detalhes LIKE ?)")
		args = append(args, padrao)
	}

	if filtro.UsuarioID != nil && *filtro.UsuarioID != "" {
		query.WriteString(" AND usuario_id = ?")
		args = append(args, *filtro.UsuarioID)
	}

	if filtro.Acao != nil && *filtro.Acao != "" {
		query.WriteString(" AND acao = ?")
		args = append(args, *filtro.Acao)
	}

	if filtro.Entidade != nil && *filtro.Entidade != "" {
		query.WriteString(" AND entidade = ?")
		args = append(args, *filtro.Entidade)
	}

	if filtro.DataInicio != nil {
		query.WriteString(" AND criado_em >= ?")
		args = append(args, *filtro.DataInicio)
	}

	if filtro.DataFim != nil {
		query.WriteString(" AND criado_em <= ?")
		args = append(args, *filtro.DataFim)
	}

	query.WriteString(" ORDER BY criado_em DESC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite, (filtro.Pagina-1)*filtro.Limite)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLLogRepository.Listar]",
			utils.LevelError,
			"erro inesperado ao listar logs no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	defer rows.Close()

	var logs []model.Log
	for rows.Next() {
		log, err := scanLog(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("[MySQLLogRepository.Listar]: %w", err)
		}
		logs = append(logs, *log)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLLogRepository.Listar]",
			utils.LevelError,
			"erro ao obter o número total de logs",
			fmt.Errorf(utils.FmtErroWrap, ErrScan, err),
		)
	}

	return logs, total, nil
}

// Métodos auxiliares

// Buscar executa uma consulta que retorna um log.
func (r *MySQLLogRepository) Buscar(ctx context.Context, query string, args ...any) (*model.Log, error) {
row := r.db.QueryRowContext(ctx, query, args...)
	log, err := scanLog(row)
	if err != nil {
		return nil, fmt.Errorf("[MySQLLogRepository.Buscar]: %w", err)
	}
	return log, nil
}

// scanUsuario mapeia os dados de um scanner (row e rows) para uma struct Log.
func scanLog(scanner interface{ Scan(dest ...any) error }) (*model.Log, error) {
	var log model.Log
	err := scanner.Scan(
		&log.ID,
		&log.UsuarioID,
		&log.Acao,
		&log.Entidade,
		&log.Detalhes,
		&log.CriadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.NewAppError(
			"[MySQLLogRepository.scanLog]",
			utils.LevelError,
			"o scanner falhou ao escanear o log",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerLog, err),
		)
	}
	return &log, nil
}
