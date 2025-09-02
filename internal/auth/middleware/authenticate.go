package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/user"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/response"
)

type ctxKey string

const UserKey ctxKey = "user"
const bearerPrefix string = "Bearer "
const unauthorizedMessage = "Você não está autorizado a acessar este recurso"

// AuthenticateUser autentica o usuário e atualiza o último login diretamente via svc
func AuthenticateUser(next http.Handler, jm *jwt.Manager, svc user.UserServiceInterface) http.Handler {
	// Retorna um handler que autentica o usuário
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Verifica o cabeçalho Authorization
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, bearerPrefix) {
			response.ErrorJSON(w, http.StatusUnauthorized, "unauthorized", unauthorizedMessage)
			return
		}

		// Divide o cabeçalho em partes
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			response.ErrorJSON(w, http.StatusUnauthorized, "unauthorized", unauthorizedMessage)
			return
		}

		// Extrai e valida o token
		token := strings.TrimSpace(parts[1])
		claims, err := jm.ParseAccess(token)
		if err != nil {
			response.ErrorJSON(w, http.StatusUnauthorized, "unauthorized", unauthorizedMessage)
			return
		}

		// Atualiza último login do usuário
		if claims.ID != "" {
			if err := svc.AtualizarUltimoLogin(r.Context(), claims.ID); err != nil {
				log.Printf("[DEBUG] erro ao atualizar último login: %v", err)
			}
		}

		// Adiciona os claims ao contexto e chama o próximo handler
		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromCtx retorna claims do usuário presente no contexto
func UserFromCtx(r *http.Request) *jwt.Claims {
	if valor := r.Context().Value(UserKey); valor != nil {
		if claims, ok := valor.(*jwt.Claims); ok {
			return claims
		}
	}
	return nil
}
