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
	acp "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/acompanhamento"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

const ACP = "ACOMPANHAMENTO"

// AcompanhamentoHandler lida com as requisições HTTP relacionadas a acompanhamentos.
type AcompanhamentoHandler struct {
	Svc    acp.Service
	LogSvc log.Service
}

// NovoAcompanhamentoHandler cria uma nova instância de AcompanhamentoHandler.
func NovoAcompanhamentoHandler(acpSvc acp.Service, logSvc log.Service) *AcompanhamentoHandler {
	return &AcompanhamentoHandler{
		Svc:    acpSvc,
		LogSvc: logSvc,
	}
}

// Criar godoc
// @Summary Cria um novo acompanhamento
// @Description Cria um acompanhamento com os dados fornecidos no corpo da requisição
// @Tags Acompanhamento
// @Accept json
// @Produce json
// @Param acompanhamento body dto.CriarAcompanhamentoReq true "Dados para criação de acompanhamento"
// @Success 201 {object} dto.AcompanhamentoResp "Acompanhamento criado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao criar acompanhamento"
// @Router /acompanhamentos [post]
func (h *AcompanhamentoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	var req dto.CriarAcompanhamentoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	acpCriado, err := h.Svc.Criar(ctx, req.ParaCriarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao criar acompanhamento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao criar acompanhamento")
			return

		default:
			slog.Error("Erro interno ao criar acompanhamento", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao criar acompanhamento")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Criar, ACP, acpCriado.String())

	ResponderJSONCreated(w, dto.ParaAcompanhamentoResp(*acpCriado))
}

// BuscarPorID godoc
// @Summary Busca acompanhamento por ID
// @Description Retorna os dados completos de um acompanhamento pelo seu ID
// @Tags Acompanhamento
// @Accept json
// @Produce json
// @Param id path string true "ID do acompanhamento"
// @Success 200 {object} dto.AcompanhamentoResp "Acompanhamento encontrado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Acompanhamento não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /acompanhamentos/buscar-por-id/{id} [get]
func (h *AcompanhamentoHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	acp, err := h.Svc.BuscarPorID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrAcompanhamentoNaoEncontrado):
			ResponderJSONNotFound(w, r, "acompanhamento não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar acompanhamento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar acompanhamento")
			return

		default:
			slog.Error("Erro interno ao buscar acompanhamento", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar acompanhamento por ID")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaAcompanhamentoResp(*acp))
}

// Atualizar godoc
// @Summary Atualiza acompanhamento
// @Description Atualiza informações de um acompanhamento pelo ID
// @Tags Acompanhamento
// @Accept json
// @Produce json
// @Param id path string true "ID do acompanhamento"
// @Param acompanhamento body dto.AtualizarAcompanhamentoReq true "Dados para atualização do acompanhamento"
// @Success 200 {object} dto.AcompanhamentoResp "Acompanhamento atualizado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Acompanhamento não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao atualizar acompanhamento"
// @Router /acompanhamentos/atualizar/{id} [patch]
func (h *AcompanhamentoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	var req dto.AtualizarAcompanhamentoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	acpAtualizado, err := h.Svc.Atualizar(ctx, id, req.ParaAtualizarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrAcompanhamentoNaoEncontrado):
			ResponderJSONNotFound(w, r, "acompanhamento não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar acompanhamento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar acompanhamento")
			return

		default:
			slog.Error("Erro interno ao atualizar acompanhamento", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar acompanhamento")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, ACP, acpAtualizado.String())

	ResponderJSONStatusOK(w, dto.ParaAcompanhamentoResp(*acpAtualizado))
}

// Deletar godoc
// @Summary Deleta acompanhamento
// @Description Deleta um acompanhamento pelo ID
// @Tags Acompanhamento
// @Accept json
// @Produce json
// @Param id path string true "ID do acompanhamento"
// @Success 204 "Acompanhamento deletado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Acompanhamento não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao deletar acompanhamento"
// @Router /acompanhamentos/deletar/{id} [delete]
func (h *AcompanhamentoHandler) Deletar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")

	if err := h.Svc.Deletar(ctx, id); err != nil {
		switch {
		case errors.Is(err, mysql.ErrAcompanhamentoNaoEncontrado):
			ResponderJSONNotFound(w, r, "acompanhamento não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao deletar acompanhamento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao deletar acompanhamento")
			return

		default:
			slog.Error("Erro interno ao deletar acompanhamento", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao deletar acompanhamento")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Desativar, ACP, fmt.Sprintf("ID(%s)", id))

	ResponderJSONNoContent(w)
}

// BuscarPorChamadoID godoc
// @Summary Busca acompanhamentos por ID do chamado
// @Description Retorna uma lista de acompanhamentos associados a um chamado pelo ID do chamado
// @Tags Acompanhamento
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Success 200 {object} []dto.AcompanhamentoResp "Acompanhamentos encontrados com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Acompanhamento não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /acompanhamentos/buscar-por-chamado/{id} [get]
func (h *AcompanhamentoHandler) BuscarPorChamadoID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	acpSlice, err := h.Svc.BuscarPorChamadoID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrAcompanhamentoNaoEncontrado):
			ResponderJSONNotFound(w, r, "acompanhamento não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar acompanhamento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar acompanhamento")
			return

		default:
			slog.Error("Erro interno ao buscar acompanhamento", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar acompanhamento")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaAcompanhamentosResp(acpSlice))
}

// ListarPaginado godoc
// @Summary Lista acompanhamentos com paginação
// @Description Retorna uma lista paginada de acompanhamentos com base nos filtros fornecidos
// @Tags Acompanhamento
// @Accept json
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param limite query int false "Número de itens por página" default(10)
// @Param chamadoId query string false "Filtra por ID do chamado"
// @Param usuarioId query string false "Filtra por ID do usuário"
// @Success 200 {object} dto.AcompanhamentoPaginado "Lista paginada de acompanhamentos retornada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /acompanhamentos/listar-paginado [get]
func (h *AcompanhamentoHandler) ListarPaginado(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro, err := h.parseFiltros(r.URL.Query())
	if err != nil {
		slog.Error("Erro ao parsear filtros de consulta", "error", err)
		ResponderJSONBadRequest(w, r, "filtros de consulta inválidos")
		return
	}

	acompanhamentos, total, filtro, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar acompanhamentos")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar acompanhamentos")
			return

		default:
			slog.Error("Erro interno ao listar acompanhamentos", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar acompanhamentos")
			return
		}
	}

	resp := dto.ParaAcompanhamentosResp(acompanhamentos)
	ResponderJSONPaginado(
		w,
		total,
		filtro.Pagina(),
		filtro.Limite(),
		resp,
	)
}

// parseFiltros extrai e converte os filtros da query string
func (*AcompanhamentoHandler) parseFiltros(query url.Values) (acp.Filtro, error) {
	pagina, limite, err := ParsePaginacaoQuery(query)
	if err != nil {
		return acp.Filtro{}, err
	}
	paginacao := dmn.NovoPaginacao(pagina, limite)
	chamadoId := ParseStringQuery(query, "chamadoId")
	usuarioId := ParseStringQuery(query, "usuarioId")

	return acp.NovoFiltro(paginacao, chamadoId, usuarioId), nil
}
