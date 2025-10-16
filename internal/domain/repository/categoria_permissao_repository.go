package repository

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// ArmazenarCategoriaPermissao define métodos para armazenamento da categoria de permissão
type ArmazenarCategoriaPermissao interface {
	// Salvar insere uma nova categoria de permissão no repositório
	Salvar(ctx context.Context, c *model.CategoriaPermissao) error

	// Atualizar modifica os dados de uma categoria de permissão existente
	Atualizar(ctx context.Context, categoriaID, usuarioID string, c *model.CategoriaPermissao) error

	// Deletar remove uma categoria de permissão do repositório
	Deletar(ctx context.Context, categoriaID, usuarioID string) error
}

// ListarCategoriaPermissao define métodos para listagem e busca filtrada
type ListarCategoriaPermissao interface {
	// Listar retorna uma lista de categorias de permissão com base em filtros e paginação
	Listar(ctx context.Context, filtro model.CategoriaPermissaoFiltro) ([]model.CategoriaPermissao, int, error)
}

// CategoriaPermissaoRepository é uma composição de todas as interfaces acima
type CategoriaPermissaoRepository interface {
	ArmazenarCategoriaPermissao
	ListarCategoriaPermissao
}
