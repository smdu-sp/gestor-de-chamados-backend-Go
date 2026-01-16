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
	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

const CTG = "CATEGORIA"

// CategoriaHandler gerencia as requisições HTTP relacionadas a categorias.
type CategoriaHandler struct {
	Svc    ctg.Service
	LogSvc log.Service
}

// NovoCategoriaHandler cria uma nova instância de CategoriaHandler.
func NovoCategoriaHandler(ctgSvc ctg.Service, logSvc log.Service) *CategoriaHandler {
	return &CategoriaHandler{
		Svc:    ctgSvc,
		LogSvc: logSvc,
	}
}

// Criar godoc
// @Summary Criar nova categoria
// @Description Cria uma nova categoria com os dados fornecidos no corpo da requisição
// @Tags Categoria
// @Accept json
// @Produce json
// @Param categoria body dto.CriarCategoriaReq true "Dados para criação de categoria"
// @Success 201 {object} dto.CategoriaResp "Categoria criada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "Categoria já existe"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /categorias [post]
func (h *CategoriaHandler) Criar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	var req dto.CriarCategoriaReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	ctgCriada, err := h.Svc.Criar(ctx, req.ParaCriarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrCategoriaJaExiste):
			ResponderJSONConflict(w, r, "já existe uma categoria com esse nome")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao criar categoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao criar categoria")
			return

		default:
			slog.Error("Erro interno ao criar categoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao criar categoria")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Criar, CTG, ctgCriada.String())

	ResponderJSONCreated(w, dto.ParaCategoriaResp(*ctgCriada))
}

// BuscarPorID godoc
// @Summary Buscar categoria por ID
// @Description Retorna uma categoria pelo ID
// @Tags Categoria
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} dto.CategoriaResp "Categoria encontrada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Categoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /categorias/buscar-por-id/{id} [get]
func (h *CategoriaHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	ctg, err := h.Svc.BuscarPorID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrCategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoria não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar categoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar categoria")
			return

		default:
			slog.Error("Erro interno ao buscar categoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar categoria")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaCategoriaResp(*ctg))
}

// BuscarPorNome godoc
// @Summary Buscar categoria por nome
// @Description Retorna uma categoria pelo nome
// @Tags Categoria
// @Accept json
// @Produce json
// @Param nome path string true "Nome da categoria"
// @Success 200 {object} dto.CategoriaResp "Categoria encontrada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Categoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /categorias/buscar-por-nome/{nome} [get]
func (h *CategoriaHandler) BuscarPorNome(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	nome := r.PathValue("nome")
	ctg, err := h.Svc.BuscarPorNome(ctx, nome)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrCategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoria não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar categoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar categoria")
			return

		default:
			slog.Error("Erro interno ao buscar categoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar categoria")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaCategoriaResp(*ctg))
}

// Atualizar godoc
// @Summary Atualizar categoria
// @Description Atualiza informações de uma categoria pelo ID
// @Tags Categoria
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Param categoria body dto.AtualizarCategoriaReq true "Dados para atualização da categoria"
// @Success 200 {object} dto.CategoriaResp "Categoria atualizada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Categoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "Categoria já existe"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /categorias/atualizar/{id} [patch]
func (h *CategoriaHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	var req dto.AtualizarCategoriaReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	ctgAtualizada, err := h.Svc.Atualizar(ctx, id, req.ParaAtualizarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrCategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoria não encontrada")
			return

		case errors.Is(err, mysql.ErrCategoriaJaExiste):
			ResponderJSONConflict(w, r, "já existe uma categoria com esse nome")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar categoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar categoria")
			return

		default:
			slog.Error("Erro interno ao atualizar categoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar categoria")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, CTG, ctgAtualizada.String())

	ResponderJSONStatusOK(w, dto.ParaCategoriaResp(*ctgAtualizada))
}

// Listar godoc
// @Summary Listar todas as categorias
// @Description Retorna lista de todas as categorias
// @Tags Categoria
// @Accept json
// @Produce json
// @Success 200 {object} []dto.CategoriaResp "Lista de categorias retornada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /categorias/listar [get]
func (h *CategoriaHandler) Listar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro := ctg.Filtro{}
	filtro.SemLimite()
	ctgSlice, _, _, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar categorias")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar categorias")
			return

		default:
			slog.Error("Erro interno ao listar categorias", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar categorias")
			return
		}
	}

	ResponderJSONStatusOK(w, dto.ParaCategoriasResp(ctgSlice))
}

// Desativar godoc
// @Summary Desativar categoria
// @Description Desativa uma categoria pelo ID
// @Tags Categoria
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} dto.CategoriaResp "Categoria desativada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Categoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /categorias/desativar/{id} [patch]
func (h *CategoriaHandler) Desativar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	ctgDesativada, err := h.Svc.Desativar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrCategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoria não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao desativar categoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao desativar categoria")
			return

		default:
			slog.Error("Erro interno ao desativar categoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao desativar categoria")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Desativar, CTG, fmt.Sprintf("ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaCategoriaResp(*ctgDesativada))
}

// Ativar godoc
// @Summary Ativar categoria
// @Description Ativa uma categoria pelo ID
// @Tags Categoria
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} dto.CategoriaResp "Categoria ativada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Categoria não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /categorias/ativar/{id} [patch]
func (h *CategoriaHandler) Ativar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	ctgAtivada, err := h.Svc.Ativar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrCategoriaNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoria não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao ativar categoria")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao ativar categoria")
			return

		default:
			slog.Error("Erro interno ao ativar categoria", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao ativar categoria")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Ativar, CTG, fmt.Sprintf("ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaCategoriaResp(*ctgAtivada))
}

// ListarPaginado godoc
// @Summary Listar categorias com paginação e filtros
// @Description Retorna uma lista paginada de categorias com base nos filtros fornecidos nos parâmetros da URL
// @Tags Categoria
// @Accept json
// @Produce json
// @Param busca query string false "Termo de busca no nome da categoria"
// @Param status query bool false "Filtrar por status (ativo/inativo)"
// @Param pagina query int false "Número da página (padrão: 1)"
// @Param limite query int false "Número de itens por página (padrão: 10)"
// @Success 200 {object} dto.CategoriaPaginado "Lista paginada de categorias retornada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /categorias/listar-paginado [get]
func (h *CategoriaHandler) ListarPaginado(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro, err := h.parseFiltros(r.URL.Query())
	if err != nil {
		slog.Error("Erro ao parsear filtros de consulta", "error", err)
		ResponderJSONBadRequest(w, r, "filtros de consulta inválidos")
		return
	}

	categorias, total, filtro, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar categorias")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar categorias")
			return

		default:
			slog.Error("Erro interno ao listar categorias", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar categorias")
			return
		}
	}

	resp := dto.ParaCategoriasResp(categorias)
	ResponderJSONPaginado(
		w,
		total,
		filtro.Pagina(),
		filtro.Limite(),
		resp,
	)
}

// parseFiltros extrai os filtros de categoria dos parâmetros da URL.
func (*CategoriaHandler) parseFiltros(query url.Values) (ctg.Filtro, error) {
	pagina, limite, err := ParsePaginacaoQuery(query)
	if err != nil {
		return ctg.Filtro{}, err
	}
	paginacao := dmn.NovoPaginacao(pagina, limite)
	busca := ParseStringQuery(query, "busca")

	status, err := ParseBoolQuery(query, "status")
	if err != nil {
		return ctg.Filtro{}, err
	}

	return ctg.NovoFiltro(paginacao, busca, status), nil
}
