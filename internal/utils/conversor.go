package utils

import (
	"fmt"
	"time"
)

var (
	// ErrDataInvalida indica que a data fornecida é inválida.
	ErrDataInvalida = fmt.Errorf("a data deve estar no formato YYYY-MM-DD")


	// ErrTimeParse indica um erro de parsing.
	ErrTimeParse = fmt.Errorf("erro de parsing ao converter valor para time.Time")
)

// StringParaTime converte uma string no formato "YYYY-MM-DD" para um objeto time.Time.
func StringParaTime(dataStr *string) (*time.Time, error) {
	if dataStr == nil || *dataStr == "" {
		return nil, NewAppError(
			"[utils.StringParaTime]",
			LevelInfo,
			"data string é nula, vazia ou inválida",
			ErrDataInvalida,
		)
	}
	
	data, err := time.Parse("2006-01-02", *dataStr) // "2006-01-02" é o layout oficial do Go para "YYYY-MM-DD"
	if err != nil {
		return nil, NewAppError(
			"[utils.StringParaTime]",
			LevelError,
			"erro ao converter string para data",
			fmt.Errorf(FmtErroWrap, ErrTimeParse, err),
		)
	}
	return &data, nil
}