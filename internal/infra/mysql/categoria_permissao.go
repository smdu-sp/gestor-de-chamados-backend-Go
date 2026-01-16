package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"

	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
)

// Erros sentinela do repositório de categoria e permissão.
var (
	ErrCategoriaPermissaoNaoEncontrada = errors.New("categoria e permissão não encontrada no banco de dados")
	ErrCategoriaPermissaoJaExiste      = errors.New("a categoria e permissão já existe no banco de dados")
)

// CategoriaPermissaoRepository é a implementação do repositório de categorias e permissões para o MySQL.
type CategoriaPermissaoRepository struct {
	db *sql.DB
}

// NovoCategoriaPermissaoRepository cria uma nova instância de CategoriaPermissaoRepository.
func NovoCategoriaPermissaoRepository(db *sql.DB) *CategoriaPermissaoRepository {
	return &CategoriaPermissaoRepository{db: db}
}

// Asserção de interface para garantir que CategoriaPermissaoRepository implementa CategoriaPermissaoRepository
var _ cpm.Repository = (*CategoriaPermissaoRepository)(nil)

// BuscarPorID recebe o ID de uma categoria e permissão e retorna a categoria e permissão correspondente.
//
// Erros sentinela possíveis: ErrCategoriaPermissaoNaoEncontrada.
func (r *CategoriaPermissaoRepository) BuscarPorID(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
	query :=
		`SELECT 
		 categoria_id, usuario_id, permissao, criado_em, atualizado_em 
		 FROM categoria_permissoes 
		 WHERE categoria_id = ? AND usuario_id = ?`

	cpm, err := r.buscar(ctx, query, ctgID, usrID)
	if err != nil {
		return nil, fmt.Errorf("buscar por ID no repo: %w", err)
	}

	if cpm == nil {
		return nil, ErrCategoriaPermissaoNaoEncontrada
	}

	return cpm, nil
}

