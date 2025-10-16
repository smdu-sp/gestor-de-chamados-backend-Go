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

var (
	ErrScannerAcompanhamento       = errors.New("erro ao escanear acompanhamento do banco de dados MySQL")
	ErrAcompanhamentoNaoEncontrado = errors.New("acompanhamento não encontrado no banco de dados MySQL")
)

// MySQLAcompanhamentoRepository é a implementação do repositório de acompanhamento para MySQL.
type MySQLAcompanhamentoRepository struct {
	db *sql.DB
}

// NewMySQLAcompanhamentoRepository cria uma nova instância do repositório de acompanhamento para MySQL.
func NewMySQLAcompanhamentoRepository(db *sql.DB) *MySQLAcompanhamentoRepository {
	return &MySQLAcompanhamentoRepository{db: db}
}

// BuscarPorID retorna um acompanhamento pelo seu ID.
func (r *MySQLAcompanhamentoRepository) BuscarPorID(ctx context.Context, id string) (*model.Acompanhamento, error) {
	acompanhamento, err := r.buscar(
		ctx,
		`SELECT id, conteudo, chamado_id, usuario_id, 
		remetente, criado_em, atualizado_em 
		FROM acompanhamentos 
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("[MySQLAcompanhamentoRepository.BuscarPorID]: %w", err)
	}

	if acompanhamento == nil {
		return nil, utils.NewAppError(
			"[MySQLAcompanhamentoRepository.BuscarPorID]",
			utils.LevelInfo,
			"a busca por ID do não retornou resultados",
			ErrAcompanhamentoNaoEncontrado,
		)
	}

	return acompanhamento, nil
}

// BuscarPorChamadoID retorna uma lista de acompanhamentos pelo ID do chamado.
func (r *MySQLAcompanhamentoRepository) BuscarPorChamadoID(ctx context.Context, chamadoID string) ([]model.Acompanhamento, error) {
	const metodo = "[MySQLAcompanhamentoRepository.BuscarPorChamadoID]"

	query := `
		SELECT 
			id, conteudo, chamado_id, usuario_id, remetente, criado_em, atualizado_em
		FROM acompanhamentos
		WHERE chamado_id = ?
		ORDER BY criado_em ASC
	`
	rows, err := r.db.QueryContext(ctx, query, chamadoID)
	if err != nil {
		return nil, utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao buscar acompanhamentos por chamadoID no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	defer rows.Close()

	var acompanhamentos []model.Acompanhamento
	for rows.Next() {
		acompanhamento, err := scanAcompanhamento(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", metodo, err)
		}
		acompanhamentos = append(acompanhamentos, *acompanhamento)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao iterar sobre os resultados de acompanhamentos",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}

	return acompanhamentos, nil
}

// Salvar insere um novo acompanhamento no repositório.
func (r *MySQLAcompanhamentoRepository) Salvar(ctx context.Context, a *model.Acompanhamento) error {
	const metodo = "[MySQLAcompanhamentoRepository.Salvar]"

	resultado, err := r.db.ExecContext(
		ctx,
		`INSERT INTO acompanhamentos (
		id, conteudo, chamado_id, usuario_id, remetente, criado_em, atualizado_em
		) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
		a.ID, a.Conteudo, a.ChamadoID, a.UsuarioID, a.Remetente,
	)
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro inesperado ao salvar o acompanhamento no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao obter o número de linhas afetadas ao salvar o acompanhamento no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrRowsAffected, err),
		)
	}
	if linhasAfetadas == 0 {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"nenhuma linha foi afetada ao salvar o acompanhamento",
			ErrExecContext,
		)
	}

	return nil
}

