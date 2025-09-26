package utils

import (
	"strings"
)

// ValidacaoErrors representa uma coleção de erros de validação
type ValidacaoErrors []error

// Error implementa a interface error
func (ve ValidacaoErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	var msgs []string
	for _, err := range ve {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// Add permite adicionar um novo erro na lista
func (ve *ValidacaoErrors) Add(err error) {
	if err != nil {
		*ve = append(*ve, err)
	}
}

// HasErrors retorna true se houver erros acumulados
func (ve ValidacaoErrors) HasErrors() bool {
	return len(ve) > 0
}

// Is permite que errors.Is reconheça ValidacaoErrors
func (ve ValidacaoErrors) Is(target error) bool {
	_, ok := target.(ValidacaoErrors)
	return ok
}

// As permite que errors.As funcione com ValidacaoErrors
func (ve ValidacaoErrors) As(target any) bool {
	vp, ok := target.(*ValidacaoErrors)
	if !ok {
		return false
	}
	*vp = ve
	return true
}
