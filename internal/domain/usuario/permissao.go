package usuario

import (
	"errors"
	"fmt"
	"strings"
)

// Permissao é o tipo que representa as permissões de um usuário.
type Permissao string

const (
	PermADM Permissao = "ADM" // Administrador
	PermTEC Permissao = "TEC" // Técnico
	PermUSR Permissao = "USR" // Usuário
	PermDEV Permissao = "DEV" // Desenvolvedor
)

var permissoesMap = map[Permissao]struct{}{
	PermADM: {},
	PermTEC: {},
	PermUSR: {},
	PermDEV: {},
}

func ValidarPermissao(permissao Permissao) error {
	if permissao == "" {
		return errors.New("a permissão do usuário não pode ser vazia")
	}

	if _, ok := permissoesMap[permissao]; ok {
		return nil
	}

	return fmt.Errorf(
		"permissão inválida %q; válidas: %s",
		permissao,
		PermissoesValidas(),
	)
}

// PermissoesValidas retorna todas as permissões válidas formatadas.
func PermissoesValidas() string {
	perms := make([]string, 0, len(permissoesMap))
	for p := range permissoesMap {
		perms = append(perms, string(p))
	}
	return strings.Join(perms, ", ")
}


// PermTecPtr retorna um ponteiro para a representação em string da permissão TEC.
func PermTecPtr() *string {
	return PermTEC.StringPtr()
}

// String retorna a representação em string da permissão.
func (p Permissao) String() string {
	return string(p)
}

// StringPtr retorna um ponteiro para a representação em string da permissão.
func (p Permissao) StringPtr() *string {
	s := string(p)
	return &s
}
