package usuario

import "context"

// CriarParams representa os dados necessários para criar um usuário.
type CriarParams struct {
	Nome      string
	Login     string
	Email     Email
	Permissao Permissao
	Avatar    *string
}

// AtualizarParams representa os dados necessários para atualizar um usuário.
type AtualizarParams struct {
	Nome      *string
	Login     *string
	Email     *Email
	Permissao *Permissao
	Status    *bool
	Avatar    *string
}

// AtualizarPermissaoParams representa os dados necessários para atualizar a permissão de um usuário.
type AtualizarPermissaoParams struct {
	Permissao Permissao
}

// Leitor define operações de consulta de usuários.
type Leitor interface {
	BuscarPorID(ctx context.Context, id string) (*Usuario, error)
	BuscarPorLogin(ctx context.Context, login string) (*Usuario, error)
	Listar(ctx context.Context, f Filtro) ([]Usuario, int, Filtro, error)
}

// Escritor define operações de criação e atualização.
type Escritor interface {
	Criar(ctx context.Context, c CriarParams) (*Usuario, error)
	Atualizar(ctx context.Context, id string, a AtualizarParams) (*Usuario, error)
	VerificarPermissao(ctx context.Context, id string, permissoes []Permissao) (bool, error)
	AtualizarUltimoLogin(ctx context.Context, id string) (*Usuario, error)
	AtualizarPermissao(ctx context.Context, id string, in AtualizarPermissaoParams) (*Usuario, error)
	Ativar(ctx context.Context, id string) (*Usuario, error)
	Desativar(ctx context.Context, id string) (*Usuario, error)
}

// Service agrega todos os casos de uso relacionados a usuários
type Service interface {
	Leitor
	Escritor
}
