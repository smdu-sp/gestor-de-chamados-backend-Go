package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// ArmazenarCategoriaPermissao define métodos para armazenamento da categoria de permissão
type ArmazenarCategoriaPermissao interface {
	// CriarCategoriaPermissao insere uma nova categoria de permissão no repositório
	CriarCategoriaPermissao(ctx context.Context, c *model.CategoriaPermissao) error

	// AtualizarCategoriaPermissao modifica os dados de uma categoria de permissão existente
	AtualizarCategoriaPermissao(ctx context.Context, categoriaID, usuarioID string, c *model.CategoriaPermissao) error

	// DeletarCategoriaPermissao remove uma categoria de permissão do repositório
	DeletarCategoriaPermissao(ctx context.Context, categoriaID, usuarioID string) error
}

// ListarCategoriaPermissao define métodos para listagem e busca filtrada
type ListarCategoriaPermissao interface {
	// ListarCategoriaPermissao retorna uma lista de categorias de permissão com base em filtros e paginação
	ListarCategoriaPermissao(ctx context.Context, filtro model.CategoriaPermissaoFiltro) ([]model.CategoriaPermissao, int, model.CategoriaPermissaoFiltro, error)
}

// CategoriaPermissaoUsecase é a interface que agrega os casos de uso relacionados a categoria de permissão
type CategoriaPermissaoUsecase interface {
	ArmazenarCategoriaPermissao
	ListarCategoriaPermissao
}