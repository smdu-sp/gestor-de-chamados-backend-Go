package response

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
)

// StatusUsuario representa a estrutura de resposta para o status do usuário
type StatusUsuario struct {
	Ativo bool `json:"ativo"`
}

// AutorizacaoUsuario representa a estrutura de resposta para autorização de usuário
type AutorizacaoUsuario struct {
	Autorizado bool `json:"autorizado"`
}

// ValidaUsuario representa a estrutura de resposta para validação de usuário
type ValidaUsuario struct {
	Valido bool `json:"valido"`
}

// BuscarNovo representa a estrutura de resposta para busca de novo usuário
type BuscarNovo struct {
	Login string `json:"login"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Tecnico representa a estrutura de resposta de um técnico com ID e Nome
type Tecnico struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

// PermissaoUsuario representa a estrutura de resposta para permissão de usuário
type PermissaoUsuario struct {
	Permissao string `json:"permissao"`
}

// UsuarioResponse representa a estrutura de resposta para dados de usuário
type UsuarioResponse struct {
	ID           string          `json:"id"`
	Nome         string          `json:"nome"`
	Login        string          `json:"login"`
	Email        string          `json:"email"`
	Permissao    model.Permissao `json:"permissao"`
	Status       bool            `json:"status"`
	Avatar       *string         `json:"avatar,omitempty"`
	UltimoLogin  time.Time       `json:"ultimoLogin"`
	CriadoEm     time.Time       `json:"criadoEm"`
	AtualizadoEm time.Time       `json:"atualizadoEm"`
}

// ToUsuarioResponse converte um modelo Usuario para UsuarioResponse
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
