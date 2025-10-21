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
	entidadeCategoria = "CATEGORIA"
)

// CategoriaHandler gerencia as requisições HTTP relacionadas a categorias.
type CategoriaHandler struct {
	Usecase    usecase.CategoriaUsecase
	UsecaseLog usecase.LogUsecase
}

// NewCategoriaHandler cria uma nova instância de CategoriaHandler.
func NewCategoriaHandler(usecase usecase.CategoriaUsecase, usecaseLog usecase.LogUsecase) *CategoriaHandler {
	return &CategoriaHandler{
		Usecase:    usecase,
		UsecaseLog: usecaseLog,
	}
}

// Criar godoc
// @Summary Criar uma nova categoria
// @Description Cria uma nova categoria com dados fornecidos no corpo da requisição.
// @Tags Categorias
// @Accept json
// @Produce json
// @Param categoria body model.Categoria true "Categoria"
// @Success 201 {object} model.Categoria
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 409 {object} any
// @Failure 500 {object} any
// @Router /categorias/criar [post]
// Criar categoria
func (h *CategoriaHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	var categoria model.Categoria
	if err := json.NewDecoder(r.Body).Decode(&categoria); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.CriarCategoria(ctx, &categoria); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao criar categoria", err.Error())
			return

		// conflitos - 409
		case errors.Is(err, repository.ErrCategoriaJaExiste):
			response.ErrorJSON(w, http.StatusConflict, "categoria já cadastrada", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, utils.ErrUUIDv7Generation),
			errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao criar categoria", err.Error())
			return

		// erro de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao criar categoria", err.Error())
			return

		// erro de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao criar categoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao criar categoria", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoCriar,
		entidadeCategoria,
		fmt.Sprintf("Categoria criada via API: %s", categoria.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, response.ToCategoriaResponse(&categoria))
}

// BuscarTudo godoc
// @Summary Listar todas as categorias
// @Description Retorna lista paginada de categorias
// @Tags Categorias
// @Accept json
// @Produce json
// @Param pagina query int false "Página"
// @Param limite query int false "Limite"
// @Param busca query string false "Busca"
// @Param status query bool false "Status"
// @Success 200 {object} []model.Categoria
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /categorias/buscar-tudo [get]
// BuscarTudo lista categorias com paginação e filtros.
func (h *CategoriaHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	query := r.URL.Query()

	filtro := model.CategoriaFiltro{}

	if pagina, err := strconv.Atoi(query.Get("pagina")); err == nil {
		filtro.Pagina = pagina
	}
	if limite, err := strconv.Atoi(query.Get("limite")); err == nil {
		filtro.Limite = limite
	}
	if busca := query.Get("busca"); busca != "" {
		filtro.Busca = &busca
	}
	if statusStr := query.Get("status"); statusStr != "" {
		if status, err := strconv.ParseBool(statusStr); err == nil {
			filtro.Status = &status
		}
	}

	items, total, filtroCorrigido, err := h.Usecase.ListarCategorias(ctx, filtro)
	if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerCategoria),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao listar categorias", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao listar categorias", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar categorias", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar categorias", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.PageResponse[model.Categoria]{
		Total:  total,
		Pagina: filtroCorrigido.Pagina,
		Limite: filtroCorrigido.Limite,
		Items:  items,
	})
}

// BuscarPorID godoc
// @Summary Buscar categoria por ID
// @Description Retorna uma categoria pelo ID
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} model.Categoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /categorias/buscar-por-id/{id} [get]
// BuscarPorID busca uma categoria pelo ID.
func (h *CategoriaHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	categoria, err := h.Usecase.BuscarCategoriaPorID(ctx, id)
	if err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrCategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao buscar categoria", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerCategoria):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao buscar categoria", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao buscar categoria", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao buscar categoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar categoria", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.ToCategoriaResponse(categoria))
}

