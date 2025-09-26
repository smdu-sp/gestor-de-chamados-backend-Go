package usecase

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// CategoriaUsecase representa a camada de caso de uso para operações relacionadas a categorias.
type CategoriaUsecase struct {
	repository repository.CategoriaRepository
}

// NewCategoriaUsecase cria uma nova instância de CategoriaUsecase.
func NewCategoriaUsecase(repository repository.CategoriaRepository) *CategoriaUsecase {
	return &CategoriaUsecase{repository: repository}
}

// BuscarCategoriaPorID busca uma categoria pelo seu ID.
func (c *CategoriaUsecase) BuscarCategoriaPorID(ctx context.Context, id string) (*model.Categoria, error) {
	if id == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarCategoriaPorID]",
			utils.LevelWarning,
			"erro ao buscar categoria por id",
			model.ErrIDInvalido,
		)
	}

	categoria, err := c.repository.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarCategoriaPorID]: %w", err)
	}
	return categoria, nil
}

// BuscarCategoriaPorNome busca uma categoria pelo seu nome.
func (c *CategoriaUsecase) BuscarCategoriaPorNome(ctx context.Context, nome string) (*model.Categoria, error) {
	if nome == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarCategoriaPorNome]",
			utils.LevelInfo,
			"erro ao buscar categoria por nome",
			model.ErrNomeInvalido,
		)
	}

	categoria, err := c.repository.BuscarPorNome(ctx, nome)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarCategoriaPorNome]: %w", err)
	}
	return categoria, nil
}

// CriarCategoria cria uma nova categoria.
func (c *CategoriaUsecase) CriarCategoria(ctx context.Context, categoria *model.Categoria) error {
	const metodo = "[usecase.CriarCategoria]: %w"

	id, err := utils.NewUUIDv7String()
	if err != nil {
		return fmt.Errorf(metodo, err)
	}
	categoria.ID = id

	categoria.Status = true

	categoria, err = model.NewCategoria(
		categoria.ID,
		categoria.Nome,
		categoria.Status,
	)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	if err := c.repository.Salvar(ctx, categoria); err != nil {
		return fmt.Errorf(metodo, err)
	}
	return nil
}

// AtualizarCategoria atualiza as informações de uma categoria existente.
func (c *CategoriaUsecase) AtualizarCategoria(ctx context.Context, id string, categoria *model.Categoria) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtualizarCategoria]",
			utils.LevelInfo,
			"erro ao atualizar categoria",
			model.ErrIDInvalido,
		)
	}

	if err := c.repository.Atualizar(ctx, id, categoria); err != nil {
		return fmt.Errorf("[usecase.AtualizarCategoria]: %w", err)
	}
	return nil
}

// DesativarCategoria desativa (soft delete) uma categoria.
func (c *CategoriaUsecase) DesativarCategoria(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.DesativarCategoria]",
			utils.LevelInfo,
			"erro ao desativar categoria",
			model.ErrIDInvalido,
		)
	}

	if err := c.repository.Desativar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.DesativarCategoria]: %w", err)
	}
	return nil
}

// AtivarCategoria ativa uma categoria.
func (c *CategoriaUsecase) AtivarCategoria(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtivarCategoria]",
			utils.LevelInfo,
			"erro ao ativar categoria",
			model.ErrIDInvalido,
		)
	}

	if err := c.repository.Ativar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.AtivarCategoria]: %w", err)
	}
	return nil
}

// ListarCategorias lista categorias com paginação e filtros opcionais.
func (c *CategoriaUsecase) ListarCategorias(ctx context.Context, filtro model.CategoriaFiltro) ([]model.Categoria, int, model.CategoriaFiltro, error) {
	if filtro.Pagina <= 0 {
		filtro.Pagina = 1
	}

	if filtro.Limite <= 0 || filtro.Limite > 100 {
		filtro.Limite = 10
	}

	categorias, total, err := c.repository.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, filtro, fmt.Errorf("[usecase.ListarCategorias]: %w", err)
	}

	return categorias, total, filtro, nil
}