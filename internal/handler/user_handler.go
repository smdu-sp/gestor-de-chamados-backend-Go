package handler

import (
	"net/http"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/httpx"
)

type UserHandler struct{}

func (h *UserHandler) AdminOnly(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"ok": "rota ADM"})
}

func (h *UserHandler) SuporteOuDev(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"ok": "rota SUP/DEV"})
}

func (h *UserHandler) Publico(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"hello": "mundo"})
}
