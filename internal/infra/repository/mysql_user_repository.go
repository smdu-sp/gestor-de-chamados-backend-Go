package repository

import (
		"context"
		"database/sql"
		"errors"
		"strings"
		"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// MySQLUserRepository implementa a interface UserRepository para MySQL.
type MySQLUserRepository struct {
    db *sql.DB
}

// NewMySQLUserRepository cria uma nova instância de MySQLUserRepository.
func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
    return &MySQLUserRepository{db: db}
}

// FindByID busca um usuário pelo ID.
func (r *MySQLUserRepository) FindByID(ctx context.Context, id string) (*model.Usuario, error) {
    return r.findOne(ctx, `SELECT id, nome, login, email, permissao, status, avatar, ultimo_login, criado_em, atualizado_em
        FROM usuarios WHERE id=?`, id)
}

// FindByLogin busca um usuário pelo login.
func (r *MySQLUserRepository) FindByLogin(ctx context.Context, login string) (*model.Usuario, error) {
    u, err := r.findOne(ctx, `SELECT id, nome, login, email, permissao, status, avatar, ultimo_login, criado_em, atualizado_em
        FROM usuarios WHERE login=?`, login)
    if err != nil {
        return nil, err
    }
    if u == nil {
        return nil, model.ErrUsuarioNaoEncontrado
    }
    return u, nil
}

// Insert insere um novo usuário no banco de dados.
func (r *MySQLUserRepository) Insert(ctx context.Context, u *model.Usuario) error {
    // Geração de ID pode ficar aqui ou na camada usecase
    _, err := r.db.ExecContext(ctx, `INSERT INTO usuarios(
        id, nome, login, email, permissao, status, avatar, ultimo_login, criado_em, atualizado_em
    ) VALUES (?,?,?,?,?,?,?,NOW(),NOW(),NOW())`,
        u.ID, u.Nome, u.Login, u.Email, u.Permissao, u.Status, u.Avatar,
    )
    return err
}

// Update atualiza os dados de um usuário existente.
func (r *MySQLUserRepository) Update(ctx context.Context, id string, u *model.Usuario) error {
    _, err := r.db.ExecContext(ctx, `UPDATE usuarios
        SET nome=?, email=?, permissao=?, status=?, avatar=?, atualizado_em=NOW()
        WHERE id=?`,
        u.Nome, u.Email, u.Permissao, u.Status, u.Avatar, id,
    )
    return err
}

// Delete remove um usuário do banco de dados pelo ID.
func (r *MySQLUserRepository) Delete(ctx context.Context, id string) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM usuarios WHERE id=?`, id)
    return err
}

// UpdateLastLogin atualiza o campo ultimo_login para o horário atual.
func (r *MySQLUserRepository) UpdateLastLogin(ctx context.Context, id string) error {
    _, err := r.db.ExecContext(ctx, `UPDATE usuarios SET ultimo_login=NOW() WHERE id=?`, id)
    return err
}

// List retorna uma lista paginada de usuários, com filtros opcionais.
func (r *MySQLUserRepository) List(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]model.Usuario, int, error) {
    var query strings.Builder
    args := []any{}

    query.WriteString(`SELECT SQL_CALC_FOUND_ROWS id, nome, login, email, permissao, status, avatar, ultimo_login, criado_em, atualizado_em
        FROM usuarios WHERE 1=1`)

    if busca != nil && *busca != "" {
        padrao := "%" + *busca + "%"
        query.WriteString(" AND (nome LIKE ? OR login LIKE ? OR email LIKE ?)")
        args = append(args, padrao, padrao, padrao)
    }

    if status != nil && *status != "" {
        query.WriteString(" AND status=?")
        args = append(args, *status == "true")
    }

    if permissao != nil && *permissao != "" {
        query.WriteString(" AND permissao=?")
        args = append(args, *permissao)
    }

    query.WriteString(" ORDER BY nome ASC LIMIT ? OFFSET ?")
    args = append(args, limite, (pagina-1)*limite)

    rows, err := r.db.QueryContext(ctx, query.String(), args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var usuarios []model.Usuario
    for rows.Next() {
        u, err := scanUsuario(rows)
        if err != nil {
            return nil, 0, err
        }
        usuarios = append(usuarios, *u)
    }

    var total int
    if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
        return nil, 0, err
    }

    return usuarios, total, nil
}

// Métodos auxiliares privados

// findOne executa uma consulta que retorna um único usuário.
func (r *MySQLUserRepository) findOne(ctx context.Context, query string, args ...any) (*model.Usuario, error) {
    row := r.db.QueryRowContext(ctx, query, args...)
    return scanUsuario(row)
}

// scanUsuario mapeia os dados de um scanner (row ou rows) para um objeto Usuario.
func scanUsuario(scanner interface {
    Scan(dest ...any) error
}) (*model.Usuario, error) {
    var u model.Usuario
    err := scanner.Scan(
        &u.ID,
        &u.Nome,
        &u.Login,
        &u.Email,
        &u.Permissao,
        &u.Status,
        &u.Avatar,
        &u.UltimoLogin,
        &u.CriadoEm,
        &u.AtualizadoEm,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &u, nil
}