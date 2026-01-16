package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	subc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

const SUBC = "SUBCATEGORIA"

// SubcategoriaHandler gerencia as requisições HTTP relacionadas a subcategorias.
type SubcategoriaHandler struct {
	Svc    subc.Service
	LogSvc log.Service
}

// NovoSubcategoriaHandler cria uma nova instância de SubcategoriaHandler.
func NovoSubcategoriaHandler(subcSvc subc.Service, logSvc log.Service) *SubcategoriaHandler {
	return &SubcategoriaHandler{
		Svc:    subcSvc,
		LogSvc: logSvc,
	}
}

// Criar godoc
// @Summary Cria uma nova subcategoria
// @Description Cria uma nova subcategoria com os dados fornecidos no corpo da requisição
// @Tags Subcategoria
// @Accept json
// @Produce json
// @Param subcategoria body dto.CriarSubcategoriaReq true "Subcategoria"
// @Success 201 {object} dto.SubcategoriaResp "Subcategoria criada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "Conflito ao criar subcategoria"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao criar subcategoria"
// @Router /subcategorias [post]
func (h *SubcategoriaHandler) Criar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	var req dto.CriarSubcategoriaReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	subcCriada, err := h.Svc.Criar(ctx, req.ParaCriarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrCategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoria não encontrada")
			return

		case errors.Is(err, mysql.ErrSubcategoriaJaExisteComNome):
			ResponderJSONConflict(w, r, "já existe uma subcategoria com o mesmo nome")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao criar subcategoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao criar subcategoria")
			return

		default:
			slog.Error("Erro interno ao criar subcategoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao criar subcategoria")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Criar, SUBC, subcCriada.String())

	ResponderJSONCreated(w, dto.ParaSubcategoriaResp(*subcCriada))
}

// BuscarPorID godoc
// @Summary Busca uma subcategoria pelo ID
// @Description Retorna os dados de uma subcategoria pelo seu ID
// @Tags Subcategoria
// @Accept json
// @Produce json
// @Param id path string true "ID da subcategoria"
// @Success 200 {object} dto.SubcategoriaResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Subcategoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /subcategorias/buscar-por-id/{id} [get]
func (h *SubcategoriaHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	subc, err := h.Svc.BuscarPorID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrSubcategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "subcategoria não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar subcategoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar subcategoria")
			return

		default:
			slog.Error("Erro interno ao buscar subcategoria por id", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar subcategoria")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaSubcategoriaResp(*subc))
}

// BuscarPorNome godoc
// @Summary Busca uma subcategoria pelo nome
// @Description Retorna os dados de uma subcategoria pelo seu nome
// @Tags Subcategoria
// @Accept json
// @Produce json
// @Param nome path string true "Nome da subcategoria"
// @Success 200 {object} dto.SubcategoriaResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Subcategoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /subcategorias/buscar-por-nome/{nome} [get]
func (h *SubcategoriaHandler) BuscarPorNome(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	nome := r.PathValue("nome")
	subc, err := h.Svc.BuscarPorNome(ctx, nome)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrSubcategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "subcategoria não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar subcategoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar subcategoria")
			return

		default:
			slog.Error("Erro interno ao buscar subcategoria por nome", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar subcategoria")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaSubcategoriaResp(*subc))
}

// Atualizar godoc
// @Summary Atualiza uma subcategoria
// @Description Atualiza informações de uma subcategoria pelo seu ID
// @Tags Subcategoria
// @Accept json
// @Produce json
// @Param id path string true "ID da subcategoria"
// @Param subcategoria body dto.AtualizarSubcategoriaReq true "Dados para atualização da subcategoria"
// @Success 200 {object} dto.SubcategoriaResp "Subcategoria atualizada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Subcategoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "Conflito ao atualizar subcategoria"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /subcategorias/atualizar/{id} [patch]
func (h *SubcategoriaHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	var req dto.AtualizarSubcategoriaReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	subcAtualizada, err := h.Svc.Atualizar(ctx, id, req.ParaAtualizarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrSubcategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "subcategoria não encontrada")
			return

		case errors.Is(err, mysql.ErrSubcategoriaJaExisteComNome):
			ResponderJSONConflict(w, r, "já existe uma subcategoria com o mesmo nome")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar subcategoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar subcategoria")
			return

		default:
			slog.Error("Erro interno ao atualizar subcategoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar subcategoria")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, SUBC, subcAtualizada.String())

	ResponderJSONStatusOK(w, dto.ParaSubcategoriaResp(*subcAtualizada))
}

