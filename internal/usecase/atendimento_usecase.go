package usecase

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// AtendimentoUsecase representa a camada de caso de uso para operações relacionadas a atendimentos.
type AtendimentoUsecase struct {
	repository repository.AtendimentoRepository
}

// NewAtendimentoUsecase cria uma nova instância de AtendimentoUsecase.
func NewAtendimentoUsecase(repository repository.AtendimentoRepository) *AtendimentoUsecase {
	return &AtendimentoUsecase{repository: repository}
}

// BuscarAtendimentoPorID busca um atendimento pelo seu ID.
func (u *AtendimentoUsecase) BuscarAtendimentoPorID(ctx context.Context, id string) (*model.Atendimento, error) {
	if id == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarAtendimentoPorID]",
			utils.LevelInfo,
			"erro ao buscar atendimento por id",
			model.ErrAtendimentoIDInvalido,
		)
	}

	atendimento, err := u.repository.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarAtendimentoPorID]: %w", err)
	}
	return atendimento, nil
}

// CriarAtendimento salva um novo atendimento.
func (u *AtendimentoUsecase) CriarAtendimento(ctx context.Context, atendimento *model.Atendimento) error {
	const metodo = "[usecase.CriarAtendimento]: %w"

	id, err := utils.NewUUIDv7String()
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	atendimento.ID = id

	// TODO 

	atendimento, err = model.NewAtendimento(
		atendimento.ID,
		atendimento.AtribuidoID,
		atendimento.ChamadoID,
	)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	if err := u.repository.Salvar(ctx, atendimento); err != nil {
		return fmt.Errorf(metodo, err)
	}
	return nil
}

// AtualizarAtendimento atualiza um atendimento existente.
func (u *AtendimentoUsecase) AtualizarAtendimento(ctx context.Context, id string, atendimento *model.Atendimento) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtualizarAtendimento]",
			utils.LevelInfo,
			"erro ao atualizar atendimento",
			model.ErrAtendimentoIDInvalido,
		)
	}

	if err := model.ValidarAtendimento(atendimento); err != nil {
		return fmt.Errorf("[usecase.AtualizarAtendimento]: %w", err)
	}

	if err := u.repository.Atualizar(ctx, id, atendimento); err != nil {
		return fmt.Errorf("[usecase.AtualizarAtendimento]: %w", err)
	}
	return nil
}

// ListarAtendimentos lista atendimentos com base em filtros.
func (u *AtendimentoUsecase) ListarAtendimentos(ctx context.Context, filtro model.AtendimentoFiltro) ([]model.Atendimento, int, model.AtendimentoFiltro, error) {
	if filtro.Pagina < 1 {
		filtro.Pagina = 1
	}

	if filtro.Limite <= 0 || filtro.Limite > 100 {
		filtro.Limite = 10
	}

	atendimentos, total, err := u.repository.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, filtro, fmt.Errorf("[usecase.ListarAtendimentos]: %w", err)
	}
	return atendimentos, total, filtro, nil
}