// Criar recebe uma categoria e permissão e a salva no repositório.
//
// Erros sentinela possíveis: ErrCategoriaPermissaoJaExiste.
func (r *CategoriaPermissaoRepository) Criar(ctx context.Context, c cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
	query :=
		`INSERT INTO categoria_permissoes (
		 categoria_id, usuario_id, permissao, criado_em, atualizado_em
		 ) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		c.CategoriaID(),
		c.UsuarioID(),
		c.Permissao(),
		c.CriadoEm(),
		c.AtualizadoEm(),
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrCategoriaPermissaoJaExiste
		}
		return nil, fmt.Errorf("criar no repo: %w", err)
	}

	cpmCriada, err := r.BuscarPorID(ctx, c.CategoriaID(), c.UsuarioID())
	if err != nil {
		return nil, fmt.Errorf("buscar após criar no repo: %w", err)
	}

	return cpmCriada, nil
}

// Atualizar recebe o ID de uma categoria e usuário, e uma categoria e permissão com os dados atualizados,
// e aplica as mudanças no banco de dados.
func (r *CategoriaPermissaoRepository) Atualizar(ctx context.Context, categoriaID, usuarioID string, c cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
	query :=
		`UPDATE categoria_permissoes 
		SET permissao = ?, atualizado_em = ?
		WHERE categoria_id = ? AND usuario_id = ?`

	_, err := r.db.ExecContext(ctx, query,
		c.Permissao(),
		c.AtualizadoEm(),
		categoriaID,
		usuarioID,
	)
	if err != nil {
		return nil, fmt.Errorf("atualizar no repo: %w", err)
	}

	cpmAtualizada, err := r.BuscarPorID(ctx, categoriaID, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("buscar após atualizar no repo: %w", err)
	}

	return cpmAtualizada, nil
}

// Deletar recebe o ID de uma categoria e usuário, e remove a categoria e permissão correspondente do banco de dados.
//
// Erros sentinela possíveis: ErrCategoriaPermissaoNaoEncontrada.
func (r *CategoriaPermissaoRepository) Deletar(ctx context.Context, categoriaID, usuarioID string) error {
	existe, err := r.ExistePorID(ctx, r.db, usuarioID, categoriaID)
	if err != nil {
		return fmt.Errorf("deletar no repo: %w", err)
	}
	if !existe {
		return ErrCategoriaPermissaoNaoEncontrada
	}

	query :=
		`DELETE FROM categoria_permissoes 
		WHERE categoria_id = ? AND usuario_id = ?`

	_, err = r.db.ExecContext(ctx, query, categoriaID, usuarioID)
	if err != nil {
		return fmt.Errorf("deletar no repo: %w", err)
	}

	return nil
}

// Listar recebe um filtro e retorna a lista de categorias e permissões correspondentes,
// o total de registros e o filtro aplicado.
func (r *CategoriaPermissaoRepository) Listar(ctx context.Context, f cpm.Filtro) ([]cpm.CategoriaPermissao, int, error) {
	// Constrói a query base
	query, args := r.construirQueryListar(f)

	// Executa a query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}
	defer rows.Close()

	// Scaneia os resultados
	cpmSlice, err := scanRows(rows, r.scanCategoriaPermissao)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	// Obtém o total de registros para paginação
	total, err := obterTotalRegistros(ctx, r.db)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	return cpmSlice, total, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// buscar recebe uma query e seus argumentos, executa a busca e retorna uma categoria e permissão.
func (r *CategoriaPermissaoRepository) buscar(ctx context.Context, query string, args ...any) (*cpm.CategoriaPermissao, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	cpm, err := r.scanCategoriaPermissao(row)
	if err != nil {
		return nil, fmt.Errorf("buscar: %w", err)
	}
	return cpm, nil
}

// ExisteCategoriaPermissaoPorID recebe o ID de um usuário e categoria, e verifica se a categoria e permissão existe no banco de dados.
func (r *CategoriaPermissaoRepository) ExistePorID(ctx context.Context, db *sql.DB, usuarioID, categoriaID string) (bool, error) {
	var existe bool
	query := "SELECT EXISTS(SELECT 1 FROM categoria_permissoes WHERE usuario_id = ? AND categoria_id = ?)"
	err := db.QueryRowContext(ctx, query, usuarioID, categoriaID).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("verificar existência: %w", err)
	}
	return existe, nil
}

// scanCategoriaPermissao recebe um scanner (row ou rows) e retorna uma categoria e permissão escaneada.
func (r *CategoriaPermissaoRepository) scanCategoriaPermissao(scanner scanner) (*cpm.CategoriaPermissao, error) {
	var cpmDB cpm.CategoriaPermissaoDB
	err := scanner.Scan(
		&cpmDB.CategoriaID,
		&cpmDB.UsuarioID,
		&cpmDB.Permissao,
		&cpmDB.CriadoEm,
		&cpmDB.AtualizadoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	return cpm.CarregarDoDB(cpmDB), nil
}

// construirQueryListarCategoriasPermissoes recebe um filtro e constrói a query SQL correspondente, retornando a query e os argumentos.
func (r *CategoriaPermissaoRepository) construirQueryListar(filtro cpm.Filtro) (string, []any) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`
		SELECT SQL_CALC_FOUND_ROWS
		 categoria_id, usuario_id, permissao, criado_em, atualizado_em 
		 FROM categoria_permissoes WHERE 1=1 
	`)

	if filtro.CategoriaID() != nil {
		query.WriteString(" AND categoria_id = ?")
		args = append(args, *filtro.CategoriaID())
	}

	if filtro.UsuarioID() != nil {
		query.WriteString(" AND usuario_id = ?")
		args = append(args, *filtro.UsuarioID())
	}

	if filtro.Permissao() != nil {
		query.WriteString(" AND permissao = ?")
		args = append(args, *filtro.Permissao())
	}

	query.WriteString(" ORDER BY criado_em DESC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite(), filtro.Offset())
	return query.String(), args
}
