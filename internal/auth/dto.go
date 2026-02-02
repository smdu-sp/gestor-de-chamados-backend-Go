package auth

import (
	"time"

	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// LoginReq representa o payload para a requisição de login.
type LoginReq struct {
	Login string `json:"login"`
	Senha string `json:"senha"`
}

// RefreshReq representa o payload para a requisição de refresh de tokens.
type RefreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

// TokenResp representa um par de tokens JWT (access e refresh)
type TokenResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// ParaTokenResp cria uma nova instância de TokenResp.
func ParaTokenResp(accessToken, refreshToken string) *TokenResp {
	return &TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// UsuarioResp representa a estrutura de resposta para dados de usuário.
type UsuarioResp struct {
	ID           string  `json:"id"`
	Nome         string  `json:"nome"`
	Login        string  `json:"login"`
	Email        string  `json:"email"`
	Permissao    string  `json:"permissao"`
	Status       bool    `json:"status"`
	Avatar       *string `json:"avatar,omitempty"`
	UltimoLogin  string  `json:"ultimoLogin"`
	CriadoEm     string  `json:"criadoEm"`
	AtualizadoEm string  `json:"atualizadoEm"`
}

// ParaUsuarioResp converte um modelo Usuario para UsuarioResp.
func ParaUsuarioResp(u usr.Usuario) UsuarioResp {
	return UsuarioResp{
		ID:           u.ID(),
		Nome:         u.Nome(),
		Login:        u.Login(),
		Email:        u.Email().String(),
		Permissao:    u.Permissao().String(),
		Status:       u.Status(),
		Avatar:       u.Avatar(),
		UltimoLogin:  u.UltimoLogin().Format(time.RFC3339),
		CriadoEm:     u.CriadoEm().Format(time.RFC3339),
		AtualizadoEm: u.AtualizadoEm().Format(time.RFC3339),
	}
}
