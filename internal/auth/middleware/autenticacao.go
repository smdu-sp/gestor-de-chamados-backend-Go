package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
)

var ErrFormatoCabecalhoInvalido = errors.New("formato do cabeçalho inválido")

type ctxKey string

const ChaveUsuario  ctxKey = "user"
const prefixoBearer  string = "Bearer "
const mensagemNaoAutorizado  = "Você não está autorizado a acessar este recurso"

// AutenticarUsuario autentica o usuário e atualiza o último login diretamente via svc
func AutenticarUsuario(next http.Handler, gJWT jwt.JWTUsecase, usecase usecase.UsuarioUsecase) http.Handler {
	// Retorna um handler que autentica o usuário
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Verifica o cabeçalho Authorization
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, prefixoBearer) {
			response.ErrorJSON(w, http.StatusUnauthorized, mensagemNaoAutorizado, nil)
			return
		}

		// Divide o cabeçalho em partes
		partes := strings.SplitN(auth, " ", 2)
		if len(partes) != 2 || strings.TrimSpace(partes[1]) == "" {
			response.ErrorJSON(w, http.StatusUnauthorized, mensagemNaoAutorizado, ErrFormatoCabecalhoInvalido.Error())
			return
		}

		// Extrai e valida o token
		token := strings.TrimSpace(partes[1])
		claims, err := gJWT.ValidarToken(token)
		if err != nil {
			response.ErrorJSON(w, http.StatusUnauthorized, mensagemNaoAutorizado, err.Error())
			return
		}

		// Atualiza último login do usuário
		if claims.ID != "" {
			_ = usecase.AtualizarUltimoLoginUsuario(r.Context(), claims.ID)
			// Se der erro, ignora (não é crítico)
		}

		// Adiciona os claims ao contexto e chama o próximo handler
		ctx := context.WithValue(r.Context(), ChaveUsuario, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UsuarioFromCtx retorna claims do usuário presente no contexto
func UsuarioFromCtx(r *http.Request) *jwt.Claims {
	if valor := r.Context().Value(ChaveUsuario); valor != nil {
		if claims, ok := valor.(*jwt.Claims); ok {
			return claims
		}
	}
	return nil
}
