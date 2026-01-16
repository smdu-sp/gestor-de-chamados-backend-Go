package log

import (
	"context"
)

// Leitor é a interface que define os métodos para obter informações de logs.
type Leitor interface {
	BuscarPorID(ctx context.Context, id string) (*Log, error)
	Listar(ctx context.Context, filtro LogFiltro) ([]Log, int, LogFiltro, error)
}

// Escritor é a interface que define os métodos para criar logs.
type Escritor interface {
	Criar(ctx context.Context, acao Acao, entidade, detalhes string)
}

// Service é a interface que agrega todos os casos de uso relacionados a logs.
type Service interface {
	Leitor
	Escritor
}
