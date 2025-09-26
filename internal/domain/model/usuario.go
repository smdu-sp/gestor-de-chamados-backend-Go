package model

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
)

// Erros de validação específicos para o modelo Usuario.
var (
	ErrEmailInvalido     = errors.New("email inválido")
	ErrPermissaoInvalida = errors.New("permissão inválida: a permissão deve ser uma das seguintes: ADM, TEC, SUP, INF, VOIP, IMP, CAD, USR, DEV")
	ErrNomeInvalido      = errors.New("nome não pode ser vazio")
	ErrLoginInvalido     = errors.New("login não pode ser vazio")
	ErrIDInvalido        = errors.New("ID não pode ser vazio")
)

// Permissao define os níveis de permissão dos usuários.
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

// Usuario representa um usuário do sistema.
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

// NewUsuario cria uma nova instância de Usuario com os dados fornecidos.
func NewUsuario(id, nome, login, email string, permissao Permissao, status bool, avatar *string) (*Usuario, error) {
	now := time.Now()
	usuario := &Usuario{
		ID:           id,
		Nome:         nome,
		Login:        login,
		Email:        email,
		Permissao:    permissao,
		Status:       status,
		Avatar:       avatar,
		UltimoLogin:  now,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
	if err := ValidarUsuario(usuario); err != nil {
		return nil, fmt.Errorf("[model.NewUsuario]: %w", err)
	}
	return usuario, nil
}

// ValidarUsuario valida os campos do usuário.
func ValidarUsuario(u *Usuario) error {
	var erros []error

	if u.Nome == "" {
		erros = append(erros, ErrNomeInvalido)
	}
	if u.Login == "" {
		erros = append(erros, ErrLoginInvalido)
	}
	if err := ValidarEmail(u.Email); err != nil {
		erros = append(erros, err)
	}
	if err := ValidarPermissao(string(u.Permissao)); err != nil {
		erros = append(erros, err)
	}
	if len(erros) > 0 {
		return fmt.Errorf("[model.ValidarUsuario] erros de validação: %v", erros)
	}
	return nil
}

// ValidarEmail valida o formato do email
func ValidarEmail(email string) error {
	email = strings.ToLower(email)

	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	if !re.MatchString(email) {
		return fmt.Errorf("[model.ValidarEmail] %w: %s", ErrEmailInvalido, email)
	}

	return nil
}

// ValidarPermissao verifica se a permissão do usuário é válida.
func ValidarPermissao(permissao string) error {
	validPermissoes := [9]string{"ADM", "TEC", "SUP", "INF", "VOIP", "IMP", "CAD", "USR", "DEV"}
	if slices.Contains(validPermissoes[:], permissao) {
		return nil
	}
	return fmt.Errorf("[model.ValidarPermissao] %w: %s", ErrPermissaoInvalida, permissao)
}

// UsuarioFiltro representa os filtros para listar usuários.
type UsuarioFiltro struct {
	Pagina    int
	Limite    int
	Busca     *string
	Status    *bool
	Permissao *string
}
