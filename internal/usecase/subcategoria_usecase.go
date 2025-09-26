package usecase

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// SubcategoriaUsecase representa a camada de caso de uso para operações relacionadas a subcategorias.
type SubcategoriaUsecase struct {
	repository repository.SubcategoriaRepository
}

// NewSubcategoriaUsecase cria uma nova instância de SubcategoriaUsecase.
func NewSubcategoriaUsecase(repository repository.SubcategoriaRepository) *SubcategoriaUsecase {
	return &SubcategoriaUsecase{repository: repository}
}

// BuscarSubcategoriaPorID busca uma subcategoria pelo seu ID.
func (s *SubcategoriaUsecase) BuscarSubcategoriaPorID(ctx context.Context, id string) (*model.Subcategoria, error) {
	if id == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarSubcategoriaPorID]",
			utils.LevelWarning,
			"erro ao buscar subcategoria por id",
			model.ErrIDInvalido,
		)
	}

	subcategoria, err := s.repository.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarSubcategoriaPorID]: %w", err)
	}
	return subcategoria, nil
}

// BuscarSubcategoriaPorNome busca uma subcategoria pelo seu nome.
func (s *SubcategoriaUsecase) BuscarSubcategoriaPorNome(ctx context.Context, nome string) (*model.Subcategoria, error) {
	if nome == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarSubcategoriaPorNome]",
			utils.LevelInfo,
			"erro ao buscar subcategoria por nome",
			model.ErrNomeInvalido,
		)
	}

	subcategoria, err := s.repository.BuscarPorNome(ctx, nome)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarSubcategoriaPorNome]: %w", err)
	}
	return subcategoria, nil
}

// CriarSubcategoria cria uma nova subcategoria.
func (s *SubcategoriaUsecase) CriarSubcategoria(ctx context.Context, subcategoria *model.Subcategoria) error {
	const metodo = "[usecase.CriarSubcategoria]: %w"

	id, err := utils.NewUUIDv7String()
	if err != nil {
		return fmt.Errorf(metodo, err)
	}
	subcategoria.ID = id

	subcategoria.Status = true

	subcategoria, err = model.NewSubcategoria(
		subcategoria.ID,
		subcategoria.Nome,
		subcategoria.CategoriaID,
		subcategoria.Status,
	)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	if err := s.repository.Salvar(ctx, subcategoria); err != nil {
		return fmt.Errorf(metodo, err)
	}
	return nil
}

// AtualizarSubcategoria atualiza uma subcategoria existente.
func (s *SubcategoriaUsecase) AtualizarSubcategoria(ctx context.Context, id string, subcategoria *model.Subcategoria) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtualizarSubcategoria]",
			utils.LevelInfo,
			"erro ao atualizar subcategoria",
			model.ErrIDInvalido,
		)
	}

	if err := s.repository.Atualizar(ctx, id, subcategoria); err != nil {
		return fmt.Errorf("[usecase.AtualizarSubcategoria]: %w", err)
	}
	return nil
}

// DesativarSubcategoria desativa uma subcategoria existente.
func (s *SubcategoriaUsecase) DesativarSubcategoria(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.DesativarSubcategoria]",
			utils.LevelInfo,
			"erro ao desativar subcategoria",
			model.ErrIDInvalido,
		)
	}

	if err := s.repository.Desativar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.DesativarSubcategoria]: %w", err)
	}
	return nil
}

// AtivarSubcategoria ativa uma subcategoria existente.
func (s *SubcategoriaUsecase) AtivarSubcategoria(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtivarSubcategoria]",
			utils.LevelInfo,
			"erro ao ativar subcategoria",
			model.ErrIDInvalido,
		)
	}

	if err := s.repository.Ativar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.AtivarSubcategoria]: %w", err)
	}
	return nil
}

// ListarSubcategorias lista todas as subcategorias.
func (s *SubcategoriaUsecase) ListarSubcategorias(ctx context.Context, filtro model.SubcategoriaFiltro) ([]model.Subcategoria, int, model.SubcategoriaFiltro, error) {
	if filtro.Pagina <= 0 {
		filtro.Pagina = 1
	}
	if filtro.Limite <= 0 || filtro.Limite > 100 {
		filtro.Limite = 10
	}

	subcategorias, total, err := s.repository.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, filtro, fmt.Errorf("[usecase.ListarSubcategorias]: %w", err)
	}
	return subcategorias, total, filtro, nil
}
