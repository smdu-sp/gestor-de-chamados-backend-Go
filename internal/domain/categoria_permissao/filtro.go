package categoria_permissao

import "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"

// Filtro representa os filtros para listar categorias de permissão.
type Filtro struct {
	paginacao   domain.Paginacao
	categoriaID *string
	usuarioID   *string
	permissao   *string
}

// NovoFiltro cria uma nova instância de Filtro com os dados fornecidos.
func NovoFiltro(paginacao domain.Paginacao, categoriaID, usuarioID, permissao *string) Filtro {
	return Filtro{
		paginacao:   paginacao,
		categoriaID: categoriaID,
		usuarioID:   usuarioID,
		permissao:   permissao,
	}
}

// Paginacao retorna o filtro de paginação do Filtro.
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

// CategoriaID retorna o valor do filtro categoriaID, se definido.
func (f *Filtro) CategoriaID() *string {
	return f.categoriaID
}

// UsuarioID retorna o valor do filtro usuarioID, se definido.
func (f *Filtro) UsuarioID() *string {
	return f.usuarioID
}

// Permissao retorna o valor do filtro permissao, se definido.
func (f *Filtro) Permissao() *string {
	return f.permissao
}