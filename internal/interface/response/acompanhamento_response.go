package response

import (
	"time"
	
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// AcompanhamentoResponse representa a estrutura de resposta para dados de acompanhamento
type AcompanhamentoResponse struct {
	ID              string    `json:"id"`
	Conteudo        string    `json:"conteudo"`
	ChamadoID       string    `json:"chamadoId"`
	UsuarioID       string    `json:"usuarioId"`
	Remetente       model.Permissao `json:"remetente"`
	CriadoEm        time.Time `json:"criadoEm"`
	AtualizadoEm    time.Time `json:"atualizadoEm"`
}

// ToAcompanhamentoResponse converte um modelo Acompanhamento para AcompanhamentoResponse
func ToAcompanhamentoResponse(a *model.Acompanhamento) *AcompanhamentoResponse {
	return &AcompanhamentoResponse{
		ID:           a.ID,
		Conteudo:     a.Conteudo,
		ChamadoID:    a.ChamadoID,
		UsuarioID:    a.UsuarioID,
		Remetente:    a.Remetente,
		CriadoEm:     a.CriadoEm,
		AtualizadoEm: a.AtualizadoEm,
	}
}