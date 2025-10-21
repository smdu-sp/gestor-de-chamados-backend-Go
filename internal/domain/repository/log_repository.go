package repository

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarLog define métodos de busca de logs
type BuscarLog interface {
	// BuscarPorID retorna um log pelo seu ID
	BuscarPorID(ctx context.Context, id string) (*model.Log, error)
}

// ArmazenarLog define métodos para armazenamento de logs
type ArmazenarLog interface {
	// Salvar insere um novo log no repositório
	Salvar(ctx context.Context, l *model.Log) error
}

// ListarLog define métodos para listagem e busca filtrada de logs
type ListarLog interface {
	// Listar retorna uma lista de logs com base em filtros e paginação
	Listar(ctx context.Context, filtro model.LogFiltro) ([]model.Log, int, error)
}

// LogRepository é uma composição de todas as interfaces acima
type LogRepository interface {
	BuscarLog
	ArmazenarLog
	ListarLog
}
