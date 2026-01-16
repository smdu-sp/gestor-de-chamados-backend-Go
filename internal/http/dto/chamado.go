package dto

import (
	"time"

	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
)

// ChamadoPaginado é um alias para uma resposta paginada de chamados,
// usado para facilitar a leitura do código e documentação Swagger.
// @Description Estrutura paginada contendo uma lista de chamados.
type ChamadoPaginado RespPaginada[ChamadoResp]

// CriarChamadoReq representa o payload para criar um chamado.
type CriarChamadoReq struct {
	Titulo         string `json:"titulo"`
	Descricao      string `json:"descricao"`
	CategoriaID    string `json:"categoriaId"`
	SubcategoriaID string `json:"subcategoriaId"`
	CriadorID      string `json:"criadorId"`
}

// ParaCriarParams converte o CriarChamadoReq para CriarParams.
func (c *CriarChamadoReq) ParaCriarParams() chm.CriarParams {
	return chm.CriarParams{
		Titulo:         c.Titulo,
		Descricao:      c.Descricao,
		CategoriaID:    c.CategoriaID,
		SubcategoriaID: c.SubcategoriaID,
		CriadorID:      c.CriadorID,
	}
}

// AtualizarChamadoReq representa o payload para atualizar um chamado.
type AtualizarChamadoReq struct {
	Titulo         *string `json:"titulo,omitempty"`
	Descricao      *string `json:"descricao,omitempty"`
	Arquivado      *bool   `json:"arquivado,omitempty"`
	CategoriaID    *string `json:"categoriaId,omitempty"`
	SubcategoriaID *string `json:"subcategoriaId,omitempty"`
}

// ParaAtualizarParams converte o AtualizarChamadoReq para AtualizarParams.
func (a *AtualizarChamadoReq) ParaAtualizarParams() chm.AtualizarParams {
	return chm.AtualizarParams{
		Titulo:         a.Titulo,
		Descricao:      a.Descricao,
		Arquivado:      a.Arquivado,
		CategoriaID:    a.CategoriaID,
		SubcategoriaID: a.SubcategoriaID,
	}
}

// AtualizarStatusReq representa o payload para atualizar o status de um chamado.
type AtualizarStatusReq struct {
	Status  string  `json:"status"`
	Solucao *string `json:"solucao,omitempty"`
}

// ParaAtualizarStatus converte o AtualizarStatusReq para AtualizarStatusParams.
func (a *AtualizarStatusReq) ParaAtualizarStatusParams() chm.AtualizarStatusParams {
	return chm.AtualizarStatusParams{
		Status:  chm.StatusChamado(a.Status),
		Solucao: a.Solucao,
	}
}

// AtualizarSolucaoReq representa o payload para atualizar a solução de um chamado.
type AtualizarSolucaoReq struct {
	Solucao string `json:"solucao"`
}

// ParaAtualizarSolucaoParams converte o AtualizarSolucaoReq para AtualizarSolucaoParams.
func (a *AtualizarSolucaoReq) ParaAtualizarSolucaoParams() chm.AtualizarSolucaoParams {
	return chm.AtualizarSolucaoParams{
		Solucao: a.Solucao,
	}
}

// ChamadoResp representa a resposta de um chamado.
type ChamadoResp struct {
	ID             string  `json:"id"`
	CategoriaID    string  `json:"categoriaId"`
	SubcategoriaID string  `json:"subcategoriaId"`
	CriadorID      string  `json:"criadorId"`
	Titulo         string  `json:"titulo"`
	Descricao      string  `json:"descricao"`
	Status         string  `json:"status"`
	Arquivado      bool    `json:"arquivado"`
	Solucao        *string `json:"solucao,omitempty"`
	SolucionadoEm  *string `json:"solucionadoEm,omitempty"`
	CriadoEm       string  `json:"criadoEm"`
	AtualizadoEm   string  `json:"atualizadoEm"`
	FechadoEm      *string `json:"fechadoEm,omitempty"`
}

// ParaChamadoResp converte o Chamado para ChamadoResp.
func ParaChamadoResp(c chm.Chamado) ChamadoResp {
	var solucionadoEm *string
	if c.SolucionadoEm() != nil {
		se := c.SolucionadoEm().Format(time.RFC3339)
		solucionadoEm = &se
	}

	var fechadoEm *string
	if c.FechadoEm() != nil {
		fe := c.FechadoEm().Format(time.RFC3339)
		fechadoEm = &fe
	}

	return ChamadoResp{
		ID:             c.ID(),
		CategoriaID:    c.CategoriaID(),
		SubcategoriaID: c.SubcategoriaID(),
		CriadorID:      c.CriadorID(),
		Titulo:         c.Titulo(),
		Descricao:      c.Descricao(),
		Status:         c.Status().String(),
		Arquivado:      c.Arquivado(),
		Solucao:        c.Solucao(),
		SolucionadoEm:  solucionadoEm,
		CriadoEm:       c.CriadoEm().Format(time.RFC3339),
		AtualizadoEm:   c.AtualizadoEm().Format(time.RFC3339),
		FechadoEm:      fechadoEm,
	}
}

// ParaChamadosResp converte uma lista de Chamado para uma lista de ChamadoResp.
func ParaChamadosResp(chamados []chm.Chamado) []ChamadoResp {
	return MapearSlice(chamados, ParaChamadoResp)
}
