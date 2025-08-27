package handlers

import (
	"encoding/json"
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
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}

// POST /usuarios/criar (ADM)
func (h *UsersHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var u user.Usuario
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "payload inválido", nil)
		return
	}
	if u.Permissao == "" {
		u.Permissao = user.PermUSR
	}
	if err := h.Svc.Criar(r.Context(), &u); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusCreated, u)
}

// GET /usuarios/buscar-tudo (ADM) — paginação e filtros
type pageResp struct {
	Total  int           `json:"total"`
	Pagina int           `json:"pagina"`
	Limite int           `json:"limite"`
	Data   []user.Usuario `json:"data"`
}

func (h *UsersHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	q := r.URL.Query()
	pagina, _ := strconv.Atoi(q.Get("pagina"))
	if pagina <= 0 {
		pagina = 1
	}
	limite, _ := strconv.Atoi(q.Get("limite"))
	if limite <= 0 || limite > 100 {
		limite = 10
	}
	busca := q.Get("busca")
	status := q.Get("status")
	permissao := q.Get("permissao")
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

	response.JSON(w, http.StatusOK, pageResp{
		Total:  total,
		Pagina: pagina,
		Limite: limite,
		Data:   items,
	})
}

// GET /usuarios/buscar-por-id/:id (ADM)
func (h *UsersHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	id := lastSegment(r.URL.Path)
	u, err := h.Svc.BuscarPorID(r.Context(), id)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, "erro", nil)
		return
	}
	if u == nil {
		response.ErrorJSON(w, http.StatusNotFound, "Usuário não encontrado", nil)
		return
	}
	response.JSON(w, http.StatusOK, u)
}

// PATCH /usuarios/atualizar/:id (ADM/TEC/USR conforme regra)
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

// GET /usuarios/lista-completa (ADM)
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

// GET /usuarios/buscar-tecnicos (ADM)
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

// DELETE /usuarios/desativar/:id (ADM) — soft delete
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

// PATCH /usuarios/autorizar/:id (ADM)
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

// GET /usuarios/valida-usuario (qualquer autenticado)
func (h *UsersHandler) ValidaUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

// GET /usuarios/buscar-novo/:login (ADM) — fluxo LDAP
func (h *UsersHandler) BuscarNovo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	login := lastSegment(r.URL.Path)

	// 1) existe ativo?
	if u, _ := h.Svc.BuscarPorLogin(r.Context(), login); u != nil && u.Status {
		response.ErrorJSON(w, http.StatusForbidden, "Login já cadastrado", nil)
		return
	}

	// 2) existe inativo? reativar e retornar
	if u, _ := h.Svc.BuscarPorLogin(r.Context(), login); u != nil && !u.Status {
		_ = h.Svc.Atualizar(r.Context(), u.ID, &user.Usuario{Status: true})
		response.JSON(w, http.StatusOK, map[string]any{
			"login": u.Login,
			"nome":  u.Nome,
			"email": u.Email,
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
