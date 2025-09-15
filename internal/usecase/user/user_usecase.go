package user

import (
	"context"
	"strings"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

type Usercase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repository repository.UserRepository) *Usercase {
	return &Usercase{repository: repository}
}

func (u *Usercase) BuscarPorID(ctx context.Context, id string) (*model.Usuario, error) {
	usuario, err := u.repository.FindByID(ctx, id)
	if usuario == nil && err == nil {
		return nil, model.ErrUsuarioNaoEncontrado
	}
	return usuario, err
}

func (u *Usercase) BuscarPorLogin(ctx context.Context, login string) (*model.Usuario, error) {
	usuario, err := u.repository.FindByLogin(ctx, login)
	if usuario == nil && err == nil {
		return nil, model.ErrUsuarioNaoEncontrado
	}
	return usuario, err
}

func (u *Usercase) Criar(ctx context.Context, user *model.Usuario) error {
	if user.Nome == "" || user.Login == "" || user.Permissao == "" {
		return model.ErrCamposObrigatorios
	}

	id, err := util.NewUUIDv7String()
	if err != nil {
		return err
	}
	
	user.ID = id

	if user.Permissao == "" {
		user.Permissao = model.PermUSR
	}

	user.Status = true

	return u.repository.Insert(ctx, user)
}

func (u *Usercase) AtualizarUltimoLogin(ctx context.Context, id string) error {
	return u.repository.UpdateLastLogin(ctx, id)
}

func (u *Usercase) Atualizar(ctx context.Context, id string, user *model.Usuario) error {
	return u.repository.Update(ctx, id, user)
}

func (u *Usercase) Deletar(ctx context.Context, id string) error {
	return u.repository.Delete(ctx, id)
}

func (u *Usercase) Listar(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]model.Usuario, int, error) {
	return u.repository.List(ctx, pagina, limite, busca, status, permissao)
}

func (u *Usercase) Permitido(ctx context.Context, id string, permissoes ...string) (bool, error) {
	usuario, err := u.BuscarPorID(ctx, id)
	if err != nil {
		return false, err
	}

	if usuario == nil {
		return false, model.ErrUsuarioNaoEncontrado
	}

	for _, permissao := range permissoes {
		if strings.EqualFold(string(usuario.Permissao), permissao) {
			return true, nil
		}
	}

	return false, nil 
}