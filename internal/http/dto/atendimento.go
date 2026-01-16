package dto

import (
	"time"

	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
)

// AtendimentoPaginado é um alias para uma resposta paginada de atendimentos,
// usado para facilitar a leitura do código e documentação Swagger.
// @Description Estrutura paginada contendo uma lista de atendimentos.
type AtendimentoPaginado RespPaginada[AtendimentoResp]

// CriarAtendimentoReq representa o payload para criar um atendimento.
type CriarAtendimentoReq struct {
	AtribuidoID string `json:"atribuidoId"`
	ChamadoID   string `json:"chamadoId"`
}

// ParaCriarParams converte um modelo CriarAtendimentoReq para CriarParams.
func (r CriarAtendimentoReq) ParaCriarParams() atd.CriarParams {
	return atd.CriarParams{
		AtribuidoID: r.AtribuidoID,
		ChamadoID:   r.ChamadoID,
	}
}

// AtualizaAtendimentoReq representa o payload para atualizar um atendimento.
type AtualizarAtendimentoReq struct {
	AtribuidoID string `json:"atribuidoId,omitempty"`
	ChamadoID   string `json:"chamadoId,omitempty"`
}

// ParaAtualizarParams converte um modelo AtualizarAtendimentoReq para AtualizarParams.
func (r AtualizarAtendimentoReq) ParaAtualizarParams() atd.AtualizarParams {
	return atd.AtualizarParams{
		AtribuidoID: r.AtribuidoID,
		ChamadoID:   r.ChamadoID,
	}
}

// AtendimentoResp representa a estrutura de resposta para dados de atendimento.
type AtendimentoResp struct {
	ID           string `json:"id"`
	AtribuidoID  string `json:"atribuidoId"`
	ChamadoID    string `json:"chamadoId"`
	CriadoEm     string `json:"criadoEm"`
	AtualizadoEm string `json:"atualizadoEm"`
}

// ParaAtendimentoResp converte um modelo Atendimento para AtendimentoResposta.
func ParaAtendimentoResp(a atd.Atendimento) AtendimentoResp {
	return AtendimentoResp{
		ID:           a.ID(),
		AtribuidoID:  a.AtribuidoID(),
		ChamadoID:    a.ChamadoID(),
		CriadoEm:     a.CriadoEm().Format(time.RFC3339),
		AtualizadoEm: a.AtualizadoEm().Format(time.RFC3339),
	}
}

// ParaAtendimentosResp converte uma lista de modelos Atendimento para uma lista de AtendimentoResp.
func ParaAtendimentosResp(atendimentos []atd.Atendimento) []AtendimentoResp {
	return MapearSlice(atendimentos, ParaAtendimentoResp)
}
