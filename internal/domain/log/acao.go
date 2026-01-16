package log

import (
	"fmt"
	"strings"
)

// Acao é o tipo que representa uma ação realizada no sistema.
type Acao string

const (
	Criar       Acao = "CRIAR"
	Atualizar   Acao = "ATUALIZAR"
	Desativar   Acao = "DESATIVAR"
	Ativar      Acao = "ATIVAR"
	Arquivar    Acao = "ARQUIVAR"
	Desarquivar Acao = "DESARQUIVAR"
	Deletar     Acao = "DELETAR"
)

// Ações válidas
var acoesValidas = map[Acao]struct{}{
	Criar:       {},
	Atualizar:   {},
	Desativar:   {},
	Ativar:      {},
	Arquivar:    {},
	Desarquivar: {},
	Deletar:     {},
}

// ValidarAcao recebe uma ação e valida se ela é válida.
func ValidarAcao(acao Acao) error {
	if acao == "" {
		return fmt.Errorf("a ação não pode ser vazia")
	}

	if _, ok := acoesValidas[acao]; ok {
		return nil
	}

	return fmt.Errorf(
		"a ação %q é inválida, ações aceitas: %s",
		acao,
		AcoesValidas(),
	)
}

// AcoesValidas retorna todas as ações válidas formatadas.
func AcoesValidas() string {
	acoes := make([]string, 0, len(acoesValidas))
	for a := range acoesValidas {
		acoes = append(acoes, string(a))
	}
	return strings.Join(acoes, ", ")
}

// String retorna a string da ação.
func (a Acao) String() string {
	return string(a)
}
