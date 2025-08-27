package user

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrCamposObrigatorios   = errors.New("campos obrigatórios ausentes")
	ErrUsuarioNaoEncontrado = errors.New("usuário não encontrado")
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// BuscarPorID retorna um usuário pelo ID.
func (s *Service) BuscarPorID(ctx context.Context, id string) (*Usuario, error) {
	u, err := s.repo.FindByID(ctx, id)
	if u == nil && err == nil {
		return nil, ErrUsuarioNaoEncontrado
	}
	return u, err
}

// BuscarPorLogin retorna um usuário pelo login.
func (s *Service) BuscarPorLogin(ctx context.Context, login string) (*Usuario, error) {
	u, err := s.repo.FindByLogin(ctx, login)
	if u == nil && err == nil {
		return nil, ErrUsuarioNaoEncontrado
	}
	return u, err
}

// BuscarPorEmail retorna um usuário pelo email.
func (s *Service) Criar(ctx context.Context, u *Usuario) error {
	if u.Nome == "" || u.Login == "" || u.Email == "" {
		return ErrCamposObrigatorios
	}
	if u.Permissao == "" {
		u.Permissao = PermUSR
	}
	u.Status = true
	return s.repo.Insert(ctx, u)
}

// Atualizar atualiza os dados de um usuário existente.
func (s *Service) Atualizar(ctx context.Context, id string, u *Usuario) error {
	return s.repo.Update(ctx, id, u)
}

// Remover exclui um usuário pelo ID.
func (s *Service) Remover(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// AtualizarUltimoLogin atualiza a data do último login do usuário.
func (s *Service) AtualizarUltimoLogin(ctx context.Context, id string) error {
	return s.repo.UpdateLastLogin(ctx, id)
}

// Permitido verifica se o usuário possui uma das permissões especificadas.
func (s *Service) Permitido(ctx context.Context, id string, permissoes []string) (bool, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, ErrUsuarioNaoEncontrado
	}

	// Verifica se a permissão do usuário está na lista de permissões permitidas
	for _, p := range permissoes {
		if strings.EqualFold(string(u.Permissao), p) {
			return true, nil
		}
	}

	return false, nil
}

// Listar retorna uma lista de usuários com paginação e filtros opcionais.
func (s *Service) Listar(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]Usuario, int, error) {
	return s.repo.List(ctx, pagina, limite, busca, status, permissao)
}
