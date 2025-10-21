package repository

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarAtendimento define métodos de busca do atendimento
type BuscarAtendimento interface {
	// BuscarPorID retorna um atendimento pelo seu ID
	BuscarPorID(ctx context.Context, id string) (*model.Atendimento, error)
}

// ArmazenarAtendimento define métodos para armazenamento do atendimento
type ArmazenarAtendimento interface {
	// Salvar insere um novo atendimento no repositório
	Salvar(ctx context.Context, a *model.Atendimento) error

	// Atualizar modifica os dados de um atendimento existente
	Atualizar(ctx context.Context, id string, a *model.Atendimento) error
}

// ListarAtendimento define métodos para listagem e busca filtrada
type ListarAtendimento interface {
	// Listar retorna uma lista de atendimentos com base em filtros e paginação
	Listar(ctx context.Context, filtro model.AtendimentoFiltro) ([]model.Atendimento, int, error)
}

// AtendimentoRepository é uma composição de todas as interfaces acima
type AtendimentoRepository interface {
	BuscarAtendimento
	ArmazenarAtendimento
	ListarAtendimento
}

