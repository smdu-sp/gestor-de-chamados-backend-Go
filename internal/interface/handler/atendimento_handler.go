package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

const (
	entidadeAtendimento = "ATENDIMENTO"
)

type AtendimentoHandler struct {
	Usecase    usecase.AtendimentoUseCase
	UsecaseLog usecase.LogUsecase
}

func NewAtendimentoHandler(usecase usecase.AtendimentoUseCase, usecaseLog usecase.LogUsecase) *AtendimentoHandler {
	return &AtendimentoHandler{
		Usecase:    usecase,
		UsecaseLog: usecaseLog,
	}
}

// Criar godoc
// @Sumary Cria um novo atendimento
// @Description Cria um novo atendimento com os dados fornecidos no corpo da requisição
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param atendimento body model.Atendimento true "Dados do atendimento"
// @Success 201 {object} any
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /atendimentos/criar [post]
// Criar atendimento
func (h *AtendimentoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	var atendimento model.Atendimento
	if err := json.NewDecoder(r.Body).Decode(&atendimento); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.CriarAtendimento(ctx, &atendimento); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao criar atendimento", err.Error())
			return

		// erros do servidor - 500
		case errors.Is(err, utils.ErrUUIDv7Generation),
			errors.Is(err, repository.ErrRowsAffected),
			errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao criar atendimento", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao criar atendimento", err.Error())
			return
		
			// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao criar atendimento", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao criar atendimento", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoCriar,
		entidadeAtendimento,
		fmt.Sprintf("Atendimento criado via API: %s", atendimento.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
	}

	response.JSON(w, http.StatusCreated, response.ToAtendimentoResponse(&atendimento))
}

// BuscarPorID godoc
// @Sumary Busca um atendimento pelo ID
// @Description Busca um atendimento pelo ID fornecido na URL
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param id path string true "ID do atendimento"
// @Success 200 {object} response.AtendimentoResponse
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /atendimentos/buscar-por-id/{id} [get]
// BuscarPorID busca um atendimento pelo ID
func (h *AtendimentoHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	atendimento, err := h.Usecase.BuscarAtendimentoPorID(ctx, id)
	if err != nil {
		switch {
		// recursos não encontrados - 404
		case errors.Is(err, model.ErrAtendimentoIDInvalido),
			errors.Is(err, repository.ErrAtendimentoNaoEncontrado):
			response.ErrorJSON(w, http.StatusBadRequest, "ID do atendimento inválido", err.Error())
			return

		// erros do servidor - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerAtendimento),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao buscar atendimento", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao buscar atendimento", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao buscar atendimento", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar atendimento", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, response.ToAtendimentoResponse(atendimento))
}

// Atualizar godoc
// @Sumary Atualiza um atendimento existente
// @Description Atualiza um atendimento existente com os dados fornecidos no corpo da requisição
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param id path string true "ID do atendimento"
// @Param atendimento body model.Atendimento true "Dados do atendimento"
// @Success 204 {object} any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /atendimentos/atualizar/{id} [put]
// Atualizar atualiza um atendimento existente
func (h *AtendimentoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPut) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	var atendimento model.Atendimento
	if err := json.NewDecoder(r.Body).Decode(&atendimento); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.AtualizarAtendimento(ctx, id, &atendimento); err != nil {
		switch {
		// recursos não encontrados - 404
		case errors.Is(err, model.ErrAtendimentoIDInvalido),
			errors.Is(err, repository.ErrAtendimentoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID do atendimento inválido ao atualizar", err.Error())
			return

		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao atualizar atendimento", err.Error())
			return

		// erros do servidor - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrQueryContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao atualizar atendimento", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao atualizar atendimento", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao atualizar atendimento", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar atendimento", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtualizar,
		entidadeAtendimento,
		fmt.Sprintf("Atendimento atualizado via API: %s", atendimento.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.ToAtendimentoResponse(&atendimento))
}

// Listar godoc
// @Sumary Lista atendimentos com filtros opcionais
// @Description Lista atendimentos com paginação e filtros opcionais para ID do chamado e ID do atribuído
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param limite query int false "Número de itens por página" default(10) maximum(100)
// @Param chamadoId query string false "ID do chamado para filtrar"
// @Param atribuidoId query string false "ID do atribuído para filtrar"
// @Success 200 {array} response.AtendimentoResponse
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /atendimentos/listar [get]
// Listar lista atendimentos com filtros opcionais
func (h *AtendimentoHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	query := r.URL.Query()

	filtro := model.AtendimentoFiltro{}

	if pagina, err := strconv.Atoi(query.Get("pagina")); err == nil {
		filtro.Pagina = pagina
	}

	if limite, err := strconv.Atoi(query.Get("limite")); err == nil {
		filtro.Limite = limite
	}

	if chamadoID := query.Get("chamadoId"); chamadoID != "" {
		filtro.ChamadoID = &chamadoID
	}

	if atribuidoID := query.Get("atribuidoId"); atribuidoID != "" {
		filtro.AtribuidoID = &atribuidoID
	}

	items, total, filtroCorrigido, err := h.Usecase.ListarAtendimentos(ctx, filtro)
		if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerAtendimento),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao listar atendimentos", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao listar atendimentos", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar atendimentos", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar atendimentos", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.PageResponse[model.Atendimento]{
		Total: total,
		Pagina: filtroCorrigido.Pagina,
		Limite: filtroCorrigido.Limite,
		Items: items,
	})
}