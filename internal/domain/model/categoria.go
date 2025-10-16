package model

import (
	"errors"
	"fmt"
	"time"
)

var ErrCategoriaIDInvalido = errors.New("ID da categoria não pode ser vazio")

// Categoria representa uma categoria de chamados no sistema.
type Categoria struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Status       bool      `json:"status"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
}

// NewCategoria cria uma nova instância de Categoria com os dados fornecidos.
func NewCategoria(id, nome string, status bool) (*Categoria, error) {
	if nome == "" {
		return nil, fmt.Errorf("[model.NewCategoria]: %w", ErrNomeInvalido)
	}

	now := time.Now()
	categoria := &Categoria{
		ID:           id,
		Nome:         nome,
		Status:       status,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
	return categoria, nil
}

// CategoriaFiltro representa os critérios de filtro para listar categorias.
type CategoriaFiltro struct {
	Pagina int
	Limite int
	Busca  *string
	Status *bool
}

// String retorna uma representação em string da Categoria para fins de logging.
func (c *Categoria) String() string {
	return fmt.Sprintf(
		"[ID=%s | Nome=%s | Status=%t]", 
		c.ID, c.Nome, c.Status,
	)
}
