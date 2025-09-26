package repository

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarSubcategoria define métodos de busca da subcategoria
type BuscarSubcategoria interface {
	// BuscarPorID busca uma subcategoria pelo seu ID.
	BuscarPorID(ctx context.Context, id string) (*model.Subcategoria, error)

	// BuscarPorNome busca uma subcategoria pelo seu nome.
	BuscarPorNome(ctx context.Context, nome string) (*model.Subcategoria, error)
}

// ArmazenarSubcategoria define métodos para salvar/atualizar/excluir subcategorias
type ArmazenarSubcategoria interface {
	// Salvar cria uma nova subcategoria.
	Salvar(ctx context.Context, c *model.Subcategoria) error

	// Atualizar atualiza as informações de uma subcategoria existente.
	Atualizar(ctx context.Context, id string, c *model.Subcategoria) error

	// Ativar ativa uma subcategoria pelo seu ID.
	Ativar(ctx context.Context, id string) error

	// Desativar desativa (soft delete) uma subcategoria pelo seu ID.
	Desativar(ctx context.Context, id string) error
}

// ListarSubcategoria define métodos para listar subcategorias
type ListarSubcategoria interface {
	// Listar todas as subcategorias com base nos critérios de filtro.
	Listar(ctx context.Context, filtro model.SubcategoriaFiltro) ([]model.Subcategoria, int, error)
}

// SubcategoriaRepository é uma composição de todas as interfaces acima
type SubcategoriaRepository interface {
	BuscarSubcategoria
	ArmazenarSubcategoria
	ListarSubcategoria
}