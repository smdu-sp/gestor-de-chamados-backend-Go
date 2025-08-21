package model

import "time"

type Permissao string

const (
	ADM  Permissao = "ADM" 	// Administrador
	SUP  Permissao = "SUP" 	// Técnico de Suporte (Help Desk)
	INF  Permissao = "INF" 	// Técnico de Infraestrutura
	VOIP Permissao = "VOIP" // Técnico de Telefonia
	CAD  Permissao = "CAD" 	// Cadastro de usuários
	USR  Permissao = "USR" 	// Usuário comum
	IMP  Permissao = "IMP" 	// Técnico de Impressoras
	DEV  Permissao = "DEV" 	// Desenvolvedor
)

type User struct {
	ID           string     `json:"id"`
	Nome         string     `json:"nome"`
	Login        string     `json:"login"`
	Email        string     `json:"email"`
	Status       bool       `json:"status"`
	Avatar       *string    `json:"avatar,omitempty"`
	UltimoLogin  *time.Time `json:"ultimoLogin,omitempty"`
	CriadoEm     time.Time  `json:"criadoEm"`
	AtualizadoEm time.Time  `json:"atualizadoEm"`
	Permissao    Permissao  `json:"permissao"`
}
