package subcategoria

import "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"

// Filtro representa os critérios de filtro para listar subcategorias.
type Filtro struct {
	paginacao domain.Paginacao
	busca     *string
	status    *bool
}

// NovoFiltro cria uma nova instância de Filtro com os dados fornecidos.
func NovoFiltro(paginacao domain.Paginacao, busca *string, status *bool) Filtro {
	return Filtro{
		paginacao: paginacao,
		busca:     busca,
		status:    status,
	}
}

// Paginacao retorna a paginação do filtro.
func (f *Filtro) Paginacao() domain.Paginacao {
	return f.paginacao
}

// Pagina retorna a página atual do filtro de paginação.
func (f *Filtro) Pagina() int {
	return f.paginacao.Pagina()
}

// Limite retorna o limite de itens por página do filtro de paginação.
func (f *Filtro) Limite() int {
	return f.paginacao.Limite()
}

// Normalizar aplica valores padrão para paginação no filtro.
func (f *Filtro) Normalizar() {
	f.paginacao.Normalizar()
}

// Offset calcula o offset para queries SQL com base no filtro de paginação.
func (f *Filtro) Offset() int {
	return f.paginacao.Offset()
}

// SemLimite define paginação com limite alto para retornar todos os registros no filtro.
func (f *Filtro) SemLimite() {
	f.paginacao.SemLimite()
}

// Busca retorna o termo de busca do filtro, se definido.
func (f *Filtro) Busca() *string {
	return f.busca
}

// Status retorna o status do filtro, se definido.
func (f *Filtro) Status() *bool {
	return f.status
}