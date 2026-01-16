package acompanhamento

import (
	"context"
	"time"
)

// AcompanhamentoDB representa a estrutura do acompanhamento no banco de dados.
type AcompanhamentoDB struct {
	ID           string
	ChamadoID    string
	UsuarioID    string
	Conteudo     string
	Remetente    string
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// Repository é a interface que define os métodos para interagir com o armazenamento de acompanhamentos.
type Repository interface {
	Criar(ctx context.Context, a Acompanhamento) (*Acompanhamento, error)
	Atualizar(ctx context.Context, id string, a Acompanhamento) (*Acompanhamento, error)
	BuscarPorID(ctx context.Context, id string) (*Acompanhamento, error)
	BuscarPorChamadoID(ctx context.Context, id string) ([]Acompanhamento, error)
	Deletar(ctx context.Context, id string) error
	Listar(ctx context.Context, f Filtro) ([]Acompanhamento, int, error)
}
