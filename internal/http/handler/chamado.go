package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

const CHM = "CHAMADO"

// ChamadoHandler gerencia as requisições HTTP relacionadas a chamados.
type ChamadoHandler struct {
	Svc    chm.Service
	LogSvc log.Service
}

// NovoChamadoHandler cria uma nova instância de ChamadoHandler.
func NovoChamadoHandler(chmSvc chm.Service, logSvc log.Service) *ChamadoHandler {
	return &ChamadoHandler{
		Svc:    chmSvc,
		LogSvc: logSvc,
	}
}

// Criar godoc
// @Summary Cria um novo chamado
// @Description Cria um novo chamado com os dados fornecidos no corpo da requisição
// @Tags Chamado
// @Accept json
// @Produce json
// @Param chamado body dto.CriarChamadoReq true "Dados para criação de chamado"
// @Success 201 {object} dto.ChamadoResp "Chamado criado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados [post]
func (h *ChamadoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	var req dto.CriarChamadoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	chmCriado, err := h.Svc.Criar(ctx, req.ParaCriarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao criar chamado")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao criar chamado")
			return

		default:
			slog.Error("Erro interno ao criar chamado", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao criar chamado")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Criar, CHM, chmCriado.String())

	ResponderJSONCreated(w, dto.ParaChamadoResp(*chmCriado))
}

// BuscarPorID godoc
// @Summary Busca um chamado pelo ID
// @Description Retorna os detalhes de um chamado específico pelo seu ID
// @Tags Chamado
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Success 200 {object} dto.ChamadoResp "Chamado encontrado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Chamado não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados/buscar/{id} [get]
func (h *ChamadoHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	chm, err := h.Svc.BuscarPorID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrChamadoNaoEncontrado):
			ResponderJSONNotFound(w, r, "chamado não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar chamado")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar chamado")
			return

		default:
			slog.Error("Erro interno ao buscar chamado", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar chamado")
			return
		}
	}

	ResponderJSONStatusOK(w, dto.ParaChamadoResp(*chm))
}

// Atualizar godoc
// @Summary Atualiza um chamado existente
// @Description Atualiza as informações de um chamado existente pelo ID
// @Tags Chamado
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Param chamado body dto.AtualizarChamadoReq true "Dados para atualização do chamado"
// @Success 200 {object} dto.ChamadoResp "Chamado atualizado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Chamado não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados/atualizar/{id} [patch]
func (h *ChamadoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	var req dto.AtualizarChamadoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	chmAtualizado, err := h.Svc.Atualizar(ctx, id, req.ParaAtualizarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrChamadoNaoEncontrado):
			ResponderJSONNotFound(w, r, "chamado não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar chamado")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar chamado")
			return

		default:
			slog.Error("erro interno ao atualizar chamado", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar chamado")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, CHM, chmAtualizado.String())

	ResponderJSONStatusOK(w, dto.ParaChamadoResp(*chmAtualizado))
}

// Arquivar godoc
// @Summary Arquiva um chamado existente
// @Description Arquiva um chamado existente pelo ID
// @Tags Chamado
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Success 200 {object} dto.ChamadoResp "Chamado arquivado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Chamado não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados/arquivar/{id} [delete]
func (h *ChamadoHandler) Arquivar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	chmArquivado, err := h.Svc.Arquivar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrChamadoNaoEncontrado):
			ResponderJSONNotFound(w, r, "chamado não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao arquivar chamado")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao arquivar chamado")
			return

		default:
			slog.Error("Erro interno ao arquivar chamado", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao arquivar chamado")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Arquivar, CHM, fmt.Sprintf("ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaChamadoResp(*chmArquivado))
}

// Desarquivar godoc
// @Summary Desarquiva um chamado existente
// @Description Desarquiva um chamado existente pelo ID
// @Tags Chamado
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Success 200 {object} dto.ChamadoResp "Chamado desarquivado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Chamado não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados/desarquivar/{id} [patch]
func (h *ChamadoHandler) Desarquivar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	chmDesarquivado, err := h.Svc.Desarquivar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrChamadoNaoEncontrado):
			ResponderJSONNotFound(w, r, "chamado não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao desarquivar chamado")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao desarquivar chamado")
			return

		default:
			slog.Error("Erro interno ao desarquivar chamado", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao desarquivar chamado")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Desarquivar, CHM, fmt.Sprintf("ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaChamadoResp(*chmDesarquivado))
}

