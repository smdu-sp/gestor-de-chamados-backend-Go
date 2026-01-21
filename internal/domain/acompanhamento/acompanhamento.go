package acompanhamento

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

const tamanhoMinimoConteudo = 1
const tamanhoMaximoConteudo = 1000

// remetentesValidos contém todas as permissões aceitas como remetentes de acompanhamentos.
var remetentesValidos = map[usr.Permissao]struct{}{
	usr.PermTEC: {},
	usr.PermUSR: {},
	usr.PermADM: {},
	usr.PermDEV: {},
}

// Acompanhamento representa um comentário ou atualização feita em um chamado.
type Acompanhamento struct {
	id           string
	chamadoID    string
	usuarioID    string
	conteudo     string
	remetente    usr.Permissao
	criadoEm     time.Time
	atualizadoEm time.Time
}

// Novo cria uma nova instância de Acompanhamento com os dados fornecidos.
func Novo(id, chamadoID, usuarioID, conteudo string, remetente usr.Permissao) (*Acompanhamento, error) {
	now := time.Now()
	acompanhamento := Acompanhamento{
		id:           strings.TrimSpace(id),
		chamadoID:    strings.TrimSpace(chamadoID),
		usuarioID:    strings.TrimSpace(usuarioID),
		conteudo:     strings.TrimSpace(conteudo),
		remetente:    remetente,
		criadoEm:     now,
		atualizadoEm: now,
	}

	if err := acompanhamento.Validar(); err != nil {
		return nil, err
	}

	return &acompanhamento, nil
}

// CarregarDoDB recebe uma estrutura AcompanhamentoDB e retorna uma instância de Acompanhamento.
func CarregarDoDB(a AcompanhamentoDB) *Acompanhamento {
	return &Acompanhamento{
		id:           a.ID,
		chamadoID:    a.ChamadoID,
		usuarioID:    a.UsuarioID,
		conteudo:     a.Conteudo,
		remetente:    usr.Permissao(a.Remetente),
		criadoEm:     a.CriadoEm,
		atualizadoEm: a.AtualizadoEm,
	}
}

// Validar valida os campos do acompanhamento.
//
// Em caso de erros de validação, retorna um erro do tipo domain.ErrosValidacao.
func (a *Acompanhamento) Validar() error {
	erros := domain.NovoErrosValidacao()

	if a.id == "" {
		erros.Add("ID", "o ID do acompanhamento não pode estar vazio")
	}

	if len(a.conteudo) < tamanhoMinimoConteudo || len(a.conteudo) > tamanhoMaximoConteudo {
		erros.Add("Conteudo", fmt.Sprintf("o conteúdo do acompanhamento deve ter entre %d e %d caracteres",
			tamanhoMinimoConteudo, tamanhoMaximoConteudo))
	}

	if a.chamadoID == "" {
		erros.Add("ChamadoID", "o ID do chamado não pode estar vazio")
	}

	if a.usuarioID == "" {
		erros.Add("UsuarioID", "o ID do usuário não pode estar vazio")
	}

	if err := ValidarRemetente(a.remetente); err != nil {
		erros.Add("Remetente", err.Error())
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// ValidarRemetente valida se o remetente é válido.
func ValidarRemetente(remetente usr.Permissao) error {
	if _, ok := remetentesValidos[remetente]; !ok {
		return errors.New("remetente inválido")
	}
	return nil
}

// AtualizarDados recebe os novos dados do acompanhamento e atualiza os campos correspondentes.
//
// Em caso de erros de validação, retorna um erro do tipo domain.ErrosValidacao.
func (a *Acompanhamento) AtualizarDados(params AtualizarParams) error {
	a.conteudo = params.Conteudo
	a.atualizadoEm = time.Now()
	return a.Validar()
}

// Metodos de acesso aos campos do acompanhamento.

func (a Acompanhamento) ID() string               { return a.id }
func (a Acompanhamento) ChamadoID() string        { return a.chamadoID }
func (a Acompanhamento) UsuarioID() string        { return a.usuarioID }
func (a Acompanhamento) Conteudo() string         { return a.conteudo }
func (a Acompanhamento) Remetente() usr.Permissao { return a.remetente }
func (a Acompanhamento) CriadoEm() time.Time      { return a.criadoEm }
func (a *Acompanhamento) AtualizadoEm() time.Time { return a.atualizadoEm }

// String retorna uma representação em string do acompanhamento.
func (a *Acompanhamento) String() string {
	return fmt.Sprintf(
		"[ID=%s | ChamadoID=%s | UsuarioID=%s | Remetente=%s | CriadoEm=%s | AtualizadoEm=%s]",
		a.id,
		a.chamadoID,
		a.usuarioID,
		a.remetente,
		a.CriadoEm().Format(time.RFC3339),
		a.AtualizadoEm().Format(time.RFC3339),
	)
}
