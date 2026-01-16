package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// LogHandler lida com requisições relacionadas a logs.
type LogHandler struct {
	Svc log.Service
}

// NovoLogHandler cria uma nova instância de LogHandler.
func NovoLogHandler(logSvc log.Service) *LogHandler {
	return &LogHandler{Svc: logSvc}
}

// BuscarPorID godoc
// @Summary Busca log por ID
// @Description Retorna os dados completos de um log pelo seu ID
// @Tags Log
// @Accept json
// @Produce json
// @Param id path string true "ID do log"
// @Success 200 {object} dto.LogResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Log não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /logs/buscar-por-id/{id} [get]
func (h *LogHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	log, err := h.Svc.BuscarPorID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrLogNaoEncontrado):
			ResponderJSONNotFound(w, r, "log não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar log")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar log")
			return

		default:
			slog.Error("erro interno ao buscar log", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar log")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaLogResp(*log))
}

// BuscarTudo godoc
// @Summary Lista logs com paginação e filtros
// @Description Retorna uma lista paginada de logs com base nos filtros fornecidos
// @Tags Log
// @Accept json
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param limite query int false "Número de itens por página" default(10)
// @Param busca query string false "Termo de busca"
// @Param usuarioId query string false "ID do usuário"
// @Param acao query string false "Ação realizada"
// @Param entidade query string false "Entidade afetada"
// @Param dataInicio query string false "Data de início no formato YYYY-MM-DD"
// @Param dataFim query string false "Data de fim no formato YYYY-MM-DD"
// @Success 200 {object} dto.LogPaginado "Lista paginada de logs"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Parâmetros de consulta inválidos"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /logs/listar-paginado [get]
func (h *LogHandler) ListarPaginado(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro, err := h.parseFiltros(r.URL.Query())
	if err != nil {
		if errors.Is(err, ErrStringParaTimeFormatoInvalido) {
			ResponderJSONBadRequest(w, r, "formato de data inválido, use YYYY-MM-DD")
			return
		}
		slog.Error("Erro ao processar filtros de log", "error", err)
		ResponderJSONInternalError(w, r, "erro ao processar filtros de log")
		return
	}

	logs, total, filtro, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar logs")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar logs")
			return

		default:
			slog.Error("Erro interno ao listar logs", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar logs")
			return
		}
	}
	
	resp := dto.ParaLogsResp(logs)
	ResponderJSONPaginado(
		w,
		total,
		filtro.Pagina(),
		filtro.Limite(),
		resp,
	)
}

// parseFiltros extrai os filtros de log dos parâmetros da query.
func (*LogHandler) parseFiltros(query url.Values) (log.LogFiltro, error) {
	pagina, limite, err := ParsePaginacaoQuery(query)
	if err != nil {
		return log.LogFiltro{}, err
	}
	paginacao := dmn.NovoPaginacao(pagina, limite)
	busca := ParseStringQuery(query, "busca")
	usuarioID := ParseStringQuery(query, "usuarioId")
	acao := ParseStringQuery(query, "acao")
	entidade := ParseStringQuery(query, "entidade")

	dataInicio, err := ParseStringParaTimeQuery(query, "dataInicio")
	if err != nil {
		return log.LogFiltro{}, err
	}

	dataFim, err := ParseStringParaTimeQuery(query, "dataFim")
	if err != nil {
		return log.LogFiltro{}, err
	}

	return log.NovoLogFiltro(
		paginacao,
		busca,
		usuarioID,
		acao,
		entidade,
		dataInicio,
		dataFim,
	), nil
}
