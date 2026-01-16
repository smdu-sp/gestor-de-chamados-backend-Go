package atendimento

import "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"

// Filtro representa os filtros para listar atendimentos.
type Filtro struct {
	paginacao   domain.Paginacao
	chamadoID   *string
	atribuidoID *string
}

// NovoFiltro cria uma nova instância de Filtro com os dados fornecidos.
func NovoFiltro(paginacao domain.Paginacao, chamadoID, atribuidoID *string) Filtro {
	return Filtro{
		paginacao:   paginacao,
		chamadoID:   chamadoID,
		atribuidoID: atribuidoID,
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

// ChamadoID retorna o ID do chamado do filtro, se definido.
func (f *Filtro) ChamadoID() *string {
	return f.chamadoID
}

// AtribuidoID retorna o ID do tecnico atribuido, se definido.
func (f *Filtro) AtribuidoID() *string {
	return f.atribuidoID
}