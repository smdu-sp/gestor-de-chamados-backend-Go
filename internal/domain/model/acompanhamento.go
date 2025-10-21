package model

import (
	"errors"
	"fmt"
	"time"
)

// Erros de validação específicos para o modelo Acompanhamento.
var (
	ErrAcompanhamentoIDInvalido = errors.New("ID do acompanhamento não pode ser vazio")
	ErrConteudoInvalido         = errors.New("conteúdo do acompanhamento não pode ser vazio")
	ErrRemetenteInvalido        = errors.New("O remetente deve ser uma das seguintes permissões: TEC, USR")
)

// rmetentesValidos contém todas as permissões aceitas como remetentes de acompanhamentos.
var remetentesValidos = map[Permissao]struct{}{
	PermTEC: {},
	PermUSR: {},
}

// Acompanhamento representa um comentário ou atualização feita em um chamado.
type Acompanhamento struct {
	ID           string    `json:"id"`
	Conteudo     string    `json:"conteudo"`
	ChamadoID    string    `json:"chamadoId"`
	UsuarioID    string    `json:"usuarioId"`
	Remetente    Permissao `json:"remetente"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
}

// NewAcompanhamento cria uma nova instância de Acompanhamento com os dados fornecidos.
func NewAcompanhamento(id, chamadoID, usuarioID, conteudo string, remetente Permissao) (*Acompanhamento, error) {
	now := time.Now()
	acompanhamento := &Acompanhamento{
		ID:           id,
		Conteudo:     conteudo,
		ChamadoID:    chamadoID,
		UsuarioID:    usuarioID,
		Remetente:    remetente,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
	if err := ValidarAcompanhamento(acompanhamento); err != nil {
		return nil, fmt.Errorf("[model.NewAcompanhamento]: %w", err)
	}
	return acompanhamento, nil
}

// ValidarAcompanhamento valida os campos do acompanhamento.
func ValidarAcompanhamento(a *Acompanhamento) error {
	var erros []error

	if a.ID == "" {
		erros = append(erros, ErrAcompanhamentoIDInvalido)
	}
	if a.ChamadoID == "" {
		erros = append(erros, ErrChamadoIDInvalido)
	}
	if a.UsuarioID == "" {
		erros = append(erros, ErrUsuarioIDInvalido)
	}
	if a.Conteudo == "" {
		erros = append(erros, ErrConteudoInvalido)
	}
	if err := ValidarRemetente(a.Remetente); err != nil {
		erros = append(erros, err)
	}
	if len(erros) > 0 {
		return fmt.Errorf("[model.ValidarAcompanhamento] erros de validação: %v", erros)
	}
	return nil
}

// ValidarRemetente valida se o remetente é válido.
func ValidarRemetente(remetente Permissao) error {
	if _, ok := remetentesValidos[remetente]; !ok {
		return fmt.Errorf("[model.ValidarRemetente] %w", ErrRemetenteInvalido)
	}
	return nil
}

// AcompanhamentoFiltro representa os filtros para listar acompanhamentos.
type AcompanhamentoFiltro struct {
	Pagina    int
	Limite    int
	ChamadoID *string
	UsuarioID *string
}

// String retorna uma representação em string do acompanhamento para fins de logging.
func (a *Acompanhamento) String() string {
	return fmt.Sprintf(
		"[ID=%s | ChamadoID=%s | UsuarioID=%s | Conteudo=%s | Remetente=%s]",
		a.ID, a.ChamadoID, a.UsuarioID, a.Conteudo, a.Remetente,
	)
}
