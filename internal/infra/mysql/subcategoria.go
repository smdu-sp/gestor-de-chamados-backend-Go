package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	subc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
)

var (
	ErrSubcategoriaNaoEncontrada   = errors.New("subcategoria não encontrada no banco de dados")
	ErrSubcategoriaJaExisteComNome = errors.New("a subcategoria já existe com este nome no banco de dados")
)

// SubcategoriaRepository é a implementação do repositório de subcategorias para o MySQL.
type SubcategoriaRepository struct {
	db *sql.DB
}

// NovoSubcategoriaRepository cria uma nova instância de SubcategoriaRepository.
func NovoSubcategoriaRepository(db *sql.DB) *SubcategoriaRepository {
	return &SubcategoriaRepository{db: db}
}

// Asserção de interface para garantir que SubcategoriaRepository implementa SubcategoriaRepository
var _ subc.Repository = (*SubcategoriaRepository)(nil)

// BuscarPorID recebe um ID e retorna a subcategoria correspondente.
//
// Erros sentinela possíveis: ErrSubcategoriaNaoEncontrada.
func (r *SubcategoriaRepository) BuscarPorID(ctx context.Context, id string) (*subc.Subcategoria, error) {
	query :=
		`SELECT id, categoria_id, nome, status, criado_em, atualizado_em 
		 FROM subcategorias 
		 WHERE id=?`

	subc, err := r.buscar(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("buscar por ID no repo: %w", err)
	}

	if subc == nil {
		return nil, ErrSubcategoriaNaoEncontrada
	}

	return subc, nil
}

// BuscarPorNome recebe um nome e retorna a subcategoria correspondente.
//
// Erros sentinela possíveis: ErrSubcategoriaNaoEncontrada.
func (r *SubcategoriaRepository) BuscarPorNome(ctx context.Context, nome string) (*subc.Subcategoria, error) {
	query :=
		`SELECT id, categoria_id, nome, status, criado_em, atualizado_em 
		 FROM subcategorias 
		 WHERE nome=?`

	subc, err := r.buscar(ctx, query, nome)
	if err != nil {
		return nil, fmt.Errorf("buscar por nome no repo: %w", err)
	}

	if subc == nil {
		return nil, ErrSubcategoriaNaoEncontrada
	}

	return subc, nil
}

// Criar recebe uma subcategoria e a insere no banco de dados.
//
// Erros sentinela possíveis: ErrSubcategoriaJaExisteComNome.
func (r *SubcategoriaRepository) Criar(ctx context.Context, s subc.Subcategoria) (*subc.Subcategoria, error) {
	query := 
		`INSERT INTO subcategorias (
		 id, categoria_id, nome, status, 
		 criado_em, atualizado_em
		 ) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx,	query, 
		s.ID(), 
		s.CategoriaID(), 
		s.Nome(), 
		s.Status(), 
		s.CriadoEm(), 
		s.AtualizadoEm(),
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return nil, ErrSubcategoriaJaExisteComNome
			}

		return nil, fmt.Errorf("criar no repo: %w", err)
	}

	subcCriada, err := r.BuscarPorID(ctx, s.ID())
	if err != nil {
		return nil, fmt.Errorf("buscar após criar no repo: %w", err)
	}

	return subcCriada, nil
}

// Atualizar recebe um ID e uma subcategoria para atualizar os dados da subcategoria correspondente.
//
// Erros sentinela possíveis: ErrSubcategoriaJaExisteComNome.
func (r *SubcategoriaRepository) Atualizar(ctx context.Context, id string, s subc.Subcategoria) (*subc.Subcategoria, error) {
	query :=
		`UPDATE subcategorias 
		 SET categoria_id=?, nome=?, status=?, atualizado_em=?
		 WHERE id=?`

	_, err := r.db.ExecContext(ctx, query, 
		s.CategoriaID(), 
		s.Nome(), 
		s.Status(), 
		s.AtualizadoEm(), 
		id,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrSubcategoriaJaExisteComNome
		}
		return nil, fmt.Errorf("atualizar no repo: %w", err)
	}

	subcAtualizada, err := r.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar após atualizar no repo: %w", err)
	}

	return subcAtualizada, nil
}

// Listar recebe um filtro e retorna uma lista de subcategorias que correspondem aos critérios do filtro,
// juntamente com o total de registros encontrados.
func (r *SubcategoriaRepository) Listar(ctx context.Context, filtro subc.Filtro) ([]subc.Subcategoria, int, error) {
	query, args := r.construirQueryListar(filtro)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}
	defer rows.Close()

	subcSclice, err := scanRows(rows, r.scanSubcategoria)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	total, err := obterTotalRegistros(ctx, r.db)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	return subcSclice, total, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// buscar recebe uma query e argumentos, executa a busca e retorna a subcategoria correspondente.
func (r *SubcategoriaRepository) buscar(ctx context.Context, query string, args ...any) (*subc.Subcategoria, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	subc, err := r.scanSubcategoria(row)
	if err != nil {
		return nil, fmt.Errorf("buscar: %w", err)
	}

	return subc, nil
}

// ExisteSubcategoriaPorID recebe um ID e verifica se a subcategoria correspondente existe no banco de dados.
func ExisteSubcategoriaPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var existe bool
	query := `SELECT EXISTS(SELECT 1 FROM subcategorias WHERE id=?)`
	err := db.QueryRowContext(ctx, query, id).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("verificar existência: %w", err)
	}

	return existe, nil
}

// scanSubcategoria recebe um scanner (row ou rows) e retorna a subcategoria correspondente.
func (r *SubcategoriaRepository) scanSubcategoria(scanner scanner) (*subc.Subcategoria, error) {
	var subcDB subc.SubcategoriaDB
	err := scanner.Scan(
		&subcDB.ID,
		&subcDB.CategoriaID,
		&subcDB.Nome,
		&subcDB.Status,
		&subcDB.CriadoEm,
		&subcDB.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	return subc.CarregarDoBD(subcDB), nil
}

// construirQueryListar recebe um filtro e constrói a query SQL correspondente, retornando a query e os argumentos.
func (r *SubcategoriaRepository) construirQueryListar(filtro subc.Filtro) (string, []any) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`SELECT SQL_CALC_FOUND_ROWS 
		id, categoria_id, nome, status, criado_em, atualizado_em 
		FROM subcategorias 
		WHERE 1=1`,
	)

	if filtro.Busca() != nil && *filtro.Busca() != "" {
		padrao := "%" + *filtro.Busca() + "%"
		query.WriteString(" AND (nome LIKE ?)")
		args = append(args, padrao)
	}

	if filtro.Status() != nil {
		query.WriteString(" AND status=?")
		args = append(args, *filtro.Status())
	}

	query.WriteString(" ORDER BY nome ASC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite(), filtro.Offset())
	return query.String(), args
}