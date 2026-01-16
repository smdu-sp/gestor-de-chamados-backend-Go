package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/dto"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

const USR = "USUARIO"

// UsuarioHandler gerencia as requisições HTTP relacionadas a usuários.
type UsuarioHandler struct {
	Svc     usr.Service
	AuthSvc auth.AuthInternoService
	LDAPSvc auth.AuthExternoService
	LogSvc  log.Service
}

// NovoUsuarioHandler cria uma nova instância de UsuarioHandler.
func NovoUsuarioHandler(usrSvc usr.Service,
	authSvc auth.AuthInternoService,
	ldapSvc auth.AuthExternoService,
	logSvc log.Service) *UsuarioHandler {

	return &UsuarioHandler{
		Svc:     usrSvc,
		AuthSvc: authSvc,
		LDAPSvc: ldapSvc,
		LogSvc:  logSvc,
	}
}

// Criar godoc
// @Summary Cria um novo usuário
// @Description Cria um usuário com os dados fornecidos no corpo da requisição
// @Tags Usuario
// @Accept json
// @Produce json
// @Param usuario body dto.CriarUsuarioReq true "Dados para criação de usuário"
// @Success 201 {object} dto.UsuarioResp "Usuário criado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "Usuário já existe"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao criar usuário"
// @Router /usuarios [post]
func (h *UsuarioHandler) Criar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	var req dto.CriarUsuarioReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	usrCriado, err := h.Svc.Criar(ctx, req.ParaCriarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrUsuarioJaExisteComEmail):
			ResponderJSONConflict(w, r, "já existe um usuário com o mesmo email")
			return

		case errors.Is(err, mysql.ErrUsuarioJaExisteComLogin):
			ResponderJSONConflict(w, r, "já existe um usuário com o mesmo login")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao criar usuário")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao criar usuário")
			return

		default:
			slog.Error("Erro interno ao criar usuário", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao criar usuário")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Criar, USR, usrCriado.String())

	ResponderJSONCreated(w, dto.ParaUsuarioResp(*usrCriado))
}

// BuscarPorID godoc
// @Summary Busca usuário por ID
// @Description Retorna os dados completos de um usuário pelo seu ID
// @Tags Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} dto.UsuarioResp "Usuário encontrado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Usuário não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/buscar-por-id/{id} [get]
func (h *UsuarioHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	usr, err := h.Svc.BuscarPorID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrUsuarioNaoEncontrado):
			ResponderJSONNotFound(w, r, "usuário não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar usuário")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar usuário")
			return

		default:
			slog.Error("Erro interno ao buscar usuário por id", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar usuário")
			return
		}
	}
	ResponderJSONStatusOK(w, dto.ParaUsuarioResp(*usr))
}

// Atualizar godoc
// @Summary Atualiza usuário
// @Description Atualiza informações de um usuário pelo ID
// @Tags Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param usuario body dto.AtualizarUsuarioReq true "Dados para atualização do usuário"
// @Success 200 {object} dto.UsuarioResp "Usuário atualizado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Usuário não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "Conflito ao atualizar usuário"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno ao atualizar usuário"
// @Router /usuarios/atualizar/{id} [patch]
func (h *UsuarioHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	var req dto.AtualizarUsuarioReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	usrAtualizado, err := h.Svc.Atualizar(ctx, id, req.ParaAtualizarParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrUsuarioNaoEncontrado):
			ResponderJSONNotFound(w, r, "usuário não encontrado")
			return

		case errors.Is(err, mysql.ErrUsuarioJaExisteComEmail) ||
			errors.Is(err, mysql.ErrUsuarioJaExisteComLogin):
			ResponderJSONConflict(w, r, "já existe um usuário com o mesmo email ou login")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar usuário")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar usuário")
			return

		default:
			slog.Error("Erro interno ao atualizar usuário", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar usuário")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, USR, usrAtualizado.String())

	ResponderJSONStatusOK(w, dto.ParaUsuarioResp(*usrAtualizado))
}

// Listar godoc
// @Summary Lista todos os usuários
// @Description Retorna uma lista de todos os usuários do sistema
// @Tags Usuario
// @Accept json
// @Produce json
// @Success 200 {object} []dto.UsuarioResp "Lista de usuários retornada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/listar [get]
func (h *UsuarioHandler) Listar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro := usr.Filtro{}
	filtro.SemLimite()
	usrSlice, _, _, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar usuários")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar usuários")
			return

		default:
			slog.Error("Erro interno ao listar usuários", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar usuários")
			return
		}
	}

	ResponderJSONStatusOK(w, dto.ParaUsuariosResp(usrSlice))
}

