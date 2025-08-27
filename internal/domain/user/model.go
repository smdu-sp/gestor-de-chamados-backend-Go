package user

import "time"

type Permissao string

const (
	PermDEV Permissao = "DEV"
	PermADM Permissao = "ADM"
	PermUSR Permissao = "USR"
	PermTEC Permissao = "TEC"
)

type Usuario struct {
	ID           string     `json:"id"`
	Nome         string     `json:"nome"`
	Login        string     `json:"login"`
	Email        string     `json:"email"`
	Permissao    Permissao  `json:"permissao"`
	Status       bool       `json:"status"`
	Avatar       *string    `json:"avatar,omitempty"`
	UltimoLogin  time.Time  `json:"ultimoLogin"`
	CriadoEm     time.Time  `json:"criadoEm"`
	AtualizadoEm time.Time  `json:"atualizadoEm"`
}

// DTOs
type LoginDTO struct {
	Login string `json:"login"`
	Senha string `json:"senha"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UsuarioResponse struct {
	ID        string    `json:"id"`
	Nome      string    `json:"nome"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Permissao Permissao `json:"permissao"`
	Status    bool      `json:"status"`
	Avatar    *string   `json:"avatar,omitempty"`
}
