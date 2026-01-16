// Package service implementa a camada de casos de uso da aplicação.
//
// O pacote service é responsável por orquestrar regras de negócio, validações e operações
// entre diferentes entidades do sistema, como chamados, categorias, subcategorias, atendimentos
// e usuários.
//
// # Visão Geral
//
// 1. Serviços por entidade:
//
//   Cada entidade possui seu serviço dedicado, responsável por coordenar operações,
//   validações e persistência via repositórios:
//
//     - ChamadoService
//     - CategoriaService
//     - SubcategoriaService
//     - UsuarioService
//     - AtendimentoService
//     - AcompanhamentoService
//     - CategoriaPermissaoService
//
// 2. Validações e regras de negócio:
//
//   - Checagem de permissões de usuários (criador, técnico, administrador).
//   - Verificação de status de entidades antes de operações (ativo/inativo, arquivado, fechado).
//   - Consistência entre entidades relacionadas (ex.: categoria e subcategoria).
//   - Regras específicas por operação, como atualizar solução ou status de chamado.
//
// 3. Auxiliares:
//
//   - Geração de IDs únicos usando UUID v7:
//       - GerarUUIDv7String()
//       - GerarUUIDv7Bytes()
//   - Métodos internos para validações de consistência entre entidades.
//
// # Exemplos de Uso
//
//	import (
//	    "context"
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/service"
//	    chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
//	)
//
//	func main() {
//	    ctx := context.Background()
//	    chamadoService := service.NovoChamadoService(
//	        repo, categoriaRepo, subcategoriaRepo, categoriaPermRepo, atendimentoRepo,
//	    )
//
//	    // Criar um chamado
//	    c, err := chamadoService.Criar(ctx, chm.CriarParams{
//	        CategoriaID:    "c123",
//	        SubcategoriaID: "sc456",
//	        CriadorID:      "u789",
//	        Titulo:         "Problema de login",
//	        Descricao:      "Usuário não consegue logar no sistema",
//	    })
//	    if err != nil {
//	        // tratar erro
//	    }
//
//	    // Atualizar status de um chamado
//	    _, err = chamadoService.AtualizarStatus(ctx, c.ID(), chm.AtualizarStatusParams{
//	        Status:  chm.StatusEmAndamento,
//	        Solucao: nil,
//	    })
//	    if err != nil {
//	        // tratar erro
//	    }
//	}
//
// # Benefícios
//
//   - Serviços centralizam regras de negócio, garantindo consistência entre entidades.
//   - Permite integração transparente com repositórios para leitura e escrita.
//   - Todas as operações usam context.Context para controle de timeout e cancelamento.
//   - Métodos seguem contratos claros de entrada e saída, retornando erros estruturados
//     do pacote domain.
//   - Facilita testes unitários e de integração, desacoplando lógica de negócio de
//     HTTP e persistência.
package service
