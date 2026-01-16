package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// Erros sentinela do repositório de usuários.
var (
	ErrUsuarioNaoEncontrado    = errors.New("usuário não encontrado no banco de dados")
	ErrUsuarioJaExisteComEmail = errors.New("erro ao tentar duplicar usuário com mesmo email no banco de dados")
	ErrUsuarioJaExisteComLogin = errors.New("erro ao tentar duplicar usuário com mesmo login no banco de dados")
)

// UsuarioRepository implementa a interface UsuarioRepository para MySQL.
type UsuarioRepository struct {
	db *sql.DB
}

// NovoUsuarioRepository cria uma nova instância de UsuarioRepository.
func NovoUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

// Asserção de interface para garantir que UsuarioRepository implementa UsuarioRepository
var _ usr.Repository = (*UsuarioRepository)(nil)

// BuscarPorID recebe um ID e retorna o usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado.
func (r *UsuarioRepository) BuscarPorID(ctx context.Context, id string) (*usr.Usuario, error) {
	query :=
		`SELECT id, nome, login, email, permissao, status, 
		 avatar, ultimo_login, criado_em, atualizado_em
     FROM usuarios 
		 WHERE id=?`

	usr, err := r.buscar(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("buscar por ID no repo: %w", err)
	}

	if usr == nil {
		return nil, ErrUsuarioNaoEncontrado
	}

	return usr, nil
}

// BuscarPorLogin recebe um login e retorna o usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado.
func (r *UsuarioRepository) BuscarPorLogin(ctx context.Context, login string) (*usr.Usuario, error) {
	query :=
		`SELECT id, nome, login, email, permissao, status,
		 avatar, ultimo_login, criado_em, atualizado_em
     FROM usuarios 
		 WHERE login=?`

	usr, err := r.buscar(ctx, query, login)
	if err != nil {
		return nil, fmt.Errorf("buscar por login no repo: %w", err)
	}

	if usr == nil {
		return nil, ErrUsuarioNaoEncontrado
	}

	return usr, nil
}

