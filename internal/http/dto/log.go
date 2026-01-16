package dto

import (
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
)

// LogPaginado é um alias para uma resposta paginada de logs,
// usado para facilitar a leitura do código e documentação Swagger.
// @Description Estrutura paginada contendo uma lista de logs.
type LogPaginado RespPaginada[LogResp]

// LogResposta representa a estrutura de resposta para logs.
type LogResp struct {
	ID        string `json:"id"`
	UsuarioID string `json:"usuarioid"`
	Acao      string `json:"acao"`
	Entidade  string `json:"entidade"`
	Detalhes  string `json:"detalhes"`
	CriadoEm  string `json:"criadoEm"`
}

// ParaLogResp converte um modelo Log para LogResp.
func ParaLogResp(l log.Log) LogResp {
	return LogResp{
		ID:        l.ID(),
		UsuarioID: l.UsuarioID(),
		Acao:      l.Acao(),
		Entidade:  l.Entidade(),
		Detalhes:  l.Detalhes(),
		CriadoEm:  l.CriadoEm().Format(time.RFC3339),
	}
}

// ParaLogsResp converte uma lista de modelos Log para uma lista de LogResp.
func ParaLogsResp(logs []log.Log) []LogResp {
	return MapearSlice(logs, ParaLogResp)
}