// Atualizar modifica os dados de um acompanhamento existente.
func (r *MySQLAcompanhamentoRepository) Atualizar(ctx context.Context, id string, a *model.Acompanhamento) error {
	const metodo = "[MySQLAcompanhamentoRepository.Atualizar]"

	existe, err := ExisteAcompanhamentoPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("%s: %w", metodo, err)
	}
	if !existe {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"o acompanhamento a ser atualizado não foi encontrado",
			ErrAcompanhamentoNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE acompanhamentos 
		SET conteudo = ?, chamado_id = ?, usuario_id = ?, remetente = ?, atualizado_em = NOW() 
		WHERE id = ?`,
		a.Conteudo, a.ChamadoID, a.UsuarioID, a.Remetente, id,
	)
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro inesperado ao atualizar o acompanhamento no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Deletar remove um acompanhamento pelo seu ID.
func (r *MySQLAcompanhamentoRepository) Deletar(ctx context.Context, id string) error {
	existe, err := ExisteAcompanhamentoPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLAcompanhamentoRepository.Deletar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLAcompanhamentoRepository.Deletar]",
			utils.LevelInfo,
			"o acompanhamento a ser deletado não foi encontrado",
			ErrAcompanhamentoNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`DELETE FROM acompanhamentos
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLAcompanhamentoRepository.Deletar]",
			utils.LevelError,
			"erro inesperado ao deletar o acompanhamento no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Listar retorna uma lista paginada de acompanhamentos, com filtros opcionais.
func (r *MySQLAcompanhamentoRepository) Listar(ctx context.Context, filtro model.AcompanhamentoFiltro) ([]model.Acompanhamento, int, error) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`
		SELECT SQL_CALC_FOUND_ROWS
			id, conteudo, chamado_id, usuario_id, remetente, criado_em, atualizado_em
		FROM acompanhamentos
		WHERE 1=1`)

	if filtro.ChamadoID != nil && *filtro.ChamadoID != "" {
		query.WriteString(" AND chamado_id = ?")
		args = append(args, *filtro.ChamadoID)
	}

	if filtro.UsuarioID != nil && *filtro.UsuarioID != "" {
		query.WriteString(" AND usuario_id = ?")
		args = append(args, *filtro.UsuarioID)
	}

	query.WriteString(" ORDER BY criado_em ASC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite, (filtro.Pagina-1)*filtro.Limite)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLAcompanhamentoRepository.Listar]",
			utils.LevelError,
			"erro ao listar acompanhamentos no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	defer rows.Close()

	var acompanhamentos []model.Acompanhamento
	for rows.Next() {
		acompanhamento, err := scanAcompanhamento(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("[MySQLAcompanhamentoRepository.Listar]: %w", err)
		}
		acompanhamentos = append(acompanhamentos, *acompanhamento)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLAcompanhamentoRepository.Listar]",
			utils.LevelError,
			"erro ao obter total de acompanhamentos",
			fmt.Errorf(utils.FmtErroWrap, ErrScan, err),
		)
	}

	return acompanhamentos, total, nil
}

// Métodos auxiliares

// buscar executa uma consulta que retorna um único acompanhamento.
func (r *MySQLAcompanhamentoRepository) buscar(ctx context.Context, query string, args ...any) (*model.Acompanhamento, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	acompanhamento, err := scanAcompanhamento(row)
	if err != nil {
		return nil, fmt.Errorf("[MySQLAcompanhamentoRepository.buscar]: %w", err)
	}
	return acompanhamento, nil
}

// ExisteAcompanhamentoPorID verifica se um acompanhamento existe pelo seu ID.
func ExisteAcompanhamentoPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM acompanhamentos WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, utils.NewAppError(
			"[MySQLAcompanhamentoRepository.ExisteAcompanhamentoPorID]",
			utils.LevelError,
			"falha ao verificar existência do acompanhamento",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerAcompanhamento, err),
		)
	}
	if !exists {
		return false, nil
	}
	return true, nil
}

// scanAcompanhamento mapeia os dados de um scanner (row ou rows) para uma struct Acompanhamento.
func scanAcompanhamento(scanner interface{ Scan(dest ...any) error }) (*model.Acompanhamento, error) {
	var acompanhamento model.Acompanhamento
	err := scanner.Scan(
		&acompanhamento.ID,
		&acompanhamento.Conteudo,
		&acompanhamento.ChamadoID,
		&acompanhamento.UsuarioID,
		&acompanhamento.Remetente,
		&acompanhamento.CriadoEm,
		&acompanhamento.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.NewAppError(
			"[MySQLAcompanhamentoRepository.scanAcompanhamento]",
			utils.LevelError,
			"o scanner falhou ao escanear o acompanhamento",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerAcompanhamento, err),
		)
	}
	return &acompanhamento, nil
}
