package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// LogHandler lida com requisições relacionadas a logs.
type LogHandler struct {
	Usecase usecase.LogUsecase
}

// NewLogHandler cria uma nova instância de LogHandler.
func NewLogHandler(usecase usecase.LogUsecase) *LogHandler {
	return &LogHandler{Usecase: usecase}
}

// BuscarTudo godoc
// @Summary      Lista todos os logs com paginação e filtros
// @Description  Retorna lista paginada de logs.
// @Tags         Logs
// @Accept       json
// @Produce      json
// @Param pagina query int false "Página"
// @Param limite query int false "Limite"
// @Param busca query string false "Busca"
// @Param usuario_id query string false "ID do usuário"
// @Param acao query string false "Ação"
// @Param entidade query string false "Entidade"
// @Param data_inicio query string false "Data de início (formato: YYYY-MM-DD)"
// @Param data_fim query string false "Data de fim (formato: YYYY-MM-DD)"
// @Success      200  {object}  []model.Log
// @Failure      400  {object}  any
// @Failure      405  {object}  any
// @Failure      408  {object}  any
// @Failure      500  {object}  any
// @Router       /logs [get]
// BuscarTudo lista todos os logs com paginação e filtros.
func (h *LogHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	query := r.URL.Query()

	filtro := model.LogFiltro{}

	if pagina, err := strconv.Atoi(query.Get("pagina")); err == nil {
		filtro.Pagina = pagina
	}

	if limite, err := strconv.Atoi(query.Get("limite")); err == nil {
		filtro.Limite = limite
	}

	if busca := query.Get("busca"); busca != "" {
		filtro.Busca = &busca
	}

	if usuarioID := query.Get("usuario_id"); usuarioID != "" {
		filtro.UsuarioID = &usuarioID
	}

	if acao := query.Get("acao"); acao != "" {
		filtro.Acao = &acao
	}

	if entidade := query.Get("entidade"); entidade != "" {
		filtro.Entidade = &entidade
	}

	if dataInicioStr := query.Get("data_inicio"); dataInicioStr != "" {
		dataInicioTime, err := utils.StringParaTime(&dataInicioStr)
		if err != nil {
			response.ErrorJSON(w, http.StatusBadRequest, "erro ao converter data de início", err.Error())
			return
		}
		filtro.DataInicio = dataInicioTime
	}

	if dataFimStr := query.Get("data_fim"); dataFimStr != "" {
		dataFimTime, err := utils.StringParaTime(&dataFimStr)
		if err != nil {
			response.ErrorJSON(w, http.StatusBadRequest, "erro ao converter data de fim", err.Error())
			return
		}
		filtro.DataFim = dataFimTime
	}

	logs, total, filtroCorrigido, err := h.Usecase.ListarLogs(ctx, filtro)
	if err != nil {
		switch {
		// Erros do repositório - status 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerLog),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao listar logs", err.Error())
			return

		// Erros de contexto - status 408 ou 400
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo limite excedido ao listar logs", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar logs", err.Error())
			return

		// Erros desconhecidos - status 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar logs", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.PageResponse[model.Log]{
		Total:  total,
		Pagina: filtroCorrigido.Pagina,
		Limite: filtroCorrigido.Limite,
		Items:  logs,
	})
}

// BuscarPorID godoc
// @Summary      Busca um log pelo ID
// @Description  Retorna um log específico pelo seu ID.
// @Tags         Logs
// @Accept       json
// @Produce      json
// @Param id path string true "ID do Log"
// @Success      200  {object}  model.Log
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      405  {object}  any
// @Failure      408  {object}  any
// @Failure      500  {object}  any
// @Router       /logs/{id} [get]
// BuscarPorID busca um log pelo seu ID.
func (h *LogHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	log, err := h.Usecase.BuscarLogPorID(ctx, id)
	if err != nil {
		switch {
		// erro de validação - status 400
		case errors.Is(err, model.ErrIDInvalido):
			response.ErrorJSON(w, http.StatusBadRequest, "ID inválido ao buscar log", err.Error())
			return

		// recurso não encontrado - status 404
		case errors.Is(err, repository.ErrLogNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao buscar log", err.Error())
			return

			// erros do repositório - status 500
			case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerLog):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao buscar log", err.Error())
			return

			// erro de contexto - status 408
			case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo limite excedido ao buscar log", err.Error())
			return

			// erro de contexto - status 400
			case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao buscar log", err.Error())
			return

			// erros desconhecidos - status 500
			default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar log", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.ToLogResponse(log))
}