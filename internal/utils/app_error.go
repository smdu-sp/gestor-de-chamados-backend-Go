package utils

import (
	"fmt"
)

// ErrorLevel representa o nível do erro (pode ser string simples ou const enum)
type ErrorLevel string

const (
	LevelError   ErrorLevel = "ERRO"
	LevelWarning ErrorLevel = "WARN"
	LevelInfo    ErrorLevel = "INFO"
)

// App é um erro com estrutura rica, útil para logging e análise
type AppError struct {
	Method  string     // Nome do método onde o erro ocorreu
	Level   ErrorLevel // ERRO, WARN, etc.
	Message string     // Mensagem legível e padronizada
	Cause   error      // Erro original (encapsulado)
}

// Error implementa a interface `error`
func (e AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s: %v", e.Method, e.Level, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Method, e.Level, e.Message)
}

// Unwrap permite `errors.Unwrap`, `errors.Is`, `errors.As`
func (e AppError) Unwrap() error {
	return e.Cause
}

// NewAppError cria um erro estruturado novo
func NewAppError(method string, level ErrorLevel, message string, cause error) error {
	return AppError{
		Method:  method,
		Level:   level,
		Message: message,
		Cause:   cause,
	}
}
