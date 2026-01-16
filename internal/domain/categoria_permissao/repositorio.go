package categoria_permissao

import (
	"context"
	"time"
)

// CategoriaPermissaoDB representa a estrutura de uma categoria de permissão no banco de dados.
type CategoriaPermissaoDB struct {
	CategoriaID  string
	UsuarioID    string
	Permissao    string
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// Repository é a interface que combina as operações de armazenamento, busca e listagem
type Repository interface {
	Criar(ctx context.Context, c CategoriaPermissao) (*CategoriaPermissao, error)
	Atualizar(ctx context.Context, categoriaID, usuarioID string, c CategoriaPermissao) (*CategoriaPermissao, error)
	Deletar(ctx context.Context, categoriaID, usuarioID string) error
	BuscarPorID(ctx context.Context, categoriaID, usuarioID string) (*CategoriaPermissao, error)
	Listar(ctx context.Context, f Filtro) ([]CategoriaPermissao, int, error)
}
