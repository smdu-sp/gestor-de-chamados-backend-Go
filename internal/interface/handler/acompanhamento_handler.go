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
	entidadeAcompanhamento = "ACOMPANHAMENTO"
)

type AcompanhamentoHandler struct {
	Usecase    usecase.AcompanhamentoUsecase
	UsecaseLog usecase.LogUsecase
}

// NewAcompanhamentoHandler cria uma nova instância de AcompanhamentoHandler.
func NewAcompanhamentoHandler(usecase usecase.AcompanhamentoUsecase, usecaseLog usecase.LogUsecase) *AcompanhamentoHandler {
	return &AcompanhamentoHandler{
		Usecase:    usecase,
		UsecaseLog: usecaseLog,
	}
}

// Criar godoc
// @Summary      Cria um novo acompanhamento
// @Description	Cria um novo acompanhamento com os dados fornecidos no corpo da requisição.
// @Tags         Acompanhamentos
// @Accept       json
// @Produce      json
// @Param        acompanhamento  body      model.Acompanhamento  true  "Dados do acompanhamento"
// @Success      201  {object}  response.AcompanhamentoResponse
// @Failure			400  {object}  any
// @Failure			405  {object}  any
// @Failure			408  {object}  any
// @Failure			500  {object}  any
// @Router			/acompanhamentos/criar [post]
// Criar acompanhamento
func (h *AcompanhamentoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	var acompanhamento model.Acompanhamento
	if err := json.NewDecoder(r.Body).Decode(&acompanhamento); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.CriarAcompanhamento(ctx, &acompanhamento); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao criar acompanhamento", err.Error())
			return

		// erros do servidor - 500
		case errors.Is(err, utils.ErrUUIDv7Generation),
			errors.Is(err, repository.ErrRowsAffected),
			errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao criar acompanhamento", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao criar acompanhamento", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao criar acompanhamento", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao criar acompanhamento", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoCriar,
		entidadeAcompanhamento,
		fmt.Sprintf("Acompanhamento criado via API: %s", acompanhamento.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
	}

	response.JSON(w, http.StatusCreated, response.ToAcompanhamentoResponse(&acompanhamento))
}

// BuscarTudo godoc
// @Sumary Lista acompanhamentos com paginação e filtros
// @Description Retorna uma lista paginada de acompanhamentos
// @Tags Acompanhamentos
// @Accept json
// @Produce json
// @Param pagina query int false "Página"
// @Param limite query int false "Limite"
// @Param chamadoId query string false "ID do Chamado"
// @Param usuarioId query string false "ID do Usuário"
// @Success 200 {object} []model.Acompanhamento
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /acompanhamentos/buscar-tudo [get]
// BuscarTudo lista acompanhamentos com paginação e filtros
func (h *AcompanhamentoHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	query := r.URL.Query()

	filtro := model.AcompanhamentoFiltro{}

	if pagina, err := strconv.Atoi(query.Get("pagina")); err == nil {
		filtro.Pagina = pagina
	}
	if limite, err := strconv.Atoi(query.Get("limite")); err == nil {
		filtro.Limite = limite
	}

	if chamadoId := query.Get("chamadoId"); chamadoId != "" {
		filtro.ChamadoID = &chamadoId
	}

	if usuarioId := query.Get("usuarioId"); usuarioId != "" {
		filtro.UsuarioID = &usuarioId
	}

	items, total, filtroCorrigido, err := h.Usecase.ListarAcompanhamentos(ctx, filtro)
	if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerAcompanhamento),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao listar acompanhamentos", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao listar acompanhamentos", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar acompanhamentos", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar acompanhamentos", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.PageResponse[model.Acompanhamento]{
		Total:  total,
		Pagina: filtroCorrigido.Pagina,
		Limite: filtroCorrigido.Limite,
		Items:  items,
	})
}

