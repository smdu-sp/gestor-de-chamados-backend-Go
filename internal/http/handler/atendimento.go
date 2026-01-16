package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

const ATD = "ATENDIMENTO"

// AtendimentoHandler gerencia as requisições HTTP relacionadas a atendimentos.
type AtendimentoHandler struct {
	Svc    atd.Service
	LogSvc log.Service
}

// NovoAtendimentoHandler cria uma nova instância de AtendimentoHandler.
func NovoAtendimentoHandler(atdSvc atd.Service, logSvc log.Service) *AtendimentoHandler {
	return &AtendimentoHandler{
		Svc:    atdSvc,
		LogSvc: logSvc,
	}
}

// Criar godoc
// @Summary Cria um novo atendimento
// @Description Cria um novo atendimento com os dados fornecidos no corpo da requisição
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param atendimento body dto.CriarAtendimentoReq true "Dados do atendimento"
// @Success 201 {object} dto.AtendimentoResp "Atendimento criado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao criar atendimento"
// @Router /atendimentos [post]
func (h *AtendimentoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	var req dto.CriarAtendimentoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	atdCriado, err := h.Svc.Criar(ctx, req.ParaCriarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao criar atendimento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao criar atendimento")
			return

		default:
			slog.Error("Erro inesperado ao criar atendimento", "error", err)
			ResponderJSONInternalError(w, r, "erro inesperado ao criar atendimento")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Criar, ATD, atdCriado.String())

	ResponderJSONCreated(w, dto.ParaAtendimentoResp(*atdCriado))
}

// BuscarPorID godoc
// @Summary Busca um atendimento pelo ID
// @Description Busca um atendimento existente pelo seu ID
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param id path string true "ID do atendimento"
// @Success 200 {object} dto.AtendimentoResp "Atendimento encontrado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Atendimento não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao buscar atendimento"
// @Router /atendimentos/buscar-por-id/{id} [get]
func (h *AtendimentoHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	atd, err := h.Svc.BuscarPorID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrAtendimentoNaoEncontrado):
			ResponderJSONNotFound(w, r, "atendimento não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar atendimento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar atendimento")
			return

		default:
			slog.Error("Erro interno ao buscar atendimento", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar atendimento")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaAtendimentoResp(*atd))
}

// BuscarPorIDs godoc
// @Summary Busca um atendimento pelo ID do chamado e do técnico atribuído
// @Description Busca um atendimento existente pelo ID do chamado e do técnico atribuído
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param chamadoId query string true "ID do chamado"
// @Param atribuidoId query string true "ID do técnico atribuído"
// @Success 200 {object} dto.AtendimentoResp "Atendimento encontrado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Atendimento não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao buscar atendimento"
// @Router /atendimentos/buscar-por-ids/chamados/{chamadoId}/tecnicos/{tecnicoId} [get]
func (h *AtendimentoHandler) BuscarPorChamadoEAtribuidoID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	query := r.URL.Query()
	chamadoID := query.Get("chamadoId")
	atribuidoID := query.Get("atribuidoId")

	atd, err := h.Svc.BuscarPorChamadoEAtribuidoID(ctx, chamadoID, atribuidoID)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrAtendimentoNaoEncontrado):
			ResponderJSONNotFound(w, r, "atendimento não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar atendimento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar atendimento")
			return

		default:
			slog.Error("Erro interno ao buscar atendimento", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar atendimento")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaAtendimentoResp(*atd))
}

// Atualizar godoc
// @Summary Atualiza um atendimento existente
// @Description Atualiza as informações de um atendimento existente pelo seu ID
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param id path string true "ID do atendimento"
// @Param atendimento body dto.AtualizarAtendimentoReq true "Dados para atualização do atendimento"
// @Success 200 {object} dto.AtendimentoResp "Atendimento atualizado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Atendimento não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao atualizar atendimento"
// @Router /atendimentos/atualizar/{id} [patch]
func (h *AtendimentoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	var req dto.AtualizarAtendimentoReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	atdAtualizado, err := h.Svc.Atualizar(ctx, id, req.ParaAtualizarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrAtendimentoNaoEncontrado):
			ResponderJSONNotFound(w, r, "atendimento não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar atendimento")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar atendimento")
			return

		default:
			slog.Error("Erro interno ao atualizar atendimento", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar atendimento")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, ATD, atdAtualizado.String())

	ResponderJSONCreated(w, dto.ParaAtendimentoResp(*atdAtualizado))
}

// Listar godoc
// @Summary Lista atendimentos com paginação e filtros
// @Description Lista atendimentos existentes com suporte a paginação e filtros opcionais
// @Tags Atendimento
// @Accept json
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param limite query int false "Número de itens por página" default(10)
// @Param chamadoId query string false "Filtra por ID do chamado"
// @Param atribuidoId query string false "Filtra por ID do usuário atribuído"
// @Success 200 {object} dto.AtendimentoPaginado "Lista paginada de atendimentos retornada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao listar atendimentos"
// @Router /atendimentos/listar-paginado [get]
func (h *AtendimentoHandler) ListarPaginado(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro, err := h.parseFiltros(r.URL.Query())
	if err != nil {
		slog.Error("Erro ao parsear filtros de consulta", "error", err)
		ResponderJSONBadRequest(w, r, "filtros de consulta inválidos")
		return
	}

	atendimentos, total, filtro, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar atendimentos")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar atendimentos")
			return

		default:
			slog.Error("Erro interno ao listar atendimentos", "error", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar atendimentos")
			return
		}
	}

	resp := dto.ParaAtendimentosResp(atendimentos)
	ResponderJSONPaginado(
		w,
		total,
		filtro.Pagina(),
		filtro.Limite(),
		resp,
	)
}

// parseFiltros extrai os filtros de atendimento dos parâmetros da URL.
func (*AtendimentoHandler) parseFiltros(query url.Values) (atd.Filtro, error) {
	pagina, limite, err := ParsePaginacaoQuery(query)
	if err != nil {
		return atd.Filtro{}, err
	}
	paginacao := dmn.NovoPaginacao(pagina, limite)
	chamadoId := ParseStringQuery(query, "chamadoId")
	atribuidoId := ParseStringQuery(query, "atribuidoId")

	return atd.NovoFiltro(paginacao, chamadoId, atribuidoId), nil
}