// AtualizarStatus godoc
// @Summary Atualiza o status de um chamado existente
// @Description Atualiza o status de um chamado existente pelo ID
// @Tags Chamado
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Param chamado body dto.AtualizarStatusReq true "Novo status do chamado"
// @Success 200 {object} dto.ChamadoResp "Chamado atualizado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Chamado não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados/atualizar-status/{id} [patch]
func (h *ChamadoHandler) AtualizarStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")

	var req dto.AtualizarStatusReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	chmAtualizado, err := h.Svc.AtualizarStatus(ctx, id, req.ParaAtualizarStatusParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrChamadoNaoEncontrado):
			ResponderJSONNotFound(w, r, "chamado não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar status do chamado")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar status do chamado")
			return

		default:
			slog.Error("Erro interno ao atualizar status do chamado", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar status do chamado")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, CHM, chmAtualizado.String())

	ResponderJSONStatusOK(w, dto.ParaChamadoResp(*chmAtualizado))
}

// AtualizarSolucao godoc
// @Summary Atualiza a solução de um chamado existente
// @Description Atualiza a solução de um chamado existente pelo ID
// @Tags Chamado
// @Accept json
// @Produce json
// @Param id path string true "ID do chamado"
// @Param chamado body dto.AtualizarSolucaoReq true "Nova solução do chamado"
// @Success 200 {object} dto.ChamadoResp "Chamado atualizado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Chamado não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados/atualizar-solucao/{id} [patch]
func (h *ChamadoHandler) AtualizarSolucao(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")

	var req dto.AtualizarSolucaoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	chmAtualizado, err := h.Svc.AtualizarSolucao(ctx, id, req.ParaAtualizarSolucaoParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrChamadoNaoEncontrado):
			ResponderJSONNotFound(w, r, "chamado não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar solução do chamado")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar solução do chamado")
			return

		default:
			slog.Error("Erro interno ao atualizar solução do chamado", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar solução do chamado")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, CHM, fmt.Sprintf("Solução do chamado atualizada ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaChamadoResp(*chmAtualizado))
}

// Listar godoc
// @Summary Lista todos os chamados
// @Description Retorna lista de todos os chamados
// @Tags Chamado
// @Accept json
// @Produce json
// @Success 200 {object} []dto.ChamadoResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados/listar [get]
func (h *ChamadoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro := chm.Filtro{}
	filtro.SemLimite()
	chmSlice, _, _, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar chamados")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar chamados")
			return

		default:
			slog.Error("Erro interno ao listar chamados", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar chamados")
			return
		}
	}

	ResponderJSONStatusOK(w, dto.ParaChamadosResp(chmSlice))
}

// ListarPaginado godoc
// @Summary Lista chamados com paginação e filtros
// @Description Retorna uma lista paginada de chamados com base nos filtros fornecidos
// @Tags Chamado
// @Accept json
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param limite query int false "Número de itens por página" default(10)
// @Param busca query string false "Termo de busca"
// @Param status query string false "Status do chamado"
// @Param categoriaId query string false "ID da categoria"
// @Param subcategoriaId query string false "ID da subcategoria"
// @Param criadorId query string false "ID do criador"
// @Success 200 {object} dto.ChamadoPaginado "Lista paginada de chamados"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /chamados/listar-paginado [get]
func (h *ChamadoHandler) ListarPaginado(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro, err := h.parseFiltros(r.URL.Query())
	if err != nil {
		slog.Error("Erro ao parsear filtros de consulta", "error", err)
		ResponderJSONBadRequest(w, r, "filtros de consulta inválidos")
		return
	}

	chamados, total, filtro, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar chamados")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar chamados")
			return

		default:
			slog.Error("Erro interno ao listar chamados", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar chamados")
			return
		}
	}

	resp := dto.ParaChamadosResp(chamados)
	ResponderJSONPaginado(
		w,
		total,
		filtro.Pagina(),
		filtro.Limite(),
		resp,
	)
}

// parseFiltros extrai os parâmetros de consulta e constrói um ChamadoFiltro.
func (*ChamadoHandler) parseFiltros(query url.Values) (chm.Filtro, error) {
	pagina, limite, err := ParsePaginacaoQuery(query)
	if err != nil {
		return chm.Filtro{}, err
	}
	paginacao := dmn.NovoPaginacao(pagina, limite)
	busca := ParseStringQuery(query, "busca")
	status := ParseStringQuery(query, "status")
	categoriaID := ParseStringQuery(query, "categoriaId")
	subcategoriaID := ParseStringQuery(query, "subcategoriaId")
	criadorID := ParseStringQuery(query, "criadorId")

	return chm.NovoFiltro(paginacao, busca, status, categoriaID, subcategoriaID, criadorID), nil
}
