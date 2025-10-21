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
	ErrChamadoNaoEncontrado = errors.New("chamado não encontrado no banco de dados MySQL")
	ErrScannerChamado       = errors.New("erro ao escanear chamado do banco de dados MySQL")
)

// MySQLChamadoRepository implementa a interface ChamadoRepository para MySQL.
type MySQLChamadoRepository struct {
	db *sql.DB
}

// NewMySQLChamadoRepository cria uma nova instância de MySQLChamadoRepository.
func NewMySQLChamadoRepository(db *sql.DB) *MySQLChamadoRepository {
	return &MySQLChamadoRepository{db: db}
}

// BuscarPorID busca um chamado pelo ID.
func (r *MySQLChamadoRepository) BuscarPorID(ctx context.Context, id string) (*model.Chamado, error) {
	chamado, err := r.buscar(
		ctx,
		`SELECT id, titulo, descricao, status, arquivado, criado_em, 
		 atualizado_em, solucionado_em, solucao, fechado_em, 
		 categoria_id, subcategoria_id, criador_id
		 FROM chamados 
		 WHERE id=?`,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("[MySQLChamadoRepository.BuscarPorID]: %w", err)
	}

	if chamado == nil {
		return nil, utils.NewAppError(
			"[MySQLChamadoRepository.BuscarPorID]",
			utils.LevelInfo,
			"a busca por ID não retornou resultados",
			ErrChamadoNaoEncontrado,
		)
	}

	return chamado, nil
}

// Salvar cria um novo chamado.
func (r *MySQLChamadoRepository) Salvar(ctx context.Context, c *model.Chamado) error {
	const metodo = "[MySQLChamadoRepository.Salvar]"

	resultado, err := r.db.ExecContext(
		ctx,
		`INSERT INTO chamados (
		 id, titulo, descricao, status, arquivado, categoria_id, 
		 subcategoria_id, criador_id, criado_em, atualizado_em
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		c.ID, c.Titulo, c.Descricao, c.Status, c.Arquivado, c.CategoriaID, c.SubcategoriaID, c.CriadorID,
	)

	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao salvar o chamado no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao obter o número de linhas afetadas ao salvar o chamado",
			fmt.Errorf(utils.FmtErroWrap, ErrRowsAffected, err),
		)
	}

	if linhasAfetadas == 0 {
		return utils.NewAppError(
			metodo,
			utils.LevelWarning,
			"nenhuma linha foi afetada ao salvar o chamado no banco de dados",
			ErrExecContext,
		)
	}

	return nil
}

// Atualizar atualiza as informações de um chamado existente.
func (r *MySQLChamadoRepository) Atualizar(ctx context.Context, id string, c *model.Chamado) error {
	existe, err := ExisteChamadoPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLChamadoRepository.Atualizar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLChamadoRepository.Atualizar]",
			utils.LevelInfo,
			"não foi possível atualizar o chamado",
			ErrChamadoNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE chamados 
		 SET titulo=?, descricao=?, status=?, arquivado=?, categoria_id=?, 
		 subcategoria_id=?, atualizado_em=NOW()
		 WHERE id=?`,
		c.Titulo, c.Descricao, c.Status, c.Arquivado, c.CategoriaID, c.SubcategoriaID, id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLChamadoRepository.Atualizar]",
			utils.LevelError,
			"erro ao atualizar o chamado no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Arquivar marca um chamado como arquivado.
func (r *MySQLChamadoRepository) Arquivar(ctx context.Context, id string) error {
	existe, err := ExisteChamadoPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLChamadoRepository.Arquivar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLChamadoRepository.Arquivar]",
			utils.LevelInfo,
			"não foi possível arquivar o chamado",
			ErrChamadoNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE chamados 
		 SET arquivado=true, atualizado_em=NOW()
		 WHERE id=?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLChamadoRepository.Arquivar]",
			utils.LevelError,
			"erro ao arquivar chamado no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Desarquivar marca um chamado como não arquivado.
