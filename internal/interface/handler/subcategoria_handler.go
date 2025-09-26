package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// SubcategoriaHandler gerencia as requisições HTTP relacionadas a subcategorias.
type SubcategoriaHandler struct {
	Usecase usecase.SubcategoriaUsecase
}

// Criar godoc
// @Summary Cria uma nova subcategoria
// @Description Cria uma nova subcategoria com os dados fornecidos no corpo da requisição.
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param subcategoria body model.Subcategoria true "Dados da subcategoria"
// @Success 201 {object} model.Subcategoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /subcategorias/criar [post]
// Criar subcategoria
func (h *SubcategoriaHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	var subcategoria model.Subcategoria
	if err := json.NewDecoder(r.Body).Decode(&subcategoria); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.CriarSubcategoria(ctx, &subcategoria); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.Is(err, model.ErrNomeInvalido):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao criar subcategoria", err.Error())

		// conflito - 409
		case errors.Is(err, repository.ErrSubcategoriaJaExiste):
			response.ErrorJSON(w, http.StatusConflict, "subcategoria já cadastrada", err.Error())

		// erros internos - 500
		case errors.Is(err, utils.ErrUUIDv7Generation),
			errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao criar subcategoria", err.Error())

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao criar subcategoria", err.Error())

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao criar subcategoria", err.Error())

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao criar subcategoria", err.Error())
		}
		return
	}

	response.JSON(w, http.StatusCreated, response.ToSubcategoriaResponse(&subcategoria))
}

// BuscarTudo godoc
// @Summary Lista todas as subcategorias
// @Description Retorna lista paginada de subcategorias
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param pagina query int false "Página"
// @Param limite query int false "Limite"
// @Param busca query string false "Busca"
// @Param status query bool false "Status"
// @Success 200 {object} []model.Subcategoria
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /subcategorias/buscar-tudo [get]
// BuscarTudo lista subcategorias com paginação e filtros.
func (h *SubcategoriaHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	query := r.URL.Query()

	filtro := model.SubcategoriaFiltro{}

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

	items, total, filtroCorrigido, err := h.Usecase.ListarSubcategorias(ctx, filtro)
	if err != nil {
		switch {
		// requisições inválidas - 400
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerSubcategoria),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusBadRequest, "erro interno ao listar subcategorias", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao listar categorias", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao listar categorias", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar categorias", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.PageResponse[model.Subcategoria]{
		Total:  total,
		Pagina: filtroCorrigido.Pagina,
		Limite: filtroCorrigido.Limite,
		Items:  items,
	})
}

// BuscarPorID godoc
// @Summary Busca uma subcategoria pelo ID
// @Description Retorna os detalhes de uma subcategoria pelo seu ID.
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param id path string true "ID da subcategoria"
// @Success 200 {object} model.Subcategoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /subcategorias/buscar-por-id/{id} [get]
// BuscarPorID busca uma subcategoria pelo ID.
func (h *SubcategoriaHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	subcategoria, err := h.Usecase.BuscarSubcategoriaPorID(ctx, id)
	if err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrSubcategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "id inválido ao buscar subcategoria", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrScannerSubcategoria):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao buscar subcategoria", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao buscar subcategoria", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao buscar subcategoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar subcategoria", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.ToSubcategoriaResponse(subcategoria))
}

// BuscarPorNome godoc
// @Summary Busca uma subcategoria pelo nome
// @Description Retorna os detalhes de uma subcategoria pelo seu nome.
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param nome path string true "Nome da subcategoria"
// @Success 200 {object} model.Subcategoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /subcategorias/buscar-por-nome/{nome} [get]
// BuscarPorNome busca uma subcategoria pelo nome.
func (h *SubcategoriaHandler) BuscarPorNome(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	nome := lastSegment(r.URL.Path)
	subcategoria, err := h.Usecase.BuscarSubcategoriaPorNome(ctx, nome)
	if err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, model.ErrNomeInvalido),
			errors.Is(err, repository.ErrSubcategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "nome inválido ao buscar subcategoria", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrScannerSubcategoria):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao buscar subcategoria", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao buscar subcategoria", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao buscar subcategoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar subcategoria", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.ToSubcategoriaResponse(subcategoria))
}

// Atualizar godoc
// @Summary Atualizar subcategoria
// @Description Atualiza os dados de uma subcategoria existente com os dados fornecidos no corpo da requisição.
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param id path string true "ID da subcategoria"
// @Param subcategoria body model.Subcategoria true "Subcategoria"
// @Success 200 {object} model.Subcategoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 409 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /subcategorias/atualizar/{id} [put]
// Atualizar atualiza uma subcategoria existente.
func (h *SubcategoriaHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPut) {
		return
	}
	// TODO (melhorar): permitir atualização parcial (PATCH)
	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	var subcategoria model.Subcategoria
	if err := json.NewDecoder(r.Body).Decode(&subcategoria); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.AtualizarSubcategoria(ctx, id, &subcategoria); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, model.ErrNomeInvalido):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao atualizar subcategoria", err.Error())

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrSubcategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "subcategoria não encontrada ao atualizar", err.Error())

		// conflito - 409
		case errors.Is(err, repository.ErrSubcategoriaJaExiste):
			response.ErrorJSON(w, http.StatusConflict, "subcategoria já cadastrada", err.Error())

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao atualizar subcategoria", err.Error())

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao atualizar subcategoria", err.Error())

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao atualizar subcategoria", err.Error())

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar subcategoria", err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "subcategoria atualizada com sucesso"})
}

// ListaCompleta godoc
// @Summary Lista todas as subcategorias sem paginação
// @Description Retorna uma lista completa de todas as subcategorias, sem paginação.
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Success 200 {array} model.Subcategoria
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /subcategorias/listar-completa [get]
// ListaCompleta lista todas as subcategorias sem paginação.
func (h *SubcategoriaHandler) ListaCompleta(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	filtro := model.SubcategoriaFiltro{
		Pagina: 1,
		Limite: 10000000,
	}

	items, _, _, err := h.Usecase.ListarSubcategorias(ctx, filtro)
	if err != nil {
		switch {
		// requisições inválidas - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerSubcategoria),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao listar subcategorias", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao listar subcategorias", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao listar subcategorias", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar subcategorias", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, items)
}

// Desativar godoc
// @Summary Desativar subcategoria
// @Description Desativa (soft delete) uma subcategoria pelo ID.
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param id path string true "ID da subcategoria"
// @Success 200 {object} model.Subcategoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /subcategorias/desativar/{id} [delete]
// Desativar desativa (soft delete) uma subcategoria.
func (h *SubcategoriaHandler) Desativar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodDelete) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	if err := h.Usecase.DesativarSubcategoria(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrSubcategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "id inválido ao desativar subcategoria", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao desativar subcategoria", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao desativar subcategoria", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao desativar subcategoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao desativar subcategoria", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.StatusSubcategoria{Ativo: false})
}

// Ativar godoc
// @Summary Ativar subcategoria
// @Description Ativa uma subcategoria pelo ID.
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param id path string true "ID da subcategoria"
// @Success 200 {object} model.Subcategoria
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /subcategorias/ativar/{id} [patch]
// Ativar ativa uma subcategoria.
func (h *SubcategoriaHandler) Ativar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPatch) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	if err := h.Usecase.AtivarSubcategoria(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrSubcategoriaNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "id inválido ao ativar subcategoria", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao ativar subcategoria", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao ativar subcategoria", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao ativar subcategoria", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao ativar subcategoria", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.StatusSubcategoria{Ativo: true})
}