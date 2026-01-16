package dto

import (
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// UsuarioPaginado é um alias para uma resposta paginada de usuários,
// usado para facilitar a leitura do código e documentação Swagger.
// @Description Estrutura paginada contendo uma lista de usuários.
type UsuarioPaginado RespPaginada[UsuarioResp]

// CriarUsuarioReq representa o payload para criar um usuário.
type CriarUsuarioReq struct {
	Nome      string  `json:"nome"`
	Login     string  `json:"login"`
	Email     string  `json:"email"`
	Permissao string  `json:"permissao"`
	Avatar    *string `json:"avatar,omitempty"`
}

// ParaCriarParams converte um modelo CriarUsuarioReq para CriarParams.
func (u CriarUsuarioReq) ParaCriarParams() usr.CriarParams {
	return usr.CriarParams{
		Nome:      u.Nome,
		Login:     u.Login,
		Email:     usr.NovoEmail(u.Email),
		Permissao: usr.Permissao(u.Permissao),
		Avatar:    u.Avatar,
	}
}

// AtualizarUsuarioReq representa o payload para atualizar um usuário.
type AtualizarUsuarioReq struct {
	Nome      *string `json:"nome,omitempty"`
	Login     *string `json:"login,omitempty"`
	Email     *string `json:"email,omitempty"`
	Permissao *string `json:"permissao,omitempty"`
	Status    *bool   `json:"status,omitempty"`
	Avatar    *string `json:"avatar,omitempty"`
}

// ParaAtualizarParams converte um modelo AtualizarUsuarioReq para AtualizarParams.
func (u AtualizarUsuarioReq) ParaAtualizarParams() usr.AtualizarParams {
	var email *usr.Email
	if u.Email != nil {
		e := usr.NovoEmail(*u.Email)
		email = &e
	}

	var permissao *usr.Permissao
	if u.Permissao != nil {
		p := usr.Permissao(*u.Permissao)
		permissao = &p
	}

	return usr.AtualizarParams{
		Nome:      u.Nome,
		Login:     u.Login,
		Email:     email,
		Permissao: permissao,
		Status:    u.Status,
		Avatar:    u.Avatar,
	}
}

// AtualizarPermissaoUsuarioReq representa o payload para atualizar a permissão de um usuário.
type AtualizarPermissaoUsuarioReq struct {
	Permissao string `json:"permissao"`
}

// ParaAtualizarPermissaoParams converte um modelo AtualizarPermissaoUsuarioReq para AtualizarPermissaoParams.
func (u AtualizarPermissaoUsuarioReq) ParaAtualizarPermissaoParams() usr.AtualizarPermissaoParams {
	return usr.AtualizarPermissaoParams{
		Permissao: usr.Permissao(u.Permissao),
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
		Email:        u.Email(),
		Permissao:    u.Permissao().String(),
		Status:       u.Status(),
		Avatar:       u.Avatar(),
		UltimoLogin:  u.UltimoLogin().Format(time.RFC3339),
		CriadoEm:     u.CriadoEm().Format(time.RFC3339),
		AtualizadoEm: u.AtualizadoEm().Format(time.RFC3339),
	}
}

// ParaUsuariosResp converte uma lista de modelos Usuario para uma lista de UsuarioResp.
func ParaUsuariosResp(usuarios []usr.Usuario) []UsuarioResp {
	return MapearSlice(usuarios, ParaUsuarioResp)
}

// Tecnico representa a estrutura de resposta para dados de técnico.
type TecnicoResp struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

// ParaTecnicoResp converte um modelo Usuario para TecnicoResp.
func ParaTecnicoResp(u usr.Usuario) TecnicoResp {
	return TecnicoResp{
		ID:   u.ID(),
		Nome: u.Nome(),
	}
}

// ParaTecnicosResp converte uma lista de modelos Usuario para uma lista de TecnicoResp.
func ParaTecnicosResp(tecnicos []usr.Usuario) []TecnicoResp {
	return MapearSlice(tecnicos, ParaTecnicoResp)
}

// UsuarioLdapResp representa a estrutura de resposta para dados de usuário LDAP.
type UsuarioLdapResp struct {
	Login string `json:"login"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// ParaUsuarioLdapResp converte dados para UsuarioLdapResp.
func ParaUsuarioLdapResp(usuarioLdap *auth.UsuarioLDAP) UsuarioLdapResp {
	return UsuarioLdapResp{
		Login: usuarioLdap.Login,
		Nome:  usuarioLdap.Nome,
		Email: strings.ToLower(usuarioLdap.Email),
	}
}
