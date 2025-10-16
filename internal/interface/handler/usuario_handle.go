package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

const (
	timeoutPadrao      = 5 * time.Second
	payloadInvalidoMsg = "verifique os dados enviados na requisição"
	erroLogMsg         = "erro ao criar log"
	entidadeUsuario    = "USUARIO"
)

// UsuarioHandler gerencia as requisições HTTP relacionadas a usuários.
type UsuarioHandler struct {
	UsecaseUsr  usecase.UsuarioUsecase
	UsecaseAuth usecase.AuthInternoUsecase
	UsecaseLDAP usecase.AuthExternoUsecase
	UsecaseLog  usecase.LogUsecase
}

// NewUsuarioHandler cria uma nova instância de UsuarioHandler.
func NewUsuarioHandler(usecaseUsr usecase.UsuarioUsecase,
	usecaseAuth usecase.AuthInternoUsecase,
	usecaseLDAP usecase.AuthExternoUsecase,
	usecaseLog usecase.LogUsecase) *UsuarioHandler {

	return &UsuarioHandler{
		UsecaseUsr:  usecaseUsr,
		UsecaseAuth: usecaseAuth,
		UsecaseLDAP: usecaseLDAP,
		UsecaseLog:  usecaseLog,
	}
}

// Helpers de path

// lastSegment extrai o último segmento de um path
func lastSegment(path string) string {
	path = strings.TrimSuffix(path, "/")
	if LastIndex := strings.LastIndex(path, "/"); LastIndex >= 0 {
		return path[LastIndex+1:]
	}
	return path
}

// metodoValido valida se o método HTTP é permitido
func metodoHttpValido(w http.ResponseWriter, r *http.Request, metodoEsperado string) bool {
	if r.Method != metodoEsperado {
		response.ErrorJSON(
			w,
			http.StatusMethodNotAllowed, // status 405
			"método não permitido",
			response.MethodErrorResponse{
				MetodoUsado:     r.Method,
				MetodoPermitido: metodoEsperado,
			},
		)
		return false
	}
	return true
}

// Criar godoc
// @Summary Cria um novo usuário
// @Description Cria um usuário com os dados fornecidos no corpo da requisição.
// @Tags usuarios
// @Accept json
// @Produce json
// @Param usuario body model.Usuario true "Dados do usuário"
// @Success 201 {object} any
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 409 {object} any
// @Failure 500 {object} any
// @Router /usuarios/criar [post]
// Cria usuário
func (h *UsuarioHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	var usuario model.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.UsecaseUsr.CriarUsuario(ctx, &usuario); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao criar usuário", err.Error())
			return

		// conflito - 409
		case errors.Is(err, repository.ErrUsuarioJaExiste):
			response.ErrorJSON(w, http.StatusConflict, "usuário já cadastrado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, utils.ErrUUIDv7Generation),
			errors.Is(err, repository.ErrRowsAffected),
			errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao criar usuário", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao criar usuário", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao criar usuário", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao criar usuário", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoCriar,
		entidadeUsuario,
		fmt.Sprintf("Usuário criado via API: %s", usuario.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, response.ToUsuarioResponse(&usuario))
}

// BuscarTudo godoc
// @Summary Lista usuários com paginação e filtros
// @Description Retorna lista paginada de usuários (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param pagina query int false "Página"
// @Param limite query int false "Limite"
// @Param busca query string false "Busca"
// @Param status query string false "Status"
// @Param permissao query string false "Permissão"
// @Success 200 {object} []model.Usuario
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /usuarios/buscar-tudo [get]
// BuscarTudo lista usuários com paginação e filtros
func (h *UsuarioHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	query := r.URL.Query()

	filtro := model.UsuarioFiltro{}

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
		statusBool, err := strconv.ParseBool(status)
		if err == nil {
			filtro.Status = &statusBool
		}
	}
	if permissao := query.Get("permissao"); permissao != "" {
		filtro.Permissao = &permissao
	}

	items, total, filtroCorrigido, err := h.UsecaseUsr.ListarUsuarios(ctx, filtro)
	if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerUsuario),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao listar usuários", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao listar usuários", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar usuários", err.Error())
			return

		default:
			// fallback de segurança - 500
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar usuários", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.PageResponse[model.Usuario]{
		Total:  total,
		Pagina: filtroCorrigido.Pagina,
		Limite: filtroCorrigido.Limite,
		Items:  items,
	})
}

