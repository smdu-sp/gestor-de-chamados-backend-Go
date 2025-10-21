package response

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// LogResponse representa a estrutura de resposta para logs.
type LogResponse struct {
	ID        string `json:"id"`
	UsuarioID string `json:"usuario_id"`
	Acao      string `json:"acao"`
	Entidade  string `json:"entidade"`
	Detalhes  string `json:"detalhes"`
	CriadoEm  time.Time `json:"criado_em"`
}

// ToLogResponse converte um modelo Log para LogResponse.
func ToLogResponse(log *model.Log) LogResponse {
	return LogResponse{
		ID:        log.ID,
		UsuarioID: log.UsuarioID,
		Acao:      string(log.Acao),
		Entidade:  log.Entidade,
		Detalhes:  log.Detalhes,
		CriadoEm:  log.CriadoEm,
	}
}