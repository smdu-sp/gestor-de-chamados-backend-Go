package chamado

import (
	"context"
)

// CriarParams representa os dados necessários para criar um chamado.
type CriarParams struct {
	Titulo         string
	Descricao      string
	CategoriaID    string
	SubcategoriaID string
	CriadorID      string
}

// AtualizarParams representa os dados necessários para um criador atualizar um chamado.
type AtualizarParams struct {
	Titulo         *string
	Descricao      *string
	Arquivado      *bool
	CategoriaID    *string
	SubcategoriaID *string
}

// AtualizarStatusParams representa os dados necessários para atualizar o status de um chamado.
type AtualizarStatusParams struct {
	Status  StatusChamado
	Solucao *string
}

// AtualizarSolucaoParams representa os dados necessários para atualizar a solução de um chamado.
type AtualizarSolucaoParams struct {
	Solucao string
}

// Leitor é a interface que define operações de consulta de chamados.
type Leitor interface {
	BuscarPorID(ctx context.Context, id string) (*Chamado, error)
	Listar(ctx context.Context, f Filtro) ([]Chamado, int, Filtro, error)
}

// Escritor é a interface operações de criação e atualização de chamados.
type Escritor interface {
	Criar(ctx context.Context, c CriarParams) (*Chamado, error)
	Atualizar(ctx context.Context, id string, a AtualizarParams) (*Chamado, error)
	Arquivar(ctx context.Context, id string) (*Chamado, error)
	Desarquivar(ctx context.Context, id string) (*Chamado, error)
	AtualizarStatus(ctx context.Context, id string, a AtualizarStatusParams) (*Chamado, error)
	AtualizarSolucao(ctx context.Context, id string, a AtualizarSolucaoParams) (*Chamado, error)
}

// Service agrega todos os casos de uso relacionados a chamados.
type Service interface {
	Leitor
	Escritor
}
