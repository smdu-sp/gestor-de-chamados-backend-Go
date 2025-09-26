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

// ChamadoHandler gerencia as requisições HTTP relacionadas a chamados.
type ChamadoHandler struct {
	Usecase usecase.ChamadoUsecase
}

// Criar godoc
// @Summary Cria um novo chamado
// @Description Cria um novo chamado com os dados fornecidos no corpo da requisição.
// @Tags chamados
// @Accept json
// @Produce json
// @Param chamado body model.Chamado true "Dados do chamado"
// @Success 201 {object} model.Chamado
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/criar [post]
// Cria chamado
func (h *ChamadoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	var chamado model.Chamado
	if err := json.NewDecoder(r.Body).Decode(&chamado); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.CriarChamado(ctx, &chamado); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao criar chamado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, utils.ErrUUIDv7Generation),
			errors.Is(err, repository.ErrRowsAffected),
			errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao criar chamado", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao criar chamado", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao criar chamado", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao criar chamado", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusCreated, 	response.ToChamadoResponse(&chamado))
}

// BuscarTudo godoc
// @Summary Lista chamados com paginação e filtros
// @Description Retorna lista paginada de chamados.
// @Tags chamados
// @Accept json
// @Produce json
// @Param pagina query int false "Pagina"
// @Param limite query int false "Limite"
// @Param busca query string false "Busca"
// @Param status query string false "Status do chamado"
// @Param categoriaId query string false "ID da categoria"
// @Param subcategoriaId query string false "ID da subcategoria"
// @Param criadorId query string false "ID do criador"
// @Param atribuidoId query string false "ID do atribuído"
// @Success 200 {object} []model.Chamado
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/buscar-tudo [get]
// BuscarTudo lista todos os chamados com paginação e filtros
func (h *ChamadoHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	query := r.URL.Query()

	filtro := model.ChamadoFiltro{}

	if pagina, err := strconv.Atoi(query.Get("pagina")); err == nil {
		filtro.Pagina = pagina
	}
	if limite, err := strconv.Atoi(query.Get("limite")); err == nil {
		filtro.Limite = limite
	}

	if busca := query.Get("busca"); busca != "" {
		filtro.Busca = &busca
	}
	if status := query.Get("status"); status != "" {
		filtro.Status = &status
	}
	if categoriaID := query.Get("categoriaId"); categoriaID != "" {
		filtro.CategoriaID = &categoriaID
	}
	if subcategoriaID := query.Get("subcategoriaId"); subcategoriaID != "" {
		filtro.SubcategoriaID = &subcategoriaID
	}
	if criadorID := query.Get("criadorId"); criadorID != "" {
		filtro.CriadorID = &criadorID
	}
	if atribuidoID := query.Get("atribuidoId"); atribuidoID != "" {
		filtro.AtribuidoID = &atribuidoID
	}

	items, total, filtroCorrigido, err := h.Usecase.ListarChamados(ctx, filtro)
	if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrScannerChamado),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao listar chamados", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao listar chamados", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao listar chamados", err.Error())
			return

			// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar chamados", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, response.PageResponse[model.Chamado]{
		Total:  total,
		Pagina: filtroCorrigido.Pagina,
		Limite: filtroCorrigido.Limite,
		Items:  items,
	})
}

// BuscarPorID godoc
// @Summary Busca um chamado por ID
// @Description Retorna chamado pelo ID.
// @Tags chamados
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Success 200 {object} model.Chamado
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/buscar-por-id/{id} [get]
// BuscarPorID busca um chamado pelo seu ID
func (h *ChamadoHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	chamado, err := h.Usecase.BuscarChamadoPorID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrChamadoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao buscar chamado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrScannerChamado):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao buscar chamado", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao buscar chamado", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao buscar chamado", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar chamado", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, response.ToChamadoResponse(chamado))
}

// Atualizar godoc
// @Summary Atualiza um chamado existente
// @Description Atualiza os dados de um chamado existente pelo ID.
// @Tags chamados
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Param chamado body model.Chamado true "Dados do chamado"
// @Success 200 {object} model.Chamado
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/atualizar/{id} [put]
// Atualizar atualiza chamado por ID
func (h *ChamadoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPut) {
		return
	}
	// TODO (melhorar): permitir atualização parcial (PATCH)
	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	var chamado model.Chamado
	if err := json.NewDecoder(r.Body).Decode(&chamado); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.AtualizarChamado(ctx, id, &chamado); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao atualizar chamado", err.Error())
			return

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrChamadoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao atualizar chamado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrQueryContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao atualizar chamado", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao atualizar chamado", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao atualizar chamado", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar chamado", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, response.ToChamadoResponse(&chamado))
}

