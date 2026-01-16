package usuario

import (
	"context"
	"time"
)

// UsuarioDB representa a estrutura do usuário no banco de dados.
type UsuarioDB struct {
	ID           string
	Nome         string
	Login        string
	Email        string
	Permissao    string
	Status       bool
	Avatar       *string
	UltimoLogin  time.Time
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// Repository é a interface que define os métodos para interagir com o armazenamento de usuários.
type Repository interface {
	Criar(ctx context.Context, u Usuario) (*Usuario, error)
	Atualizar(ctx context.Context, id string, u Usuario) (*Usuario, error)
	BuscarPorID(ctx context.Context, id string) (*Usuario, error)
	BuscarPorLogin(ctx context.Context, login string) (*Usuario, error)
	Listar(ctx context.Context, filtro Filtro) ([]Usuario, int, error)
}