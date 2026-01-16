package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

// scanner é uma interface que abstrai o método Scan de sql.Row e sql.Rows.
// Isso permite que funções genéricas trabalhem com ambos os tipos.
type scanner interface {
	Scan(dest ...any) error
}

// scanFunc é um tipo para funções que escaneiam dados de um scanner para um tipo T.
// Ela é usada para que funções genéricas possam receber diferentes funções de escaneamento.
type scanFunc[T any] func(scanner) (*T, error)

//==============================================================================
// Funções auxiliares para manipulação de linhas do banco de dados
//==============================================================================

// obterTotalRegistros recebe um contexto e uma conexão com o banco de dados,
// e retorna o total de registros encontrados na última consulta SQL que utilizou SQL_CALC_FOUND_ROWS.
func obterTotalRegistros(ctx context.Context, db *sql.DB) (int, error) {
	var total int
	if err := db.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&total); err != nil {
		return 0, fmt.Errorf("obter total de registros: %w", err)
	}
	return total, nil
}

// scanRows recebe um conjunto de linhas e uma função de escaneamento,
// e retorna uma fatia de itens do tipo T escaneados a partir das linhas.
func scanRows[T any](rows *sql.Rows, scan scanFunc[T]) ([]T, error) {
	var items []T

	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanRows: %w", err)
		}
		if item != nil { // evita nil caso o scanFunc retorne nil, nil (como no ErrNoRows)
			items = append(items, *item)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanRows: %w", err)
	}

	return items, nil
}

