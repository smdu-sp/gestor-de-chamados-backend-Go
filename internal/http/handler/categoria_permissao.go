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
	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

const CPM = "CATEGORIA_PERMISSAO"

// CategoriaPermissaoHandler gerencia as requisições HTTP relacionadas a categorias e permissões.
type CategoriaPermissaoHandler struct {
	Svc    cpm.Service
	LogSvc log.Service
}

// NovoCategoriaPermissaoHandler cria uma nova instância de CategoriaPermissaoHandler.
func NovoCategoriaPermissaoHandler(cpmSvc cpm.Service, logSvc log.Service) *CategoriaPermissaoHandler {
	return &CategoriaPermissaoHandler{
		Svc:    cpmSvc,
		LogSvc: logSvc,
	}
}

// Criar godoc
// @Summary Cria uma nova categoria e permissão
// @Description Adiciona uma nova categoria e permissão ao sistema
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param categoria_permissao body dto.CriarCategoriaPermissaoReq true "Dados da nova categoria e permissão"
// @Success 201 {object} dto.CategoriaPermissaoResp "CategoriaPermissao criada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "CategoriaPermissao já existe"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao criar usuário"
// @Router /categoria-permissoes [post]
func (h *CategoriaPermissaoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	var req dto.CriarCategoriaPermissaoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	cpmCriada, err := h.Svc.Criar(ctx, req.ParaCriarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrCategoriaPermissaoJaExiste):
			ResponderJSONConflict(w, r, "já existe uma categoriaPermissao com os dados fornecidos")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao criar categoriaPermissao")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao criar categoriaPermissao")
			return

		default:
			slog.Error("Erro interno ao criar categoriaPermissao", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao criar categoriaPermissao")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Criar, CPM, cpmCriada.String())

	ResponderJSONCreated(w, dto.ParaCategoriaPermissaoResp(*cpmCriada))
}

// BuscarPorID godoc
// @Summary Busca uma categoria e permissão por ID
// @Description Recupera os detalhes de uma categoria e permissão específica usando seu ID
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param categoriaId path string true "ID da categoria"
// @Param usuarioId path string true "ID do usuário"
// @Success 200 {object} dto.CategoriaPermissaoResp "CategoriaPermissao encontrada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "CategoriaPermissao não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao buscar categoriaPermissao"
// @Router /categorias-permissoes/{categoriaId}/usuarios/{usuarioId} [get]
func (h *CategoriaPermissaoHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	usuarioID := r.PathValue("usuarioId")
	cpm, err := h.Svc.BuscarPorID(ctx, id, usuarioID)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrCategoriaPermissaoNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoriaPermissao não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar categoriaPermissao")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar categoriaPermissao")
			return

		default:
			slog.Error("Erro interno ao buscar categoriaPermissao por id", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar categoriaPermissao")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaCategoriaPermissaoResp(*cpm))
}

// Atualizar godoc
// @Summary Atualiza uma categoria e permissão existente
// @Description Modifica os detalhes de uma categoria e permissão específica
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param categoriaId path string true "ID da categoria"
// @Param usuarioId path string true "ID do usuário"
// @Param categoriaPermissao body dto.AtualizarCategoriaPermissaoReq true "Dados para atualizar a categoria e permissão"
// @Success 200 {object} dto.CategoriaPermissaoResp "CategoriaPermissao atualizada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "CategoriaPermissao não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "CategoriaPermissao já existe"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao atualizar categoriaPermissao"
// @Router /categorias-permissoes/atualizar/{categoriaId}/usuarios/{usuarioId} [patch]
func (h *CategoriaPermissaoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	usuarioID := r.PathValue("usuarioId")
	var req dto.AtualizarCategoriaPermissaoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	cpmAtualizada, err := h.Svc.Atualizar(ctx, id, usuarioID, req.ParaAtualizarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrCategoriaPermissaoNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoriaPermissao não encontrada")
			return

		case errors.Is(err, mysql.ErrCategoriaPermissaoJaExiste):
			ResponderJSONConflict(w, r, "já existe uma categoriaPermissao com os dados fornecidos")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar categoriaPermissao")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar categoriaPermissao")
			return

		default:
			slog.Error("Erro interno ao atualizar categoriaPermissao", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar categoriaPermissao")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, CPM, cpmAtualizada.String())

	ResponderJSONStatusOK(w, dto.ParaCategoriaPermissaoResp(*cpmAtualizada))
}

// Deletar godoc
// @Summary Deleta uma categoria e permissão
// @Description Remove uma categoria e permissão específica do sistema
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param categoriaId path string true "ID da categoria"
// @Param usuarioId path string true "ID do usuário"
// @Success 204 "CategoriaPermissao deletada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "CategoriaPermissao não encontrada"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao deletar categoriaPermissao"
// @Router /categorias-permissoes/deletar/{categoriaId}/usuarios/{usuarioId} [delete]
func (h *CategoriaPermissaoHandler) Deletar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	usuarioID := r.PathValue("usuarioId")
	if err := h.Svc.Deletar(ctx, id, usuarioID); err != nil {
		switch {
		case errors.Is(err, mysql.ErrCategoriaPermissaoNaoEncontrada):
			ResponderJSONNotFound(w, r, "categoriaPermissao não encontrada")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao deletar categoriaPermissao")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao deletar categoriaPermissao")
			return

		default:
			slog.Error("Erro interno ao deletar categoriaPermissao", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao deletar categoriaPermissao")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Deletar, CPM, fmt.Sprintf("categoriaID=%s, usuarioID=%s", id, usuarioID))

	ResponderJSONNoContent(w)
}

// ListarPaginado godoc
// @Summary Lista categorias e permissões com paginação
// @Description Recupera uma lista paginada de categorias e permissões com base nos filtros fornecidos
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param limite query int false "Número de itens por página" default(10)
// @Param categoriaId query string false "Filtra por ID da categoria"
// @Param usuarioId query string false "Filtra por ID do usuário"
// @Param permissao query string false "Filtra por permissão"
// @Success 200 {object} dto.CategoriaPermissaoPaginado "Lista paginada de categorias e permissões retornada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao listar categorias e permissões"
// @Router /categorias-permissoes/listar-paginado [get]
func (h *CategoriaPermissaoHandler) ListarPaginado(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro, err := h.parseFiltros(r.URL.Query())
	if err != nil {
		slog.Error("Erro ao parsear filtros de consulta", "error", err)
		ResponderJSONBadRequest(w, r, "filtros de consulta inválidos")
		return
	}

	categoriasPermissoes, total, filtro, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar categoria permissao")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar categoria permissao")
			return

		default:
			slog.Error("Erro interno ao listar categoria permissao", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar categoria permissao")
			return
		}
	}

	resp := dto.ParaCategoriasPermissoesResp(categoriasPermissoes)
	ResponderJSONPaginado(
		w,
		total,
		filtro.Pagina(),
		filtro.Limite(),
		resp,
	)
}

// parseFiltrosCategoriaPermissao extrai os filtros da query string e retorna um Filtro.
func (*CategoriaPermissaoHandler) parseFiltros(query url.Values) (cpm.Filtro, error) {
	pagina, limite, err := ParsePaginacaoQuery(query)
	if err != nil {
		return cpm.Filtro{}, err
	}
	paginacao := dmn.NovoPaginacao(pagina, limite)
	categoriaID := ParseStringQuery(query, "categoriaId")
	usuarioID := ParseStringQuery(query, "usuarioId")
	permissao := ParseStringQuery(query, "permissao")

	return cpm.NovoFiltro(paginacao, categoriaID, usuarioID, permissao), nil
}
