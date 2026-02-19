package acompanhamento

import "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"

// Filtro representa os filtros para listar acompanhamentos.
type Filtro struct {
	paginacao domain.Paginacao
	chamadoID *string
	usuarioID *string
}

// NovoFiltro cria uma nova instância de Filtro com os dados fornecidos.
func NovoFiltro(paginacao domain.Paginacao, chamadoID, usuarioID *string) Filtro {
	return Filtro{
		paginacao: paginacao,
		chamadoID: chamadoID,
		usuarioID: usuarioID,
	}
}

// Paginacao retorna a estrutura de paginação do filtro.
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

// ChamadoID retorna o ID do chamado do filtro, se definido.
func (f *Filtro) ChamadoID() *string {
	return f.chamadoID
}

// UsuarioID retorna o ID do usuário do filtro, se definido.
func (f *Filtro) UsuarioID() *string {
	return f.usuarioID
}
