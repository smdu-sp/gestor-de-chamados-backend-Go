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
	ErrScannerCategoria       = errors.New("erro ao scanear categoria do banco de dados MySQL")
	ErrCategoriaNaoEncontrada = errors.New("categoria não encontrada no banco de dados MySQL")
	ErrCategoriaJaExiste      = errors.New("a categoria já existe no banco de dados MySQL")
)

// MySQLCategoriaRepository é a implementação do repositório de categorias para o MySQL.
type MySQLCategoriaRepository struct {
	db *sql.DB
}

// NewMySQLCategoriaRepository cria uma nova instância de MySQLCategoriaRepository.
func NewMySQLCategoriaRepository(db *sql.DB) *MySQLCategoriaRepository {
	return &MySQLCategoriaRepository{db: db}
}

// BuscarPorID busca uma categoria pelo seu ID.
func (r *MySQLCategoriaRepository) BuscarPorID(ctx context.Context, id string) (*model.Categoria, error) {
	categoria, err := r.buscar(
		ctx,
		`SELECT id, nome, status, criado_em, atualizado_em 
		FROM categorias 
		WHERE id=?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("[MySQLCategoriaRepository.BuscarPorID]: %w", err)
	}

	if categoria == nil {
		return nil, utils.NewAppError(
			"[MySQLCategoriaRepository.BuscarPorID]",
			utils.LevelInfo,
			"a busca por ID não retornou resultados",
			ErrCategoriaNaoEncontrada,
		)
	}

	return categoria, nil
}

// BuscarPorNome busca uma categoria pelo seu nome.
func (r *MySQLCategoriaRepository) BuscarPorNome(ctx context.Context, nome string) (*model.Categoria, error) {
	categoria, err := r.buscar(
		ctx,
		`SELECT id, nome, status, criado_em, atualizado_em 
		FROM categorias 
		WHERE nome=?`,
		nome,
	)

	if err != nil {
		return nil, fmt.Errorf("[MySQLCategoriaRepository.BuscarPorNome]: %w", err)
	}

	if categoria == nil {
		return nil, utils.NewAppError(
			"[MySQLCategoriaRepository.BuscarPorNome]",
			utils.LevelInfo,
			"a busca por nome não retornou resultados",
			ErrCategoriaNaoEncontrada,
		)
	}

	return categoria, nil
}

// Salvar cria uma nova categoria.
func (r *MySQLCategoriaRepository) Salvar(ctx context.Context, c *model.Categoria) error {
	const metodo = "[MySQLCategoriaRepository.Salvar]"

	resultado, err := r.db.ExecContext(
		ctx,
		`INSERT INTO categorias (
		id, nome, status, criado_em, atualizado_em
		) VALUES (?, ?, ?, NOW(), NOW())`,
		c.ID, c.Nome, c.Status,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			errorMensagem := err.Error()

			if strings.Contains(errorMensagem, "nome") {
				return utils.NewAppError(
					metodo,
					utils.LevelInfo,
					"erro ao salvar a categoria com nome duplicado no banco de dados",
					fmt.Errorf(utils.FmtErroWrap, ErrCategoriaJaExiste, err),
				)
			}
		}
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"falha ao obter o número de linhas afetadas ao salvar categoria no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrRowsAffected, err),
		)
	}

	if linhasAfetadas == 0 {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"nenhuma linha foi afetada ao salvar a categoria no banco de dados",
			ErrExecContext,
		)
	}

	return nil
}

// Atualizar atualiza as informações de uma categoria existente.
func (r *MySQLCategoriaRepository) Atualizar(ctx context.Context, id string, c *model.Categoria) error {
	const metodo = "[MySQLCategoriaRepository.Atualizar]"

	existe, err := ExisteCategoriaPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("%s: %w", metodo, err)
	}
	if !existe {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"não foi possível atualizar a categoria",
			ErrCategoriaNaoEncontrada,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE categorias 
		SET nome=?, status=?, atualizado_em=NOW() 
		WHERE id=?`,
		c.Nome, c.Status, id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			errorMensagem := err.Error()

			if strings.Contains(errorMensagem, "nome") {
				return utils.NewAppError(
					metodo,
					utils.LevelWarning,
					"erro ao atualizar a categoria com nome duplicado no banco de dados",
					fmt.Errorf(utils.FmtErroWrap, ErrCategoriaJaExiste, err),
				)
			}
		}
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"falha ao atualizar a categoria no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Ativar ativa uma categoria pelo seu ID.
func (r *MySQLCategoriaRepository) Ativar(ctx context.Context, id string) error {
	existe, err := ExisteCategoriaPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLCategoriaRepository.Ativar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLCategoriaRepository.Ativar]",
			utils.LevelInfo,
			"não foi possível ativar a categoria",
			ErrCategoriaNaoEncontrada,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE categorias 
		SET status=true, atualizado_em=NOW() 
		WHERE id=?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLCategoriaRepository.Ativar]",
			utils.LevelError,
			"falha ao ativar a categoria no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Desativar desativa (soft delete) uma categoria pelo seu ID.
