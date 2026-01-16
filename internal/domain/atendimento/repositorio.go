package atendimento

import (
	"context"
	"time"
)

// AtendimentoDB representa a estrutura do atendimento no banco de dados.
type AtendimentoDB struct {
	ID           string    
	AtribuidoID  string
	ChamadoID    string
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// Repository é a interface que define os métodos para interagir com o armazenamento de atendimentos.
type Repository interface {
	Criar(ctx context.Context, a Atendimento) (*Atendimento, error)
	Atualizar(ctx context.Context, id string, a Atendimento) (*Atendimento, error)
	BuscarPorID(ctx context.Context, id string) (*Atendimento, error)
	BuscarPorChamadoEAtribuidoID(ctx context.Context, chamadoID, atribuidoID string) (*Atendimento, error)
	Listar(ctx context.Context, filtro Filtro) ([]Atendimento, int, error)
}
