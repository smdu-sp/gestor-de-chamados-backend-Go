package repository

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarAcompanhamento define métodos de busca do acompanhamento
type BuscarAcompanhamento interface {
	// BuscarPorID retorna um acompanhamento pelo seu ID
	BuscarPorID(ctx context.Context, id string) (*model.Acompanhamento, error)

	// BuscarPorChamadoID retorna uma lista de acompanhamentos pelo ID do chamado
	BuscarPorChamadoID(ctx context.Context, chamadoID string) ([]model.Acompanhamento, error)
}

// ArmazenarAcompanhamento define métodos para armazenamento do acompanhamento
type ArmazenarAcompanhamento interface {
	// Salvar insere um novo acompanhamento no repositório
	Salvar(ctx context.Context, a *model.Acompanhamento) error

	// Atualizar modifica os dados de um acompanhamento existente
	Atualizar(ctx context.Context, id string, a *model.Acompanhamento) error

	// Deletar remove um acompanhamento do repositório pelo seu ID
	Deletar(ctx context.Context, id string) error
}

// ListarAcompanhamento define métodos para listagem e busca filtrada
type ListarAcompanhamento interface {
	// Listar retorna uma lista de acompanhamentos com base em filtros e paginação
	Listar(ctx context.Context, filtro model.AcompanhamentoFiltro) ([]model.Acompanhamento, int, error)
}

// AcompanhamentoRepository é uma composição de todas as interfaces acima
type AcompanhamentoRepository interface {
	BuscarAcompanhamento
	ArmazenarAcompanhamento
	ListarAcompanhamento
}