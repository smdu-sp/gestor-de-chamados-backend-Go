package chamado

import (
	"errors"
	"fmt"
	"strings"
)

// StatusChamado é o tipo que representa o status de um chamado
type StatusChamado string

const (
	StatusAberto    StatusChamado = "ABERTO"
	StatusAtribuido StatusChamado = "ATRIBUIDO"
	StatusResolvido StatusChamado = "RESOLVIDO"
	StatusRejeitado StatusChamado = "REJEITADO"
	StatusFechado   StatusChamado = "FECHADO"
)

// Status válidos para um chamado
var statusMap = map[StatusChamado]struct{}{
	StatusAberto:    {},
	StatusAtribuido: {},
	StatusResolvido: {},
	StatusRejeitado: {},
	StatusFechado:   {},
}

// ValidarStatus recebe um status de chamado e verifica se ele é válido.
func ValidarStatus(status StatusChamado) error {
	if status == "" {
		return errors.New("o status do chamado não pode ser vazio")
	}

	if _, ok := statusMap[status]; ok {
		return nil
	}

	return fmt.Errorf(
		"o status do chamado %q é inválido, status aceitos: %s",
		status,
		StatusValidos(),
	)
}

// Métodos de conveniência para obter os status

func (s StatusChamado) Aberto() StatusChamado    { return StatusAberto }
func (s StatusChamado) Atribuido() StatusChamado { return StatusAtribuido }
func (s StatusChamado) Resolvido() StatusChamado { return StatusResolvido }
func (s StatusChamado) Rejeitado() StatusChamado { return StatusRejeitado }
func (s StatusChamado) Fechado() StatusChamado   { return StatusFechado }

// StatusValidos retorna todos os status válidos formatados
func StatusValidos() string {
	status := make([]string, 0, len(statusMap))
	for s := range statusMap {
		status = append(status, string(s))
	}
	return strings.Join(status, ", ")
}

// String retorna a representação em string do status do chamado.
func (s StatusChamado) String() string {
	return string(s)
}
