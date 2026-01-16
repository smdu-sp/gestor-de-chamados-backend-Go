package usuario

import (
	"fmt"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

const tamanhoMinNome = 3
const tamanhoMaxNome = 100

// Usuario representa um usuário do sistema.
type Usuario struct {
	id           string
	nome         string
	login        string
	email        Email
	permissao    Permissao
	status       bool
	avatar       *string
	ultimoLogin  time.Time
	criadoEm     time.Time
	atualizadoEm time.Time
}

// Novo cria uma nova instância de Usuario com os dados fornecidos.
func Novo(id, nome, login string, email Email, permissao Permissao, avatar *string) (*Usuario, error) {
	now := time.Now()
	usuario := Usuario{
		id:           strings.TrimSpace(id),
		nome:         strings.TrimSpace(nome),
		login:        strings.TrimSpace(login),
		email:        email,
		permissao:    permissao,
		status:       true,
		avatar:       avatar,
		ultimoLogin:  now,
		criadoEm:     now,
		atualizadoEm: now,
	}

	if err := usuario.Validar(); err != nil {
		return nil, err
	}

	return &usuario, nil
}

// CarregarDoBD recebe uma estrutura UsuarioDB e retorna uma instância de Usuario.
func CarregarDoBD(u UsuarioDB) *Usuario {
	return &Usuario{
		id:           u.ID,
		nome:         u.Nome,
		login:        u.Login,
		email:        NovoEmail(u.Email),
		permissao:    Permissao(u.Permissao),
		status:       u.Status,
		avatar:       u.Avatar,
		ultimoLogin:  u.UltimoLogin,
		criadoEm:     u.CriadoEm,
		atualizadoEm: u.AtualizadoEm,
	}
}

// Validar valida os campos do usuário.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (u *Usuario) Validar() error {
	erros := domain.NovoErrosValidacao()

	if u.id == "" {
		erros.Add("ID", "o ID do usuário não pode estar vazio")
	}

	if len(u.nome) < tamanhoMinNome || len(u.nome) > tamanhoMaxNome {
		erros.Add("Nome", fmt.Sprintf("o nome do usuário deve ter entre %d e %d caracteres",
			tamanhoMinNome, tamanhoMaxNome))
	}

	if u.login == "" {
		erros.Add("Login", "o login não pode estar vazio")
	}

	if err := u.email.ValidarEmail(); err != nil {
		erros.Add("Email", err.Error())
	}

	if err := ValidarPermissao(u.permissao); err != nil {
		erros.Add("Permissao", err.Error())
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// AtualizarDados recebe os parâmetros para atualizar os dados do usuário.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (u *Usuario) AtualizarDados(params AtualizarParams) error {
	if params.Nome != nil {
		u.nome = *params.Nome
	}

	if params.Login != nil {
		u.login = *params.Login
	}

	if params.Email != nil {
		u.email = *params.Email
	}

	if params.Permissao != nil {
		u.permissao = *params.Permissao
	}

	if params.Status != nil {
		u.status = *params.Status
	}

	if params.Avatar != nil {
		u.avatar = params.Avatar
	}

	u.atualizadoEm = time.Now()

	return u.Validar()
}

// Ativar ativa o usuário.
func (u *Usuario) Ativar() {
	u.status = true
	u.atualizadoEm = time.Now()
}

// Desativar desativa o usuário.
func (u *Usuario) Desativar() {
	u.status = false
	u.atualizadoEm = time.Now()
}

// AtualizarPermissao recebe uma nova permissão e a atribui ao usuário.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (u *Usuario) AtualizarPermissao(novaPermissao Permissao) error {
	u.permissao = novaPermissao
	u.atualizadoEm = time.Now()

	return u.Validar()
}

// AtualizarUltimoLogin atualiza o timestamp do último login do usuário.
func (u *Usuario) AtualizarUltimoLogin() {
	now := time.Now()
	u.ultimoLogin = now
}

// --- Métodos de acesso aos campos do usuário ---

func (u Usuario) ID() string              { return u.id }
func (u Usuario) Nome() string            { return u.nome }
func (u Usuario) Login() string           { return u.login }
func (u Usuario) Email() string           { return strings.ToLower(u.email.String()) }
func (u Usuario) Permissao() Permissao    { return u.permissao }
func (u Usuario) Status() bool            { return u.status }
func (u Usuario) Avatar() *string         { return u.avatar }
func (u Usuario) UltimoLogin() time.Time  { return u.ultimoLogin }
func (u Usuario) CriadoEm() time.Time     { return u.criadoEm }
func (u Usuario) AtualizadoEm() time.Time { return u.atualizadoEm }

// String retorna uma representação em string do usuário para fins de logging.
func (u *Usuario) String() string {
	return fmt.Sprintf(
		"[ID=%s | Nome=%s | Login=%s | Email=%s | Permissao=%s"+
			"| Status=%t | UltimoLogin=%s | CriadoEm=%s | AtualizadoEm=%s]",
		u.id,
		u.nome,
		u.login,
		u.Email(),
		u.permissao,
		u.status,
		u.UltimoLogin().Format(time.RFC3339),
		u.CriadoEm().Format(time.RFC3339),
		u.AtualizadoEm().Format(time.RFC3339),
	)
}
