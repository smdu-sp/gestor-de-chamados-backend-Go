package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarAcompanhamento é a interface que define os métodos para obter informações de acompanhamentos.
type BuscarAcompanhamento interface {
	// BuscarAcompanhamentoPorID busca um acompanhamento pelo ID.
	BuscarAcompanhamentoPorID(ctx context.Context, id string) (*model.Acompanhamento, error)

	// BuscarAcompanhamentosPorChamadoID busca acompanhamentos pelo ID do chamado.
	BuscarAcompanhamentosPorChamadoID(ctx context.Context, chamadoID string) ([]model.Acompanhamento, error)
}

// ArmazenarAcompanhamento é a interface que define os métodos para criar e atualizar acompanhamentos.
type ArmazenarAcompanhamento interface {
	// CriarAcompanhamento cria um novo acompanhamento.
	CriarAcompanhamento(ctx context.Context, a *model.Acompanhamento) error

	// AtualizarAcompanhamento atualiza as informações de um acompanhamento existente.
	AtualizarAcompanhamento(ctx context.Context, id string, a *model.Acompanhamento) error

	// DeletarAcompanhamento remove um acompanhamento pelo ID.
	DeletarAcompanhamento(ctx context.Context, id string) error
}

// ListarAcompanhamentos é a interface que define os métodos para listar e buscar acompanhamentos com filtros.
type ListarAcompanhamentos interface {
	// ListarAcompanhamentos lista acompanhamentos com paginação e filtros opcionais.
	ListarAcompanhamentos(ctx context.Context, filtro model.AcompanhamentoFiltro) ([]model.Acompanhamento, int, model.AcompanhamentoFiltro, error)
}

// AcompanhamentoUsecase é a interface que agrega os casos de uso relacionados a acompanhamentos.
type AcompanhamentoUsecase interface {
	BuscarAcompanhamento
	ArmazenarAcompanhamento
	ListarAcompanhamentos
}
