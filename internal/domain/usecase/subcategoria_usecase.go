package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarSubcategoria é a interface que define os métodos para obter informações de subcategorias.
type BuscarSubcategoria interface {
	// BuscarSubcategoriaPorID busca uma subcategoria pelo ID.
	BuscarSubcategoriaPorID(ctx context.Context, id string) (*model.Subcategoria, error)

	// BuscarSubcategoriaPorNome busca uma subcategoria pelo nome.
	BuscarSubcategoriaPorNome(ctx context.Context, nome string) (*model.Subcategoria, error)
}

// ArmazenarSubcategoria é a interface que define os métodos para criar, atualizar, ativar e desativar subcategorias.
type ArmazenarSubcategoria interface {
	// CriarSubcategoria cria uma nova subcategoria.
	CriarSubcategoria(ctx context.Context, s *model.Subcategoria) error

	// AtualizarSubcategoria atualiza as informações de uma subcategoria existente.
	AtualizarSubcategoria(ctx context.Context, id string, s *model.Subcategoria) error

	// AtivarSubcategoria ativa uma subcategoria pelo seu ID.
	AtivarSubcategoria(ctx context.Context, id string) error

	// DesativarSubcategoria desativa (soft delete) uma subcategoria pelo seu ID.
	DesativarSubcategoria(ctx context.Context, id string) error
}

// ListarSubcategorias é a interface que define os métodos para listar e buscar subcategorias com filtros.
type ListarSubcategorias interface {
	// ListarSubcategorias lista subcategorias com paginação e filtros opcionais.
	ListarSubcategorias(ctx context.Context, filtro model.SubcategoriaFiltro) ([]model.Subcategoria, int, model.SubcategoriaFiltro, error)
}

// SubcategoriaUsecase é a interface que agrega os casos de uso relacionados a subcategorias.
type SubcategoriaUsecase interface {
	BuscarSubcategoria
	ArmazenarSubcategoria
	ListarSubcategorias
}