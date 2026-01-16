package handler

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

var ErrStringParaTimeFormatoInvalido = errors.New("a data deve estar no formato YYYY-MM-DD")

// ParseIntQuery extrai um valor int dos query params.
// Retorna nil se não existir.
func ParseIntQuery(query url.Values, key string) (*int, error) {
	valStr := query.Get(key)
	if valStr == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return nil, fmt.Errorf("parse int query param %s: %w", key, err)
	}

	return &val, nil
}

// ParseBoolQuery extrai um valor bool dos query params.
// Retorna nil se não existir.
func ParseBoolQuery(query url.Values, key string) (*bool, error) {
	valStr := query.Get(key)
	if valStr == "" {
		return nil, nil
	}

	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return nil, fmt.Errorf("parse bool query param %s: %w", key, err)
	}

	return &val, nil
}

// ParseStringQuery extrai um valor string dos query params.
// Retorna nil se vazio ou inexistente.
func ParseStringQuery(query url.Values, key string) *string {
	val := query.Get(key)
	if val == "" {
		return nil
	}
	return &val
}

// ParseStringParaTimeQuery extrai uma data no formato YYYY-MM-DD dos query params.
// Retorna nil se não existir.
func ParseStringParaTimeQuery(query url.Values, key string) (*time.Time, error) {
	valStr := query.Get(key)
	if valStr == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", valStr)
	if err != nil {
		return nil, ErrStringParaTimeFormatoInvalido
	}

	return &t, nil
}

// ParsePaginacao extrai pagina e limite dos query params
func ParsePaginacaoQuery(query url.Values) (pagina, limite int, err error) {
	p, err := ParseIntQuery(query, "pagina")
	if err != nil {
		return 0, 0, fmt.Errorf("parse pagina query param: %w", err)
	}
	if p != nil {
		pagina = *p
	}

	l, err := ParseIntQuery(query, "limite")
	if err != nil {
		return 0, 0, fmt.Errorf("parse limite query param: %w", err)
	}
	if l != nil {
		limite = *l
	}

	return pagina, limite, nil
}
