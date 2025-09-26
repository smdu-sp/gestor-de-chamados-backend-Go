package response

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// StatusCategoria representa o status de uma categoria.
type StatusCategoria struct {
	Ativo bool `json:"ativo"`
}

// CategoriaResponse representa a estrutura de resposta para uma categoria.
type CategoriaResponse struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Status       bool      `json:"status"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// ToCategoriaResponse converte um modelo Categoria para CategoriaResponse
func ToCategoriaResponse(c *model.Categoria) *CategoriaResponse {
	return &CategoriaResponse{
		ID:           c.ID,
		Nome:         c.Nome,
		Status:       c.Status,
		CriadoEm:     c.CriadoEm,
		AtualizadoEm: c.AtualizadoEm,
	}
}
