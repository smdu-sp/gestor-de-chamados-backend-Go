package model

import (
	"errors"
	"fmt"
	"time"
)

// Erros de validação específicos para o modelo Atendimento.
var (
	ErrAtribuidoIDInvalido = errors.New("ID do técnico atribuído não pode ser vazio")
	ErrChamadoIDInvalido   = errors.New("ID do chamado não pode ser vazio")
	ErrAtendimentoIDInvalido = errors.New("ID do atendimento não pode ser vazio")
)

// Atendimento representa a atribuição de um técnico a um chamado.
type Atendimento struct {
	ID           string    `json:"id"`
	AtribuidoID  string    `json:"atribuidoId"`
	ChamadoID    string    `json:"chamadoId"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
}

// NewAtendimento cria uma nova instância de Atendimento com os dados fornecidos.
func NewAtendimento(id, atribuidoID, chamadoID string) (*Atendimento, error) {
	now := time.Now()
	atendimento := &Atendimento{
		ID:           id,
		AtribuidoID:  atribuidoID,
		ChamadoID:    chamadoID,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
	if err := ValidarAtendimento(atendimento); err != nil {
		return nil, fmt.Errorf("[model.NewAtendimento]: %w", err)
	}
	return atendimento, nil
}

// ValidarAtendimento valida os campos do atendimento.
func ValidarAtendimento(a *Atendimento) error {
	var erros []error

	if a.AtribuidoID == "" {
		erros = append(erros, ErrAtribuidoIDInvalido)
	}
	if a.ChamadoID == "" {
		erros = append(erros, ErrChamadoIDInvalido)
	}
	if len(erros) > 0 {
		return fmt.Errorf("[model.ValidarAtendimento] erros de validação: %v", erros)
	}
	return nil
}

// AtendimentoFiltro representa os filtros para listar atendimentos.
type AtendimentoFiltro struct {
	Pagina      int
	Limite      int
	ChamadoID   *string
	AtribuidoID *string
}

// String retorna uma representação em string do atendimento para fins de logging.
func (a *Atendimento) String() string {
	return fmt.Sprintf(
		"[ID=%s | AtribuidoID=%s | ChamadoID=%s]",
		a.ID, a.AtribuidoID, a.ChamadoID,
	)
}