// Arquivar godoc
// @Summary Arquiva um chamado existente
// @Description Arquiva um chamado existente pelo ID.
// @Tags chamados
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Success 200 {object} map[string]any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/arquivar/{id} [patch]
// Arquivar arquiva chamado por ID
func (h *ChamadoHandler) Arquivar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodDelete) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	if err := h.Usecase.ArquivarChamado(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrChamadoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao arquivar chamado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrQueryContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao arquivar chamado", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao arquivar chamado", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao arquivar chamado", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao arquivar chamado", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, response.ChamadoStatus{Status: "Arquivado"})
}

// Desarquivar godoc
// @Summary Desarquiva um chamado existente
// @Description Desarquiva um chamado existente pelo ID.
// @Tags chamados
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Success 200 {object} map[string]any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/desarquivar/{id} [patch]
// Desarquivar desarquiva chamado por ID
func (h *ChamadoHandler) Desarquivar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPatch) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	if err := h.Usecase.DesarquivarChamado(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrChamadoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao desarquivar chamado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrQueryContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao desarquivar chamado", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao desarquivar chamado", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao desarquivar chamado", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao desarquivar chamado", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.ChamadoStatus{Status: "Aberto"})
}

// AtualizarStatus godoc
// @Summary Atualiza o status de um chamado existente
// @Description Atualiza o status de um chamado existente pelo ID.
// @Tags chamados
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Param chamado body object true "Status e solução do chamado" { "status": "string", "solucao": "string (opcional)" }
// @Success 200 {object} map[string]any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/atualizar-status/{id} [patch]
// AtualizarStatus atualiza o status do chamado por ID
func (h *ChamadoHandler) AtualizarStatus(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPatch) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	var requisicao struct {
		Status  string  `json:"status"`
		Solucao *string `json:"solucao,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requisicao); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.AtualizarStatusChamado(ctx, id, requisicao.Status, requisicao.Solucao); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao atualizar status do chamado", err.Error())
			return

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrChamadoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao atualizar status do chamado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrQueryContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao atualizar status do chamado", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao atualizar status do chamado", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao atualizar status do chamado", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar status do chamado", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, response.ChamadoStatus{Status: requisicao.Status})
}

// AtribuirTecnico godoc
// @Summary Atribui um técnico a um chamado existente
// @Description Atribui um técnico a um chamado existente pelo ID.
// @Tags chamados
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Param tecnico body object true "ID do técnico" { "tecnicoId": "string" }
// @Success 200 {object} map[string]any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/atribuir-tecnico/{id} [patch]
// AtribuirTecnico atribui um técnico ao chamado por ID
func (h *ChamadoHandler) AtribuirTecnico(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPatch) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	var requisicao struct {
		TecnicoID string `json:"tecnicoId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requisicao); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.AtribuirTecnicoChamado(ctx, id, requisicao.TecnicoID); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrChamadoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao atribuir técnico ao chamado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrQueryContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao atribuir técnico ao chamado", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao atribuir técnico ao chamado", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao atribuir técnico ao chamado", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atribuir técnico ao chamado", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, response.TecnicoAtribuido{Atribuido: true})
}

// RemoverTecnicoChamado godoc
// @Summary Remove o técnico atribuído de um chamado existente
// @Description Remove o técnico atribuído de um chamado existente pelo ID.
// @Tags chamados
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Success 200 {object} map[string]any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/remover-tecnico/{id} [delete]
// RemoverTecnicoChamado remove o técnico atribuído do chamado por ID
func (h *ChamadoHandler) RemoverTecnicoChamado(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodDelete) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	if err := h.Usecase.RemoverTecnicoChamado(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrChamadoNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao remover técnico do chamado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrQueryContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao remover técnico do chamado", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao remover técnico do chamado", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao remover técnico do chamado", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao remover técnico do chamado", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, response.TecnicoAtribuido{Atribuido: false})
}

// ListaCompleta godoc
// @Summary Retorna todos os chamados sem paginação
// @Description Retorna todos os chamados sem paginação, útil para relatórios ou exportação de dados.
// @Tags chamados
// @Accept json
// @Produce json
// @Success 200 {object} []model.Chamado
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Failure 504 {object} any
// @Router /chamados/lista-completa [get]
// ListaCompleta retorna todos os chamados sem paginação
func (h *ChamadoHandler) ListaCompleta(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	filtro := model.ChamadoFiltro{
		Pagina: 1,
		Limite: 10000000, 
	}
	items, _, _, err := h.Usecase.ListarChamados(ctx, filtro)
	if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrScannerChamado),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro ao listar chamados", err.Error())
			return

		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusGatewayTimeout, "tempo de requisição excedido ao listar chamados", err.Error())
			return

		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao listar chamados", err.Error())
			return

			// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar chamados", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, items)
}