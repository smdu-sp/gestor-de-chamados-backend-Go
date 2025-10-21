package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarChamado é a interface que define os métodos para obter informações de chamados.
type BuscarChamado interface {
	// BuscarChamadoPorID busca um chamado pelo ID.
	BuscarChamadoPorID(ctx context.Context, id string) (*model.Chamado, error)
}

// ArmazenarChamado é a interface que define os métodos para criar e atualizar chamados.
type ArmazenarChamado interface {
	// CriarChamado cria um novo chamado.
	CriarChamado(ctx context.Context, c *model.Chamado) error

	// AtualizarChamado atualiza as informações de um chamado existente.
	AtualizarChamado(ctx context.Context, id string, c *model.Chamado) error

	// ArquivarChamado marca um chamado como arquivado.
	ArquivarChamado(ctx context.Context, id string) error

	// DesarquivarChamado marca um chamado como não arquivado.
	DesarquivarChamado(ctx context.Context, id string) error
}

// AtualizarChamado é a interface que define os métodos específicos de atualização de chamados.
type AtualizarChamado interface {
	// AtualizarStatusChamado atualiza o status de um chamado, podendo incluir uma solução.
	AtualizarStatusChamado(ctx context.Context, id string, status string, solucao *string) error
}

// ListarChamados é a interface que define os métodos para listar e buscar chamados com filtros.
type ListarChamados interface {
	// ListarChamados lista chamados com paginação e filtros opcionais.
	ListarChamados(ctx context.Context, filtro model.ChamadoFiltro) ([]model.Chamado, int, model.ChamadoFiltro, error)
}

// ChamadoUsecase é a interface que agrega os casos de uso relacionados a chamados.
type ChamadoUsecase interface {
	BuscarChamado
	ArmazenarChamado
	AtualizarChamado
	ListarChamados
}