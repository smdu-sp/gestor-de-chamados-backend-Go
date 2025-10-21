package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarAtendimento define métodos de busca do atendimento
type BuscarAtendimento interface {
	// BuscarAtendimentoPorID retorna um atendimento pelo seu ID
	BuscarAtendimentoPorID(ctx context.Context, id string) (*model.Atendimento, error)
}

// ArmazenarAtendimento define métodos para armazenamento do atendimento
type ArmazenarAtendimento interface {
	// CriarAtendimento insere um novo atendimento no repositório
	CriarAtendimento(ctx context.Context, a *model.Atendimento) error

	// AtualizarAtendimento modifica os dados de um atendimento existente
	AtualizarAtendimento(ctx context.Context, id string, a *model.Atendimento) error
}

// ListarAtendimento define métodos para listagem e busca filtrada
type ListarAtendimento interface {
	// ListarAtendimentos retorna uma lista de atendimentos com base em filtros e paginação
	ListarAtendimentos(ctx context.Context, filtro model.AtendimentoFiltro) ([]model.Atendimento, int, model.AtendimentoFiltro, error)
}

// AtendimentoUseCase é uma composição de todas as interfaces acima
type AtendimentoUseCase interface {
	BuscarAtendimento
	ArmazenarAtendimento
	ListarAtendimento
}