// Criar recebe um usuário e o insere no banco de dados.
//
// Erros sentinela possíveis: ErrUsuarioJaExisteComLogin, ErrUsuarioJaExisteComEmail.
func (r *UsuarioRepository) Criar(ctx context.Context, u usr.Usuario) (*usr.Usuario, error) {
	query :=
		`INSERT INTO usuarios(
     id, nome, login, email, permissao, status, 
	   avatar, ultimo_login, criado_em, atualizado_em
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		u.ID(),
		u.Nome(),
		u.Login(),
		u.Email(),
		u.Permissao(),
		u.Status(),
		u.Avatar(),
		u.UltimoLogin(),
		u.CriadoEm(),
		u.AtualizadoEm(),
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			msg := mysqlErr.Message

			switch {
			case strings.Contains(msg, "login"):
				return nil, ErrUsuarioJaExisteComLogin
			case strings.Contains(msg, "email"):
				return nil, ErrUsuarioJaExisteComEmail
			}
		}
		return nil, fmt.Errorf("criar no repo: %w", err)
	}

	usrCriado, err := r.BuscarPorID(ctx, u.ID())
	if err != nil {
		return nil, fmt.Errorf("buscar após criar no repo: %w", err)
	}

	return usrCriado, nil
}

// Atualizar recebe um ID e um usuário, e atualiza os dados do usuário no banco de dados.
//
// Erros sentinela possíveis: ErrUsuarioJaExisteComLogin, ErrUsuarioJaExisteComEmail.
func (r *UsuarioRepository) Atualizar(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
	query :=
		`UPDATE usuarios
     SET nome=?, login=?, email=?, permissao=?, status=?, avatar=?, ultimo_login=?, atualizado_em=?
     WHERE id=?`

	_, err := r.db.ExecContext(ctx, query,
		u.Nome(),
		u.Login(),
		u.Email(),
		u.Permissao(),
		u.Status(),
		u.Avatar(),
		u.UltimoLogin(),
		u.AtualizadoEm(),
		id,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			msg := mysqlErr.Message

			switch {
			case strings.Contains(msg, "login"):
				return nil, ErrUsuarioJaExisteComLogin
			case strings.Contains(msg, "email"):
				return nil, ErrUsuarioJaExisteComEmail
			}
		}
		return nil, fmt.Errorf("atualizar no repo: %w", err)
	}

	usrAtualizado, err := r.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar após atualizar no repo: %w", err)
	}

	return usrAtualizado, nil
}

// Listar recebe um filtro e retorna uma lista de usuários que correspondem aos critérios do filtro,
// juntamente com o total de registros encontrados.
func (r *UsuarioRepository) Listar(ctx context.Context, filtro usr.Filtro) ([]usr.Usuario, int, error) {
	// Constrói a query base
	query, args := r.construirQueryListar(filtro)

	// Executa a query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}
	defer rows.Close()

	// Scaneia os resultados
	usrSlice, err := scanRows(rows, r.scanUsuario)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	// Obtém o total de registros
	total, err := obterTotalRegistros(ctx, r.db)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	return usrSlice, total, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// buscar recebe uma query e argumentos, executa a consulta e retorna o usuário correspondente.
func (r *UsuarioRepository) buscar(ctx context.Context, query string, args ...any) (*usr.Usuario, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	usr, err := r.scanUsuario(row)
	if err != nil {
		return nil, fmt.Errorf("buscar: %w", err)
	}

	return usr, nil
}

// ExisteUsuarioPorID recebe um ID e verifica se um usuário com esse ID existe no banco de dados.
func ExisteUsuarioPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var existe bool
	query := `SELECT EXISTS(SELECT 1 FROM usuarios WHERE id=?)`
	err := db.QueryRowContext(ctx, query, id).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("verificar existência: %w", err)
	}

	return existe, nil
}

// ExistePorLogin recebe um login e verifica se um usuário com esse login existe no banco de dados.
func (r *UsuarioRepository) ExistePorLogin(ctx context.Context, login string) (bool, error) {
	var existe bool
	query := `SELECT EXISTS(SELECT 1 FROM usuarios WHERE login=?)`
	err := r.db.QueryRowContext(ctx, query, login).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("verificar existência: %w", err)
	}

	return existe, nil
}

// scanUsuario recebe um scanner (row ou rows) e retorna o usuário escaneado.
func (r *UsuarioRepository) scanUsuario(scanner scanner) (*usr.Usuario, error) {
	var usrDB usr.UsuarioDB
	err := scanner.Scan(
		&usrDB.ID,
		&usrDB.Nome,
		&usrDB.Login,
		&usrDB.Email,
		&usrDB.Permissao,
		&usrDB.Status,
		&usrDB.Avatar,
		&usrDB.UltimoLogin,
		&usrDB.CriadoEm,
		&usrDB.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}

	return usr.CarregarDoBD(usrDB), nil
}

// construirQueryListarUsuarios recebe um filtro e constrói a query SQL correspondente.
func (r *UsuarioRepository) construirQueryListar(filtro usr.Filtro) (string, []any) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`
		SELECT SQL_CALC_FOUND_ROWS 
			id, nome, login, email, permissao, status, avatar, 
			ultimo_login, criado_em, atualizado_em 
		FROM usuarios 
		WHERE 1=1
	`)

	if filtro.Busca() != nil && *filtro.Busca() != "" {
		padrao := "%" + *filtro.Busca() + "%"
		query.WriteString(" AND (nome LIKE ? OR login LIKE ? OR email LIKE ?)")
		args = append(args, padrao, padrao, padrao)
	}

	if filtro.Status() != nil {
		query.WriteString(" AND status = ?")
		args = append(args, *filtro.Status())
	}

	if filtro.Permissao() != nil && *filtro.Permissao() != "" {
		query.WriteString(" AND permissao = ?")
		args = append(args, *filtro.Permissao())
	}

	query.WriteString(" ORDER BY nome ASC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite(), filtro.Offset())
	return query.String(), args
}
