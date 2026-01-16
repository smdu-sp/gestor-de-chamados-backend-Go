package dto

import (
	"time"

	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// CategoriaPermissaoPaginado é um alias para uma resposta paginada de categorias e permissões,
// usado para facilitar a leitura do código e documentação Swagger.
// @Description Estrutura paginada contendo uma lista de categorias e permissões.
type CategoriaPermissaoPaginado RespPaginada[CategoriaPermissaoResp]

// CriarCategoriaPermissaoReq representa o payload para criar uma nova categoria/permissão.
type CriarCategoriaPermissaoReq struct {
	CategoriaID string `json:"categoriaId"`
	UsuarioID   string `json:"usuarioId"`
	Permissao   string `json:"permissao"`
}

// ParaCriarParams converte um modelo CriarCategoriaPermissaoReq para CriarParams.
func (c CriarCategoriaPermissaoReq) ParaCriarParams() cpm.CriarParams {
	return cpm.CriarParams{
		CategoriaID: c.CategoriaID,
		UsuarioID:   c.UsuarioID,
		Permissao:   usr.Permissao(c.Permissao),
	}
}

// AtualizarCategoriaPermissaoReq representa o payload para atualizar uma categoria/permissão existente.
type AtualizarCategoriaPermissaoReq struct {
	Permissao string `json:"permissao"`
}

// ParaAtualizarParams converte um modelo AtualizarCategoriaPermissaoReq para AtualizarParams.
func (a AtualizarCategoriaPermissaoReq) ParaAtualizarParams() cpm.AtualizarParams {
	return cpm.AtualizarParams{
		Permissao: usr.Permissao(a.Permissao),
	}
}

// CategoriaPermissaoResp representa a resposta para operações relacionadas a categoria/permissão.
type CategoriaPermissaoResp struct {
	CategoriaID  string `json:"categoriaId"`
	UsuarioID    string `json:"usuarioId"`
	Permissao    string `json:"permissao"`
	CriadoEm     string `json:"criadoEm"`
	AtualizadoEm string `json:"atualizadoEm"`
}

// ParaCategoriaPermissaoResp converte um modelo CategoriaPermissao para CategoriaPermissaoResp.
func ParaCategoriaPermissaoResp(c cpm.CategoriaPermissao) CategoriaPermissaoResp {
	return CategoriaPermissaoResp{
		CategoriaID:  c.CategoriaID(),
		UsuarioID:    c.UsuarioID(),
		Permissao:    c.Permissao(),
		CriadoEm:     c.CriadoEm().Format(time.RFC3339),
		AtualizadoEm: c.AtualizadoEm().Format(time.RFC3339),
	}
}

// ParaCategoriasPermissoesResp converte um slice de modelos CategoriaPermissao para um slice de CategoriaPermissaoResp.
func ParaCategoriasPermissoesResp(categoriasPermissoes []cpm.CategoriaPermissao) []CategoriaPermissaoResp {
	return MapearSlice(categoriasPermissoes, ParaCategoriaPermissaoResp)
}
