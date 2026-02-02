package auth

import (
	"context"
	"fmt"
	"net/http"

	goJwt "github.com/golang-jwt/jwt/v5"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

var ErrClaimsInvalidasOuAusentes = fmt.Errorf("claims inválidas ou ausentes")

//============================================================================//
// Definições e funções relacionadas às claims do JWT
//============================================================================//

// ctxKey é um tipo personalizado para evitar colisões de chaves no contexto
type ctxKey string

// chaveUsuario é a chave usada para armazenar as claims do usuário no contexto
const chaveUsuario ctxKey = "user"

// Claims define as claims personalizadas para o token JWT
type Claims struct {
	ID        string `json:"sub"`
	Login     string `json:"login"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Permissao string `json:"permissao"`
	goJwt.RegisteredClaims
}

// criarClaims cria as claims do JWT a partir do usuário.
func novoClaims(u *usr.Usuario) Claims {
	return Claims{
		ID:        u.ID(),
		Login:     u.Login(),
		Nome:      u.Nome(),
		Email:     u.Email().String(),
		Permissao: u.Permissao().String(),
	}
}

// ClaimsdoContexto retorna claims do usuário presente no contexto.
func ClaimsdoContexto(ctx context.Context) (*Claims, error) {
	valor := ctx.Value(chaveUsuario)
	claims, ok := valor.(*Claims)
	if !ok || claims == nil {
		return nil, ErrClaimsInvalidasOuAusentes
	}
	return claims, nil
}

// ClaimsDoRequest extrai as claims do contexto da requisição HTTP.
func ClaimsDoRequest(r *http.Request) (*Claims, error) {
	return ClaimsdoContexto(r.Context())
}

//============================================================================//
// Definições e funções relacionadas ao usuário autenticado
//============================================================================//

// UsuarioAutenticado representa o usuário autenticado com suas claims.
type UsuarioAutenticado struct {
	id        string
	login     string
	nome      string
	email     string
	permissao usr.Permissao
}

// NovoUsuarioAutenticado cria uma nova instância de UsuarioAutenticado a partir das claims.
func NovoUsuarioAutenticado(claims *Claims) *UsuarioAutenticado {
	return &UsuarioAutenticado{
		id:        claims.ID,
		login:     claims.Login,
		nome:      claims.Nome,
		email:     claims.Email,
		permissao: usr.Permissao(claims.Permissao),
	}
}

// --- Métodos de acesso aos campos de UsuarioAutenticado ---

func (u *UsuarioAutenticado) ID() string               { return u.id }
func (u *UsuarioAutenticado) Login() string            { return u.login }
func (u *UsuarioAutenticado) Nome() string             { return u.nome }
func (u *UsuarioAutenticado) Email() string            { return u.email }
func (u *UsuarioAutenticado) Permissao() usr.Permissao { return u.permissao }

func (u *UsuarioAutenticado) EhTEC() bool {
	return u.permissao == usr.PermTEC
}

func (u *UsuarioAutenticado) EhADM() bool {
	return u.permissao == usr.PermADM
}

func (u *UsuarioAutenticado) EhDEV() bool {
	return u.permissao == usr.PermDEV
}

// UsuarioAutenticadoDoContexto extrai o usuário autenticado do contexto.
func UsuarioAutenticadoDoContexto(ctx context.Context) (*UsuarioAutenticado, error) {
	claims, err := ClaimsdoContexto(ctx)
	if err != nil {
		return nil, err
	}
	return NovoUsuarioAutenticado(claims), nil
}

// UsuarioAutenticadoDoRequest extrai o usuário autenticado do contexto da requisição HTTP.
func UsuarioAutenticadoDoRequest(r *http.Request) (*UsuarioAutenticado, error) {
	return UsuarioAutenticadoDoContexto(r.Context())
}

// ContextoComClaims adiciona as claims ao contexto. Útil para testes.
func ContextoComClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, chaveUsuario, claims)
}
