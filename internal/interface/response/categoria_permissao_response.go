package response

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

type CategoriaPermissaoResponse struct {
	CategoriaID  string          `json:"categoriaId"`
	UsuarioID    string          `json:"usuarioId"`
	Permissao    model.Permissao `json:"permissao"`
	CriadoEm     time.Time       `json:"criadoEm"`
	AtualizadoEm time.Time       `json:"atualizadoEm"`
}

// ToCategoriaPermissaoResponse converte um modelo CategoriaPermissao para CategoriaPermissaoResponse
func ToCategoriaPermissaoResponse(c *model.CategoriaPermissao) *CategoriaPermissaoResponse {
	return &CategoriaPermissaoResponse{
		CategoriaID:  c.CategoriaID,
		UsuarioID:    c.UsuarioID,
		Permissao:    c.Permissao,
		CriadoEm:     c.CriadoEm,
		AtualizadoEm: c.AtualizadoEm,
	}
}