// Listar godoc
// @Summary Lista todas as subcategorias
// @Description Retorna todas as subcategorias sem paginação
// @Tags Subcategoria
// @Accept json
// @Produce json
// @Success 200 {object} []dto.SubcategoriaResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /subcategorias/listar [get]
func (h *SubcategoriaHandler) Listar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro := subc.Filtro{}
	filtro.SemLimite()
	items, _, _, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar subcategorias")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar subcategorias")
			return

		default:
			slog.Error("Erro interno ao listar subcategorias", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar subcategorias")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaSubcategoriasResp(items))
}

// Desativar godoc
// @Summary Desativar subcategoria
// @Description Desativa uma subcategoria pelo ID
// @Tags Subcategoria
// @Accept json
// @Produce json
// @Param id path string true "ID da subcategoria"
// @Success 200 {object} dto.SubcategoriaResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta " subcategoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /subcategorias/desativar/{id} [patch]
func (h *SubcategoriaHandler) Desativar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	subcDesativada, err := h.Svc.Desativar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrSubcategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "subcategoria não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao desativar subcategoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao desativar subcategoria")
			return

		default:
			slog.Error("Erro interno ao desativar subcategoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao desativar subcategoria")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Desativar, SUBC, fmt.Sprintf("ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaSubcategoriaResp(*subcDesativada))
}

// Ativar godoc
// @Summary Ativar subcategoria
// @Description Ativa uma subcategoria pelo ID
// @Tags Subcategoria
// @Accept json
// @Produce json
// @Param id path string true "ID da subcategoria"
// @Success 200 {object} dto.SubcategoriaResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta " subcategoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /subcategorias/ativar/{id} [patch]
func (h *SubcategoriaHandler) Ativar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	subcAtivada, err := h.Svc.Ativar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrSubcategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "subcategoria não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao ativar subcategoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao ativar subcategoria")
			return

		default:
			slog.Error("Erro interno ao ativar subcategoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao ativar subcategoria")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Ativar, SUBC, fmt.Sprintf("ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaSubcategoriaResp(*subcAtivada))
}

// ListarPaginado godoc
// @Summary Lista subcategorias com paginação
// @Description Retorna uma lista paginada de subcategorias com base nos parâmetros de filtro fornecidos
// @Tags Subcategoria
// @Accept json
// @Produce json
// @Param pagina query int false "Número da página (padrão: 1)"
// @Param limite query int false "Número de itens por página (padrão: 10)"
// @Param busca query string false "Termo de busca no nome da subcategoria"
// @Param status query boolean false "Filtrar por status ativo/inativo"
// @Success 200 {object} dto.SubcategoriaPaginado "Lista paginada de subcategorias"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /subcategorias/listar-paginado [get]
func (h *SubcategoriaHandler) ListarPaginado(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro, err := h.parseFiltros(r.URL.Query())
	if err != nil {
		slog.Error("Erro ao parsear filtros de consulta", "error", err)
		ResponderJSONBadRequest(w, r, "filtros de consulta inválidos")
		return
	}

	subcategorias, total, filtro, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar subcategorias")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar subcategorias")
			return

		default:
			slog.Error("Erro interno ao listar subcategorias", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar subcategorias")
			return
		}
	}

	resp := dto.ParaSubcategoriasResp(subcategorias)
	ResponderJSONPaginado(
		w,
		total,
		filtro.Pagina(),
		filtro.Limite(),
		resp,
	)
}

// parseFiltros extrai os parâmetros de filtro da query string.
func (*SubcategoriaHandler) parseFiltros(query url.Values) (subc.Filtro, error) {
	pagina, limite, err := ParsePaginacaoQuery(query)
	if err != nil {
		return subc.Filtro{}, err
	}
	paginacao := dmn.NovoPaginacao(pagina, limite)
	busca := ParseStringQuery(query, "busca")

	status, err := ParseBoolQuery(query, "status")
	if err != nil {
		return subc.Filtro{}, err
	}

	return subc.NovoFiltro(paginacao, busca, status), nil
}
