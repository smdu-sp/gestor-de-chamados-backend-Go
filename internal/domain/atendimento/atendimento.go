package atendimento

import (
	"fmt"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

// Atendimento representa a atribuição de um técnico a um chamado.
type Atendimento struct {
	id           string
	atribuidoID  string
	chamadoID    string
	criadoEm     time.Time
	atualizadoEm time.Time
}

// NewAtendimento cria uma nova instância de Atendimento com os dados fornecidos.
func Novo(id, atribuidoID, chamadoID string) (*Atendimento, error) {
	now := time.Now()
	atendimento := Atendimento{
		id:           strings.TrimSpace(id),
		atribuidoID:  strings.TrimSpace(atribuidoID),
		chamadoID:    strings.TrimSpace(chamadoID),
		criadoEm:     now,
		atualizadoEm: now,
	}

	if err := atendimento.Validar(); err != nil {
		return nil, err
	}

	return &atendimento, nil
}

// CarregarDoBD recebe uma estrutura AtendimentoDB e retorna uma instância de Atendimento.
func CarregarDoBD(a AtendimentoDB) *Atendimento {
	return &Atendimento{
		id:           a.ID,
		atribuidoID:  a.AtribuidoID,
		chamadoID:    a.ChamadoID,
		criadoEm:     a.CriadoEm,
		atualizadoEm: a.AtualizadoEm,
	}
}

// Validar valida os campos do atendimento.
//
// Em caso de erros de validação, retorna um erro do tipo domain.ErrosValidacao.
func (a *Atendimento) Validar() error {
	erros := domain.NovoErrosValidacao()
	
	if a.id == "" {
		erros.Add("ID", "o ID do atendimento não pode estar vazio")
	}

	if a.atribuidoID == "" {
		erros.Add("AtribuidoID", "o ID do técnico atribuído não pode estar vazio")
	}

	if a.chamadoID == "" {
		erros.Add("ChamadoID", "o ID do chamado não pode estar vazio")
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// AtualizarDados recebe os novos dados para atualizar o atendimento.
//
// Em caso de erros de validação, retorna um erro do tipo domain.ErrosValidacao.
func (a *Atendimento) AtualizarDados(params AtualizarParams) error {
	if params.AtribuidoID != "" {
		a.atribuidoID = params.AtribuidoID
	}
	if params.ChamadoID != "" {
		a.chamadoID = params.ChamadoID
	}
	a.atualizadoEm = time.Now()
	return a.Validar()
}

// Metodos de acesso aos campos do atendimento.

func (a Atendimento) ID() string              { return a.id }
func (a Atendimento) AtribuidoID() string     { return a.atribuidoID }
func (a Atendimento) ChamadoID() string       { return a.chamadoID }
func (a Atendimento) CriadoEm() time.Time     { return a.criadoEm }
func (a Atendimento) AtualizadoEm() time.Time { return a.atualizadoEm }

// String retorna uma representação em string do atendimento.
func (a *Atendimento) String() string {
	return fmt.Sprintf(
		"[ID=%s | AtribuidoID=%s | ChamadoID=%s | CriadoEm=%s | AtualizadoEm=%s]",
		a.id,
		a.atribuidoID,
		a.chamadoID,
		a.CriadoEm().Format(time.RFC3339),
		a.AtualizadoEm().Format(time.RFC3339),
	)
}
