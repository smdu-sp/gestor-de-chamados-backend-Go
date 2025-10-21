	package usecase

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// AcompanhamentoUsecase representa a camada de caso de uso para operações relacionadas a acompanhamentos.
type AcompanhamentoUsecase struct {
	repository repository.AcompanhamentoRepository
}

// NewAcompanhamentoUsecase cria uma nova instância de AcompanhamentoUsecase.
func NewAcompanhamentoUsecase(repository repository.AcompanhamentoRepository) *AcompanhamentoUsecase {
	return &AcompanhamentoUsecase{repository: repository}
}

// BuscarAcompanhamentoPorID busca um acompanhamento pelo seu ID.
func (u *AcompanhamentoUsecase) BuscarAcompanhamentoPorID(ctx context.Context, id string) (*model.Acompanhamento, error) {
	if id == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarAcompanhamentoPorID]",
			utils.LevelInfo,
			"erro ao buscar acompanhamento por id",
			model.ErrAcompanhamentoIDInvalido,
		)
	}

	acompanhamento, err := u.repository.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarAcompanhamentoPorID]: %w", err)
	}
	return acompanhamento, nil
}

// BuscarAcompanhamentosPorChamadoID busca acompanhamentos pelo ID do chamado.
func (u *AcompanhamentoUsecase) BuscarAcompanhamentosPorChamadoID(ctx context.Context, chamadoID string) ([]model.Acompanhamento, error) {
	if chamadoID == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarAcompanhamentosPorChamadoID]",
			utils.LevelInfo,
			"erro ao buscar acompanhamentos por chamado id",
			model.ErrChamadoIDInvalido,
		)
	}

	acompanhamentos, err := u.repository.BuscarPorChamadoID(ctx, chamadoID)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarAcompanhamentosPorChamadoID]: %w", err)
	}
	return acompanhamentos, nil
}

// CriarAcompanhamento cria um novo acompanhamento.
func (u *AcompanhamentoUsecase) CriarAcompanhamento(ctx context.Context, acompanhamento *model.Acompanhamento) error {
	const metodo = "[usecase.CriarAcompanhamento]: %w"

	id, err := utils.NewUUIDv7String()
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	acompanhamento.ID = id

	acompanhamento, err = model.NewAcompanhamento(
		acompanhamento.ID,
		acompanhamento.ChamadoID,
		acompanhamento.UsuarioID,
		acompanhamento.Conteudo,
		acompanhamento.Remetente,
	)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	if err := u.repository.Salvar(ctx, acompanhamento); err != nil {
		return fmt.Errorf(metodo, err)
	}
	return nil
}

// AtualizarAcompanhamento atualiza as informações de um acompanhamento existente.
func (u *AcompanhamentoUsecase) AtualizarAcompanhamento(ctx context.Context, id string, acompanhamento *model.Acompanhamento) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtualizarAcompanhamento]",
			utils.LevelInfo,
			"erro ao atualizar acompanhamento",
			model.ErrAcompanhamentoIDInvalido,
		)
	}

	if err := model.ValidarAcompanhamento(acompanhamento); err != nil {
		return fmt.Errorf("[usecase.AtualizarAcompanhamento]: %w", err)
	}

	if err := u.repository.Atualizar(ctx, id, acompanhamento); err != nil {
		return fmt.Errorf("[usecase.AtualizarAcompanhamento]: %w", err)
	}
	return nil
}

// DeletarAcompanhamento remove um acompanhamento pelo ID.
func (u *AcompanhamentoUsecase) DeletarAcompanhamento(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.DeletarAcompanhamento]",
			utils.LevelInfo,
			"erro ao deletar acompanhamento",
			model.ErrAcompanhamentoIDInvalido,
		)
	}

	if err := u.repository.Deletar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.DeletarAcompanhamento]: %w", err)
	}
	return nil
}

// ListarAcompanhamentos lista acompanhamentos com paginação e filtros opcionais.
func (u *AcompanhamentoUsecase) ListarAcompanhamentos(ctx context.Context, filtro model.AcompanhamentoFiltro) ([]model.Acompanhamento, int, model.AcompanhamentoFiltro, error) {
	if filtro.Pagina < 1 {
		filtro.Pagina = 1
	}

	if filtro.Limite <= 0 || filtro.Limite > 100 {
		filtro.Limite = 10
	}

	acompanhamentos, total, err := u.repository.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, filtro, fmt.Errorf("[usecase.ListarAcompanhamentos]: %w", err)
	}
	return acompanhamentos, total, filtro, nil
}