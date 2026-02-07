package log

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

// LogFiltro representa os critérios de filtragem para listar logs.
type LogFiltro struct {
	paginacao  domain.Paginacao
	busca      *string
	usuarioID  *string
	acao       *string
	entidade   *string
	dataInicio *time.Time
	dataFim    *time.Time
}

// NovoLogFiltro cria uma nova instância de LogFiltro com os dados fornecidos.
func NovoLogFiltro(
	paginacao domain.Paginacao,
	busca, usuarioID, acao, entidade *string,
	dataInicio, dataFim *time.Time,
) LogFiltro {
	return LogFiltro{
		paginacao:  paginacao,
		busca:      busca,
		usuarioID:  usuarioID,
		acao:       acao,
		entidade:   entidade,
		dataInicio: dataInicio,
		dataFim:    dataFim,
	}
}

// Paginacao retorna a estrutura de paginação do filtro.
func (f *LogFiltro) Paginacao() domain.Paginacao {
	return f.paginacao
}

// Pagina retorna a página atual do filtro de paginação.
func (f *LogFiltro) Pagina() int {
	return f.paginacao.Pagina()
}

// Limite retorna o limite de itens por página do filtro de paginação.
func (f *LogFiltro) Limite() int {
	return f.paginacao.Limite()
}

// Normalizar aplica valores padrão para paginação no filtro.
func (f *LogFiltro) Normalizar() {
	f.paginacao.Normalizar()
}

// Offset calcula o offset para queries SQL com base no filtro de paginação.
func (f *LogFiltro) Offset() int {
	return f.paginacao.Offset()
}

// SemLimite define paginação com limite alto para retornar todos os registros no filtro.
func (f *LogFiltro) SemLimite() {
	f.paginacao.SemLimite()
}

// Busca retorna o termo de busca do filtro, se definido.
func (f *LogFiltro) Busca() *string {
	return f.busca
}

// UsuarioID retorna o ID do usuário do filtro, se definido.
func (f *LogFiltro) UsuarioID() *string {
	return f.usuarioID
}

// Acao retorna a ação do filtro, se definida.
func (f *LogFiltro) Acao() *string {
	return f.acao
}

// Entidade retorna a entidade do filtro, se definida.
func (f *LogFiltro) Entidade() *string {
	return f.entidade
}

// DataInicio retorna a data de início do filtro, se definida.
func (f *LogFiltro) DataInicio() *time.Time {
	return f.dataInicio
}

// DataFim retorna a data de fim do filtro, se definida.
func (f *LogFiltro) DataFim() *time.Time {
	return f.dataFim
}
