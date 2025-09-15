package response

import (
	"time"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// UsuarioResponse representa a estrutura de resposta para dados de usuário
type UsuarioResponse struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Login        string    `json:"login"`
	Email        string    `json:"email"`
	Permissao    model.Permissao `json:"permissao"`
	Status       bool      `json:"status"`
	Avatar       *string   `json:"avatar,omitempty"`
	UltimoLogin  time.Time `json:"ultimoLogin"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
}

func ToUsuarioResponse(u *model.Usuario) *UsuarioResponse {
	return &UsuarioResponse{
		ID:           u.ID,
		Nome:         u.Nome,
		Login:        u.Login,
		Email:        u.Email,
		Permissao:    u.Permissao,
		Status:       u.Status,
		Avatar:       u.Avatar,
		UltimoLogin:  u.UltimoLogin,
		CriadoEm:     u.CriadoEm,
		AtualizadoEm: u.AtualizadoEm,
	}
}