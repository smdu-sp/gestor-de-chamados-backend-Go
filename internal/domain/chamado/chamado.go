package chamado

import (
	"fmt"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

const tamanhoMinTitulo = 10
const tamanhoMaxTitulo = 255
const tamanhoMinDescricao = 10
const tamanhoMaxDescricao = 1000

// Chamado representa um chamado no sistema
type Chamado struct {
	id             string
	categoriaID    string
	subcategoriaID string
	criadorID      string
	titulo         string
	descricao      string
	status         StatusChamado
	arquivado      bool
	solucao        *string
	solucionadoEm  *time.Time
	fechadoEm      *time.Time
	criadoEm       time.Time
	atualizadoEm   time.Time
}

// Novo cria uma nova instância de Chamado com os dados fornecidos
func Novo(id, categoriaID, subcategoriaID, criadorID, titulo, descricao string) (*Chamado, error) {
	now := time.Now()
	chamado := Chamado{
		id:             strings.TrimSpace(id),
		categoriaID:    strings.TrimSpace(categoriaID),
		subcategoriaID: strings.TrimSpace(subcategoriaID),
		criadorID:      strings.TrimSpace(criadorID),
		titulo:         strings.TrimSpace(titulo),
		descricao:      strings.TrimSpace(descricao),
		status:         StatusAberto,
		arquivado:      false,
		criadoEm:       now,
		atualizadoEm:   now,
	}

	if err := chamado.Validar(); err != nil {
		return nil, err
	}

	return &chamado, nil
}

// CarregarDoDB recebe um ChamadoDB e retorna uma instância de Chamado.
func CarregarDoDB(c ChamadoDB) *Chamado {
	return &Chamado{
		id:             c.ID,
		categoriaID:    c.CategoriaID,
		subcategoriaID: c.SubcategoriaID,
		criadorID:      c.CriadorID,
		titulo:         c.Titulo,
		descricao:      c.Descricao,
		status:         StatusChamado(c.Status),
		arquivado:      c.Arquivado,
		solucao:        c.Solucao,
		solucionadoEm:  c.SolucionadoEm,
		fechadoEm:      c.FechadoEm,
		criadoEm:       c.CriadoEm,
		atualizadoEm:   c.AtualizadoEm,
	}
}

// Validar valida os campos do chamado.
//
// Em caso de erros de validação, retorna um erro do tipo domain.ErrosValidacao.
func (c Chamado) Validar() error {
	erros := domain.NovoErrosValidacao()

	if c.id == "" {
		erros.Add("ID", "o ID do chamado não pode estar vazio")
	}

	if len(c.titulo) < tamanhoMinTitulo || len(c.titulo) > tamanhoMaxTitulo {
		erros.Add("Titulo", fmt.Sprintf("o título do chamado deve ter entre %d e %d caracteres",
			tamanhoMinTitulo, tamanhoMaxTitulo))
	}

	if len(c.descricao) < tamanhoMinDescricao || len(c.descricao) > tamanhoMaxDescricao {
		erros.Add("Descricao", fmt.Sprintf("a descrição do chamado deve ter entre %d e %d caracteres",
			tamanhoMinDescricao, tamanhoMaxDescricao))
	}

	if err := ValidarStatus(c.status); err != nil {
		erros.Add("Status", err.Error())
	}

	if c.categoriaID == "" {
		erros.Add("CategoriaID", "o ID da categoria não pode estar vazio")
	}

	if c.subcategoriaID == "" {
		erros.Add("SubcategoriaID", "o ID da subcategoria não pode estar vazio")
	}

	if c.criadorID == "" {
		erros.Add("CriadorID", "o ID do criador não pode estar vazio")
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// AtualizarSolucao recebe uma solução e atualiza o chamado como resolvido.
//
// Em caso de erro de validação, retorna uma instância de domain.ErrosValidacao.
func (c Chamado) AtualizarSolucao(solucao string) (Chamado, error) {
	now := time.Now()
	c.solucao = &solucao
	c.solucionadoEm = &now
	c.status = StatusResolvido
	c.atualizadoEm = now

	if err := c.Validar(); err != nil {
		return Chamado{}, err
	}

	return c, nil
}

// AtualizarStatus recebe um status e atualiza o chamado.
//
// Em caso de erro de validação, retorna uma instância de domain.ErrosValidacao.
func (c Chamado) AtualizarStatus(status StatusChamado, solucao *string) (Chamado, error) {
	now := time.Now()
	c.status = status
	c.atualizadoEm = now

	if status == StatusResolvido && solucao != nil {
		c.solucao = solucao
		c.solucionadoEm = &now
	}

	if status == StatusFechado {
		c.fechadoEm = &now
	}

	if err := c.Validar(); err != nil {
		return Chamado{}, err
	}

	return c, nil
}

// ComDadosAtualizados recebe os parâmetros para atualização do chamado
// e retorna uma nova instância com os dados atualizados.
//
// Em caso de erro de validação, retorna uma instância de domain.ErrosValidacao.
func (c Chamado) ComDadosAtualizados(params AtualizarParams) (Chamado, error) {
	if params.Titulo != nil {
		c.titulo = *params.Titulo
	}

	if params.Descricao != nil {
		c.descricao = *params.Descricao
	}

	if params.Arquivado != nil {
		c.arquivado = *params.Arquivado
	}

	if params.CategoriaID != nil {
		c.categoriaID = *params.CategoriaID
	}

	if params.SubcategoriaID != nil {
		c.subcategoriaID = *params.SubcategoriaID
	}

	c.atualizadoEm = time.Now()

	if err := c.Validar(); err != nil {
		return Chamado{}, err
	}

	return c, nil
}

// NaoPodeSerModificado verifica se o chamado está em um status que não permite modificações.
func (c Chamado) NaoPodeSerModificado() bool {
	return c.arquivado || c.status == StatusFechado || c.status == StatusRejeitado
}

// StatusAtribuido marca o chamado como atribuído a um técnico.
func (c Chamado) StatusAtribuido() Chamado {
	c.status = StatusAtribuido
	c.atualizadoEm = time.Now()
	return c
}

// Arquivar marca o chamado como arquivado.
func (c Chamado) Arquivar() Chamado {
	c.arquivado = true
	c.atualizadoEm = time.Now()
	return c
}

// Desarquivar marca o chamado como não arquivado.
func (c Chamado) Desarquivar() Chamado {
	c.arquivado = false
	c.atualizadoEm = time.Now()
	return c
}

//--- Métodos de acesso aos campos do Chamado ---

func (c Chamado) ID() string                { return c.id }
func (c Chamado) CategoriaID() string       { return c.categoriaID }
func (c Chamado) SubcategoriaID() string    { return c.subcategoriaID }
func (c Chamado) CriadorID() string         { return c.criadorID }
func (c Chamado) Titulo() string            { return c.titulo }
func (c Chamado) Descricao() string         { return c.descricao }
func (c Chamado) Status() StatusChamado     { return c.status }
func (c Chamado) Arquivado() bool           { return c.arquivado }
func (c Chamado) Solucao() *string          { return c.solucao }
func (c Chamado) SolucionadoEm() *time.Time { return c.solucionadoEm }
func (c Chamado) FechadoEm() *time.Time     { return c.fechadoEm }
func (c Chamado) CriadoEm() time.Time       { return c.criadoEm }
func (c Chamado) AtualizadoEm() time.Time   { return c.atualizadoEm }

// String retorna uma representação em string do Chamado para fins de logging.
func (c Chamado) String() string {
	formatTime := func(t *time.Time) string {
		if t == nil {
			return "<nil>"
		}
		return t.Format(time.RFC3339)
	}

	return fmt.Sprintf(
		"[ID: %s | CategoriaID: %s | SubcategoriaID: %s | CriadorID: %s | Titulo: %s"+
			"| Descricao: %s | Status: %s | Arquivado: %t | Solucao: %v | SolucionadoEm: %s"+
			"| FechadoEm: %s | CriadoEm: %s | AtualizadoEm: %s]",
		c.id,
		c.categoriaID,
		c.subcategoriaID,
		c.criadorID,
		c.titulo,
		c.descricao,
		c.status.String(),
		c.arquivado,
		c.solucao,
		formatTime(c.solucionadoEm),
		formatTime(c.fechadoEm),
		formatTime(&c.criadoEm),
		formatTime(&c.atualizadoEm),
	)
}