// BuscarPorID godoc
// @Summary Busca usuário por ID
// @Description Retorna usuário pelo ID (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} model.Usuario
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /usuarios/buscar-por-id/{id} [get]
// BuscarPorID busca usuário por ID
func (h *UsuarioHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	usuario, err := h.UsecaseUsr.BuscarUsuarioPorID(ctx, id)
	if err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, model.ErrIDInvalido),
			errors.Is(err, repository.ErrUsuarioNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao buscar usuário", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerUsuario):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao buscar usuário", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao buscar usuário", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao buscar usuário", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao buscar usuário", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, response.ToUsuarioResponse(usuario))
}

// Atualizar godoc
// @Summary Atualiza usuário
// @Description Atualiza dados do usuário (ADM/TEC/USR conforme regra)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param usuario body model.Usuario true "Dados do usuário"
// @Success 200 {object} model.Usuario
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /usuarios/atualizar/{id} [put]
// Atualizar atualiza usuário por ID
func (h *UsuarioHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPut) {
		return
	}
	// TODO (melhorar): permitir atualização parcial (PATCH)
	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	var usuario model.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.UsecaseUsr.AtualizarUsuario(ctx, id, &usuario); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao atualizar usuário", err.Error())
			return

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrUsuarioNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao atualizar usuário", err.Error())
			return

		// Conflito - 409
		case errors.Is(err, repository.ErrUsuarioJaExiste):
			response.ErrorJSON(w, http.StatusConflict, "usuário já cadastrado", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrQueryContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao atualizar usuário", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao atualizar usuário", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao atualizar usuário", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar usuário", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtualizar,
		entidadeUsuario,
		fmt.Sprintf("Usuário atualizado via API: %s", usuario.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.ToUsuarioResponse(&usuario))
}

// ListaCompleta godoc
// @Summary Lista completa de usuários
// @Description Retorna todos os usuários (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Success 200 {object} []model.Usuario
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /usuarios/lista-completa [get]
// ListaCompleta retorna todos os usuários sem paginação
func (h *UsuarioHandler) ListaCompleta(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	filtro := model.UsuarioFiltro{
		Pagina: 1,
		Limite: 10000000,
	}
	items, _, _, err := h.UsecaseUsr.ListarUsuarios(ctx, filtro)
	if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerUsuario),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao listar usuários", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao listar usuários", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar usuários", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar usuários", err.Error())
			return
		}
	}
	response.JSON(w, http.StatusOK, items)
}

// BuscarTecnicos godoc
// @Summary Lista técnicos
// @Description Retorna lista de técnicos (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Success 200 {object} []model.Usuario
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /usuarios/buscar-tecnicos [get]
func (h *UsuarioHandler) BuscarTecnicos(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	perm := string(model.PermTEC)
	filtro := model.UsuarioFiltro{
		Pagina:    1,
		Limite:    10000,
		Permissao: &perm,
	}
	items, _, _, err := h.UsecaseUsr.ListarUsuarios(ctx, filtro)
	if err != nil {
		switch {
		// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerUsuario),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao listar técnicos", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao listar técnicos", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar técnicos", err.Error())
			return

		default:
			// fallback de segurança - 500
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar técnicos", err.Error())
			return
		}
	}

	resp := make([]response.Tecnico, 0, len(items))
	for _, u := range items {
		resp = append(resp, response.Tecnico{ID: u.ID, Nome: u.Nome})
	}

	response.JSON(w, http.StatusOK, resp)
}

// Desativar godoc
// @Summary Desativa usuário
// @Description Desativa (soft delete) usuário pelo ID (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} map[string]any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /usuarios/desativar/{id} [delete]
func (h *UsuarioHandler) Desativar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodDelete) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	if err := h.UsecaseUsr.DesativarUsuario(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrUsuarioNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao desativar usuário", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao desativar usuário", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao desativar usuário", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao desativar usuário", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao desativar usuário", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoDesativar,
		entidadeUsuario,
		fmt.Sprintf("Usuário desativado via API: usuário ID(%s)", id),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.StatusUsuario{Ativo: false})
}

// Ativar godoc
// @Summary Ativa usuário
// @Description Ativa (reativa) usuário pelo ID (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} map[string]any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 500 {object} any
// @Router /usuarios/ativar/{id} [patch]
func (h *UsuarioHandler) Ativar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPatch) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	if err := h.UsecaseUsr.AtivarUsuario(ctx, id); err != nil {
		switch {
		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrUsuarioNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao ativar usuário", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao ativar usuário", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao ativar usuário", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao ativar usuário", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao ativar usuário", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtivar,
		entidadeUsuario,
		fmt.Sprintf("Usuário ativado via API: usuário ID(%s)", id),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.StatusUsuario{Ativo: true})
}

