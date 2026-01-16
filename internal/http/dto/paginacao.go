package dto

// RespostaPaginada representa a resposta de paginação para listagens
type RespPaginada[T any] struct {
	Total        int `json:"total"`
	Pagina       int `json:"pagina"`
	Limite       int `json:"limite"`
	TotalPaginas int `json:"totalPaginas"`
	Items        []T `json:"items"`
}

// NovoRespPaginada cria uma nova instância de RespPaginada
func NovoRespPaginada[T any](total, pagina, limite int, items []T) *RespPaginada[T] {
	totalPaginas := total / limite
	if total%limite != 0 {
		totalPaginas++
	}

	return &RespPaginada[T]{
		Total:        total,
		Pagina:       pagina,
		Limite:       limite,
		TotalPaginas: totalPaginas,
		Items:        items,
	}
}