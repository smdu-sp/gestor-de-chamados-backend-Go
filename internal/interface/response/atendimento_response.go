package response

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

type AtendimentoResponse struct {
	ID           string    `json:"id"`
	AtribuidoID  string    `json:"atribuidoId"`
	ChamadoID    string    `json:"chamadoId"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
}

// ToAtendimentoResponse converte um modelo Atendimento para AtendimentoResponse
func ToAtendimentoResponse(a *model.Atendimento) *AtendimentoResponse {
	return &AtendimentoResponse{
		ID:           a.ID,
		AtribuidoID:  a.AtribuidoID,
		ChamadoID:    a.ChamadoID,
		CriadoEm:     a.CriadoEm,
		AtualizadoEm: a.AtualizadoEm,
	}
}
