package dto

// MapearSlice percorre cada elemento do slice "items",
// aplica a função "mapper" a ele,
// e devolve um novo slice contendo os resultados.
//
// É como um "for" que transforma um tipo em outro.
func MapearSlice[T any, R any](items []T, mapper func(T) R) []R {

	// 1 - Criamos um novo slice com o mesmo tamanho do original.
	resultados := make([]R, len(items))

	// 2 - Percorremos todos os itens.
	for indice := range items {

		// 3 - Pegamos o item da posição atual.
		itemAtual := items[indice]

		// 4 - Aplicamos a função de conversão.
		resultado := mapper(itemAtual)

		// 5 - Salvamos o resultado no novo slice.
		resultados[indice] = resultado
	}

	// 6 - Devolvemos o slice preenchido.
	return resultados
}
