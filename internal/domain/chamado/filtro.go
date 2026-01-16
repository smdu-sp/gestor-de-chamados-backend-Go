package chamado

import "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"

// Filtro representa os filtros possíveis para buscar chamados
type Filtro struct {
	paginacao      domain.Paginacao
	busca          *string
	status         *string
	categoriaID    *string
	subcategoriaID *string
	criadorID      *string
}

// NovoFiltro cria uma nova instância de Filtro com os dados fornecidos.
func NovoFiltro(paginacao domain.Paginacao, busca, status, categoriaID, subcategoriaID, criadorID *string) Filtro {
	return Filtro{
		paginacao:      paginacao,
		busca:          busca,
		status:         status,
		categoriaID:    categoriaID,
		subcategoriaID: subcategoriaID,
		criadorID:      criadorID,
	}
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
func (f *Filtro) Status() *string {
	return f.status
}

// CategoriaID retorna o ID da categoria do filtro, se definido.
func (f *Filtro) CategoriaID() *string {
	return f.categoriaID
}

// SubcategoriaID retorna o ID da subcategoria do filtro, se definido.
func (f *Filtro) SubcategoriaID() *string {
	return f.subcategoriaID
}

// CriadorID retorna o ID do criador do filtro, se definido.
func (f *Filtro) CriadorID() *string {
	return f.criadorID
}