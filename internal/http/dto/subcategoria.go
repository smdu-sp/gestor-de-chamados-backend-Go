package dto

import (
	"time"

	subc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
)
// SubcategoriaPaginado é um alias para uma resposta paginada de subcategorias,
// usado para facilitar a leitura do código e documentação Swagger.
// @Description Estrutura paginada contendo uma lista de subcategorias.
type SubcategoriaPaginado RespPaginada[SubcategoriaResp]

// CriarSubcategoriaReq representa o payload para criar uma subcategoria.
type CriarSubcategoriaReq struct {
	Nome        string `json:"nome"`
	CategoriaID string `json:"categoriaid"`
}

// ParaCriarParams converte um modelo CriarSubcategoriaReq para CriarParams.
func (c CriarSubcategoriaReq) ParaCriarParams() subc.CriarParams {
	return subc.CriarParams{
		Nome:        c.Nome,
		CategoriaID: c.CategoriaID,
	}
}

// AtualizarSubcategoriaReq representa o payload para atualizar uma subcategoria.
type AtualizarSubcategoriaReq struct {
	Nome        *string `json:"nome,omitempty"`
	CategoriaID *string `json:"categoriaid,omitempty"`
	Status      *bool   `json:"status,omitempty"`
}

// ParaAtualizarParams converte um modelo AtualizarSubcategoriaReq para AtualizarParams.
func (a AtualizarSubcategoriaReq) ParaAtualizarParams() subc.AtualizarParams {
	return subc.AtualizarParams{
		Nome:        a.Nome,
		CategoriaID: a.CategoriaID,
		Status:      a.Status,
	}
}

// SubcategoriaResp representa a estrutura de resposta para dados de subcategoria.
type SubcategoriaResp struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	CategoriaID  string `json:"categoriaid"`
	Status       bool   `json:"status"`
	CriadoEm     string `json:"criadoEm"`
	AtualizadoEm string `json:"atualizadoEm"`
}

// paraSubcategoriaResp converte um modelo Subcategoria para SubcategoriaResp.
func ParaSubcategoriaResp(s subc.Subcategoria) SubcategoriaResp {
	return SubcategoriaResp{
		ID:           s.ID(),
		Nome:         s.Nome(),
		CategoriaID:  s.CategoriaID(),
		Status:       s.Status(),
		CriadoEm:     s.CriadoEm().Format(time.RFC3339),
		AtualizadoEm: s.AtualizadoEm().Format(time.RFC3339),
	}
}

// ParaSubcategoriasResp converte uma lista de modelos Subcategoria para uma lista de SubcategoriaResp.
func ParaSubcategoriasResp(subcategorias []subc.Subcategoria) []SubcategoriaResp {
	return MapearSlice(subcategorias, ParaSubcategoriaResp)
}
