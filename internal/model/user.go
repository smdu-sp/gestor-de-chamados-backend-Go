package model

import "time"

type Permissao string

const (
	ADM  Permissao = "ADM"
	SUP  Permissao = "SUP"
	INF  Permissao = "INF"
	VOIP Permissao = "VOIP"
	CAD  Permissao = "CAD"
	USR  Permissao = "USR"
	IMP  Permissao = "IMP"
	DEV  Permissao = "DEV"
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
