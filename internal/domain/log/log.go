package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

// Log representa uma entrada de log no sistema.
type Log struct {
	id        string
	usuarioID string
	acao      Acao
	entidade  string
	detalhes  string
	criadoEm  time.Time
}

// Novo cria uma nova instância de Log com os dados fornecidos.
func Novo(id, usuarioID string, acao Acao, entidade, detalhes string) (*Log, error) {
	log := Log{
		id:        strings.TrimSpace(id),
		usuarioID: strings.TrimSpace(usuarioID),
		acao:      acao,
		entidade:  strings.TrimSpace(entidade),
		detalhes:  strings.TrimSpace(detalhes),
		criadoEm:  time.Now(),
	}

	if err := log.Validar(); err != nil {
		return nil, err
	}

	return &log, nil
}

// CarregarDoBD recebe uma estrutura LogDB e retorna uma instância de Log.
func CarregarDoBD(l LogDB) *Log {
	return &Log{
		id:        l.ID,
		usuarioID: l.UsuarioID,
		acao:      Acao(l.Acao),
		entidade:  l.Entidade,
		detalhes:  l.Detalhes,
		criadoEm:  l.CriadoEm,
	}
}

// Validar valida os campos do log.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (l *Log) Validar() error {
	erros := domain.NovoErrosValidacao()

	if l.id == "" {
		erros.Add("ID", "o ID não pode estar vazio")
	}

	if l.usuarioID == "" {
		erros.Add("UsuarioID", "o UsuarioID não pode estar vazio")
	}

	if err := ValidarAcao(l.acao); err != nil {
		erros.Add("Acao", err.Error())
	}

	if l.entidade == "" {
		erros.Add("Entidade", "a entidade não pode estar vazia")
	}

	if l.detalhes == "" {
		erros.Add("Detalhes", "os detalhes não podem estar vazios")
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// -- Métodos de acesso aos campos do Log ---

func (l Log) ID() string          { return l.id }
func (l Log) UsuarioID() string   { return l.usuarioID }
func (l Log) Acao() string        { return l.acao.String() }
func (l Log) Entidade() string    { return l.entidade }
func (l Log) Detalhes() string    { return l.detalhes }
func (l Log) CriadoEm() time.Time { return l.criadoEm }

// String retorna uma representação em string do log para fins de logging.
func (l *Log) String() string {
	return fmt.Sprintf(
		"[ID=%s | UsuarioID=%s | Acao=%s | Entidade=%s | Detalhes=%s | CriadoEm=%s]",
		l.id,
		l.usuarioID,
		l.acao.String(),
		l.entidade,
		l.detalhes,
		l.CriadoEm().Format(time.RFC3339),
	)
}
