package subcategoria

import (
	"fmt"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

const tamanhoMinNome = 3
const tamanhoMaxNome = 100

// Subcategoria representa uma subcategoria de chamados no sistema.
type Subcategoria struct {
	id           string
	nome         string
	status       bool
	categoriaID  string
	criadoEm     time.Time
	atualizadoEm time.Time
}

// Novo cria uma nova instância de Subcategoria com os dados fornecidos.
func Novo(id, nome, categoriaID string) (*Subcategoria, error) {
	now := time.Now()
	subcategoria := Subcategoria{
		id:           strings.TrimSpace(id),
		nome:         strings.TrimSpace(nome),
		status:       true,
		categoriaID:  strings.TrimSpace(categoriaID),
		criadoEm:     now,
		atualizadoEm: now,
	}

	if err := subcategoria.Validar(); err != nil {
		return nil, err
	}

	return &subcategoria, nil
}

// CarregarDoBD cria uma nova instância de Subcategoria a partir dos dados do banco de dados.
func CarregarDoBD(s SubcategoriaDB) *Subcategoria {
	return &Subcategoria{
		id:           s.ID,
		nome:         s.Nome,
		status:       s.Status,
		categoriaID:  s.CategoriaID,
		criadoEm:     s.CriadoEm,
		atualizadoEm: s.AtualizadoEm,
	}
}

// Validar valida os campos da subcategoria.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (s Subcategoria) Validar() error {
	erros := domain.NovoErrosValidacao()

	if s.id == "" {
		erros.Add("ID", "o ID da subcategoria não pode estar vazio")
	}

	if len(s.nome) < tamanhoMinNome || len(s.nome) > tamanhoMaxNome {
		erros.Add("Nome", fmt.Sprintf("o nome da subcategoria deve ter entre %d e %d caracteres",
			tamanhoMinNome, tamanhoMaxNome))
	}

	if s.categoriaID == "" {
		erros.Add("CategoriaID", "o ID da categoria não pode estar vazio")
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// Desativar desativa a subcategoria.
func (s Subcategoria) Desativar() Subcategoria {
	s.status = false
	s.atualizadoEm = time.Now()
	return s
}

// Ativar ativa a subcategoria.
func (s Subcategoria) Ativar() Subcategoria {
	s.status = true
	s.atualizadoEm = time.Now()
	return s
}

// ComDadosAtualizados recebe os novos dados para a subcategoria 
// e retorna uma nova instância com os dados atualizados.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (s Subcategoria) ComDadosAtualizados(params AtualizarParams) (Subcategoria, error) {
	if params.Nome != nil {
		s.nome = *params.Nome
	}

	if params.CategoriaID != nil {
		s.categoriaID = *params.CategoriaID
	}

	if params.Status != nil {
		s.status = *params.Status
	}

	s.atualizadoEm = time.Now()

	if err := s.Validar(); err != nil {
		return Subcategoria{}, err
	}

	return s, nil
}

// --- Métodos de acesso aos campos da subcategoria ---

func (s Subcategoria) ID() string              { return s.id }
func (s Subcategoria) Nome() string            { return s.nome }
func (s Subcategoria) Status() bool            { return s.status }
func (s Subcategoria) CategoriaID() string     { return s.categoriaID }
func (s Subcategoria) CriadoEm() time.Time     { return s.criadoEm }
func (s Subcategoria) AtualizadoEm() time.Time { return s.atualizadoEm }

// String retorna uma representação em string da subcategoria para fins de logging.
func (s Subcategoria) String() string {
	return fmt.Sprintf(
		"[ID=%s | Nome=%s | Status=%t | CategoriaID=%s | CriadoEm=%s | AtualizadoEm=%s]",
		s.id,
		s.nome,
		s.status,
		s.categoriaID,
		s.CriadoEm().Format(time.RFC3339),
		s.AtualizadoEm().Format(time.RFC3339),
	)
}
