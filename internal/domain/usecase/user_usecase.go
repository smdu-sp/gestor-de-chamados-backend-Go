package usecase

import (
	"context"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// UserGetter é a interface que define os métodos para obter informações de usuários.
type UserGetter interface {
	BuscarPorID(ctx context.Context, id string) (*model.Usuario, error)
	BuscarPorLogin(ctx context.Context, login string) (*model.Usuario, error)
	Listar(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]model.Usuario, int, error)
	Permitido(ctx context.Context, id string, permissoes ...string) (bool, error)
}

// UserCreator é a interface que define os métodos para criar, atualizar e deletar usuários.
type UserCreator interface {
	Criar(ctx context.Context, u *model.Usuario) error	
	Atualizar(ctx context.Context, id string, u *model.Usuario) error
	Deletar(ctx context.Context, id string) error
}

// UserUpdater é a interface que define os métodos para atualizar informações específicas de usuários.
type UserUpdater interface {
	AtualizarUltimoLogin(ctx context.Context, id string) error
}

// UserUsecase é a interface que agrega os casos de uso relacionados a usuários.
type UserUsecase interface {
	UserGetter
	UserCreator
	UserUpdater
}