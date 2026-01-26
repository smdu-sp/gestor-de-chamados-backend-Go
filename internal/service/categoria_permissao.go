package service

import (
	"context"
	"fmt"

	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
)

// CategoriaPermissaoService representa a camada de caso de uso para operações relacionadas a categorias de permissão.
type CategoriaPermissaoService struct {
	repo cpm.Repository
}

// NovoCategoriaPermissaoService cria uma nova instância de CategoriaPermissaoService.
func NovoCategoriaPermissaoService(repo cpm.Repository) *CategoriaPermissaoService {
	return &CategoriaPermissaoService{repo: repo}
}

// Asserção de interface para garantir que CategoriaPermissaoService implementa CategoriaPermissaoService
var _ cpm.Service = (*CategoriaPermissaoService)(nil)

// Criar recebe os parâmetros para criar uma nova categoria de permissão e a salva no repositório.
//
// Erros sentinela possíveis: ErrCategoriaPermissaoJaExiste.
func (c *CategoriaPermissaoService) Criar(ctx context.Context, criar cpm.CriarParams) (*cpm.CategoriaPermissao, error) {
	// 1 - Criar a categoria de permissão
	cpm, err := cpm.Novo(
		criar.CategoriaID,
		criar.UsuarioID,
		criar.Permissao,
	)
	if err != nil {
		return nil, err
	}

	// 2 - Salvar a categoria de permissão no repositório
	cpmCriada, err := c.repo.Criar(ctx, *cpm)
	if err != nil {
		return nil, fmt.Errorf("criar categoria/permissão: %w", err)
	}

	return cpmCriada, nil
}

// BuscarPorID recebe o ID composto de categoria e usuário, e retorna a categoria de permissão correspondente.
//
// Erros sentinela possíveis: ErrCategoriaPermissaoNaoEncontrada.
func (c *CategoriaPermissaoService) BuscarPorID(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
	cpm, err := c.repo.BuscarPorID(ctx, ctgID, usrID)
	if err != nil {
		return nil, fmt.Errorf("buscar categoriaPermissão por ID composto: %w", err)
	}
	return cpm, nil
}

// Atualizar recebe o ID composto de categoria e usuário, e os parâmetros para atualizar a categoria de permissão existente.
//
// Erros sentinela possíveis: ErrCategoriaPermissaoNaoEncontrada.
func (c *CategoriaPermissaoService) Atualizar(ctx context.Context, ctgID, usrID string, atualizar cpm.AtualizarParams) (*cpm.CategoriaPermissao, error) {
	// 1 - Buscar a categoria de permissão existente
	categoriaPermissaoAtual, err := c.repo.BuscarPorID(ctx, ctgID, usrID)
	if err != nil {
		return nil, fmt.Errorf("atualizar categoriaPermissão: %w", err)
	}

	// 2 - Atualizar os dados da categoria de permissão
	categoriaPermissaoAtualizada, err := categoriaPermissaoAtual.ComDadosAtualizados(atualizar)
	if err != nil {
		return nil, err
	}

	// 3 - Salvar as alterações no repositório
	categoriaPermissaoSalva, err := c.repo.Atualizar(ctx, ctgID, usrID, categoriaPermissaoAtualizada)
	if err != nil {
		return nil, fmt.Errorf("atualizar categoriaPermissão: %w", err)
	}

	return categoriaPermissaoSalva, nil
}

// Deletar recebe o ID composto de categoria e usuário, e remove a categoria de permissão correspondente do repositório.
//
// Erros sentinela possíveis: ErrCategoriaPermissaoNaoEncontrada.
func (c *CategoriaPermissaoService) Deletar(ctx context.Context, ctgID, usrID string) error {
	// 1 - Verificar se a categoria de permissão existe
	_, err := c.repo.BuscarPorID(ctx, ctgID, usrID)
	if err != nil {
		return fmt.Errorf("deletar categoriaPermissão: %w", err)
	}

	// 2 - Deletar a categoria de permissão do repositório
	if err := c.repo.Deletar(ctx, ctgID, usrID); err != nil {
		return fmt.Errorf("deletar categoriaPermissão: %w", err)
	}
	return nil
}

// Listar recebe um filtro e retorna a lista de categorias e permissões correspondentes,
// o total de registros e o filtro aplicado.
func (c *CategoriaPermissaoService) Listar(ctx context.Context, f cpm.Filtro) ([]cpm.CategoriaPermissao, int, cpm.Filtro, error) {
	f.Normalizar()
	cpmSlice, total, err := c.repo.Listar(ctx, f)
	if err != nil {
		return nil, 0, f, fmt.Errorf("listar categoriaPermissão: %w", err)
	}

	return cpmSlice, total, f, nil
}
