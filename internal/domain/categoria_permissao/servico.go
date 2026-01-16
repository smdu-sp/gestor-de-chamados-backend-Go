package categoria_permissao

import (
	"context"

	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// CriarParams representa o payload para criar uma categoria de permissão.
type CriarParams struct {
	CategoriaID string
	UsuarioID   string
	Permissao   usr.Permissao
}

// AtualizarParams representa os dados necessários para atualizar uma categoria de permissão.
type AtualizarParams struct {
	Permissao usr.Permissao
}

// Leitor define operações de consulta de categoria de permissão.
type Leitor interface {
	BuscarPorID(ctx context.Context, categoriaID, usuarioID string) (*CategoriaPermissao, error)
	Listar(ctx context.Context, f Filtro) ([]CategoriaPermissao, int, Filtro, error)
}

// Escritor define operações de criação e atualização de categoria de permissão.
type Escritor interface {
	Criar(ctx context.Context, c CriarParams) (*CategoriaPermissao, error)
	Atualizar(ctx context.Context, categoriaID, usuarioID string, a AtualizarParams) (*CategoriaPermissao, error)
	Deletar(ctx context.Context, categoriaID, usuarioID string) error
}

// Service é a interface que agrega os casos de uso relacionados a categoria de permissão
type Service interface {
	Leitor
	Escritor
}
