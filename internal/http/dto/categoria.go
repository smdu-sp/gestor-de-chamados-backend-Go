package dto

import (
	"time"

	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
)

// CategoriaPaginado é um alias para uma resposta paginada de categorias,
// usado para facilitar a leitura do código e documentação Swagger.
// @Description Estrutura paginada contendo uma lista de categorias.
type CategoriaPaginado RespPaginada[CategoriaResp]

// CriarCategoriaReq representa o payload para criar uma categoria.
type CriarCategoriaReq struct {
	Nome string `json:"nome"`
}

// ParaCriarParams converte um modelo CriarCategoriaReq para CriarParams.
func (c CriarCategoriaReq) ParaCriarParams() ctg.CriarParams {
	return ctg.CriarParams{
		Nome: c.Nome,
	}
}

// AtualizarCategoriaReq representa o payload para atualizar uma categoria.
type AtualizarCategoriaReq struct {
	Nome   *string `json:"nome,omitempty"`
	Status *bool   `json:"status,omitempty"`
}

// ParaAtualizarParams converte um modelo AtualizarCategoriaReq para AtualizarParams.
func (a AtualizarCategoriaReq) ParaAtualizarParams() ctg.AtualizarParams {
	return ctg.AtualizarParams{
		Nome:   a.Nome,
		Status: a.Status,
	}
}

// CategoriaResp representa a estrutura de resposta para dados de categoria.
type CategoriaResp struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	Status       bool   `json:"status"`
	CriadoEm     string `json:"criadoEm"`
	AtualizadoEm string `json:"atualizadoEm"`
}

// ParaCategoriaResp converte um modelo Categoria para CategoriaResposta.
func ParaCategoriaResp(c ctg.Categoria) CategoriaResp {
	return CategoriaResp{
		ID:           c.ID(),
		Nome:         c.Nome(),
		Status:       c.Status(),
		CriadoEm:     c.CriadoEm().Format(time.RFC3339),
		AtualizadoEm: c.AtualizadoEm().Format(time.RFC3339),
	}
}

// ParaCategoriasResp converte uma lista de Categoria para uma lista de CategoriaResp.
func ParaCategoriasResp(categorias []ctg.Categoria) []CategoriaResp {
	return MapearSlice(categorias, ParaCategoriaResp)
}