// BuscarTecnicos godoc
// @Summary Lista todos os técnicos
// @Description Retorna uma lista de todos os usuários com permissão de técnico
// @Tags Usuario
// @Accept json
// @Produce json
// @Success 200 {object} []dto.TecnicoResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/buscar-tecnicos [get]
func (h *UsuarioHandler) BuscarTecnicos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro := usr.Filtro{}
	filtro.DefinirPermissao(usr.PermTecPtr())
	filtro.SemLimite()

	tecSlice, _, _, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar técnicos")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar técnicos")
			return

		default:
			slog.Error("Erro interno ao listar técnicos", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar técnicos")
			return
		}
	}

	ResponderJSONStatusOK(w, dto.ParaTecnicosResp(tecSlice))
}

// Desativar godoc
// @Summary Desativa usuário
// @Description Desativa usuário pelo ID
// @Tags Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} dto.UsuarioResp "Usuário desativado com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Usuário não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/desativar/{id} [patch]
func (h *UsuarioHandler) Desativar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	usrAtualizado, err := h.Svc.Desativar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrUsuarioNaoEncontrado):
			ResponderJSONNotFound(w, r, "usuário não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao desativar usuário")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao desativar usuário")
			return

		default:
			slog.Error("Erro interno ao desativar usuário", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao desativar usuário")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Desativar, USR, fmt.Sprintf("ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaUsuarioResp(*usrAtualizado))
}

// Ativar godoc
// @Summary Ativa usuário
// @Description Ativa usuário pelo ID
// @Tags Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} dto.UsuarioResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Usuário não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/ativar/{id} [patch]
func (h *UsuarioHandler) Ativar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	usrAtualizado, err := h.Svc.Ativar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrUsuarioNaoEncontrado):
			ResponderJSONNotFound(w, r, "usuário não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao ativar usuário")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao ativar usuário")
			return

		default:
			slog.Error("Erro interno ao ativar usuário", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao ativar usuário")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Ativar, USR, fmt.Sprintf("ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaUsuarioResp(*usrAtualizado))
}

// Autorizar godoc
// @Summary Autoriza usuário
// @Description Autoriza (reativa) usuário pelo ID
// @Tags Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} dto.UsuarioResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Usuário não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/autorizar/{id} [patch]
func (h *UsuarioHandler) Autorizar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")
	usrAtualizado, err := h.Svc.Ativar(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrUsuarioNaoEncontrado):
			ResponderJSONNotFound(w, r, "usuário não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao autorizar usuário")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao autorizar usuário")
			return

		default:
			slog.Error("Erro interno ao autorizar usuário", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao autorizar usuário")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, USR, fmt.Sprintf("Autorizado ID(%s)", id))

	ResponderJSONStatusOK(w, dto.ParaUsuarioResp(*usrAtualizado))
}

// AtualizarPermissao godoc
// @Summary Atualiza permissão do usuário
// @Description Atualiza a permissão de um usuário pelo ID
// @Tags Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param permissao body dto.AtualizarPermissaoUsuarioReq true "Dados para atualização da permissão do usuário"
// @Success 200 {object} dto.UsuarioResp "Permissão do usuário atualizada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 400 {object} handler.ErroResposta "Payload inválido ou JSON malformado"
// @Failure 404 {object} handler.ErroResposta "Usuário não encontrado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 422 {object} handler.ErroResposta "Erro de validação nos dados enviados"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/atualizar-permissao/{id} [patch]
func (h *UsuarioHandler) AtualizarPermissao(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	id := r.PathValue("id")

	var req dto.AtualizarPermissaoUsuarioReq
	if !DecodificarJSON(w, r, &req) {
		return
	}

	usrAtualizado, err := h.Svc.AtualizarPermissao(ctx, id, req.ParaAtualizarPermissaoParams())
	if err != nil {
		var errosValidacao dmn.ErrosValidacao
		switch {
		case errors.As(err, &errosValidacao):
			ResponderJSONValidationError(w, r, errosValidacao)
			return

		case errors.Is(err, mysql.ErrUsuarioNaoEncontrado):
			ResponderJSONNotFound(w, r, "usuário não encontrado")
			return

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao atualizar permissão do usuário")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao atualizar permissão do usuário")
			return

		default:
			slog.Error("Erro interno ao atualizar permissão do usuário", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao atualizar permissão do usuário")
			return
		}
	}

	h.LogSvc.Criar(ctx, log.Atualizar, USR, fmt.Sprintf("Permissão(%s) ID(%s)", req.Permissao, id))

	ResponderJSONStatusOK(w, dto.ParaUsuarioResp(*usrAtualizado))
}

// ValidarUsuario godoc
// @Summary Valida usuário autenticado
// @Description Valida se o usuário autenticado está ativo no sistema
// @Tags Usuario
// @Accept json
// @Produce json
// @Success 200 {object} map[string]bool "Usuário válido"
// @Router /usuarios/validar-usuario [get]
func (h *UsuarioHandler) ValidarUsuario(w http.ResponseWriter, r *http.Request) {
	type Resp struct {
		Valido bool `json:"valido"`
	}

	ResponderJSONStatusOK(w, Resp{Valido: true})
}

