package categoria

import (
	"fmt"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

const tamanhoMinimoNome = 3
const tamanhoMaximoNome = 100

// Categoria representa uma categoria de chamados no sistema.
type Categoria struct {
	id           string
	nome         string
	status       bool
	criadoEm     time.Time
	atualizadoEm time.Time
}

// Novo cria uma nova instância de Categoria com os dados fornecidos.
func Novo(id, nome string) (*Categoria, error) {
	now := time.Now()
	categoria := Categoria{
		id:           strings.TrimSpace(id),
		nome:         strings.TrimSpace(nome),
		status:       true,
		criadoEm:     now,
		atualizadoEm: now,
	}

	if err := categoria.Validar(); err != nil {
		return nil, err
	}

	return &categoria, nil
}

// CarregarDoBD recebe os dados de uma categoria do banco de dados e retorna uma instância de Categoria.
func CarregarDoBD(c CategoriaDB) *Categoria {
	return &Categoria{
		id:           c.ID,
		nome:         c.Nome,
		status:       c.Status,
		criadoEm:     c.CriadoEm,
		atualizadoEm: c.AtualizadoEm,
	}
}

// Validar valida os campos da categoria.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao
func (c Categoria) Validar() error {
	erros := domain.NovoErrosValidacao()

	if c.id == "" {
		erros.Add("ID", "o ID da categoria não pode estar vazio")
	}

	if len(c.nome) < tamanhoMinimoNome || len(c.nome) > tamanhoMaximoNome {
		erros.Add("Nome", fmt.Sprintf("o nome da categoria deve ter entre %d e %d caracteres",
			tamanhoMinimoNome, tamanhoMaximoNome))
	}

	if erros.HaErros() {
		return erros
	}
	return nil
}

// Ativar ativa a categoria.
func (c Categoria) Ativar() Categoria {
	c.status = true
	c.atualizadoEm = time.Now()
	return c
}

// Desativar desativa a categoria.
func (c Categoria) Desativar() Categoria {
	c.status = false
	c.atualizadoEm = time.Now()
	return c
}

// ComDadosAtualizados recebe os parâmetros para atualizar os dados da categoria 
// e retorna uma nova instância de Categoria com os dados atualizados.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (c Categoria) ComDadosAtualizados(params AtualizarParams) (Categoria, error) {
	if params.Nome != nil {
		c.nome = *params.Nome
	}

	if params.Status != nil {
		c.status = *params.Status
	}

	c.atualizadoEm = time.Now()

	if err := c.Validar(); err != nil {
		return Categoria{}, err
	}

	return c, nil
}

// Métodos de acesso aos campos da Categoria.

func (c Categoria) ID() string              { return c.id }
func (c Categoria) Nome() string            { return c.nome }
func (c Categoria) Status() bool            { return c.status }
func (c Categoria) CriadoEm() time.Time     { return c.criadoEm }
func (c Categoria) AtualizadoEm() time.Time { return c.atualizadoEm }

// String retorna uma representação em string da Categoria.
func (c Categoria) String() string {
	return fmt.Sprintf(
		"[ID=%s | Nome=%s | Status=%t | CriadoEm=%s | AtualizadoEm=%s]",
		c.id,
		c.nome,
		c.status,
		c.CriadoEm().Format(time.RFC3339),
		c.AtualizadoEm().Format(time.RFC3339),
	)
}
