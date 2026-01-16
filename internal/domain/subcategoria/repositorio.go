package subcategoria

import (
	"context"
	"time"
)

// SubcategoriaDB representa a estrutura da subcategoria no banco de dados.
type SubcategoriaDB struct {
	ID           string
	Nome         string
	Status       bool
	CategoriaID  string
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// Repository é a interface que define os métodos para interagir com o armazenamento de subcategorias.
type Repository interface {
	Criar(ctx context.Context, c Subcategoria) (*Subcategoria, error)
	Atualizar(ctx context.Context, id string, c Subcategoria) (*Subcategoria, error)
	BuscarPorID(ctx context.Context, id string) (*Subcategoria, error)
	BuscarPorNome(ctx context.Context, nome string) (*Subcategoria, error)
	Listar(ctx context.Context, filtro Filtro) ([]Subcategoria, int, error)
}
