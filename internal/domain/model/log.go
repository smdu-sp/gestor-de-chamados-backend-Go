package model

import (
	"errors"
	"fmt"
	"time"
)

// Erros relacionados ao modelo de Log.
var (
	ErrUsuarioIDInvalido = errors.New("usuario_id não pode ser vazio")
	ErrEntidadeInvalida  = errors.New("entidade não pode ser vazia")
	ErrAcaoInvalida      = errors.New("ação inválida: a ação deve ser uma das seguintes: CRIACAO, ATUALIZACAO, REMOCAO, ATIVACAO, DESATIVACAO")
)

// Acao define os tipos de ações que podem ser registradas nos logs.
type Acao string

const (
	AcaoCriar       Acao = "CRIAR"
	AcaoAtualizar   Acao = "ATUALIZAR"
	AcaoDesativar   Acao = "DESATIVAR"
	AcaoAtivar      Acao = "ATIVAR"
	AcaoArquivar    Acao = "ARQUIVAR"
	AcaoDesarquivar Acao = "DESARQUIVAR"
	AcaoDeletar     Acao = "DELETAR"
)

var acoesValidas = map[Acao]struct{}{
	AcaoCriar:       {},
	AcaoAtualizar:   {},
	AcaoAtivar:      {},
	AcaoDesativar:   {},
	AcaoArquivar:    {},
	AcaoDesarquivar: {},
}

// Log representa uma entrada de log no sistema.
type Log struct {
	ID        string    `json:"id"`
	UsuarioID string    `json:"usuario_id"`
	Acao      Acao      `json:"acao"`
	Entidade  string    `json:"entidade"`
	Detalhes  string    `json:"detalhes,omitempty"`
	CriadoEm  time.Time `json:"criado_em"`
}

// NewLog cria uma nova instância de Log com os dados fornecidos.
func NewLog(id, usuarioID string, acao Acao, entidade, detalhes string) (*Log, error) {
	now := time.Now()
	log := &Log{
		ID:        id,
		UsuarioID: usuarioID,
		Acao:      acao,
		Entidade:  entidade,
		Detalhes:  detalhes,
		CriadoEm:  now,
	}
	if err := ValidarLog(log); err != nil {
		return nil, err
	}
	return log, nil
}

// ValidarLog valida os campos do log.
func ValidarLog(l *Log) error {
	var erros []error

	if l.UsuarioID == "" {
		erros = append(erros, ErrUsuarioIDInvalido)
	}
	if l.Entidade == "" {
		erros = append(erros, ErrEntidadeInvalida)
	}
	if err := ValidarAcao(l.Acao); err != nil {
		erros = append(erros, err)
	}
	if len(erros) > 0 {
		return fmt.Errorf("[model.ValidarLog] erros de validação: %v", erros)
	}
	return nil
}

// ValidarAcao valida se a ação é uma das ações permitidas.
func ValidarAcao(acao Acao) error {
	if _, ok := acoesValidas[acao]; ok {
		return nil
	}
	return fmt.Errorf("[model.ValidarAcao]: %w", ErrAcaoInvalida)
}

// LogFiltro representa os critérios de filtragem para listar logs.
type LogFiltro struct {
	Pagina     int
	Limite     int
	Busca      *string
	UsuarioID  *string
	Acao       *string
	Entidade   *string
	DataInicio *time.Time
	DataFim    *time.Time
}
