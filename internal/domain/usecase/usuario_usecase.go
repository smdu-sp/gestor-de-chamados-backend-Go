package usecase

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarUsuario é a interface que define os métodos para obter informações de usuários.
type BuscarUsuario interface {
	// BuscarUsuarioPorID busca um usuário pelo ID.
	BuscarUsuarioPorID(ctx context.Context, id string) (*model.Usuario, error)

	// BuscarUsuarioPorLogin busca um usuário pelo login.
	BuscarUsuarioPorLogin(ctx context.Context, login string) (*model.Usuario, error)
}

// ArmazenarUsuario é a interface que define os métodos para criar, atualizar, ativar e desativar usuários.
type ArmazenarUsuario interface {
	// CriarUsuario cria um novo usuário.
	CriarUsuario(ctx context.Context, u *model.Usuario) error

	// AtualizarUsuario atualiza as informações de um usuário existente.
	AtualizarUsuario(ctx context.Context, id string, u *model.Usuario) error

	// DesativarUsuario desativa um usuário.
	DesativarUsuario(ctx context.Context, id string) error

	// AtivarUsuario ativa um usuário.
	AtivarUsuario(ctx context.Context, id string) error

	// VerificarPermissao verifica se o usuário possui uma das permissões especificadas.
	VerificarPermissao(ctx context.Context, id string, permissoes ...string) (bool, error)
}

// AtualizarUsuario é a interface que define os métodos para atualizar informações específicas de usuários.
type AtualizarUsuario interface {
	// AtualizarUltimoLoginUsuario atualiza o campo de último login do usuário.
	AtualizarUltimoLoginUsuario(ctx context.Context, id string) error

	// AtualizarPermissaoUsuario atualiza a permissão do usuário.
	AtualizarPermissaoUsuario(ctx context.Context, id string, permissao string) error
}

// ListarUsuarios é a interface que define os métodos para listar e buscar usuários com filtros.
type ListarUsuarios interface {
	// Listar lista usuários com paginação e filtros opcionais.
	ListarUsuarios(ctx context.Context, filtro model.UsuarioFiltro) ([]model.Usuario, int, model.UsuarioFiltro, error)
}

// UsuarioUsecase é a interface que agrega os casos de uso relacionados a usuários.
type UsuarioUsecase interface {
	BuscarUsuario
	ArmazenarUsuario
	AtualizarUsuario
	ListarUsuarios
}
