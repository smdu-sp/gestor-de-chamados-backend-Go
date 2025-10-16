package usecase

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// CategoriaPermissaoUsecase representa a camada de caso de uso para operações relacionadas a categorias de permissão.
type CategoriaPermissaoUsecase struct {
	repository repository.CategoriaPermissaoRepository
}

// NewCategoriaPermissaoUsecase cria uma nova instância de CategoriaPermissaoUsecase.
func NewCategoriaPermissaoUsecase(repository repository.CategoriaPermissaoRepository) *CategoriaPermissaoUsecase {
	return &CategoriaPermissaoUsecase{repository: repository}
}

// CriarCategoriaPermissao cria uma nova categoria de permissão.
func (c *CategoriaPermissaoUsecase) CriarCategoriaPermissao(ctx context.Context, categoriaPermissao *model.CategoriaPermissao) error {
	const metodo = "[usecase.CriarCategoriaPermissao]: %w"

	categoriaPermissao, err := model.NewCategoriaPermissao(
		categoriaPermissao.CategoriaID,
		categoriaPermissao.UsuarioID,
		categoriaPermissao.Permissao,
	)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	if err := c.repository.Salvar(ctx, categoriaPermissao); err != nil {
		return fmt.Errorf(metodo, err)
	}
	return nil
}

// AtualizarCategoriaPermissao atualiza uma categoria de permissão existente.
func (c *CategoriaPermissaoUsecase) AtualizarCategoriaPermissao(ctx context.Context, categoriaID, usuarioID string, categoriaPermissao *model.CategoriaPermissao) error {
	if categoriaID == "" {
		return utils.NewAppError(
			"[usecase.AtualizarCategoriaPermissao]",
			utils.LevelInfo,
			"erro ao atualizar categoriaPermissao",
			model.ErrCategoriaPermissaoIDInvalido,
		)
	}
	if usuarioID == "" {
		return utils.NewAppError(
			"[usecase.AtualizarCategoriaPermissao]",
			utils.LevelInfo,
			"erro ao atualizar categoriaPermissao",
			model.ErrUsuarioIDInvalido,
		)
	}

	if err := c.repository.Atualizar(ctx, categoriaID, usuarioID, categoriaPermissao); err != nil {
		return fmt.Errorf("[usecase.AtualizarCategoriaPermissao]: %w", err)
	}
	return nil
}

// DeletarCategoriaPermissao remove uma categoria de permissão existente.
func (c *CategoriaPermissaoUsecase) DeletarCategoriaPermissao(ctx context.Context, categoriaID, usuarioID string) error {
	if categoriaID == "" {
		return utils.NewAppError(
			"[usecase.DeletarCategoriaPermissao]",
			utils.LevelInfo,
			"erro ao deletar categoriaPermissao",
			model.ErrCategoriaPermissaoIDInvalido,
		)
	}
	if usuarioID == "" {
		return utils.NewAppError(
			"[usecase.DeletarCategoriaPermissao]",
			utils.LevelInfo,
			"erro ao deletar categoriaPermissao",
			model.ErrUsuarioIDInvalido,
		)
	}

	if err := c.repository.Deletar(ctx, categoriaID, usuarioID); err != nil {
		return fmt.Errorf("[usecase.DeletarCategoriaPermissao]: %w", err)
	}
	return nil
}

// ListarCategoriaPermissao lista categorias de permissão com base em filtros e paginação.
func (c *CategoriaPermissaoUsecase) ListarCategoriaPermissao(ctx context.Context, filtro model.CategoriaPermissaoFiltro) ([]model.CategoriaPermissao, int, model.CategoriaPermissaoFiltro, error) {
	if filtro.Pagina <= 0 {
		filtro.Pagina = 1
	}

	if filtro.Limite <= 0 || filtro.Limite > 100 {
		filtro.Limite = 10
	}

	categoriasPermissao, total, err := c.repository.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, filtro, fmt.Errorf("[usecase.ListarCategoriaPermissao]: %w", err)
	}

	return categoriasPermissao, total, filtro, nil
}