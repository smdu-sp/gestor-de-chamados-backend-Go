package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	ldapauth "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/ldap"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/response"
)

type UsersHandler struct {
	Svc  *user.Service
	LDAP *ldapauth.Client
}

// Helpers de path
func lastSegment(path string) string {
	path = strings.TrimSuffix(path, "/")
	if LastIndex := strings.LastIndex(path, "/"); LastIndex >= 0 {
		return path[LastIndex+1:]
	}
	return path
}

// Criar godoc
// @Summary Cria um novo usuário
// @Description Cria um usuário no sistema (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param usuario body user.Usuario true "Dados do usuário"
// @Success 201 {object} user.Usuario
// @Failure 400 {object} response.ErrorResponse
// @Router /usuarios/criar [post]
// Cria usuário
func (h *UsersHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var usuario user.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "payload inválido", nil)
		return
	}
	if usuario.Permissao == "" {
		usuario.Permissao = user.PermUSR
	}
	if err := h.Svc.Criar(r.Context(), &usuario); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusCreated, usuario)
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
// @Success 200 {object} response.PageResp
// @Failure 500 {object} response.ErrorResponse
// @Router /usuarios/buscar-tudo [get]
func (h *UsersHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	query := r.URL.Query()
	pagina, _ := strconv.Atoi(query.Get("pagina"))
	if pagina <= 0 {
		pagina = 1
	}

	limite, _ := strconv.Atoi(query.Get("limite"))
	if limite <= 0 || limite > 100 {
		limite = 10
	}

	busca := query.Get("busca")
	status := query.Get("status")
	permissao := query.Get("permissao")
	var bPtr, sPtr, pPtr *string

	if busca != "" {
		bPtr = &busca
	}
	if status != "" {
		sPtr = &status
	}
	if permissao != "" {
		pPtr = &permissao
	}

	items, total, err := h.Svc.Listar(r.Context(), pagina, limite, bPtr, sPtr, pPtr)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, "erro na listagem", nil)
		return
	}

	response.JSON(w, http.StatusOK, response.PageResponse{
		Total:  total,
		Pagina: pagina,
		Limite: limite,
		Data:   items,
	})
}

// BuscarPorID godoc
// @Summary Busca usuário por ID
// @Description Retorna usuário pelo ID (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} user.Usuario
// @Failure 404 {object} response.ErrorResponse
// @Router /usuarios/buscar-por-id/{id} [get]
func (h *UsersHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	id := lastSegment(r.URL.Path)
	usuario, err := h.Svc.BuscarPorID(r.Context(), id)
	if errors.Is(err, user.ErrUsuarioNaoEncontrado) {
		response.ErrorJSON(w, http.StatusNotFound, "Usuário não encontrado", nil)
		return
	}
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, "erro", nil)
		return
	}
	response.JSON(w, http.StatusOK, usuario)
}

// Atualizar godoc
// @Summary Atualiza usuário
// @Description Atualiza dados do usuário (ADM/TEC/USR conforme regra)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param usuario body user.Usuario true "Dados do usuário"
// @Success 200 {object} map[string]any
// @Failure 400 {object} response.ErrorResponse
// @Router /usuarios/atualizar/{id} [patch]
func (h *UsersHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.NotFound(w, r)
		return
	}

	// extrai id da URL
	id := lastSegment(r.URL.Path)
	var u user.Usuario
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "payload inválido", nil)
		return
	}

	// atualiza
	if err := h.Svc.Atualizar(r.Context(), id, &u); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

// ListaCompleta godoc
// @Summary Lista completa de usuários
// @Description Retorna todos os usuários (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Success 200 {array} user.Usuario
// @Failure 500 {object} response.ErrorResponse
// @Router /usuarios/lista-completa [get]
func (h *UsersHandler) ListaCompleta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	busca := ""
	status := ""
	permissao := ""
	items, _, err := h.Svc.Listar(r.Context(), 1, 10000, &busca, &status, &permissao)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, "erro na listagem", nil)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

