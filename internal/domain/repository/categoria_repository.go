package repository

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarCategoria define métodos de busca da categoria
type BuscarCategoria interface {
	// BuscarPorID busca uma categoria pelo seu ID.
	BuscarPorID(ctx context.Context, id string) (*model.Categoria, error)

	// BuscarPorNome busca uma categoria pelo seu nome.
	BuscarPorNome(ctx context.Context, nome string) (*model.Categoria, error)
}

// ArmazenarCategoria define métodos para salvar/atualizar/excluir categorias
type ArmazenarCategoria interface {
	// Salvar cria uma nova categoria.
	Salvar(ctx context.Context, c *model.Categoria) error

	// Atualizar atualiza as informações de uma categoria existente.
	Atualizar(ctx context.Context, id string, c *model.Categoria) error

	// Ativar ativa uma categoria pelo seu ID.
	Ativar(ctx context.Context, id string) error

	// Desativar desativa (soft delete) uma categoria pelo seu ID.
	Desativar(ctx context.Context, id string) error
}

// ListarCategoria define métodos para listagem e busca filtrada
type ListarCategoria interface {
	// Listar lista categorias com paginação e filtros opcionais.
	Listar(ctx context.Context, filtro model.CategoriaFiltro) ([]model.Categoria, int, error)
}

// CategoriaRepository é uma composição de todas as interfaces acima
type CategoriaRepository interface {
	BuscarCategoria
	ArmazenarCategoria
	ListarCategoria
}
