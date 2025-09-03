package user

import "time"

type Permissao string

const (
	PermADM  Permissao = "ADM"  // Administrador
	PermTEC  Permissao = "TEC"  // Técnico
	PermSUP  Permissao = "SUP"  // Técnico de Suporte (Help Desk)
	PermINF  Permissao = "INF"  // Técnico de Infraestrutura
	PermVOIP Permissao = "VOIP" // Técnico de Telefonia
	PermIMP  Permissao = "IMP"  // Técnico de Impressoras
	PermCAD  Permissao = "CAD"  // Cadastro de usuários
	PermUSR  Permissao = "USR"  // Usuário comum (pode apenas abrir chamados)
	PermDEV  Permissao = "DEV"  // Desenvolvedor
)

type Usuario struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Login        string    `json:"login"`
	Email        string    `json:"email"`
	Permissao    Permissao `json:"permissao"`
	Status       bool      `json:"status"`
	Avatar       *string   `json:"avatar,omitempty"`
	UltimoLogin  time.Time `json:"ultimoLogin"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
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
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Login        string    `json:"login"`
	Email        string    `json:"email"`
	Permissao    Permissao `json:"permissao"`
	Status       bool      `json:"status"`
	Avatar       *string   `json:"avatar,omitempty"`
	UltimoLogin  time.Time `json:"ultimoLogin"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
}
