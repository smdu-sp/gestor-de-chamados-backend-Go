package repository

import (
		"context"
		"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// UserGetter define métodos de busca do usuário
type UserGetter interface {
    FindByID(ctx context.Context, id string) (*model.Usuario, error)
    FindByLogin(ctx context.Context, login string) (*model.Usuario, error)
}

// UserSaver define métodos para salvar/atualizar/excluir usuários
type UserSaver interface {
    Insert(ctx context.Context, u *model.Usuario) error
    Update(ctx context.Context, id string, u *model.Usuario) error
    Delete(ctx context.Context, id string) error
}

// UserUpdater define métodos específicos de atualização
type UserUpdater interface {
    UpdateLastLogin(ctx context.Context, id string) error
}

// UserLister define métodos para listagem e busca filtrada
type UserLister interface {
    List(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]model.Usuario, int, error)
}

// UserRepository é uma composição de todas as interfaces acima
type UserRepository interface {
    UserGetter
    UserSaver
    UserUpdater
    UserLister
}