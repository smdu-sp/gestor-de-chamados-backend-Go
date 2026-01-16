package usuario

import "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"

// Filtro representa os filtros para listar usuários.
type Filtro struct {
	paginacao domain.Paginacao
	busca     *string
	status    *bool
	permissao *string
}

// NovoFiltro cria uma nova instância de Filtro com os dados fornecidos.
func NovoFiltro(paginacao domain.Paginacao, busca *string, status *bool, permissao *string) Filtro {
	return Filtro{
		paginacao: paginacao,
		busca:     busca,
		status:    status,
		permissao: permissao,
	}
}

// Pagina retorna a página atual do filtro de paginação.
func (f *Filtro) Pagina() int {
	return f.paginacao.Pagina()
}

// GetLimite retorna o limite de itens por página do filtro de paginação.
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

// Permissao retorna a permissão do filtro, se definida.
func (f *Filtro) Permissao() *string {
	return f.permissao
}

// DefinirPermissao define a permissão no filtro.
func (f *Filtro) DefinirPermissao(permissao *string) {
	f.permissao = permissao
}