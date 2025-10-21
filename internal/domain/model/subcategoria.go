package model

import (
	"fmt"
	"time"
)

// Subcategoria representa uma subcategoria de chamados no sistema.
type Subcategoria struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Status       bool      `json:"status"`
	CategoriaID  string    `json:"categoriaId"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
}

// NewSubcategoria cria uma nova instância de Subcategoria com os dados fornecidos.
func NewSubcategoria(id, nome, categoriaID string, status bool) (*Subcategoria, error) {
	if nome == "" {
		return nil, fmt.Errorf("[model.NewSubcategoria]: %w", ErrNomeInvalido)
	}
	if categoriaID == "" {
		return nil, fmt.Errorf("[model.NewSubcategoria]: %w", ErrCategoriaIDInvalido)
	}

	now := time.Now()
	subcategoria := &Subcategoria{
		ID:           id,
		Nome:         nome,
		Status:       status,
		CategoriaID:  categoriaID,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
	return subcategoria, nil
}

// SubcategoriaFiltro representa os critérios de filtro para listar subcategorias.
type SubcategoriaFiltro struct {
	Pagina      int
	Limite      int
	Busca       *string
	Status      *bool
}

// String retorna uma representação em string da subcategoria para fins de logging.
func (s *Subcategoria) String() string {
    return fmt.Sprintf(
        "[ID=%s | Nome=%s | Status=%t | CategoriaID=%s]",
        s.ID, s.Nome, s.Status, s.CategoriaID,
    )
}
