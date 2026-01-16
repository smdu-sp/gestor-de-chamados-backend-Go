// Package router fornece a configuração centralizada das rotas HTTP da aplicação.
//
// Este pacote organiza o registro de rotas por entidade, permitindo separar responsabilidades
// e aplicar middlewares específicos como autenticação JWT, controle de permissões e CORS.
//
// # Visão Geral
//
// 1. Registro de rotas por entidade:
//
// Cada entidade possui uma função dedicada para registrar suas rotas, por exemplo:
//
//   - RegistrarRotasChamado: rotas para operações de chamados
//   - RegistrarRotasUsuario: rotas para operações de usuários
//   - RegistrarRotasCategoria, RegistrarRotasSubcategoria, etc.
//
// Essas funções recebem um http.ServeMux, os handlers correspondentes e serviços
// de autenticação/permissão, garantindo controle de acesso por rota.
//
// 2. Rotas públicas:
//
// Algumas rotas não exigem autenticação, incluindo:
//
//   - /login
//   - /refresh
//   - /swagger/*
//   - /docs/*
//   - /health
//
// 3. Roteador principal:
//
// A função NovoRoteador cria o http.Handler principal que combina:
//
//   - ServeMux para rotas públicas
//   - ServeMux para rotas protegidas
//   - Middlewares globais, como CORS e recuperação de pânico
//
// 4. Rotas auxiliares:
//
//   - Swagger: documentação da API via /swagger/
//   - Health check: endpoint /health para monitoramento
//   - Rotas de autenticação: /login, /refresh, /eu
//
// # Boas Práticas
//
//   - Manter o registro de rotas por entidade para maior modularidade;
//   - Rotas protegidas devem sempre validar JWT e permissões;
//   - Rotas públicas devem ser limitadas a endpoints que não exigem autenticação;
//   - Middlewares globais devem ser aplicados a todas as rotas para consistência;
//
// # Exemplo de Uso
//
//	import (
//	    "net/http"
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/container"
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/router"
//	)
//
//	func main() {
//	    cfg := config.Carregar()
//	    cont := container.NovoContainer(cfg)
//
//	    roteador := router.NovoRouter(cont, cfg)
//	    http.ListenAndServe(":8080", roteador)
//	}
//
// # Resumo
//
// O pacote router fornece uma configuração centralizada e modular das rotas HTTP,
// aplicando middlewares, autenticação e controle de permissões de forma consistente.
// Ele separa claramente rotas públicas e protegidas, permitindo fácil manutenção
// e extensão do sistema.
package router
