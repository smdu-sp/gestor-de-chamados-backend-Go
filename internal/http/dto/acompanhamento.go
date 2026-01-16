package dto

import (
	"time"

	acp "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/acompanhamento"
)

// AcompanhamentoPaginado é um alias para uma resposta paginada de acompanhamentos,
// usado para facilitar a leitura do código e documentação Swagger.
// @Description Estrutura paginada contendo uma lista de acompanhamentos.
type AcompanhamentoPaginado RespPaginada[AcompanhamentoResp]

// CriarAcompanhamentoReq representa o payload para criar um acompanhamento.
type CriarAcompanhamentoReq struct {
	Conteudo  string `json:"conteudo"`
	ChamadoID string `json:"chamadoId"`
}

// ParaCriarParams converte a requisição para os parâmetros de criação de acompanhamento.
func (r CriarAcompanhamentoReq) ParaCriarParams() acp.CriarParams {
	return acp.CriarParams{
		Conteudo:  r.Conteudo,
		ChamadoID: r.ChamadoID,
	}
}

// AtualizaAcompanhamentoReq representa o payload para atualizar um acompanhamento.
type AtualizarAcompanhamentoReq struct {
	Conteudo  string `json:"conteudo"`
}

// ParaAtualizarParams converte a requisição para os parâmetros de atualização de acompanhamento.
func (r AtualizarAcompanhamentoReq) ParaAtualizarParams() acp.AtualizarParams {
	return acp.AtualizarParams{
		Conteudo:  r.Conteudo,
	}
}

// AcompanhamentoResp representa a estrutura de resposta para um acompanhamento.
type AcompanhamentoResp struct {
	ID           string `json:"id"`
	Conteudo     string `json:"conteudo"`
	ChamadoID    string `json:"chamadoId"`
	UsuarioID    string `json:"usuarioId"`
	Remetente    string `json:"remetente"`
	CriadoEm     string `json:"criadoEm"`
	AtualizadoEm string `json:"atualizadoEm"`
}

// ParaAcompanhamentoResp converte o modelo de acompanhamento para a estrutura de resposta.
func ParaAcompanhamentoResp(a acp.Acompanhamento) AcompanhamentoResp {
	return AcompanhamentoResp{
		ID:           a.ID(),
		Conteudo:     a.Conteudo(),
		ChamadoID:    a.ChamadoID(),
		UsuarioID:    a.UsuarioID(),
		Remetente:    a.Remetente().String(),
		CriadoEm:     a.CriadoEm().Format(time.RFC3339),
		AtualizadoEm: a.AtualizadoEm().Format(time.RFC3339),
	}
}

// ParaAcompanhamentosResp converte um slice de acompanhamentos para um slice de respostas.
func ParaAcompanhamentosResp(acompanhamentos []acp.Acompanhamento) []AcompanhamentoResp {
	return MapearSlice(acompanhamentos, ParaAcompanhamentoResp)
}
