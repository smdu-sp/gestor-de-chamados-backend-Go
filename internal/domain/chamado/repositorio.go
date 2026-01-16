package chamado

import (
	"context"
	"time"
)

// ChamadoDB representa a estrutura de um chamado no banco de dados.
type ChamadoDB struct {
	ID             string
	CategoriaID    string
	SubcategoriaID string
	CriadorID      string
	Titulo         string
	Descricao      string
	Status         string
	Arquivado      bool
	Solucao        *string
	SolucionadoEm  *time.Time
	FechadoEm      *time.Time
	CriadoEm       time.Time
	AtualizadoEm   time.Time
}

// ChamadoRepository é a interface que define os métodos para interagir com o armazenamento de chamados.
type Repository interface {
	Criar(ctx context.Context, c Chamado) (*Chamado, error)
	Atualizar(ctx context.Context, id string, c Chamado) (*Chamado, error)
	BuscarPorID(ctx context.Context, id string) (*Chamado, error)
	Listar(ctx context.Context, f Filtro) ([]Chamado, int, error)
}
