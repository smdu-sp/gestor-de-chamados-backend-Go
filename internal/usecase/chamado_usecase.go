package usecase

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// ChamadoUsecase representa a camada de caso de uso para operações relacionadas a chamados.
type ChamadoUsecase struct {
	repository repository.ChamadoRepository
}

// NewChamadoUsecase cria uma nova instância de ChamadoUsecase.
func NewChamadoUsecase(repository repository.ChamadoRepository) *ChamadoUsecase {
	return &ChamadoUsecase{repository: repository}
}

// BuscarChamadoPorID busca um chamado pelo seu ID.
func (c *ChamadoUsecase) BuscarChamadoPorID(ctx context.Context, id string) (*model.Chamado, error) {
	if id == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarChamadoPorID]",
			utils.LevelInfo,
			"erro ao buscar chamado por id",
			model.ErrIDInvalido,
		)
	}

	chamado, err := c.repository.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarChamadoPorID]: %w", err)
	}
	return chamado, nil
}

// CriarChamado cria um novo chamado.
func (c *ChamadoUsecase) CriarChamado(ctx context.Context, chamado *model.Chamado) error {
	const metodo = "[usecase.CriarChamado]: %w"

	id, err := utils.NewUUIDv7String()
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	chamado.ID = id

	if chamado.Status == "" {
		chamado.Status = model.StatusAberto
	}

	chamado, err = model.NewChamado(
		chamado.ID,
		chamado.Titulo,
		chamado.Descricao,
		chamado.Status,
		chamado.CategoriaID,
		chamado.SubcategoriaID,
		chamado.CriadorID,
	)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	if err := c.repository.Salvar(ctx, chamado); err != nil {
		return fmt.Errorf(metodo, err)
	}

	return nil
}

// AtualizarChamado atualiza um chamado existente.
func (c *ChamadoUsecase) AtualizarChamado(ctx context.Context, id string, chamado *model.Chamado) error {
	if err := model.ValidarChamado(chamado); err != nil {
		return fmt.Errorf("[usecase.AtualizarChamado] %w", err)
	}

	if err := c.repository.Atualizar(ctx, id, chamado); err != nil {
		return fmt.Errorf("[usecase.AtualizarChamado]: %w", err)
	}

	return nil
}

// ArquivarChamado arquiva um chamado existente.
func (c *ChamadoUsecase) ArquivarChamado(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.ArquivarChamado]",
			utils.LevelInfo,
			"erro ao arquivar chamado",
			model.ErrIDInvalido,
		)
	}

	if err := c.repository.Arquivar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.ArquivarChamado]: %w", err)
	}

	return nil
}

// DesarquivarChamado desarquiva um chamado existente.
func (c *ChamadoUsecase) DesarquivarChamado(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.DesarquivarChamado]",
			utils.LevelInfo,
			"erro ao desarquivar chamado",
			model.ErrIDInvalido,
		)
	}

	if err := c.repository.Desarquivar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.DesarquivarChamado]: %w", err)
	}

	return nil
}

// AtualizarStatusChamado atualiza o status de um chamado existente.
func (c *ChamadoUsecase) AtualizarStatusChamado(ctx context.Context, id string, status string, solucao *string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtualizarStatusChamado]",
			utils.LevelInfo,
			"erro ao atualizar status do chamado",
			model.ErrIDInvalido,
		)
	}

	if err := model.ValidarStatusChamado(status); err != nil {
		return fmt.Errorf("[usecase.AtualizarStatusChamado]: %w", err)
	}

	if err := c.repository.AtualizarStatus(ctx, id, status, solucao); err != nil {
		return fmt.Errorf("[usecase.AtualizarStatusChamado]: %w", err)
	}

	return nil
}

// AtribuirTecnicoChamado atribui um técnico a um chamado existente.
func (c *ChamadoUsecase) AtribuirTecnicoChamado(ctx context.Context, id string, tecnicoID string) error {
// TODO adicionar verificação se permissão do usuário é correta em relação a categoria do chamado

	if id == "" || tecnicoID == "" {
		return utils.NewAppError(
			"[usecase.AtribuirTecnicoChamado]",
			utils.LevelInfo,
			"erro ao atribuir técnico ao chamado",
			model.ErrIDInvalido,
		)
	}

	if err := c.repository.AtribuirTecnico(ctx, id, tecnicoID); err != nil {
		return fmt.Errorf("[usecase.AtribuirTecnicoChamado]: %w", err)
	}

	return nil
}

// RemoverTecnicoChamado remove o técnico atribuído de um chamado existente.
func (c *ChamadoUsecase) RemoverTecnicoChamado(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.RemoverTecnicoChamado]",
			utils.LevelInfo,
			"erro ao remover técnico do chamado",
			model.ErrIDInvalido,
		)
	}

	if err := c.repository.RemoverTecnico(ctx, id); err != nil {
		return fmt.Errorf("[usecase.RemoverTecnicoChamado]: %w", err)
	}

	return nil
}

// ListarChamados lista todos os chamados com paginação.
func (c *ChamadoUsecase) ListarChamados(ctx context.Context, filtro model.ChamadoFiltro) ([]model.Chamado, int, model.ChamadoFiltro, error) {
	if filtro.Pagina <= 0 {
		filtro.Pagina = 1
	}

	if filtro.Limite <= 0 || filtro.Limite > 100 {
		filtro.Limite = 10
	}

	chamados, total, err := c.repository.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, filtro, fmt.Errorf("[usecase.ListarChamados]: %w", err)
	}

	return chamados, total, filtro, nil
}