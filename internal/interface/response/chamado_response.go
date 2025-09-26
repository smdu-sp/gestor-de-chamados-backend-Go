package response

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

type ChamadoStatus struct {
	Status string `json:"status"`
}

type TecnicoAtribuido struct {
	Atribuido bool `json:"atribuido"`
}

type ChamadoResponse struct {
	ID             string     `json:"id"`
	Titulo         string     `json:"titulo"`
	Descricao      string     `json:"descricao"`
	Status         string     `json:"status"`
	CriadoEm       time.Time  `json:"criado_em"`
	AtualizadoEm   time.Time  `json:"atualizado_em"`
	SolucionadoEm  *time.Time `json:"solucionado_em"`
	Solucao        *string    `json:"solucao"`
	FechadoEm      *time.Time `json:"fechado_em"`
	CategoriaID    string     `json:"categoria_id"`
	SubcategoriaID string     `json:"subcategoria_id"`
	CriadorID      string     `json:"criador_id"`
	AtribuidoID    *string     `json:"atribuido_id"`
}

// ToChamadoResponse converte um modelo Chamado para ChamadoResponse
func ToChamadoResponse(c *model.Chamado) *ChamadoResponse {
	return &ChamadoResponse{
		ID:             c.ID,
		Titulo:         c.Titulo,
		Descricao:      c.Descricao,
		Status:         string(c.Status),
		CriadoEm:       c.CriadoEm,
		AtualizadoEm:   c.AtualizadoEm,
		SolucionadoEm:  c.SolucionadoEm,
		Solucao:        c.Solucao,
		FechadoEm:      c.FechadoEm,
		CategoriaID:    c.CategoriaID,
		SubcategoriaID: c.SubcategoriaID,
		CriadorID:      c.CriadorID,
		AtribuidoID:    c.AtribuidoID,
	}
}

