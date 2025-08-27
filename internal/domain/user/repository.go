package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// Repository encapsula operações de banco de dados para usuários
type Repository struct {
	db *sql.DB
}

// NewRepository cria um novo Repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FindByID retorna um usuário pelo ID ou nil se não encontrado
func (r *Repository) FindByID(ctx context.Context, id string) (*Usuario, error) {
	return r.findOne(ctx, `SELECT id, nome, login, email, permissao, status, avatar, ultimo_login, criado_em, atualizado_em
		FROM usuarios WHERE id=?`, id)
}

// FindByLogin retorna um usuário pelo login ou nil se não encontrado
func (r *Repository) FindByLogin(ctx context.Context, login string) (*Usuario, error) {
	return r.findOne(ctx, `SELECT id, nome, login, email, permissao, status, avatar, ultimo_login, criado_em, atualizado_em
		FROM usuarios WHERE login=?`, login)
}

// Insert adiciona um novo usuário, gerando UUID no Go
func (r *Repository) Insert(ctx context.Context, u *Usuario) error {
	id := uuid.New().String()
	_, err := r.db.ExecContext(ctx, `INSERT INTO usuarios(
		id, nome, login, email, permissao, status, avatar, ultimo_login, criado_em, atualizado_em
	) VALUES (?,?,?,?,?,?,?,NOW(),NOW(),NOW())`,
		id, u.Nome, u.Login, u.Email, u.Permissao, u.Status, u.Avatar,
	)
	return err
}

// Update altera os dados de um usuário existente
func (r *Repository) Update(ctx context.Context, id string, u *Usuario) error {
	_, err := r.db.ExecContext(ctx, `UPDATE usuarios
		SET nome=?, email=?, permissao=?, status=?, avatar=?, atualizado_em=NOW()
		WHERE id=?`,
		u.Nome, u.Email, u.Permissao, u.Status, u.Avatar, id,
	)
	return err
}

// Delete remove um usuário pelo ID
func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM usuarios WHERE id=?`, id)
	return err
}

// UpdateLastLogin atualiza a data de último login
func (r *Repository) UpdateLastLogin(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE usuarios SET ultimo_login=NOW() WHERE id=?`, id)
	return err
}

// List retorna usuários paginados, filtrando por busca, status e permissão
func (r *Repository) List(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]Usuario, int, error) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`SELECT SQL_CALC_FOUND_ROWS id, nome, login, email, permissao, status, avatar, ultimo_login, criado_em, atualizado_em
		FROM usuarios WHERE 1=1`)

	if busca != nil && *busca != "" {
		q := "%" + *busca + "%"
		query.WriteString(" AND (nome LIKE ? OR login LIKE ? OR email LIKE ?)")
		args = append(args, q, q, q)
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

	// executa consulta
	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// varre resultados
	var usuarios []Usuario
	for rows.Next() {
		u, err := scanUsuario(rows)
		if err != nil {
			return nil, 0, err
		}
		usuarios = append(usuarios, *u)
	}

	// total de registros encontrados
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return nil, 0, err
	}

	return usuarios, total, nil
}

// --- helpers internos ---

// findOne é helper para retornar um único usuário ou nil
func (r *Repository) findOne(ctx context.Context, query string, args ...any) (*Usuario, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	return scanUsuario(row)
}

// scanUsuario lê dados de um usuário de Row ou Rows
func scanUsuario(scanner interface {
	Scan(dest ...any) error
}) (*Usuario, error) {
	var u Usuario
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
