package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
)

// Erros sentinela do repositório de categorias.
var (
	ErrCategoriaNaoEncontrada = errors.New("categoria não encontrada no banco de dados")
	ErrCategoriaJaExiste      = errors.New("a categoria já existe no banco de dados")
)

// CategoriaRepository é a implementação do repositório de categorias para o MySQL.
type CategoriaRepository struct {
	db *sql.DB
}

// NovoCategoriaRepository cria uma nova instância de CategoriaRepository.
func NovoCategoriaRepository(db *sql.DB) *CategoriaRepository {
	return &CategoriaRepository{db: db}
}

// Asserção de interface para garantir que CategoriaRepository implementa CategoriaService
var _ ctg.Repository = (*CategoriaRepository)(nil)

// BuscarPorID recebe o ID de uma categoria e retorna a categoria correspondente.
//
// Erros sentinela possíveis: ErrCategoriaNaoEncontrada.
func (r *CategoriaRepository) BuscarPorID(ctx context.Context, id string) (*ctg.Categoria, error) {
	query :=
		`SELECT id, nome, status, criado_em, atualizado_em 
	   FROM categorias 
	   WHERE id=?`

	ctg, err := r.buscar(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("buscar por ID no repo: %w", err)
	}

	if ctg == nil {
		return nil, ErrCategoriaNaoEncontrada
	}

	return ctg, nil
}

// BuscarPorNome recebe o nome de uma categoria e retorna a categoria correspondente.
//
// Erros sentinela possíveis: ErrCategoriaNaoEncontrada.
func (r *CategoriaRepository) BuscarPorNome(ctx context.Context, nome string) (*ctg.Categoria, error) {
	query :=
		`SELECT id, nome, status, criado_em, atualizado_em 
	   FROM categorias 
	   WHERE nome=?`
	ctg, err := r.buscar(ctx, query, nome)

	if err != nil {
		return nil, fmt.Errorf("buscar por nome no repo: %w", err)
	}

	if ctg == nil {
		return nil, ErrCategoriaNaoEncontrada
	}

	return ctg, nil
}

// Criar recebe uma categoria e a persiste no banco de dados.
//
// Erros sentinela possíveis: ErrCategoriaJaExiste.
func (r *CategoriaRepository) Criar(ctx context.Context, c ctg.Categoria) (*ctg.Categoria, error) {
	query := `
		INSERT INTO categorias (
			id, nome, status, criado_em, atualizado_em
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		c.ID(),
		c.Nome(),
		c.Status(),
		c.CriadoEm(),
		c.AtualizadoEm(),
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrCategoriaJaExiste
		}

		return nil, fmt.Errorf("criar no repo: %w", err)
	}

	ctgSalva, err := r.BuscarPorID(ctx, c.ID())
	if err != nil {
		return nil, fmt.Errorf("buscar após criar no repo: %w", err)
	}

	return ctgSalva, nil
}

// Atualizar recebe o ID de uma categoria e uma categoria com os dados atualizados,
// e aplica as mudanças no banco de dados.
//
// Erros sentinela possíveis: ErrCategoriaJaExiste.
func (r *CategoriaRepository) Atualizar(ctx context.Context, id string, c ctg.Categoria) (*ctg.Categoria, error) {
	query := `
		UPDATE categorias 
		SET nome=?, status=?, atualizado_em=? 
		WHERE id=?
	`
	_, err := r.db.ExecContext(ctx, query,
		c.Nome(),
		c.Status(),
		c.AtualizadoEm(),
		id,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrCategoriaJaExiste
		}
		return nil, fmt.Errorf("atualizar no repo: %w", err)
	}

	ctgAtualizada, err := r.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar após atualizar no repo: %w", err)
	}

	return ctgAtualizada, nil
}

// Listar recebe um filtro e retorna a lista de categorias correspondentes,
// o total de registros e o filtro aplicado.
func (r *CategoriaRepository) Listar(ctx context.Context, filtro ctg.Filtro) ([]ctg.Categoria, int, error) {
	// Constrói a query base
	query, args := r.construirQueryListar(filtro)

	// Executa a query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}
	defer rows.Close()

	// Scaneia os resultados
	ctgSlice, err := scanRows(rows, scanCategoria)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	// Obtém o total de registros para paginação
	total, err := obterTotalRegistros(ctx, r.db)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	return ctgSlice, total, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// buscar recebe uma query e argumentos, executa a busca e retorna uma categoria.
func (r *CategoriaRepository) buscar(ctx context.Context, query string, args ...any) (*ctg.Categoria, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	ctg, err := scanCategoria(row)
	if err != nil {
		return nil, fmt.Errorf("buscar: %w", err)
	}
	return ctg, nil
}

// ExisteCategoriaPorID recebe o ID de uma categoria e verifica se ela existe no banco de dados.
func ExisteCategoriaPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var existe bool
	query := "SELECT EXISTS(SELECT 1 FROM categorias WHERE id = ?)"
	err := db.QueryRowContext(ctx, query, id).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("verificar existência: %w", err)
	}
	return existe, nil
}

// scanCategoria recebe um scanner (row ou rows) e retorna uma categoria escaneada.
func scanCategoria(scanner scanner) (*ctg.Categoria, error) {
	var ctgDB ctg.CategoriaDB
	err := scanner.Scan(
		&ctgDB.ID,
		&ctgDB.Nome,
		&ctgDB.Status,
		&ctgDB.CriadoEm,
		&ctgDB.AtualizadoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	
	return ctg.CarregarDoBD(ctgDB), nil
}

// construirQueryListar recebe um filtro e constrói a query SQL para listar categorias.
func (r *CategoriaRepository) construirQueryListar(filtro ctg.Filtro) (string, []any) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`
		SELECT SQL_CALC_FOUND_ROWS
			id, nome, status, criado_em, atualizado_em
		FROM categorias
		WHERE 1=1
	`)

	if filtro.Busca() != nil {
		query.WriteString(" AND nome LIKE ?")
		args = append(args, "%"+*filtro.Busca()+"%")
	}

	if filtro.Status() != nil {
		query.WriteString(" AND status = ?")
		args = append(args, *filtro.Status())
	}

	query.WriteString(" ORDER BY criado_em DESC")
	args = append(args, filtro.Limite(), filtro.Offset())
	return query.String(), args
}
