package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarLog é a interface que define os métodos para obter informações de logs.
type BuscarLog interface {
	// BuscarLogPorID busca um log pelo ID.
	BuscarLogPorID(ctx context.Context, id string) (*model.Log, error)
}

// ArmazenarLog é a interface que define os métodos para criar logs.
type ArmazenarLog interface {
	// CriarLog cria um novo log.
	CriarLog(ctx context.Context, acao model.Acao, entidade, detalhes string) error
}

// ListarLogs é a interface que define os métodos para listar e buscar logs com filtros.
type ListarLogs interface {
	// ListarLogs lista logs com paginação e filtros opcionais.
	ListarLogs(ctx context.Context, filtro model.LogFiltro) ([]model.Log, int, model.LogFiltro, error)
}

// LogUsecase é a interface que agrega os casos de uso relacionados a logs.
type LogUsecase interface {
	BuscarLog
	ArmazenarLog
	ListarLogs
}