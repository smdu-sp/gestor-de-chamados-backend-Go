package service

import (
	"context"
	"fmt"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	sub "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
)

// SubcategoriaService representa a camada de caso de uso para operações relacionadas a subcategorias.
type SubcategoriaService struct {
	id   dmn.GeradorID
	repo sub.Repository
}

// NewSubcategoriaService cria uma nova instância de SubcategoriaService.
func NovoSubcategoriaService(id dmn.GeradorID, repo sub.Repository) *SubcategoriaService {
	return &SubcategoriaService{id: id, repo: repo}
}

// Asserção de interface para garantir que SubcategoriaService implementa SubcategoriaService
var _ sub.Service = (*SubcategoriaService)(nil)

// BuscarPorID recebe um ID e retorna a subcategoria correspondente.
//
// Erros possíveis: ErrSubcategoriaNaoEncontrada.
func (s *SubcategoriaService) BuscarPorID(ctx context.Context, id string) (*sub.Subcategoria, error) {
	subc, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("subcategoria por ID: %w", err)
	}
	return subc, nil
}

// BuscarPorNome recebe um nome e retorna a subcategoria correspondente.
//
// Erros possíveis: ErrSubcategoriaNaoEncontrada.
func (s *SubcategoriaService) BuscarPorNome(ctx context.Context, nome string) (*sub.Subcategoria, error) {
	subc, err := s.repo.BuscarPorNome(ctx, nome)
	if err != nil {
		return nil, fmt.Errorf("subcategoria por nome: %w", err)
	}
	return subc, nil
}

// Criar recebe os parâmetros para criar uma nova subcategoria e a insere no repositório.
//
// Erros possíveis: ErrSubcategoriaJaExisteComNome.
func (s *SubcategoriaService) Criar(ctx context.Context, criar sub.CriarParams) (*sub.Subcategoria, error) {
	// 1 - Gerar um novo ID para a subcategoria
	id, err := s.id.NovoID()
	if err != nil {
		return nil, fmt.Errorf("criar subcategoria: %w", err)
	}

	// 2 - Criar a subcategoria
	subc, err := sub.Novo(
		id,
		criar.Nome,
		criar.CategoriaID,
	)
	if err != nil {
		return nil, err
	}

	// 3 - Salvar a subcategoria no repositório
	subcCriada, err := s.repo.Criar(ctx, *subc)
	if err != nil {
		return nil, fmt.Errorf("criar subcategoria: %w", err)
	}
	return subcCriada, nil
}

// Atualizar recebe um ID e uma subcategoria para atualizar os dados da subcategoria correspondente.
//
// Erros possíveis: ErrSubcategoriaNaoEncontrada, ErrSubcategoriaJaExisteComNome, ErrosValidacao.
func (s *SubcategoriaService) Atualizar(ctx context.Context, id string, atualizar sub.AtualizarParams) (*sub.Subcategoria, error) {
	subcategoriaAtual, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar subcategoria por ID: %w", err)
	}

	subcategoriaAtualizada, err := subcategoriaAtual.ComDadosAtualizados(atualizar)
	if err != nil {
		return nil, err
	}

	subcategoriaSalva, err := s.repo.Atualizar(ctx, id, subcategoriaAtualizada)
	if err != nil {
		return nil, fmt.Errorf("atualizar subcategoria: %w", err)
	}
	return subcategoriaSalva, nil
}

// Desativar desativa uma subcategoria existente.
//
// Erros possíveis: ErrSubcategoriaNaoEncontrada.
func (s *SubcategoriaService) Desativar(ctx context.Context, id string) (*sub.Subcategoria, error) {
	subcategoriaAtual, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("desativar subcategoria: %w", err)
	}

	subcategoriaDesativada := subcategoriaAtual.Desativar()

	subcategoriaSalva, err := s.repo.Atualizar(ctx, id, subcategoriaDesativada)
	if err != nil {
		return nil, fmt.Errorf("desativar subcategoria: %w", err)
	}
	return subcategoriaSalva, nil
}

// Ativar recebe um ID e ativa a subcategoria correspondente.
//
// Erros possíveis: ErrSubcategoriaNaoEncontrada.
func (s *SubcategoriaService) Ativar(ctx context.Context, id string) (*sub.Subcategoria, error) {
	subcategoriaAtual, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ativar subcategoria: %w", err)
	}

	subcategoriaAtivada := subcategoriaAtual.Ativar()

	subcategoriaSalva, err := s.repo.Atualizar(ctx, id, subcategoriaAtivada)
	if err != nil {
		return nil, fmt.Errorf("ativar subcategoria: %w", err)
	}
	return subcategoriaSalva, nil
}

// Listar recebe um filtro e retorna uma lista de subcategorias que correspondem aos critérios do filtro,
// juntamente com o total de registros encontrados.
func (s *SubcategoriaService) Listar(ctx context.Context, filtro sub.Filtro) ([]sub.Subcategoria, int, sub.Filtro, error) {
	filtro.Normalizar()
	subcategorias, total, err := s.repo.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, sub.Filtro{}, fmt.Errorf("listar subcategorias: %w", err)
	}

	return subcategorias, total, filtro, nil
}
