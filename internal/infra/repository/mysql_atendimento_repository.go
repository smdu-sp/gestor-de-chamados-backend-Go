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
	ErrScannerAtendimento       = errors.New("erro ao escanear atendimento do banco de dados MySQL")
	ErrAtendimentoNaoEncontrado = errors.New("atendimento não encontrado no banco de dados MySQL")
)

// MySQLAtendimentoRepository é a implementação do repositório de atendimentos para MySQL.
type MySQLAtendimentoRepository struct {
	db *sql.DB
}

// NewMySQLAtendimentoRepository cria uma nova instância de MySQLAtendimentoRepository.
func NewMySQLAtendimentoRepository(db *sql.DB) *MySQLAtendimentoRepository {
	return &MySQLAtendimentoRepository{db: db}
}

// BuscarPorID retorna um atendimento pelo seu ID.
func (r *MySQLAtendimentoRepository) BuscarPorID(ctx context.Context, id string) (*model.Atendimento, error) {
	atendimento, err := r.buscar(
		ctx,
		`SELECT id, atribuido_id, chamado_id, criado_em, atualizado_em
		FROM atendimentos 
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("[MySQLAtendimentoRepository.BuscarPorID]: %w", err)
	}

	if atendimento == nil {
		return nil, utils.NewAppError(
			"[MySQLAtendimentoRepository.BuscarPorID]",
			utils.LevelInfo,
			"a busca por ID do atendimento não retornou resultados",
			ErrAtendimentoNaoEncontrado,
		)
	}

	return atendimento, nil
}

// Salvar insere um novo atendimento no repositório.
func (r *MySQLAtendimentoRepository) Salvar(ctx context.Context, a *model.Atendimento) error {
	const metodo = "[MySQLAtendimentoRepository.Salvar]: %w"

	existeChamado, err := ExisteChamadoPorID(ctx, r.db, a.ChamadoID)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}
	if !existeChamado {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"o chamado associado ao atendimento não existe",
			ErrChamadoNaoEncontrado,
		)
	}

	existeTecnico, err := ExisteUsuarioPorID(ctx, r.db, a.AtribuidoID)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}
	if !existeTecnico {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"o técnico atribuído ao atendimento não existe",
			ErrUsuarioNaoEncontrado,
		)
	}

	resultado, err := r.db.ExecContext(
		ctx,
		`INSERT INTO atendimentos (
		id, atribuido_id, chamado_id, criado_em, atualizado_em
		) VALUES (?, ?, ?, NOW(), NOW())`,
		a.ID, a.AtribuidoID, a.ChamadoID,
	)
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao salvar o atendimento no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao obter número de linhas afetadas após salvar o atendimento",
			fmt.Errorf(utils.FmtErroWrap, ErrRowsAffected, err),
		)
	}
	if linhasAfetadas == 0 {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"nenhuma linha foi afetada ao salvar o atendimento",
			ErrExecContext,
		)
	}

	return nil
}

// Atualizar modifica os dados de um atendimento existente.
func (r *MySQLAtendimentoRepository) Atualizar(ctx context.Context, id string, a *model.Atendimento) error {
	const metodo = "[MySQLAtendimentoRepository.Atualizar]"

	existe, err := ExisteAtendimentoPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("%s: %w", metodo, err)
	}
	if !existe {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"o atendimento que se tentou atualizar não existe",
			ErrAtendimentoNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE atendimentos
		SET atribuido_id = ?, chamado_id = ?, atualizado_em = NOW()
		WHERE id = ?`,
		a.AtribuidoID, a.ChamadoID, id,
	)
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao atualizar o atendimento no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Listar retorna uma lista de atendimentos com base em filtros e paginação.
func (r *MySQLAtendimentoRepository) Listar(ctx context.Context, filtro model.AtendimentoFiltro) ([]model.Atendimento, int, error) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`
		SELECT SQL_CALC_FOUND_ROWS
		id, atribuido_id, chamado_id, criado_em, atualizado_em
		FROM atendimentos
		WHERE 1=1
	`)

	if filtro.ChamadoID != nil && *filtro.ChamadoID != "" {
		query.WriteString(" AND chamado_id = ?")
		args = append(args, *filtro.ChamadoID)
	}

	if filtro.AtribuidoID != nil && *filtro.AtribuidoID != "" {
		query.WriteString(" AND atribuido_id = ?")
		args = append(args, *filtro.AtribuidoID)
	}

	query.WriteString(" ORDER BY criado_em DESC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite, (filtro.Pagina-1)*filtro.Limite)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLAtendimentoRepository.Listar]",
			utils.LevelError,
			"erro ao listar os atendimentos no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	defer rows.Close()

	var atendimentos []model.Atendimento
	for rows.Next() {
		atendimento, err := scanAtendimento(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("[MySQLAtendimentoRepository.Listar]: %w", err)
		}
		atendimentos = append(atendimentos, *atendimento)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLAtendimentoRepository.Listar]",
			utils.LevelError,
			"erro ao obter total de atendimentos",
			fmt.Errorf(utils.FmtErroWrap, ErrScan, err),
		)
	}
	return atendimentos, total, nil
}

// Metodos auxiliares

// buscar executa uma consulta que retorna um único atendimento.
func (r *MySQLAtendimentoRepository) buscar(ctx context.Context, query string, args ...any) (*model.Atendimento, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	atendimento, err := scanAtendimento(row)
	if err != nil {
		return nil, fmt.Errorf("[MySQLAtendimentoRepository.buscar]: %w", err)
	}
	return atendimento, nil
}

// ExisteAtendimentoPorID verifica se um atendimento existe pelo seu ID.
func ExisteAtendimentoPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM atendimentos WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, utils.NewAppError(
			"[MySQLAtendimentoRepository.ExisteAtendimentoPorID]",
			utils.LevelError,
			"erro ao verificar existência do atendimento",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerAtendimento, err),
		)
	}
	if !exists {
		return false, nil
	}
	return true, nil
}

// scanAtendimento mapeia os dados de uma linha do resultado da consulta para um modelo de Atendimento.
func scanAtendimento(scanner interface{ Scan(dest ...any) error }) (*model.Atendimento, error) {
	var acompanhamento model.Atendimento
	err := scanner.Scan(
		&acompanhamento.ID,
		&acompanhamento.AtribuidoID,
		&acompanhamento.ChamadoID,
		&acompanhamento.CriadoEm,
		&acompanhamento.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.NewAppError(
			"[MySQLAtendimentoRepository.scanAtendimento]",
			utils.LevelError,
			"o scanner falhou ao escanear o atendimento",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerAtendimento, err),
		)
	}
	return &acompanhamento, nil
}
