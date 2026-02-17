package service

import (
	"context"
	"fmt"
	"slices"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// UsuarioService representa a camada de caso de uso para operações relacionadas a usuários.
type UsuarioService struct {
	id   dmn.GeradorID
	repo usr.Repository
}

// NovoUsuarioService cria uma nova instância de UsuarioService.
func NovoUsuarioService(id dmn.GeradorID, repo usr.Repository) *UsuarioService {
	return &UsuarioService{id: id, repo: repo}
}

// Asserção de interface para garantir que UsuarioService implementa usr.Service
var _ usr.Service = (*UsuarioService)(nil)

// BuscarPorID recebe um ID e retorna o usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado.
func (u *UsuarioService) BuscarPorID(ctx context.Context, id string) (*usr.Usuario, error) {
	usuario, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar usuário por ID: %w", err)
	}
	return usuario, nil
}

// BuscarPorLogin recebe um login e retorna o usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado.
func (u *UsuarioService) BuscarPorLogin(ctx context.Context, login string) (*usr.Usuario, error) {
	usuario, err := u.repo.BuscarPorLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("buscar usuário por login: %w", err)
	}
	return usuario, nil
}

// Criar recebe os parâmetros para criação de um usuário e retorna o usuário criado.
//
// Erros sentinela possíveis: ErrUsuarioJaExisteComEmail, ErrUsuarioJaExisteComLogin, ErrosValidacao.
func (u *UsuarioService) Criar(ctx context.Context, criar usr.CriarParams) (*usr.Usuario, error) {
	// 1 - Gerar um novo ID para o usuário
	id, err := u.id.NovoID()
	if err != nil {
		return nil, fmt.Errorf("criar usuário: %w", err)
	}

	// 2 - Criar o usuário
	usuario, err := usr.Novo(
		id,
		criar.Nome,
		criar.Login,
		usr.NovoEmail(criar.Email.String()),
		usr.Permissao(criar.Permissao),
		criar.Avatar,
	)
	if err != nil {
		return nil, err
	}

	// 3 - Salvar o usuário no repositório
	usuarioCriado, err := u.repo.Criar(ctx, *usuario)
	if err != nil {
		return nil, fmt.Errorf("criar usuário: %w", err)
	}
	return usuarioCriado, nil
}

// AtualizarUltimoLogin recebe um ID e atualiza o último login do usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado.
func (u *UsuarioService) AtualizarUltimoLogin(ctx context.Context, id string) (*usr.Usuario, error) {
	// 1 - Buscar o usuário pelo ID
	usuario, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar último login: %w", err)
	}

	// 2 - Atualizar o último login
	usuario.AtualizarUltimoLogin()

	// 3 - Salvar a atualização no repositório
	usuarioAtualizado, err := u.repo.Atualizar(ctx, id, *usuario)
	if err != nil {
		return nil, fmt.Errorf("atualizar último login: %w", err)
	}
	return usuarioAtualizado, nil
}

// Atualizar recebe um ID e parâmetros para atualizar os dados do usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado,  ErrUsuarioJaExisteComLogin, 
// ErrUsuarioJaExisteComEmail, ErrosValidacao.
func (u *UsuarioService) Atualizar(ctx context.Context, id string, atualizar usr.AtualizarParams) (*usr.Usuario, error) {
	// 1 - Buscar o usuário pelo ID
	usuarioAtual, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar usuário: %w", err)
	}

	// 2 - Atualizar os dados do usuário
	usuarioAtualizado, err := usuarioAtual.ComDadosAtualizados(atualizar)
	if err != nil {
		return nil, err
	}

	// 3 - Salvar a atualização no repositório
	usuarioSalvo, err := u.repo.Atualizar(ctx, id, usuarioAtualizado)
	if err != nil {
		return nil, fmt.Errorf("atualizar usuário: %w", err)
	}

	return usuarioSalvo, nil
}

// AtualizarPermissao recebe um ID e parâmetros para atualizar a permissão do usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado, ErrosValidacao.
func (u *UsuarioService) AtualizarPermissao(ctx context.Context, id string, a usr.AtualizarPermissaoParams) (*usr.Usuario, error) {
	// 1 - Buscar o usuário pelo ID
	usuarioAtual, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar permissão: %w", err)
	}

	// 2 - Atualizar a permissão do usuário
	usuarioAtualizado, err := usuarioAtual.AtualizarPermissao(a.Permissao)
	if err != nil {
		return nil, err
	}

	// 3 - Salvar a atualização no repositório
	usuarioSalvo, err := u.repo.Atualizar(ctx, id, usuarioAtualizado)
	if err != nil {
		return nil, fmt.Errorf("atualizar permissão: %w", err)
	}

	return usuarioSalvo, nil
}

// Desativar recebe um ID e desativa o usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado.
func (u *UsuarioService) Desativar(ctx context.Context, id string) (*usr.Usuario, error) {
	// 1 - Buscar o usuário pelo ID
	usuarioAtual, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("desativar usuário: %w", err)
	}

	// 2 - Desativar o usuário
	usuarioDesativado := usuarioAtual.Desativar()

	// 3 - Salvar a atualização no repositório
	usuarioSalvo, err := u.repo.Atualizar(ctx, id, usuarioDesativado)
	if err != nil {
		return nil, fmt.Errorf("desativar usuário: %w", err)
	}

	return usuarioSalvo, nil
}

// Ativar recebe um ID e ativa o usuário correspondente.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado.
func (u *UsuarioService) Ativar(ctx context.Context, id string) (*usr.Usuario, error) {
	// 1 - Buscar o usuário pelo ID
	usuarioAtual, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ativar usuário: %w", err)
	}

	// 2 - Ativar o usuário
	usuarioAtivado := usuarioAtual.Ativar()

	// 3 - Salvar a atualização no repositório
	usuarioSalvo, err := u.repo.Atualizar(ctx, id, usuarioAtivado)
	if err != nil {
		return nil, fmt.Errorf("ativar usuário: %w", err)
	}
	return usuarioSalvo, nil
}

// Listar recebe um filtro e retorna uma lista de usuários que correspondem aos critérios do filtro,
// juntamente com o total de registros encontrados e o filtro aplicado.
//
// Erros sentinela possíveis: ErrosValidacao.
func (u *UsuarioService) Listar(ctx context.Context, filtro usr.Filtro) ([]usr.Usuario, int, usr.Filtro, error) {
	filtro.Normalizar()
	usuarios, total, err := u.repo.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, usr.Filtro{}, fmt.Errorf("listar usuários: %w", err)
	}

	return usuarios, total, filtro, nil
}

// VerificarPermissao recebe um ID e uma lista de permissões, 
//e verifica se o usuário correspondente possui alguma das permissões fornecidas.
//
// Erros sentinela possíveis: ErrUsuarioNaoEncontrado.
func (u *UsuarioService) VerificarPermissao(ctx context.Context, id string, p []usr.Permissao) (bool, error) {
	usuario, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("verificar permissão: %w", err)
	}

	if slices.Contains(p, usuario.Permissao()) {
		return true, nil
	}

	return false, nil
}
