package model

import (
	"time"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
	"errors"
)

// Erros sentinela do domínio de usuário
var (
	ErrCamposObrigatorios   = errors.New("campos obrigatórios ausentes")
	ErrUsuarioNaoEncontrado = errors.New("usuário não encontrado")
)

// Permissao define os níveis de permissão dos usuários
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

// Usuario representa um usuário do sistema
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

// NewUser cria uma nova instância de Usuario com os dados fornecidos
func NewUser(nome, login, email string, permissao Permissao, status bool, avatar *string) (*Usuario, error) {
	validUUID, err := util.NewUUIDv7String()
	if err != nil {
		return nil, util.NewAppError(
			"NewUser",
			util.LevelError,
			"erro ao gerar ID para novo usuário",
			err,
		)
	}
	
	now := time.Now()
	return &Usuario{
		ID:           validUUID,
		Nome:         nome,
		Login:        login,
		Email:        email,
		Permissao:    permissao,
		Status:       status,
		Avatar:       avatar,
		UltimoLogin:  now,
		CriadoEm:     now,
		AtualizadoEm: now,
	}, nil
}


