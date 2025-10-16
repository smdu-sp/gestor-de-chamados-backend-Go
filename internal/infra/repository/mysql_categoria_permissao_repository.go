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
	ErrScannerCategoriaPermissao       = errors.New("erro ao scanear categoria e permissão do banco de dados MySQL")
	ErrCategoriaPermissaoNaoEncontrada = errors.New("categoria e permissão não encontrada no banco de dados MySQL")
	ErrCategoriaPermissaoJaExiste      = errors.New("a categoria e permissão já existe no banco de dados MySQL")
)

// MySQLCategoriaPermissaoRepository é a implementação do repositório de categorias e permissões para o MySQL.
type MySQLCategoriaPermissaoRepository struct {
	db *sql.DB
}

// NewMySQLCategoriaPermissaoRepository cria uma nova instância de MySQLCategoriaPermissaoRepository.
func NewMySQLCategoriaPermissaoRepository(db *sql.DB) *MySQLCategoriaPermissaoRepository {
	return &MySQLCategoriaPermissaoRepository{db: db}
}

// Salvar insere uma nova categoria e permissão no repositório.
func (r *MySQLCategoriaPermissaoRepository) Salvar(ctx context.Context, c *model.CategoriaPermissao) error {
	const metodo = "[MySQLCategoriaPermissaoRepository.Salvar]"

	resultado, err := r.db.ExecContext(
		ctx,
		`INSERT INTO categoria_permissoes (
		categoria_id, usuario_id, permissao, criado_em, atualizado_em
		) VALUES (?, ?, ?, NOW(), NOW())`,
		c.CategoriaID, c.UsuarioID, c.Permissao,
	)
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro inesperado ao salvar categoriaPermissao no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"falha ao obter o número de linhas afetadas ao salvar categoriaPermissao no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrRowsAffected, err),
		)
	}

	if linhasAfetadas == 0 {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"nenhuma linha foi inserida ao tentar salvar a categoriaPermissao",
			ErrCategoriaPermissaoNaoEncontrada,
		)
	}

	return nil
}

// Atualizar modifica os dados de uma categoria e permissão existente.
func (r *MySQLCategoriaPermissaoRepository) Atualizar(ctx context.Context, categoriaID, usuarioID string, c *model.CategoriaPermissao) error {
	const metodo = "[MySQLCategoriaPermissaoRepository.Atualizar]"

	existe, err := ExisteCategoriaPermissaoPorID(ctx, r.db, usuarioID, categoriaID)
	if err != nil {
		return fmt.Errorf("%s: %w", metodo, err)
	}
	if !existe {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"não foi possível atualizar a categoriaPermissao",
			ErrCategoriaPermissaoNaoEncontrada,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE categoria_permissoes
		SET permissao = ?, atualizado_em = NOW() 
		WHERE categoria_id = ? AND usuario_id = ?`,
		c.Permissao, categoriaID, usuarioID,
	)
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"falha ao atualizar a categoriaPermissao no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Deletar remove uma categoria e permissão do repositório.
func (r *MySQLCategoriaPermissaoRepository) Deletar(ctx context.Context, categoriaID, usuarioID string) error {
	existe, err := ExisteCategoriaPermissaoPorID(ctx, r.db, usuarioID, categoriaID)
	if err != nil {
		return fmt.Errorf("[MySQLCategoriaPermissaoRepository.Deletar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLCategoriaPermissaoRepository.Deletar]",
			utils.LevelInfo,
			"não foi possível deletar a categoriaPermissao",
			ErrCategoriaPermissaoNaoEncontrada,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`DELETE FROM categoria_permissoes 
		WHERE categoria_id = ? AND usuario_id = ?`,
		categoriaID, usuarioID,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLCategoriaPermissaoRepository.Deletar]",
			utils.LevelError,
			"falha ao deletar a categoriaPermissao no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Listar retorna uma lista de categorias e permissões com base em filtros e paginação.
func (r *MySQLCategoriaPermissaoRepository) Listar(ctx context.Context, filtro model.CategoriaPermissaoFiltro) ([]model.CategoriaPermissao, int, error) {
	var query strings.Builder
	args := []any{}

	query.WriteString(
		`SELECT SQL_CALC_FOUND_ROWS
		categoria_id, usuario_id, permissao, criado_em, atualizado_em 
		FROM categoria_permissoes 
		WHERE 1=1`,
	)

	if filtro.CategoriaID != nil && *filtro.CategoriaID != "" {
		query.WriteString(" AND categoria_id = ?")
		args = append(args, *filtro.CategoriaID)
	}

	if filtro.UsuarioID != nil && *filtro.UsuarioID != "" {
		query.WriteString(" AND usuario_id = ?")
		args = append(args, *filtro.UsuarioID)
	}

	if filtro.Permissao != nil && *filtro.Permissao != "" {
		query.WriteString(" AND permissao = ?")
		args = append(args, *filtro.Permissao)
	}

	query.WriteString(" ORDER BY nome ASC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite, (filtro.Pagina-1)*filtro.Limite)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLCategoriaPermissaoRepository.Listar]",
			utils.LevelError,
			"falha ao listar categorias e permissões no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}

	defer rows.Close()

	var categoriasPermissoes []model.CategoriaPermissao
	for rows.Next() {
		categoriaPermissao, err := scanCategoriaPermissao(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("[MySQLCategoriaPermissaoRepository.Listar]: %w", err)
		}
		categoriasPermissoes = append(categoriasPermissoes, *categoriaPermissao)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLCategoriaPermissaoRepository.Listar]",
			utils.LevelError,
			"falha ao obter o número total de categorias e permissões",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}

	return categoriasPermissoes, total, nil
}

// metodos auxiliares

// ExisteCategoriaPermissaoPorID verifica se uma categoria e permissão existe pelo seu ID de usuário e ID de categoria.
func ExisteCategoriaPermissaoPorID(ctx context.Context, db *sql.DB, usuarioID, categoriaID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categoria_permissoes WHERE usuario_id = ? AND categoria_id = ?)", usuarioID, categoriaID).Scan(&exists)
	if err != nil {
		return false, utils.NewAppError(
			"[MySQLCategoriaPermissaoRepository.ExisteCategoriaPermissaoPorID]",
			utils.LevelError,
			"Falha ao verificar existência da categoriaPermissao por ID",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	return exists, nil
}

// scanCategoriaPermissao mapeia os dados de um scanner (row ou rows) para uma struct CategoriaPermissao.
func scanCategoriaPermissao(scanner interface{ Scan(dest ...any) error }) (*model.CategoriaPermissao, error) {
	var categoriaPermissao model.CategoriaPermissao
	err := scanner.Scan(
		&categoriaPermissao.CategoriaID,
		&categoriaPermissao.UsuarioID,
		&categoriaPermissao.Permissao,
		&categoriaPermissao.CriadoEm,
		&categoriaPermissao.AtualizadoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, utils.NewAppError(
			"[MySQLCategoriaPermissaRepository.scanCategoriaPermissao]",
			utils.LevelError,
			"O scanner falhou ao scanear a categoriaPermissao",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerCategoriaPermissao, err),
		)
	}
	return &categoriaPermissao, nil
}
