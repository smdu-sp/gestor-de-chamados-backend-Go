package repository

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarChamado define métodos de busca do chamado
type BuscarChamado interface {
	// BuscarPorID busca um chamado pelo seu ID.
	BuscarPorID(ctx context.Context, id string) (*model.Chamado, error)
}

// ArmazenarChamado define métodos para salvar/atualizar/excluir chamados
type ArmazenarChamado interface {
	// Salvar cria um novo chamado.
	Salvar(ctx context.Context, c *model.Chamado) error

	// Atualizar atualiza as informações de um chamado existente.
	Atualizar(ctx context.Context, id string, c *model.Chamado) error

	// Arquivar marca um chamado como arquivado.
	Arquivar(ctx context.Context, id string) error

	// Desarquivar marca um chamado como não arquivado.
	Desarquivar(ctx context.Context, id string) error
}

// AtualizarChamado define métodos específicos de atualização
type AtualizarChamado interface {
	// AtualizarStatus atualiza o status de um chamado, podendo incluir uma solução.
	AtualizarStatus(ctx context.Context, id string, status string, solucao *string) error
}

// ListarChamado define métodos para listagem e busca filtrada
type ListarChamado interface {
	// Listar lista chamados com paginação e filtros opcionais.
	Listar(ctx context.Context, filtro model.ChamadoFiltro) ([]model.Chamado, int, error)
}

// ChamadoRepository é uma composição de todas as interfaces acima
type ChamadoRepository interface {
	BuscarChamado
	ArmazenarChamado	
	AtualizarChamado
	ListarChamado
}
