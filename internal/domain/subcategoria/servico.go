package subcategoria

import (
	"context"
)

// CriarParams representa os dados necessários para criar uma subcategoria.
type CriarParams struct {
	Nome        string
	CategoriaID string
}

// AtualizarParams representa os dados necessários para atualizar uma subcategoria.
type AtualizarParams struct {
	Nome        *string
	CategoriaID *string
	Status      *bool
}

// Leitor define operações de consulta de subcategorias.
type Leitor interface {
	BuscarPorID(ctx context.Context, id string) (*Subcategoria, error)
	BuscarPorNome(ctx context.Context, nome string) (*Subcategoria, error)
	Listar(ctx context.Context, f Filtro) ([]Subcategoria, int, Filtro, error)
}

// Escritor define operações de criação e atualização de subcategorias.
type Escritor interface {
	Criar(ctx context.Context, c CriarParams) (*Subcategoria, error)
	Atualizar(ctx context.Context, id string, a AtualizarParams) (*Subcategoria, error)
	Ativar(ctx context.Context, id string) (*Subcategoria, error)
	Desativar(ctx context.Context, id string) (*Subcategoria, error)
}

// Service é a interface que agrega os casos de uso relacionados a subcategorias.
type Service interface {
	Leitor
	Escritor
}