func (r *MySQLCategoriaRepository) Desativar(ctx context.Context, id string) error {
	existe, err := ExisteCategoriaPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLCategoriaRepository.Desativar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLCategoriaRepository.Desativar]",
			utils.LevelInfo,
			"não foi possível desativar a categoria",
			ErrCategoriaNaoEncontrada,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE categorias 
		SET status=false, atualizado_em=NOW() 
		WHERE id=?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLCategoriaRepository.Desativar]",
			utils.LevelError,
			"falha ao desativar a categoria no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Listar lista categorias com paginação e filtros opcionais.
func (r *MySQLCategoriaRepository) Listar(ctx context.Context, filtro model.CategoriaFiltro) ([]model.Categoria, int, error) {
	var query strings.Builder
	args := []any{}

	// TODO nao trazer os arquivados, incluir flag para exibir ou nao status false
	query.WriteString(
		`SELECT SQL_CALC_FOUND_ROWS
		id, nome, status, criado_em, atualizado_em 
		FROM categorias 
		WHERE 1=1`,
	)

	if filtro.Busca != nil && *filtro.Busca != "" {
		padrao := "%" + *filtro.Busca + "%"
		query.WriteString(" AND (nome LIKE ?)")
		args = append(args, padrao)
	}

	if filtro.Status != nil {
		query.WriteString(" AND status = ?")
		args = append(args, *filtro.Status)
	}

	query.WriteString(" ORDER BY nome ASC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite, (filtro.Pagina-1)*filtro.Limite)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLCategoriaRepository.Listar]",
			utils.LevelError,
			"falha ao listar categorias no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	defer rows.Close()

	var categorias []model.Categoria
	for rows.Next() {
		categoria, err := scanCategoria(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("[MySQLCategoriaRepository.Listar]: %w", err)
		}
		categorias = append(categorias, *categoria)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLCategoriaRepository.Listar]",
			utils.LevelError,
			"erro ao obter total de categorias",
			fmt.Errorf(utils.FmtErroWrap, ErrScan, err),
		)
	}

	return categorias, total, nil
}

// Metodos auxiliares

// buscar executa uma consulta que retorna uma única categoria.
func (r *MySQLCategoriaRepository) buscar(ctx context.Context, query string, args ...any) (*model.Categoria, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	categoria, err := scanCategoria(row)
	if err != nil {
		return nil, fmt.Errorf("[MySQLCategoriaRepository.buscar]: %w", err)
	}
	return categoria, nil
}

// ExisteCategoriaPorID verifica se uma categoria existe pelo seu ID.
func ExisteCategoriaPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categorias WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, utils.NewAppError(
			"[MySQLCategoriaRepository.ExisteCategoriaPorID]",
			utils.LevelError,
			"Falha ao verificar existência da categoria por ID",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	return exists, nil
}

// scanCategoria mapeia os dados de um scanner (row ou rows) para uma struct Categoria.
func scanCategoria(scanner interface{ Scan(dest ...any) error }) (*model.Categoria, error) {
	var categoria model.Categoria
	err := scanner.Scan(
		&categoria.ID,
		&categoria.Nome,
		&categoria.Status,
		&categoria.CriadoEm,
		&categoria.AtualizadoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, utils.NewAppError(
			"[MySQLCategoriaRepository.scanCategoria]",
			utils.LevelError,
			"O scanner falhou ao scanear a categoria",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerCategoria, err),
		)
	}
	return &categoria, nil
}
