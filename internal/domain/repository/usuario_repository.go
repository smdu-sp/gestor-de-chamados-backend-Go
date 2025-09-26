package repository

import (
	"context"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// BuscarUsuario define métodos de busca do usuário
type BuscarUsuario interface {
	// BuscarPorID retorna um usuário pelo seu ID
	BuscarPorID(ctx context.Context, id string) (*model.Usuario, error)

	// BuscarPorLogin retorna um usuário pelo seu login
	BuscarPorLogin(ctx context.Context, login string) (*model.Usuario, error)
}

// ArmazenarUsuario define métodos para armazenamento do usuário
type ArmazenarUsuario interface {
	// Salvar insere um novo usuário no repositório
	Salvar(ctx context.Context, u *model.Usuario) error

	// Atualizar modifica os dados de um usuário existente
	Atualizar(ctx context.Context, id string, u *model.Usuario) error

	// Desativar desativa um usuário pelo seu ID
	Desativar(ctx context.Context, id string) error

	// Ativar reativa um usuário desativado pelo seu ID
	Ativar(ctx context.Context, id string) error
}

// AtualizarUsuario define métodos específicos de atualização
type AtualizarUsuario interface {
	// AtualizarSenha atualiza a senha de um usuário
	AtualizarUltimoLogin(ctx context.Context, id string) error

	// AtualizarPermissao atualiza a permissão de um usuário
	AtualizarPermissao(ctx context.Context, id string, permissao string) error
}

// ListarUsuario define métodos para listagem e busca filtrada
type ListarUsuario interface {
	// Listar retorna uma lista de usuários com base em filtros e paginação
	Listar(ctx context.Context, filtro model.UsuarioFiltro) ([]model.Usuario, int, error)
}

// UsuarioRepository é uma composição de todas as interfaces acima
type UsuarioRepository interface {
	BuscarUsuario
	ArmazenarUsuario
	AtualizarUsuario
	ListarUsuario
}
