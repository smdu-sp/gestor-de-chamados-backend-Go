package acompanhamento

import (
	"context"
)

// CriarParams representa os dados necessários para criar um acompanhamento.
type CriarParams struct {
	Conteudo  string
	ChamadoID string
}

// AtualizarParams representa os dados necessários para atualizar um acompanhamento.
type AtualizarParams struct {
	Conteudo string
}

// Leitor é a interface que define as operações de leitura para acompanhamentos.
type Leitor interface {
	BuscarPorID(ctx context.Context, id string) (*Acompanhamento, error)
	BuscarPorChamadoID(ctx context.Context, id string) ([]Acompanhamento, error)
	Listar(ctx context.Context, f Filtro) ([]Acompanhamento, int, Filtro, error)
}

// Escritor é a interface que define as operações de criação e atualização para acompanhamentos.
type Escritor interface {
	Criar(ctx context.Context, c CriarParams) (*Acompanhamento, error)
	Atualizar(ctx context.Context, id string, a AtualizarParams) (*Acompanhamento, error)
	Deletar(ctx context.Context, id string) error
}

// Service é a interface que agrega os casos de uso relacionados a acompanhamentos.
type Service interface {
	Leitor
	Escritor
}
