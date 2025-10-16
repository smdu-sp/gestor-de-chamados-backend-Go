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

// Erros sentinela do repositório de usuários.
var (
	ErrUsuarioNaoEncontrado = errors.New("usuário não encontrado no banco de dados MySQL")
	ErrUsuarioJaExiste      = errors.New("erro ao tentar duplicar usuário no banco de dados MySQL")
	ErrScannerUsuario       = errors.New("erro ao escanear usuário do banco de dados MySQL")
	ErrExecContext          = errors.New("erro ao executar comando no banco de dados MySQL")
	ErrRowsAffected         = errors.New("erro ao obter número de linhas afetadas no banco de dados MySQL")
	ErrQueryContext         = errors.New("erro ao executar consulta no banco de dados MySQL")
	ErrScan                 = errors.New("erro ao escanear resultado da consulta no banco de dados MySQL")
)

// MySQLUsuarioRepository implementa a interface UserRepository para MySQL.
type MySQLUsuarioRepository struct {
	db *sql.DB
}

// NewMySQLUsuarioRepository cria uma nova instância de MySQLUsuarioRepository.
func NewMySQLUsuarioRepository(db *sql.DB) *MySQLUsuarioRepository {
	return &MySQLUsuarioRepository{db: db}
}

// BuscarPorID busca um usuário pelo ID.
func (r *MySQLUsuarioRepository) BuscarPorID(ctx context.Context, id string) (*model.Usuario, error) {
	usuario, err := r.buscar(
		ctx,
		`SELECT id, nome, login, email, permissao, status, 
		 avatar, ultimo_login, criado_em, atualizado_em
     FROM usuarios 
		 WHERE id=?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("[MySQLUsuarioRepository.BuscarPorID]: %w", err)
	}

	if usuario == nil {
		return nil, utils.NewAppError(
			"[MySQLUsuarioRepository.BuscarPorID]",
			utils.LevelInfo,
			"a busca por ID não retornou resultados",
			ErrUsuarioNaoEncontrado,
		)
	}

	return usuario, nil
}

// BuscarPorLogin busca um usuário pelo login.
func (r *MySQLUsuarioRepository) BuscarPorLogin(ctx context.Context, login string) (*model.Usuario, error) {
	usuario, err := r.buscar(
		ctx,
		`SELECT id, nome, login, email, permissao, status,
		 avatar, ultimo_login, criado_em, atualizado_em
     FROM usuarios 
		 WHERE login=?`,
		login,
	)
	if err != nil {
		return nil, fmt.Errorf("[MySQLUsuarioRepository.BuscarPorLogin]: %w", err)
	}

	if usuario == nil {
		return nil, utils.NewAppError(
			"[MySQLUsuarioRepository.BuscarPorLogin]",
			utils.LevelWarning,
			"a busca por login não retornou resultados",
			ErrUsuarioNaoEncontrado,
		)
	}

	return usuario, nil
}

// Salvar insere um novo usuário no banco de dados.
func (r *MySQLUsuarioRepository) Salvar(ctx context.Context, u *model.Usuario) error {
	const metodo = "[MySQLUsuarioRepository.Salvar]"

	resultado, err := r.db.ExecContext(
		ctx,
		`INSERT INTO usuarios(
     id, nome, login, email, permissao, status, 
	   avatar, ultimo_login, criado_em, atualizado_em
    ) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), NOW())`,
		u.ID, u.Nome, u.Login, u.Email, u.Permissao, u.Status, u.Avatar,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			errorMensagem := err.Error()

			if strings.Contains(errorMensagem, "login") {
				return utils.NewAppError(
					metodo,
					utils.LevelWarning,
					"erro ao salvar usuário com login duplicado no banco de dados",
					ErrUsuarioJaExiste,
				)
			}

			if strings.Contains(errorMensagem, "email") {
				return utils.NewAppError(
					metodo,
					utils.LevelWarning,
					"erro ao salvar usuário com email duplicado no banco de dados",
					ErrUsuarioJaExiste,
				)
			}
		}
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro inesperado ao salvar usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao obter número de linhas afetadas ao salvar usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrRowsAffected, err),
		)
	}
	if linhasAfetadas == 0 {
		return utils.NewAppError(
			metodo,
			utils.LevelWarning,
			"nenhuma linha foi afetada ao salvar usuário no banco de dados",
			ErrExecContext,
		)
	}

	return nil
}

// Atualizar atualiza os dados de um usuário existente.
func (r *MySQLUsuarioRepository) Atualizar(ctx context.Context, id string, u *model.Usuario) error {
	const metodo = "[MySQLUsuarioRepository.Atualizar]"

	existe, err := ExisteUsuarioPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("%s: %w", metodo, err)
	}
	if !existe {
		return utils.NewAppError(
			metodo,
			utils.LevelInfo,
			"não foi possível atualizar o usuário",
			ErrUsuarioNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE usuarios
     SET nome=?, email=?, permissao=?, 
		 status=?, avatar=?, atualizado_em=NOW()
     WHERE id=?`,
		u.Nome, u.Email, u.Permissao, u.Status, u.Avatar, id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			errorMensagem := err.Error()

			if strings.Contains(errorMensagem, "email") {
				return utils.NewAppError(
					metodo,
					utils.LevelWarning,
					"erro ao atualizar usuário com email duplicado no banco de dados",
					fmt.Errorf(utils.FmtErroWrap, ErrUsuarioJaExiste, err),
				)
			}
		}
		return utils.NewAppError(
			metodo,
			utils.LevelError,
			"erro ao atualizar usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// AtualizarPermissao atualiza a permissão de um usuário.
func (r *MySQLUsuarioRepository) AtualizarPermissao(ctx context.Context, id string, permissao string) error {
	existe, err := ExisteUsuarioPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLUsuarioRepository.AtualizarPermissao]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLUsuarioRepository.AtualizarPermissao]",
			utils.LevelInfo,
			"não foi possível atualizar a permissão do usuário",
			ErrUsuarioNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE usuarios 
		 SET permissao=?, atualizado_em=NOW()
     WHERE id=?`,
		permissao, id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLUsuarioRepository.AtualizarPermissao]",
			utils.LevelError,
			"erro ao atualizar permissão do usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Desativar desativa um usuário do banco de dados pelo ID.
func (r *MySQLUsuarioRepository) Desativar(ctx context.Context, id string) error {
	existe, err := ExisteUsuarioPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLUsuarioRepository.Desativar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLUsuarioRepository.Desativar]",
			utils.LevelInfo,
			"não foi possível desativar o usuário",
			ErrUsuarioNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE usuarios 
		 SET status=false, atualizado_em=NOW() 
		 WHERE id=?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLUsuarioRepository.Desativar]",
			utils.LevelError,
			"erro ao desativar usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Ativar ativa um usuário do banco de dados pelo ID.
func (r *MySQLUsuarioRepository) Ativar(ctx context.Context, id string) error {
	existe, err := ExisteUsuarioPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLUsuarioRepository.Ativar]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLUsuarioRepository.Ativar]",
			utils.LevelInfo,
			"não foi possível ativar o usuário",
			ErrUsuarioNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE usuarios 
		 SET status=true 
		 WHERE id=?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLUsuarioRepository.Ativar]",
			utils.LevelError,
			"erro ao ativar usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// AtualizarUltimoLogin atualiza o campo ultimo_login para o horário atual.
func (r *MySQLUsuarioRepository) AtualizarUltimoLogin(ctx context.Context, id string) error {
	existe, err := ExisteUsuarioPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("[MySQLUsuarioRepository.AtualizarUltimoLogin]: %w", err)
	}
	if !existe {
		return utils.NewAppError(
			"[MySQLUsuarioRepository.AtualizarUltimoLogin]",
			utils.LevelInfo,
			"não foi possível atualizar o último login do usuário",
			ErrUsuarioNaoEncontrado,
		)
	}

	_, err = r.db.ExecContext(ctx,
		`UPDATE usuarios 
		 SET ultimo_login=NOW()
		 WHERE id=?`,
		id,
	)
	if err != nil {
		return utils.NewAppError(
			"[MySQLUsuarioRepository.AtualizarUltimoLogin]",
			utils.LevelError,
			"erro ao atualizar último login do usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrExecContext, err),
		)
	}

	return nil
}

// Listar retorna uma lista paginada de usuários, com filtros opcionais.
func (r *MySQLUsuarioRepository) Listar(ctx context.Context, filtro model.UsuarioFiltro) ([]model.Usuario, int, error) {
	var query strings.Builder
	args := []any{}

	// TODO nao trazer os arquivados, incluir flag para exibir ou nao status false
	query.WriteString(
		`SELECT SQL_CALC_FOUND_ROWS 
     id, nome, login, email, permissao, status, 
		 avatar, ultimo_login, criado_em, atualizado_em
     FROM usuarios 
		 WHERE 1=1`,
	)

	if filtro.Busca != nil && *filtro.Busca != "" {
		padrao := "%" + *filtro.Busca + "%"
		query.WriteString(" AND (nome LIKE ? OR login LIKE ? OR email LIKE ?)")
		args = append(args, padrao, padrao, padrao)
	}

	if filtro.Status != nil {
		query.WriteString(" AND status=?")
		args = append(args, *filtro.Status)
	}

	if filtro.Permissao != nil && *filtro.Permissao != "" {
		query.WriteString(" AND permissao=?")
		args = append(args, *filtro.Permissao)
	}

	query.WriteString(" ORDER BY nome ASC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite, (filtro.Pagina-1)*filtro.Limite)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLUsuarioRepository.Listar]",
			utils.LevelError,
			"erro ao listar usuários no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}
	defer rows.Close()

	var usuarios []model.Usuario
	for rows.Next() {
		usuario, err := scanUsuario(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("[MySQLUsuarioRepository.Listar]: %w", err)
		}
		usuarios = append(usuarios, *usuario)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return nil, 0, utils.NewAppError(
			"[MySQLUsuarioRepository.Listar]",
			utils.LevelError,
			"erro ao obter total de usuários",
			fmt.Errorf(utils.FmtErroWrap, ErrScan, err),
		)
	}

	return usuarios, total, nil
}

// Métodos auxiliares

// buscar executa uma consulta que retorna um único usuário.
func (r *MySQLUsuarioRepository) buscar(ctx context.Context, query string, args ...any) (*model.Usuario, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	usuario, err := scanUsuario(row)
	if err != nil {
		return nil, fmt.Errorf("[MySQLUsuarioRepository.buscar]: %w", err)
	}
	return usuario, nil
}

// ExisteUsuarioPorID verifica se um usuário existe com base no ID.
func ExisteUsuarioPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var existe bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM usuarios WHERE id=?)`, id).Scan(&existe)
	if err != nil {
		return false, utils.NewAppError(
			"[MySQLUsuarioRepository.ExisteUsuarioPorID]",
			utils.LevelError,
			"erro ao verificar existência de usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}

	if !existe {
		return false, nil
	}
	return true, nil
}

// ExistePorLogin verifica se um usuário existe com base no login.
func (r *MySQLUsuarioRepository) ExistePorLogin(ctx context.Context, login string) (bool, error) {
	var existe bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM usuarios WHERE login=?)`, login).Scan(&existe)
	if err != nil {
		return false, utils.NewAppError(
			"[MySQLUsuarioRepository.ExistePorLogin]",
			utils.LevelError,
			"erro ao verificar existência de usuário no banco de dados",
			fmt.Errorf(utils.FmtErroWrap, ErrQueryContext, err),
		)
	}

	if !existe {
		return false, nil
	}
	return true, nil
}

// scanUsuario mapeia os dados de um scanner (row ou rows) para uma struct Usuario.
func scanUsuario(scanner interface{ Scan(dest ...any) error }) (*model.Usuario, error) {
	var usuario model.Usuario
	err := scanner.Scan(
		&usuario.ID,
		&usuario.Nome,
		&usuario.Login,
		&usuario.Email,
		&usuario.Permissao,
		&usuario.Status,
		&usuario.Avatar,
		&usuario.UltimoLogin,
		&usuario.CriadoEm,
		&usuario.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, utils.NewAppError(
			"[MySQLUsuarioRepository.scanUsuario]",
			utils.LevelError,
			"o scanner falhou ao escanear o usuário",
			fmt.Errorf(utils.FmtErroWrap, ErrScannerUsuario, err),
		)
	}
	return &usuario, nil
}
