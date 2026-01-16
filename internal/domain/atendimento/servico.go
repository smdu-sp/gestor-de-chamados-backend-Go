package atendimento

import (
	"context"
)

// CriarParams representa os dados necessários para criar um atendimento.
type CriarParams struct {
	AtribuidoID string
	ChamadoID   string
}

// AtualizarParams representa os dados necessários para atualizar um atendimento.
type AtualizarParams struct {
	AtribuidoID string
	ChamadoID   string
}

// Leitor define operações de consulta de atendimentos.
type Leitor interface {
	BuscarPorID(ctx context.Context, id string) (*Atendimento, error)
	BuscarPorChamadoEAtribuidoID(ctx context.Context, chamadoID, atribuidoID string) (*Atendimento, error)
	Listar(ctx context.Context, filtro Filtro) ([]Atendimento, int, Filtro, error)
}

// Escritor define operações de criação e atualização de atendimentos.
type Escritor interface {
	Criar(ctx context.Context, c CriarParams) (*Atendimento, error)
	Atualizar(ctx context.Context, id string, a AtualizarParams) (*Atendimento, error)
}

// Service é uma composição de todas as interfaces acima
type Service interface {
	Leitor
	Escritor
}
