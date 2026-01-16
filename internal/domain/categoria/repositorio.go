package categoria

import (
	"context"
	"time"
)

// CategoriaDB representa a estrutura da categoria no banco de dados.
type CategoriaDB struct {
	ID           string
	Nome         string
	Status       bool
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// Repository é a interface que define os métodos para interagir com o armazenamento de categorias.
type Repository interface {
	Criar(ctx context.Context, c Categoria) (*Categoria, error)
	Atualizar(ctx context.Context, id string, c Categoria) (*Categoria, error)
	BuscarPorID(ctx context.Context, id string) (*Categoria, error)
	BuscarPorNome(ctx context.Context, nome string) (*Categoria, error)
	Listar(ctx context.Context, f Filtro) ([]Categoria, int, error)
}
