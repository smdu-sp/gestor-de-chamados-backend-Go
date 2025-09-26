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
	ErrScannerSubcategoria       = errors.New("erro ao scanear subcategoria do banco de dados MySQL")
	ErrSubcategoriaNaoEncontrada = errors.New("subcategoria não encontrada no banco de dados MySQL")
	ErrSubcategoriaJaExiste      = errors.New("a subcategoria já existe no banco de dados MySQL")
)

// MySQLSubcategoriaRepository é a implementação do repositório de subcategorias para o MySQL.
type MySQLSubcategoriaRepository struct {
	db *sql.DB
}

// NewMySQLSubcategoriaRepository cria uma nova instância de MySQLSubcategoriaRepository.
func NewMySQLSubcategoriaRepository(db *sql.DB) *MySQLSubcategoriaRepository {
	return &MySQLSubcategoriaRepository{db: db}
}

// BuscarPorID busca uma subcategoria pelo seu ID.
func (r *MySQLSubcategoriaRepository) BuscarPorID(ctx context.Context, id string) (*model.Subcategoria, error) {
	subcategoria, err := r.buscar(
		ctx,
		`SELECT id, categoria_id, nome, status, criado_em, atualizado_em 
		FROM subcategorias 
		WHERE id=?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("[MySQLSubcategoriaRepository.BuscarPorID]: %w", err)
	}

	if subcategoria == nil {
		return nil, utils.NewAppError(
			"[MySQLSubcategoriaRepository.BuscarPorID]",
			utils.LevelInfo,
			"a busca por ID não retornou resultados",
			ErrSubcategoriaNaoEncontrada,
		)
	}

	return subcategoria, nil
}

// BuscarPorNome busca uma subcategoria pelo seu nome.
func (r *MySQLSubcategoriaRepository) BuscarPorNome(ctx context.Context, nome string) (*model.Subcategoria, error) {
	subcategoria, err := r.buscar(
		ctx,
		`SELECT id, categoria_id, nome, status, criado_em, atualizado_em 
		FROM subcategorias 
		WHERE nome=?`,
		nome,
	)
	if err != nil {
		return nil, fmt.Errorf("[MySQLSubcategoriaRepository.BuscarPorNome]: %w", err)
	}

	if subcategoria == nil {
		return nil, utils.NewAppError(
			"[MySQLSubcategoriaRepository.BuscarPorNome]",
			utils.LevelInfo,
			"a busca por nome não retornou resultados",
			ErrSubcategoriaNaoEncontrada,
		)
	}

	return subcategoria, nil
}

// Salvar cria uma nova subcategoria.
func (r *MySQLSubcategoriaRepository) Salvar(ctx context.Context, s *model.Subcategoria) error {
	const metodo = "[MySQLSubcategoriaRepository.Salvar]"

	resultado, err := r.db.ExecContext(
		ctx,
		`INSERT INTO subcategorias (
		id, categoria_id, nome, status, criado_em, atualizado_em
		) VALUES (?, ?, ?, ?, NOW(), NOW())`,
		s.ID,
		s.CategoriaID,
		s.Nome,
		s.Status,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			errorMensagem := err.Error()

			if strings.Contains(errorMensagem, "nome") {
				return utils.NewAppError(
					metodo,
					utils.LevelInfo,
					"erro ao salvar a subcategoria com nome duplicado no banco de dados",
					ErrSubcategoriaJaExiste,
				)
			}
		}
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao obter o número de linhas afetadas após inserir a subcategoria no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrRowsAffected, err),
		)
	}

	if linhasAfetadas == 0 {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"nenhuma linha foi afetada ao inserir a subcategoria no banco de dados",
			ErrExecContext,
		)
	}

	return nil
}

// Atualizar atualiza uma subcategoria existente.
func (r *MySQLSubcategoriaRepository) Atualizar(ctx context.Context, id string, s *model.Subcategoria) error {
	const metodo = "[MySQLSubcategoriaRepository.Atualizar]"

	existe, err := r.ExistePorID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", metodo, err)
	}
	if !existe {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"não foi possível atualizar a subcategoria",
			ErrSubcategoriaNaoEncontrada,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE subcategorias 
		SET categoria_id=?, nome=?, status=?, atualizado_em=NOW() 
		WHERE id=?`,
		s.CategoriaID,
		s.Nome,
		s.Status,
		id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			errorMensagem := err.Error()

			if strings.Contains(errorMensagem, "nome") {
				return utils.NewAppError(
					metodo,
					utils.LevelInfo,
					"erro ao atualizar a subcategoria com nome duplicado no banco de dados",
					ErrSubcategoriaJaExiste,
				)
			}
		}
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao atualizar a subcategoria no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Ativar ativa uma subcategoria.
func (r *MySQLSubcategoriaRepository) Ativar(ctx context.Context, id string) error {
	existe, err := r.ExistePorID(ctx, id)
	if err != nil {
		return fmt.Errorf("[MySQLSubcategoriaRepository.Ativar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLSubcategoriaRepository.Ativar]",
			utils.LevelInfo,
			"não foi possível ativar a subcategoria",
			ErrSubcategoriaNaoEncontrada,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE subcategorias 
		SET status=true, atualizado_em=NOW() 
		WHERE id=?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLSubcategoriaRepository.Ativar]",
			utils.LevelError,
			"falha ao ativar a subcategoria no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Desativar desativa uma subcategoria.
func (r *MySQLSubcategoriaRepository) Desativar(ctx context.Context, id string) error {
	existe, err := r.ExistePorID(ctx, id)
	if err != nil {
		return fmt.Errorf("[MySQLSubcategoriaRepository.Desativar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLSubcategoriaRepository.Desativar]",
			utils.LevelInfo,
			"não foi possível desativar a subcategoria",
			ErrSubcategoriaNaoEncontrada,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE subcategorias 
		SET status=false, atualizado_em=NOW() 
		WHERE id=?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLSubcategoriaRepository.Desativar]",
			utils.LevelError,
			"falha ao desativar a subcategoria no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Listar lista subcategorias com paginação e filtro opcionais.
func (r *MySQLSubcategoriaRepository) Listar(ctx context.Context, filtro model.SubcategoriaFiltro) ([]model.Subcategoria, int, error) {
	var query strings.Builder
	args := []any{}

	// TODO nao trazer os arquivados, incluir flag para exibir ou nao status false
	query.WriteString(`SELECT SQL_CALC_FOUND_ROWS 
		id, categoria_id, nome, status, criado_em, atualizado_em 
		FROM subcategorias 
		WHERE 1=1`,
	)

	if filtro.Busca != nil && *filtro.Busca != "" {
		padrao := "%" + *filtro.Busca + "%"
		query.WriteString(" AND (nome LIKE ?)")
		args = append(args, padrao)
	}

	if filtro.Status != nil {
		query.WriteString(" AND status=?")
		args = append(args, *filtro.Status)
	}

	query.WriteString(" ORDER BY nome ASC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite, (filtro.Pagina-1)*filtro.Limite)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLSubcategoriaRepository.Listar]",
			utils.LevelError,
			"falha ao listar as subcategorias no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	defer rows.Close()

	var subcategorias []model.Subcategoria
	for rows.Next() {
		subcategoria, err := scanSubcategoria(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("[MySQLSubcategoriaRepository.Listar]: %w", err)
		}
		subcategorias = append(subcategorias, *subcategoria)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLSubcategoriaRepository.Listar]",
			utils.LevelError,
			"erro ao obter total de subcategorias",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}

	return subcategorias, total, nil
}

// Metodos auxiliares

// buscar executa uma consulta que retorna uma única subcategoria.
func (r *MySQLSubcategoriaRepository) buscar(ctx context.Context, query string, args ...interface{}) (*model.Subcategoria, error) {
	row := r.db.QueryRowContext(ctx, query, args...)

	subcategoria, err := scanSubcategoria(row)
	if err != nil {
		return nil, fmt.Errorf("[MySQLSubcategoriaRepository.buscar]: %w", err)
	}

	return subcategoria, nil
}

// ExistePorID verifica se uma subcategoria existe pelo seu ID.
func (r *MySQLSubcategoriaRepository) ExistePorID(ctx context.Context, id string) (bool, error) {
	var existe bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM subcategorias WHERE id=?)`, id).Scan(&existe)
	if err != nil {
		return false, utils.NewAppError(
			"[MySQLSubcategoriaRepository.ExistePorID]",
			utils.LevelError,
			"erro ao verificar a existência da subcategoria pelo ID no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	return existe, nil
}

// scanSubcategoria mapeia os dados de um scanner (row ou rows) para uma struct Subcategoria.
func scanSubcategoria(scanner interface{ Scan(dest ...any) error }) (*model.Subcategoria, error) {
	var subcategoria model.Subcategoria
	err := scanner.Scan(
		&subcategoria.ID,
		&subcategoria.CategoriaID,
		&subcategoria.Nome,
		&subcategoria.Status,
		&subcategoria.CriadoEm,
		&subcategoria.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.NewAppError(
			"[MySQLSubcategoriaRepository.scanSubcategoria]",
			utils.LevelError,
			"O scanner falhou ao scanear a subcategoria",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerSubcategoria, err),
		)
	}
	return &subcategoria, nil
}
