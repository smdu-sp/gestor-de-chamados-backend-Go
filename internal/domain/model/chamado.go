package model

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// Erros de validação específicos para o modelo Chamado
var (
	ErrStatusChamadoInvalido       = errors.New("status inválido: o status deve ser uma das seguintes opções: ABERTO, ATRIBUIDO, RESOLVIDO, REJEITADO, FECHADO, ARQUIVADO")
	ErrCriadorChamadoInvalido      = errors.New("criador do chamado não pode ser vazio")
	ErrSubcategoriaChamadoInvalido = errors.New("subcategoria do chamado não pode ser vazia")
	ErrCategoriaChamadoInvalido    = errors.New("categoria do chamado não pode ser vazia")
	ErrDescricaoChamadoInvalido    = errors.New("descrição do chamado não pode ser vazia")
	ErrTituloChamadoInvalido       = errors.New("título do chamado não pode ser vazio")
)

// StatusChamado define os possíveis status de um chamado
type StatusChamado string

const (
	StatusAberto    StatusChamado = "ABERTO"
	StatusAtribuido StatusChamado = "ATRIBUIDO"
	StatusResolvido StatusChamado = "RESOLVIDO"
	StatusRejeitado StatusChamado = "REJEITADO"
	StatusFechado   StatusChamado = "FECHADO"
	StatusArquivado StatusChamado = "ARQUIVADO"
)

// Chamado representa um chamado no sistema
type Chamado struct {
	ID             string        `json:"id"`
	Titulo         string        `json:"titulo"`
	Descricao      string        `json:"descricao"`
	Status         StatusChamado `json:"status"`
	CriadoEm       time.Time     `json:"criadoEm"`
	AtualizadoEm   time.Time     `json:"atualizadoEm"`
	SolucionadoEm  *time.Time    `json:"solucionadoEm,omitempty"`
	Solucao        *string       `json:"solucao,omitempty"`
	FechadoEm      *time.Time    `json:"fechadoEm,omitempty"`
	CategoriaID    string        `json:"categoriaId"`
	SubcategoriaID string        `json:"subcategoriaId"`
	CriadorID      string        `json:"criadorId"`
	AtribuidoID    *string        `json:"atribuidoId"`
}

// NewChamado cria uma nova instância de Chamado com os dados fornecidos
func NewChamado(id, titulo, descricao string, status StatusChamado, categoriaID, subcategoriaID, criadorID string) (*Chamado, error) {
	now := time.Now()
	chamado := &Chamado{
		ID:             id,
		Titulo:         titulo,
		Descricao:      descricao,
		Status:         status,
		CriadoEm:       now,
		AtualizadoEm:   now,
		CategoriaID:    categoriaID,
		SubcategoriaID: subcategoriaID,
		CriadorID:      criadorID,
	}

	if err := ValidarChamado(chamado); err != nil {
		return nil, fmt.Errorf("[model.NewChamado] %w", err)
	}
	return chamado, nil
}

// ValidarChamado valida os campos do chamado
func ValidarChamado(c *Chamado) error {
	var erros []error

	if c.Titulo == "" {
		erros = append(erros, ErrTituloChamadoInvalido)
	}
	if c.Descricao == "" {
		erros = append(erros, ErrDescricaoChamadoInvalido)
	}
	if err := ValidarStatusChamado(string(c.Status)); err != nil {
		erros = append(erros, err)
	}
	if c.CategoriaID == "" {
		erros = append(erros, ErrCategoriaChamadoInvalido)
	}
	if c.SubcategoriaID == "" {
		erros = append(erros, ErrSubcategoriaChamadoInvalido)
	}
	if c.CriadorID == "" {
		erros = append(erros, ErrCriadorChamadoInvalido)
	}
	if len(erros) > 0 {
		return fmt.Errorf("[model.ValidarChamado] erros de validação: %v", erros)
	}
	return nil
}

// ValidarStatusChamado verifica se o status do chamado é válido
func ValidarStatusChamado(status string) error {
	validStatus := [6]string{"ABERTO", "ATRIBUIDO", "RESOLVIDO", "REJEITADO", "FECHADO", "ARQUIVADO"}
	if slices.Contains(validStatus[:], status) {
		return nil
	}
	return fmt.Errorf(
		"[validator.ValidarStatusChamado] erro ao validar status do chamado: %w",
		ErrStatusChamadoInvalido,
	)
}

// AdiconarSolucao adiciona uma solução ao chamado e atualiza seu status
func (c *Chamado) AdiconarSolucao(solucao string) {
	now := time.Now()
	c.Solucao = &solucao
	c.SolucionadoEm = &now
	c.Status = StatusResolvido
	c.AtualizadoEm = now
}

// AtribuirTecnico atribui um técnico ao chamado e atualiza seu status
func (c *Chamado) AtribuirTecnico(tecnicoID string) {
	now := time.Now()
	c.AtribuidoID = &tecnicoID
	c.Status = StatusAtribuido
	c.AtualizadoEm = now
}

// ChamadoFiltro representa os filtros possíveis para buscar chamados
type ChamadoFiltro struct {
	Pagina         int
	Limite         int
	Busca          *string
	Status         *string
	CategoriaID    *string
	SubcategoriaID *string
	CriadorID      *string
	AtribuidoID    *string
}
