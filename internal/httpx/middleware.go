package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"
)

// Chave para o context
type ctxKey string

// Chave usada para armazenar os claims JWT no context da requisição
const ctxClaims ctxKey = "claims"

// Envia um objeto como resposta HTTP em formato JSON
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Envia uma resposta de erro em JSON com o status e mensagem fornecidos
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

// Verifica se o header Authorization está presente e válido (Bearer <token>)
// Decodifica o token e adiciona os claims no context da requisição
func AuthMiddleware(tm *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authz := r.Header.Get("Authorization")

			if authz == "" {
				Error(w, http.StatusUnauthorized, "token ausente")
				return
			}

			parts := strings.SplitN(authz, " ", 2)

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				Error(w, http.StatusUnauthorized, "formato de autorização inválido")
				return
			}

			claims, err := tm.Parse(parts[1])

			if err != nil {
				Error(w, http.StatusUnauthorized, "token inválido")
				return
			}

			// Armazena os claims no context e passa para o próximo handler
			ctx := context.WithValue(r.Context(), ctxClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRoles é um middleware de RBAC (controle de acesso por papéis)
// Permite apenas usuários com um dos papéis especificados
func RequireRoles(allowed ...model.Permissao) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := r.Context().Value(ctxClaims)

			if v == nil {
				Error(w, http.StatusUnauthorized, "não autenticado")
				return
			}

			c := v.(*auth.Claims)

			if !auth.IsAllowed(c.Permissao, allowed...) {
				Error(w, http.StatusForbidden, "acesso negado")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Middleware que captura panics e evita crash do servidor
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {

			if rec := recover(); rec != nil {
				Error(w,
					http.StatusInternalServerError, "erro interno")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Cria um middleware que define tempo máximo de execução de uma requisição
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, "request timeout")
	}
}

// Permite encadear múltiplos middlewares de forma simples
// Aplica os middlewares da última para a primeira (como decoradores)
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Helper para extrair os claims JWT do context da requisição
func ClaimsFrom(r *http.Request) (*auth.Claims, error) {
	v := r.Context().Value(ctxClaims)

	if v == nil {
		return nil, errors.New("sem claims")
	}

	c, ok := v.(*auth.Claims)
	
	if !ok {
		return nil, errors.New("claims inválidas")
	}
	return c, nil
}
