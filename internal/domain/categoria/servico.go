package categoria

import (
	"context"
)

// CriarParams representa os dados necessários para criar uma categoria.
type CriarParams struct {
	Nome string
}

// AtualizarParams representa os dados necessários para atualizar uma categoria.
type AtualizarParams struct {
	Nome   *string
	Status *bool
}

// Leitor define os métodos de consulta para categorias.
type Leitor interface {
	BuscarPorID(ctx context.Context, id string) (*Categoria, error)
	BuscarPorNome(ctx context.Context, nome string) (*Categoria, error)
	Listar(ctx context.Context, filtro Filtro) ([]Categoria, int, Filtro, error)
}

// Escritor define os métodos de escrita para categorias.
type Escritor interface {
	Criar(ctx context.Context, c CriarParams) (*Categoria, error)
	Atualizar(ctx context.Context, id string, c AtualizarParams) (*Categoria, error)
	Desativar(ctx context.Context, id string) (*Categoria, error)
	Ativar(ctx context.Context, id string) (*Categoria, error)
}

// Service é a interface que agrega os casos de uso relacionados a categorias.
type Service interface {
	Leitor
	Escritor
}