// BuscarPorNome godoc
// @Summary Buscar categoria por nome
// @Description Retorna uma categoria pelo nome
// @Tags Categorias
// @Accept json
// @Produce json
// @Param nome path string true "Nome da categoria"
// @Success 200 {object} model.Categoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /categorias/buscar-por-nome/{nome} [get]
// BuscarPorNome busca uma categoria pelo nome.
func (h *CategoriaHandler) BuscarPorNome(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	nome := lastSegment(r.URL.Path)
	categoria, err := h.Usecase.BuscarCategoriaPorNome(ctx, nome)
	if err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, model.ErrNomeInvalido),
			errors.Is(err, repository.ErrCategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "nome inválido ao buscar categoria", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerCategoria):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao buscar categoria", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao buscar categoria", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao buscar categoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar categoria", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.ToCategoriaResponse(categoria))
}

// Atualizar godoc
// @Summary Atualizar categoria
// @Description Atualiza uma categoria existente pelo ID com os dados fornecidos no corpo da requisição.
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Param categoria body model.Categoria true "Categoria"
// @Success 200 {object} map[string]string
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 409 {object} any
// @Failure 500 {object} any
// @Router /categorias/atualizar/{id} [put]
// Atualizar atualiza uma categoria existente.
func (h *CategoriaHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPut) {
		return
	}

	// TODO (melhorar): permitir atualização parcial (PATCH)
	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	var categoria model.Categoria
	if err := json.NewDecoder(r.Body).Decode(&categoria); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.AtualizarCategoria(ctx, id, &categoria); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.Is(err, model.ErrCategoriaIDInvalido),
			errors.Is(err, model.ErrNomeInvalido):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao atualizar categoria", err.Error())

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrCategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao atualizar categoria", err.Error())

		// conflitos - 409
		case errors.Is(err, repository.ErrCategoriaJaExiste):
			response.ErrorJSON(w, http.StatusConflict, "categoria já cadastrada", err.Error())

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao atualizar categoria", err.Error())

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao atualizar categoria", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao atualizar categoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar categoria", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtualizar,
		entidadeCategoria,
		fmt.Sprintf("Categoria atualizada via API: %s", categoria.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "categoria atualizada com sucesso"})
}

// ListaCompleta godoc
// @Summary Listar todas as categorias (completa)
// @Description Retorna lista completa de categorias sem paginação
// @Tags Categorias
// @Accept json
// @Produce json
// @Success 200 {array} model.Categoria
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /categorias/listar-completa [get]
// ListaCompleta lista todas as categorias sem paginação.
func (h *CategoriaHandler) ListaCompleta(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	filtro := model.CategoriaFiltro{
		Pagina: 1,
		Limite: 10000000,
	}
	items, _, _, err := h.Usecase.ListarCategorias(ctx, filtro)
	if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerCategoria),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao listar categorias", err.Error())
			return

		// erro de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao listar categorias", err.Error())
			return

		// erro de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar categorias", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar categorias", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, items)
}

// Desativar godoc
// @Summary Desativar categoria
// @Description Desativa (soft delete) uma categoria pelo ID.
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} model.Categoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /categorias/desativar/{id} [delete]
// Desativar desativa (soft delete) uma categoria.
func (h *CategoriaHandler) Desativar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodDelete) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	if err := h.Usecase.DesativarCategoria(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrCategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao desativar categoria", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao desativar categoria", err.Error())

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao desativar categoria", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao desativar categoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao desativar categoria", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoDesativar,
		entidadeCategoria,
		fmt.Sprintf("Categoria desativada via API: categoria ID(%s)", id),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.StatusCategoria{Ativo: false})
}

// Ativar godoc
// @Summary Ativar categoria
// @Description Ativa uma categoria pelo ID.
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} model.Categoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /categorias/ativar/{id} [patch]
// Ativar ativa uma categoria.
func (h *CategoriaHandler) Ativar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPatch) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	if err := h.Usecase.AtivarCategoria(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrCategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao ativar categoria", err.Error())

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao ativar categoria", err.Error())

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao ativar categoria", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao ativar categoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao ativar categoria", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtivar,
		entidadeCategoria,
		fmt.Sprintf("Categoria ativada via API: categoria ID(%s)", id),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.StatusCategoria{Ativo: true})
}
