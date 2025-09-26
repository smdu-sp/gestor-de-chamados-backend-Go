package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarCategoria é a interface que define os métodos para obter informações de categorias.
type BuscarCategoria interface {
	// BuscarCategoriaPorID busca uma categoria pelo ID.
	BuscarCategoriaPorID(ctx context.Context, id string) (*model.Categoria, error)

	// BuscarCategoriaPorNome busca uma categoria pelo nome.
	BuscarCategoriaPorNome(ctx context.Context, nome string) (*model.Categoria, error)
}

// ArmazenarCategoria é a interface que define os métodos para criar, atualizar, ativar e desativar categorias.
type ArmazenarCategoria interface {
	// CriarCategoria cria uma nova categoria.
	CriarCategoria(ctx context.Context, c *model.Categoria) error

	// AtualizarCategoria atualiza as informações de uma categoria existente.
	AtualizarCategoria(ctx context.Context, id string, c *model.Categoria) error

	// DesativarCategoria desativa (soft delete) uma categoria.
	DesativarCategoria(ctx context.Context, id string) error

	// AtivarCategoria ativa uma categoria.
	AtivarCategoria(ctx context.Context, id string) error
}

// ListarCategorias é a interface que define os métodos para listar e buscar categorias com filtros.
type ListarCategorias interface {
	// ListarCategorias lista categorias com paginação e filtros opcionais.
	ListarCategorias(ctx context.Context, filtro model.CategoriaFiltro) ([]model.Categoria, int, model.CategoriaFiltro, error)
}

// CategoriaUsecase é a interface que agrega os casos de uso relacionados a categorias.
type CategoriaUsecase interface {
	BuscarCategoria
	ArmazenarCategoria
	ListarCategorias
}