// BuscarPorID godoc
// @Summary      Busca um acompanhamento pelo ID
// @Description	Busca um acompanhamento específico pelo ID
// @Tags         Acompanhamentos
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do acompanhamento"
// @Success      200  {object}  model.AcompanhamentoResponse
// @Failure			404  {object}  any
// @Failure			405  {object}  any
// @Failure			408  {object}  any
// @Failure			500  {object}  any
// @Router			/acompanhamentos/buscar-por-id/{id} [get]
// BuscarPorID busca um acompanhamento pelo ID
func (h *AcompanhamentoHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	acompanhamento, err := h.Usecase.BuscarAcompanhamentoPorID(ctx, id)
	if err != nil {
		switch {
		// recursos não encontrados - 404
		case errors.Is(err, model.ErrAcompanhamentoIDInvalido),
			errors.Is(err, repository.ErrAcompanhamentoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao buscar acompanhamento", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerAcompanhamento),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao buscar acompanhamento", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao buscar acompanhamento", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao buscar acompanhamento", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar acompanhamento", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.ToAcompanhamentoResponse(acompanhamento))
}

// Atualizar godoc
// @Summary      Atualiza um acompanhamento
// @Description	Atualiza dados do acompanhamento
// @Tags         Acompanhamentos
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do acompanhamento"
// @Param        acompanhamento body model.Acompanhamento true "Dados do acompanhamento"
// @Success      200  {object}  model.Acompanhamento
// @Failure			400  {object}  any
// @Failure			404  {object}  any
// @Failure			405  {object}  any
// @Failure			408  {object}  any
// @Failure			500  {object}  any
// @Router			/acompanhamentos/atualizar/{id} [put]
// Atualizar acompanhamento
func (h *AcompanhamentoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPut) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	var acompanhamento model.Acompanhamento
	if err := json.NewDecoder(r.Body).Decode(&acompanhamento); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.AtualizarAcompanhamento(ctx, id, &acompanhamento); err != nil {
		switch {
		// recursos não encontrados - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrAcompanhamentoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao atualizar acompanhamento", err.Error())
			return

		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao atualizar acompanhamento", err.Error())
			return

		// erros do servidor - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao atualizar acompanhamento", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao atualizar acompanhamento", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao atualizar acompanhamento", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar acompanhamento", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtualizar,
		entidadeAcompanhamento,
		fmt.Sprintf("Acompanhamento atualizado via API: %s", acompanhamento.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.ToAcompanhamentoResponse(&acompanhamento))
}

// Deletar godoc
// @Summary      Deleta um acompanhamento
// @Description	Deleta um acompanhamento pelo ID
// @Tags         Acompanhamentos
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do acompanhamento"
// @Success      200  {object}  map[string]any
// @Failure			404  {object}  any
// @Failure			405  {object}  any
// @Failure			408  {object}  any
// @Failure			500  {object}  any
// @Router			/acompanhamentos/deletar/{id} [delete]
// Deletar acompanhamento
func (h *AcompanhamentoHandler) Deletar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodDelete) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	if err := h.Usecase.DeletarAcompanhamento(ctx, id); err != nil {
		switch {
		// recursos não encontrados - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrAcompanhamentoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao deletar acompanhamento", err.Error())
			return

		// erros do servidor - 500
		case errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao deletar acompanhamento", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao deletar acompanhamento", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao deletar acompanhamento", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao deletar acompanhamento", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoDesativar,
		entidadeAcompanhamento,
		fmt.Sprintf("Acompanhamento deletado via API: ID %s", id),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "acompanhamento deletado com sucesso"})
}

// BuscarPorChamadoID godoc
// @Summary      Busca acompanhamentos pelo ID do chamado
// @Description	Busca todos os acompanhamentos associados a um chamado específico pelo ID do chamado
// @Tags         Acompanhamentos
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do chamado"
// @Success      200  {object}  []model.AcompanhamentoResponse
// @Failure			404  {object}  any
// @Failure			405  {object}  any
// @Failure			408  {object}  any
// @Failure			500  {object}  any
// @Router			/acompanhamentos/buscar-por-chamado-id/{id} [get]
// BuscarPorChamadoID busca acompanhamentos pelo ID do chamado
func (h *AcompanhamentoHandler) BuscarPorChamadoID(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	acompanhamentos, err := h.Usecase.BuscarAcompanhamentosPorChamadoID(ctx, id)
	if err != nil {
		switch {
		// recursos não encontrados - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrAcompanhamentoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao buscar acompanhamento", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerAcompanhamento),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao buscar acompanhamento", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao buscar acompanhamento", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao buscar acompanhamento", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar acompanhamento", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, acompanhamentos)
}