// Autorizar godoc
// @Summary Autoriza usuário
// @Description Autoriza (reativa) usuário pelo ID (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} map[string]any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /usuarios/autorizar/{id} [patch]
func (h *UsuarioHandler) Autorizar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPatch) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)
	if err := h.UsecaseUsr.AtualizarUsuario(ctx, id, &model.Usuario{Status: true}); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos", err.Error())
			return

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrUsuarioNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao autorizar usuário", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao autorizar usuário", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao autorizar usuário", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao autorizar usuário", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtualizar,
		entidadeUsuario,
		fmt.Sprintf("Usuário autorizado via API: usuário ID(%s)", id),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.AutorizacaoUsuario{Autorizado: true})
}

// AtualizarPermissao godoc
// @Summary Atualiza permissão do usuário
// @Description Atualiza permissão do usuário pelo ID (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param permissao body object true "Permissão do usuário" example({"permissao": "ADM"})
// @Success 200 {object} map[string]any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /usuarios/atualizar-permissao/{id} [patch]
func (h *UsuarioHandler) AtualizarPermissao(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPatch) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := lastSegment(r.URL.Path)

	var requisicao struct {
		Permissao string `json:"permissao"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requisicao); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}
	if err := h.UsecaseUsr.AtualizarPermissaoUsuario(ctx, id, requisicao.Permissao); err != nil {
		switch {
		// requisições inválidas - 400
		case errors.As(err, &utils.ValidacaoErrors{}):
			response.ErrorJSON(w, http.StatusBadRequest, "permissão inválida", err.Error())
			return

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrUsuarioNaoEncontrado):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido", err.Error())
			return

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao atualizar permissão do usuário", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao atualizar permissão do usuário", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusInternalServerError, "requisição cancelada ao atualizar permissão do usuário", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar permissão do usuário", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtualizar,
		entidadeUsuario,
		fmt.Sprintf("Permissão do usuário atualizada via API: usuário ID(%s), nova permissão(%s)", id, requisicao.Permissao),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.PermissaoUsuario{Permissao: requisicao.Permissao})
}

// ValidaUsuario godoc
// @Summary Valida usuário autenticado
// @Description Verifica se o usuário está autenticado
// @Tags usuarios
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Router /usuarios/valida-usuario [get]
func (h *UsuarioHandler) ValidaUsuario(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}
	response.JSON(w, http.StatusOK, response.ValidaUsuario{Valido: true})
}

// BuscarNovo godoc
// @Summary Busca usuário novo no LDAP
// @Description Busca usuário no LDAP e retorna dados (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param login path string true "Login do usuário"
// @Success 200 {object} map[string]any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Router /usuarios/buscar-novo/{login} [get]
func (h *UsuarioHandler) BuscarNovo(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	login := lastSegment(r.URL.Path)

	// 1) existe ativo?
	if usuario, _ := h.UsecaseUsr.BuscarUsuarioPorLogin(ctx, login); usuario != nil && usuario.Status {
		response.ErrorJSON(w, http.StatusForbidden, "Login já cadastrado", nil)
		return
	}

	// 2) existe inativo? reativar e retornar
	if usuario, _ := h.UsecaseUsr.BuscarUsuarioPorLogin(r.Context(), login); usuario != nil && !usuario.Status {
		_ = h.UsecaseUsr.AtualizarUsuario(r.Context(), usuario.ID, &model.Usuario{Status: true})
		response.JSON(w, http.StatusOK, response.BuscarNovo{
			Login: usuario.Login,
			Nome:  usuario.Nome,
			Email: usuario.Email,
		})

		err := h.UsecaseLog.CriarLog(
			ctx,
			model.AcaoAtivar,
			entidadeUsuario,
			fmt.Sprintf("Usuário reativado via API: usuário ID(%s)", usuario.ID),
		)
		if err != nil {
			response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		}
		
		return
	}

	// 3) consulta LDAP
	if h.UsecaseLDAP == nil {
		response.ErrorJSON(w, http.StatusInternalServerError, "LDAP não configurado", nil)
		return
	}

	name, mail, sLogin, err := h.UsecaseLDAP.PesquisarPorLogin(login)
	if err != nil || sLogin == "" {
		response.ErrorJSON(w, http.StatusNotFound, "Usuário não encontrado no LDAP", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, response.BuscarNovo{
		Login: sLogin,
		Nome:  name,
		Email: mail,
	})
}