func (r *MySQLChamadoRepository) Desarquivar(ctx context.Context, id string) error {
	existe, err := ExisteChamadoPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLChamadoRepository.Desarquivar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLChamadoRepository.Desarquivar]",
			utils.LevelInfo,
			"não foi possível desarquivar o chamado",
			ErrChamadoNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE chamados 
		 SET arquivado=false, atualizado_em=NOW()
		 WHERE id=?`,
		id,
	)

	if err != nil {
		return utils.NewAppError(
			"[MySQLChamadoRepository.Desarquivar]",
			utils.LevelError,
			"erro ao desarquivar chamado no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// AtualizarStatus atualiza o status de um chamado, podendo incluir uma solução.
func (r *MySQLChamadoRepository) AtualizarStatus(ctx context.Context, id string, status string, solucao *string) error {
	existe, err := ExisteChamadoPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLChamadoRepository.AtualizarStatus]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLChamadoRepository.AtualizarStatus]",
			utils.LevelInfo,
			"não foi possível atualizar o status do chamado",
			ErrChamadoNaoEncontrado,
		)
	}

	var query string
	var args []any

	if solucao != nil {
		query = `UPDATE chamados 
		SET status=?, solucao=?, solucionado_em=NOW(), atualizado_em=NOW() 
		WHERE id=?`
		args = []any{status, *solucao, id}
	} else {
		query = `UPDATE chamados 
		SET status=?, atualizado_em=NOW() 
		WHERE id=?`
		args = []any{status, id}
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return utils.NewAppError(
			"[MySQLChamadoRepository.AtualizarStatus]",
			utils.LevelError,
			"erro ao atualizar status do chamado no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Listar lista chamados com paginação e filtros opcionais.
func (r *MySQLChamadoRepository) Listar(ctx context.Context, filtro model.ChamadoFiltro) ([]model.Chamado, int, error) {
	var query strings.Builder
	args := []any{}

	// TODO nao trazer os arquivados, incluir flag para exibir ou nao arquivados
	query.WriteString(
		`SELECT SQL_CALC_FOUND_ROWS
		id, titulo, descricao, status, criado_em, 
		atualizado_em, solucionado_em, solucao, fechado_em, 
		categoria_id, subcategoria_id, criador_id
		FROM chamados WHERE 1=1`,
	)

	if filtro.Busca != nil && *filtro.Busca != "" {
		padrao := "%" + *filtro.Busca + "%"
		query.WriteString(" AND (titulo LIKE ? OR descricao LIKE ?)")
		args = append(args, padrao, padrao)
	}

	if filtro.Status != nil && *filtro.Status != "" {
		query.WriteString(" AND status = ?")
		args = append(args, *filtro.Status)
	}

	if filtro.CategoriaID != nil && *filtro.CategoriaID != "" {
		query.WriteString(" AND categoria_id = ?")
		args = append(args, *filtro.CategoriaID)
	}

	if filtro.SubcategoriaID != nil && *filtro.SubcategoriaID != "" {
		query.WriteString(" AND subcategoria_id = ?")
		args = append(args, *filtro.SubcategoriaID)
	}

	if filtro.CriadorID != nil && *filtro.CriadorID != "" {
		query.WriteString(" AND criador_id = ?")
		args = append(args, *filtro.CriadorID)
	}

	query.WriteString(" ORDER BY criado_em DESC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite, (filtro.Pagina-1)*filtro.Limite)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLChamadoRepository.Listar]",
			utils.LevelError,
			"erro ao listar chamados no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	defer rows.Close()

	var chamados []model.Chamado
	for rows.Next() {
		chamado, err := scanChamado(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("[MySQLChamadoRepository.Listar]: %w", err)
		}
		chamados = append(chamados, *chamado)
	}

	var total int
	err = r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLChamadoRepository.Listar]",
			utils.LevelError,
			"erro ao contar total de chamados no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}

	return chamados, total, nil
}

// Métodos auxiliares

// buscar é um método auxiliar para buscar um chamado com base em uma consulta SQL.
func (r *MySQLChamadoRepository) buscar(ctx context.Context, query string, args ...any) (*model.Chamado, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	chamado, err := scanChamado(row)
	if err != nil {
		return nil, fmt.Errorf("[MySQLChamadoRepository.buscar]: %w", err)
	}
	return chamado, nil
}

// ExisteChamadoPorID verifica se um chamado existe pelo ID.
func ExisteChamadoPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var existe bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM chamados WHERE id=?)`, id).Scan(&existe)
	if err != nil {
		return false, utils.NewAppError(
			"[MySQLChamadoRepository.ExisteChamadoPorID]",
			utils.LevelError,
			"erro ao verificar existência do chamado no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}

	if !existe {
		return false, nil
	}

	return true, nil
}

// scanChamado mapeia os dados de um scanner (row ou rows) para uma struct Chamado.
func scanChamado(scanner interface{ Scan(dest ...any) error }) (*model.Chamado, error) {
	var chamado model.Chamado
	err := scanner.Scan(
		&chamado.ID,
		&chamado.Titulo,
		&chamado.Descricao,
		&chamado.Status,
		&chamado.CriadoEm,
		&chamado.AtualizadoEm,
		&chamado.SolucionadoEm,
		&chamado.Solucao,
		&chamado.FechadoEm,
		&chamado.CategoriaID,
		&chamado.SubcategoriaID,
		&chamado.CriadorID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.NewAppError(
			"[MySQLChamadoRepository.scanChamado]",
			utils.LevelError,
			"o scanner falhou ao escanear o chamado",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerChamado, err),
		)
	}
	return &chamado, nil
}
