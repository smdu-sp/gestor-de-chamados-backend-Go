package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// Erros relacionados à autenticação
var (
	ErrCabecalhoAusente         = errors.New("cabeçalho Authorization não informado")
	ErrFormatoCabecalhoInvalido = errors.New("formato do cabeçalho inválido")
	ErrUsuarioNaoAutenticado    = errors.New("usuário não autenticado")
	ErrUsuarioNaoAutorizado     = errors.New("usuário não autorizado")
	ErrTokenInvalido            = errors.New("token inválido")
)

const (
	prefixoBearer string = "Bearer " // com espaço no final!
	cabecalhoAuth string = "Authorization"
)

// AutenticarUsuario é um middleware que recebe um http.Handler, um JWTService e um usr.Service,
// e retorna um http.Handler que autentica o usuário usando um token JWT.
//
// Erros sentinela que podem ser retornados: ErrCabecalhoAusente, ErrFormatoCabecalhoInvalido, ErrTokenInvalido.
func AutenticarUsuario(next http.Handler, jwtSvc JWTService, usuarioSvc usr.Service) http.Handler {

	// Retorna um handler que autentica o usuário
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extrai o token do cabeçalho Authorization
		token, err := ExtrairTokenDoRequest(r)
		if err != nil {
			ResponderJSONUnauthorized(w, r, err.Error())
			return
		}

		// Valida o token JWT
		claims, err := jwtSvc.ValidarToken(token)
		if err != nil {
			ResponderJSONUnauthorized(w, r, ErrTokenInvalido.Error())
			return
		}

		// Atualiza último login do usuário
		if claims.ID != "" {
			if _, err = usuarioSvc.AtualizarUltimoLogin(r.Context(), claims.ID); err != nil {
				slog.Error("Autenticar usuário:", "error", err)
			}
		}

		// Adiciona os claims ao contexto e chama o próximo handler
		ctx := context.WithValue(r.Context(), chaveUsuario, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ExtrairTokenDoRequest recebe um *http.Request, extrai o token JWT do cabeçalho Authorization e o retorna como string.
//
// Erros sentinela possíveis: ErrCabecalhoAusente ou ErrFormatoCabecalhoInvalido.
func ExtrairTokenDoRequest(r *http.Request) (string, error) {
	// Obtém o valor do cabeçalho Authorization
	auth := r.Header.Get(cabecalhoAuth)
	if !strings.HasPrefix(auth, prefixoBearer) {
		return "", ErrCabecalhoAusente
	}

	// Separa o token do prefixo "Bearer "
	partes := strings.SplitN(auth, " ", 2)
	if len(partes) != 2 || strings.TrimSpace(partes[1]) == "" {
		return "", ErrFormatoCabecalhoInvalido
	}

	// Extrai o token do cabeçalho Authorization
	token := strings.TrimSpace(partes[1])
	return token, nil
}
