package response

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// StatusSubcategoria representa o status de uma subcategoria.
type StatusSubcategoria struct {
	Ativo bool `json:"ativo"`
}

// SubcategoriaResponse representa a estrutura de resposta para uma subcategoria.
type SubcategoriaResponse struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Status       bool      `json:"status"`
	CategoriaID  string    `json:"categoria_id"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// ToSubcategoriaResponse converte um modelo Subcategoria para SubcategoriaResponse
func ToSubcategoriaResponse(s *model.Subcategoria) *SubcategoriaResponse {
	return &SubcategoriaResponse{
		ID:           s.ID,
		CategoriaID:  s.CategoriaID,
		Nome:         s.Nome,
		Status:       s.Status,
		CriadoEm:     s.CriadoEm,
		AtualizadoEm: s.AtualizadoEm,
	}
}
