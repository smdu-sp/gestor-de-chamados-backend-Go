// Package middleware fornece middlewares HTTP reutilizáveis para aplicações Go.
//
// Middlewares são funções que interceptam requisições HTTP antes ou depois do handler principal,
// permitindo adicionar funcionalidades transversais como autenticação, logging, CORS ou tratamento de erros.
//
// # Visão Geral
//
// Este pacote inclui middlewares para:
//
//   - CORS:
//       - Adiciona cabeçalhos CORS às respostas HTTP;
//       - Permite configurar origens, métodos e headers aceitos;
//       - Responde automaticamente a requisições OPTIONS com status 204.
//
//   - Recuperação de pânico:
//       - Captura pânicos (panic) em handlers HTTP;
//       - Registra stack trace e detalhes do erro;
//       - Retorna resposta JSON padronizada com status 500.
//
// # Uso
//
// Para aplicar middlewares:
//
//	import (
//	    "net/http"
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/middleware"
//	)
//
//	func main() {
//	    mux := http.NewServeMux()
//	    mux.HandleFunc("/usuarios", usuariosHandler.Listar)
//
//	    // Aplica middlewares
//	    handler := middleware.RecuperarDePanico(mux)
//	    handler = middleware.CORS("https://meusite.com")(handler)
//
//	    http.ListenAndServe(":8080", handler)
//	}
//
// # Boas Práticas
//
//   - Middlewares devem ser aplicados na ordem correta (ex.: recuperação de pânico primeiro, CORS depois);
//   - Evitar lógica de negócio dentro de middlewares;
//   - Sempre retornar respostas JSON padronizadas para erros capturados;
//   - Garantir que middlewares não alterem indevidamente o contexto da requisição.
//
// # Considerações Finais
//
// O pacote middleware oferece uma camada de interceptação HTTP modular e reutilizável,
// facilitando a adição de funcionalidades transversais como CORS e tratamento de pânicos,
// mantendo os handlers limpos e focados na lógica de negócio.
package middleware
