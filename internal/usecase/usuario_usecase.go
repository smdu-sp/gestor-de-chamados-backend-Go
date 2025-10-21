package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// UsuarioUsecase representa a camada de caso de uso para operações relacionadas a usuários.
type UsuarioUsecase struct {
	repository repository.UsuarioRepository
}

// NewUsuarioUsecase cria uma nova instância de UsuarioUsecase.
func NewUsuarioUsecase(repository repository.UsuarioRepository) *UsuarioUsecase {
	return &UsuarioUsecase{repository: repository}
}

// BuscarUsuarioPorID busca um usuário pelo seu ID.
func (u *UsuarioUsecase) BuscarUsuarioPorID(ctx context.Context, id string) (*model.Usuario, error) {
	if id == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarUsuarioPorID]",
			utils.LevelInfo,
			"erro ao buscar usuário por id",
			model.ErrIDInvalido,
		)
	}

	usuario, err := u.repository.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarUsuarioPorID]: %w", err)
	}
	return usuario, nil
}

// BuscarUsuarioPorLogin busca um usuário pelo seu login.
func (u *UsuarioUsecase) BuscarUsuarioPorLogin(ctx context.Context, login string) (*model.Usuario, error) {
	if login == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarUsuarioPorLogin]",
			utils.LevelInfo,
			"erro ao buscar usuário por login",
			model.ErrLoginInvalido,
		)
	}

	usuario, err := u.repository.BuscarPorLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarUsuarioPorLogin]: %w", err)
	}
	return usuario, nil
}

// CriarUsuario cria um novo usuário.
func (u *UsuarioUsecase) CriarUsuario(ctx context.Context, usuario *model.Usuario) error {
	const metodo = "[usecase.CriarUsuario]: %w"

	id, err := utils.NewUUIDv7String()
	if err != nil {
		return fmt.Errorf(metodo, err)
	}
	usuario.ID = id

	if usuario.Permissao == "" {
		usuario.Permissao = model.PermUSR
	}

	usuario.Status = true

	usuario, err = model.NewUsuario(
		usuario.ID,
		usuario.Nome,
		usuario.Login,
		usuario.Email,
		usuario.Permissao,
		usuario.Status,
		usuario.Avatar,
	)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	if err = u.repository.Salvar(ctx, usuario); err != nil {
		return fmt.Errorf(metodo, err)
	}
	return nil
}

// AtualizarUltimoLoginUsuario atualiza a data do último login do usuário.
func (u *UsuarioUsecase) AtualizarUltimoLoginUsuario(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtualizarUltimoLoginUsuario]",
			utils.LevelInfo,
			"erro ao atualizar último login do usuário",
			model.ErrIDInvalido,
		)
	}
	if err := u.repository.AtualizarUltimoLogin(ctx, id); err != nil {
		return fmt.Errorf("[usecase.AtualizarUltimoLoginUsuario]: %w", err)
	}
	return nil
}

// AtualizarUsuario atualiza as informações de um usuário.
func (u *UsuarioUsecase) AtualizarUsuario(ctx context.Context, id string, usuario *model.Usuario) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtualizarUsuario]",
			utils.LevelInfo,
			"erro ao atualizar usuário",
			model.ErrIDInvalido,
		)
	}

	if err := model.ValidarUsuario(usuario); err != nil {
		return fmt.Errorf("[usecase.AtualizarUsuario]: %w", err)
	}

	if err := u.repository.Atualizar(ctx, id, usuario); err != nil {
		return fmt.Errorf("[usecase.AtualizarUsuario]: %w", err)
	}
	return nil
}

// AtualizarPermissaoUsuario atualiza a permissão do usuário.
func (u *UsuarioUsecase) AtualizarPermissaoUsuario(ctx context.Context, id string, permissao string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtualizarPermissaoUsuario]",
			utils.LevelInfo,
			"erro ao atualizar permissão do usuário",
			model.ErrIDInvalido,
		)
	}

	if err := model.ValidarPermissao(model.Permissao(permissao)); err != nil {
		return fmt.Errorf("[usecase.AtualizarPermissaoUsuario]: %w", err)
	}

	if err := u.repository.AtualizarPermissao(ctx, id, permissao); err != nil {
		return fmt.Errorf("[usecase.AtualizarPermissaoUsuario]: %w", err)
	}
	return nil
}

// DesativarUsuario desativa um usuário.
func (u *UsuarioUsecase) DesativarUsuario(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.DesativarUsuario]",
			utils.LevelInfo,
			"erro ao desativar usuário",
			model.ErrIDInvalido,
		)
	}

	if err := u.repository.Desativar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.DesativarUsuario]: %w", err)
	}
	return nil
}

// AtivarUsuario ativa um usuário.
func (u *UsuarioUsecase) AtivarUsuario(ctx context.Context, id string) error {
	if id == "" {
		return utils.NewAppError(
			"[usecase.AtivarUsuario]",
			utils.LevelInfo,
			"erro ao ativar usuário",
			model.ErrIDInvalido,
		)
	}

	if err := u.repository.Ativar(ctx, id); err != nil {
		return fmt.Errorf("[usecase.AtivarUsuario]: %w", err)
	}
	return nil
}

// ListarUsuarios lista usuários com paginação e filtros opcionais.
func (u *UsuarioUsecase) ListarUsuarios(ctx context.Context, filtro model.UsuarioFiltro) ([]model.Usuario, int, model.UsuarioFiltro, error) {
	if filtro.Pagina <= 0 {
		filtro.Pagina = 1
	}

	if filtro.Limite <= 0 || filtro.Limite > 100 {
		filtro.Limite = 10
	}

	usuarios, total, err := u.repository.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, filtro, fmt.Errorf("[usecase.ListarUsuarios]: %w", err)
	}

	return usuarios, total, filtro, nil
}

// VerificarPermissao verifica se o usuário possui uma das permissões especificadas.
func (u *UsuarioUsecase) VerificarPermissao(ctx context.Context, id string, permissoes ...string) (bool, error) {
	if id == "" {
		return false, utils.NewAppError(
			"[usecase.VerificarPermissao]",
			utils.LevelInfo,
			"erro ao verificar permissão do usuário",
			model.ErrIDInvalido,
		)
	}

	usuario, err := u.BuscarUsuarioPorID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[usecase.VerificarPermissao]: %w", err)
	}

	for _, permissao := range permissoes {
		if strings.EqualFold(string(usuario.Permissao), permissao) {
			return true, nil
		}
	}

	return false, nil
}