// BuscarTecnicos godoc
// @Summary Lista técnicos
// @Description Retorna lista de técnicos (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Success 200 {array} map[string]string
// @Failure 500 {object} response.ErrorResponse
// @Router /usuarios/buscar-tecnicos [get]
func (h *UsersHandler) BuscarTecnicos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	perm := string(user.PermTEC)
	items, _, err := h.Svc.Listar(r.Context(), 1, 10000, nil, nil, &perm)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, "erro", nil)
		return
	}

	// reduzir payload
	type mini struct {
		ID, Nome string
	}
	out := make([]mini, 0, len(items))
	for _, u := range items {
		out = append(out, mini{ID: u.ID, Nome: u.Nome})
	}

	response.JSON(w, http.StatusOK, out)
}

// Desativar godoc
// @Summary Desativa usuário
// @Description Desativa (soft delete) usuário pelo ID (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} map[string]any
// @Failure 400 {object} response.ErrorResponse
// @Router /usuarios/desativar/{id} [delete]
func (h *UsersHandler) Desativar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.NotFound(w, r)
		return
	}
	id := lastSegment(r.URL.Path)

	// status=false
	if err := h.Svc.Atualizar(r.Context(), id, &user.Usuario{Status: false}); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"desativado": true})
}

// Autorizar godoc
// @Summary Autoriza usuário
// @Description Autoriza (reativa) usuário pelo ID (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} map[string]any
// @Failure 400 {object} response.ErrorResponse
// @Router /usuarios/autorizar/{id} [patch]
func (h *UsersHandler) Autorizar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.NotFound(w, r)
		return
	}
	id := lastSegment(r.URL.Path)
	if err := h.Svc.Atualizar(r.Context(), id, &user.Usuario{Status: true}); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"autorizado": true})
}

// ValidaUsuario godoc
// @Summary Valida usuário autenticado
// @Description Verifica se o usuário está autenticado
// @Tags usuarios
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 404 {object} response.ErrorResponse
// @Router /usuarios/valida-usuario [get]
func (h *UsersHandler) ValidaUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

// BuscarNovo godoc
// @Summary Busca usuário novo no LDAP
// @Description Busca usuário no LDAP e retorna dados (apenas ADM)
// @Tags usuarios
// @Accept json
// @Produce json
// @Param login path string true "Login do usuário"
// @Success 200 {object} map[string]any
// @Failure 404 {object} response.ErrorResponse
// @Router /usuarios/buscar-novo/{login} [get]
func (h *UsersHandler) BuscarNovo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	login := lastSegment(r.URL.Path)

	// 1) existe ativo?
	if usuario, _ := h.Svc.BuscarPorLogin(r.Context(), login); usuario != nil && usuario.Status {
		response.ErrorJSON(w, http.StatusForbidden, "Login já cadastrado", nil)
		return
	}

	// 2) existe inativo? reativar e retornar
	if usuario, _ := h.Svc.BuscarPorLogin(r.Context(), login); usuario != nil && !usuario.Status {
		_ = h.Svc.Atualizar(r.Context(), usuario.ID, &user.Usuario{Status: true})
		response.JSON(w, http.StatusOK, map[string]any{
			"login": usuario.Login,
			"nome":  usuario.Nome,
			"email": usuario.Email,
		})
		return
	}

	// 3) consulta LDAP
	if h.LDAP == nil {
		response.ErrorJSON(w, http.StatusInternalServerError, "LDAP não configurado", nil)
		return
	}

	name, mail, sLogin, err := h.LDAP.SearchByLogin(login)
	if err != nil || sLogin == "" {
		response.ErrorJSON(w, http.StatusNotFound, "Usuário não encontrado no LDAP", nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"login": sLogin,
		"nome":  name,
		"email": mail,
	})
}
