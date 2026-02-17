package service

import (
	"context"
	"fmt"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
)

// CategoriaService representa a camada de caso de uso para operações relacionadas a categorias.
type CategoriaService struct {
	id   dmn.GeradorID
	repo ctg.Repository
}

// NovoCategoriaService cria uma nova instância de CategoriaService.
func NovoCategoriaService(id dmn.GeradorID, repo ctg.Repository) *CategoriaService {
	return &CategoriaService{id: id, repo: repo}
}

// Asserção de interface para garantir que CategoriaService implementa CategoriaService
var _ ctg.Service = (*CategoriaService)(nil)

// BuscarPorID recebe o ID de uma categoria e retorna a categoria correspondente.
//
// Erros sentinela possíveis: ErrCategoriaNaoEncontrada.
func (c *CategoriaService) BuscarPorID(ctx context.Context, id string) (*ctg.Categoria, error) {
	ctg, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar categoria por ID: %w", err)
	}
	return ctg, nil
}

// BuscarPorNome recebe o nome de uma categoria e retorna a categoria correspondente.
//
// Erros sentinela possíveis: ErrCategoriaNaoEncontrada.
func (c *CategoriaService) BuscarPorNome(ctx context.Context, nome string) (*ctg.Categoria, error) {
	ctg, err := c.repo.BuscarPorNome(ctx, nome)
	if err != nil {
		return nil, fmt.Errorf("buscar categoria por nome: %w", err)
	}
	return ctg, nil
}

// Criar recebe os parâmetros para criar uma nova categoria e retorna a categoria criada.
func (c *CategoriaService) Criar(ctx context.Context, criar ctg.CriarParams) (*ctg.Categoria, error) {
	// 1 - Gerar ID
	id, err := c.id.NovoID()
	if err != nil {
		return nil, fmt.Errorf("criar categoria: %w", err)
	}

	// 2 - Criar a nova categoria
	ctg, err := ctg.Novo(
		id,
		criar.Nome,
	)
	if err != nil {
		return nil, err
	}

	// 3 - Salvar a categoria no repositório
	ctgCriada, err := c.repo.Criar(ctx, *ctg)
	if err != nil {
		return nil, fmt.Errorf("criar categoria: %w", err)
	}
	return ctgCriada, nil
}

// Atualizar recebe o ID de uma categoria e os parâmetros para atualização,
// e retorna a categoria atualizada.
//
// Erros sentinela possíveis: ErrCategoriaNaoEncontrada, ErrCategoriaJaExiste, ErrosValidacao.
func (c *CategoriaService) Atualizar(ctx context.Context, id string, atualizar ctg.AtualizarParams) (*ctg.Categoria, error) {
	// 1 - Buscar a categoria existente
	categoriaAtual, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar categoria: %w", err)
	}

	// 2 - Atualizar os dados da categoria
	categoriaAtualizada, err := categoriaAtual.ComDadosAtualizados(atualizar)
	if err != nil {
		return nil, err
	}

	// 3 - Salvar as alterações no repositório
	categoriaSalva, err := c.repo.Atualizar(ctx, id, categoriaAtualizada)
	if err != nil {
		return nil, fmt.Errorf("atualizar categoria: %w", err)
	}
	return categoriaSalva, nil
}

// Desativar desativa (soft delete) uma categoria.
//
// Erros sentinela possíveis: ErrCategoriaNaoEncontrada.
func (c *CategoriaService) Desativar(ctx context.Context, id string) (*ctg.Categoria, error) {
	// 1 - Buscar a categoria existente
	categoriaAtual, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("desativar categoria: %w", err)
	}

	// 2 - Desativar a categoria
	categoriaAtualizada := categoriaAtual.Desativar()

	// 3 - Salvar as alterações no repositório
	categoriaDesativada, err := c.repo.Atualizar(ctx, id, categoriaAtualizada)
	if err != nil {
		return nil, fmt.Errorf("desativar categoria: %w", err)
	}

	return categoriaDesativada, nil
}

// Ativar recebe o ID de uma categoria e a ativa.
//
// Erros sentinela possíveis: ErrCategoriaNaoEncontrada.
func (c *CategoriaService) Ativar(ctx context.Context, id string) (*ctg.Categoria, error) {
	// 1 - Buscar a categoria existente
	categoriaAtual, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ativar categoria: %w", err)
	}

	// 2 - Ativar a categoria
	categoriaAtualizada := categoriaAtual.Ativar()

	// 3 - Salvar as alterações no repositório
	categoriaAtivada, err := c.repo.Atualizar(ctx, id, categoriaAtualizada)
	if err != nil {
		return nil, fmt.Errorf("ativar categoria: %w", err)
	}
	
	return categoriaAtivada, nil
}

// Listar recebe um filtro e retorna a lista de categorias correspondentes,
// o total de registros e o filtro aplicado.
func (c *CategoriaService) Listar(ctx context.Context, f ctg.Filtro) ([]ctg.Categoria, int, ctg.Filtro, error) {
	f.Normalizar()
	ctgSlice, total, err := c.repo.Listar(ctx, f)
	if err != nil {
		return nil, 0, ctg.Filtro{}, fmt.Errorf("listar categoria: %w", err)
	}

	return ctgSlice, total, f, nil
}