// BuscarNovo godoc
// @Summary Busca novo usuário no LDAP
// @Description Busca um novo usuário pelo login no LDAP
// @Tags Usuario
// @Accept json
// @Produce json
// @Param login path string true "Login do usuário"
// @Success 200 {object} dto.UsuarioLdapResp
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 404 {object} handler.ErroResposta "Usuário não encontrado no LDAP"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 409 {object} handler.ErroResposta "Conflito ao atualizar usuário"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/buscar-novo/{login} [get]
func (h *UsuarioHandler) BuscarNovo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	login := r.PathValue("login")

	usr, err := h.Svc.BuscarPorLogin(ctx, login)
	if err != nil {
		switch {
		case errors.Is(err, mysql.ErrUsuarioNaoEncontrado):
			usr = nil // garantir que usr é nil e proceguir para buscar no LDAP

		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar usuário pelo login")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar usuário pelo login")
			return

		default:
			slog.Error("Erro interno ao buscar usuário pelo login", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar usuário pelo login")
			return
		}
	}

	// 1- existe já ativo
	if usr != nil && usr.Status() {
		ResponderJSONConflict(w, r, "Login já cadastrado")
		return
	}

	// 2- existe inativo -> reativa e retorna
	if usr != nil && !usr.Status() {
		if _, err := h.Svc.Ativar(r.Context(), usr.ID()); err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				ResponderJSONTimeout(w, r, "tempo de requisição excedido ao reativar usuário")
				return

			case errors.Is(err, context.Canceled):
				ResponderJSONContextCanceled(w, r, "requisição cancelada ao reativar usuário")
				return
			}
			slog.Error("Erro interno ao reativar usuário", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao reativar usuário")
			return
		}

		h.LogSvc.Criar(ctx, log.Ativar, USR, fmt.Sprintf("ID(%s)", usr.ID()))

		ResponderJSONStatusOK(w, dto.ParaUsuarioResp(*usr))
		return
	}

	// 3- consulta LDAP por novo usuário
	usuarioLDAP, err := h.LDAPSvc.PesquisarPorLogin(login)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao buscar usuário no LDAP")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao buscar usuário no LDAP")
			return

		default:
			slog.Error("Erro interno ao buscar usuário no LDAP", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao buscar usuário no LDAP")
			return
		}
	}

	if usuarioLDAP.Login == "" {
		ResponderJSONNotFound(w, r, "Usuário não encontrado no LDAP")
		return
	}

	ResponderJSONStatusOK(w, dto.ParaUsuarioLdapResp(usuarioLDAP))
}

// ListarPaginado godoc
// @Summary Lista usuários com paginação e filtros
// @Description Retorna uma lista paginada de usuários com base nos filtros fornecidos
// @Tags Usuario
// @Accept json
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param limite query int false "Número de itens por página" default(10)
// @Param busca query string false "Termo de busca no nome ou email do usuário"
// @Param status query bool false "Filtra por status do usuário (ativo/inativo)"
// @Param permissao query string false "Filtra por permissão do usuário"
// @Success 200 {object} dto.UsuarioPaginado "Lista paginada de usuários retornada com sucesso"
// @Failure 400 {object} handler.ErroResposta "Contexto cancelado"
// @Failure 408 {object} handler.ErroResposta "Tempo de requisição excedido"
// @Failure 500 {object} handler.ErroResposta "Erro interno do servidor"
// @Router /usuarios/listar-paginado [get]
func (h *UsuarioHandler) ListarPaginado(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeoutPadrao)
	defer cancel()

	filtro, err := h.parseFiltros(r.URL.Query())
	if err != nil {
		slog.Error("Erro ao parsear filtros de consulta", "erro", err)
		ResponderJSONBadRequest(w, r, "filtros de consulta inválidos")
		return
	}

	usuarios, total, filtro, err := h.Svc.Listar(ctx, filtro)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			ResponderJSONTimeout(w, r, "tempo de requisição excedido ao listar usuários")
			return

		case errors.Is(err, context.Canceled):
			ResponderJSONContextCanceled(w, r, "requisição cancelada ao listar usuários")
			return

		default:
			slog.Error("Erro interno ao listar usuários", "erro", err)
			ResponderJSONInternalError(w, r, "erro interno ao listar usuários")
			return
		}
	}

	resp := dto.ParaUsuariosResp(usuarios)
	ResponderJSONPaginado(
		w,
		total,
		filtro.Pagina(),
		filtro.Limite(),
		resp,
	)
}

// parseFiltros extrai e converte os filtros da query string
func (*UsuarioHandler) parseFiltros(query url.Values) (usr.Filtro, error) {
	pagina, limite, err := ParsePaginacaoQuery(query)
	if err != nil {
		return usr.Filtro{}, err
	}
	paginacao := dmn.NovoPaginacao(pagina, limite)
	busca := ParseStringQuery(query, "busca")

	status, err := ParseBoolQuery(query, "status")
	if err != nil {
		return usr.Filtro{}, err
	}
	permissao := ParseStringQuery(query, "permissao")

	return usr.NovoFiltro(paginacao, busca, status, permissao), nil
}
