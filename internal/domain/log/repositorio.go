package log

import (
	"context"
	"time"
)

// LogDB representa a estrutura do log no banco de dados.
type LogDB struct {
	ID        string
	UsuarioID string
	Acao      string
	Entidade  string
	Detalhes  string
	CriadoEm  time.Time
}

// Repository é a interface que define os métodos para interagir com o armazenamento de logs.
type Repository interface {
	Criar(ctx context.Context, l Log) error
	BuscarPorID(ctx context.Context, id string) (*Log, error)
	Listar(ctx context.Context, filtro LogFiltro) ([]Log, int, error)
}
