// Package log gerencia o registro de logs de ações realizadas no sistema.
//
// Um Log representa uma entrada de auditoria que registra ações de usuários sobre entidades do sistema,
// incluindo o tipo de ação, a entidade afetada, detalhes e timestamp de criação.
//
// O pacote fornece:
//
//   - A entidade principal `Log` com validações e campos de auditoria (usuarioID, acao, entidade, detalhes, criadoEm);
//   - Tipos de ações definidas em `Acao` (CRIAR, ATUALIZAR, DESATIVAR, ATIVAR, ARQUIVAR, DESARQUIVAR, DELETAR);
//   - Estruturas de filtro (`LogFiltro`) para listagem de logs, com suporte a paginação, busca textual, intervalo de datas,
//     usuário, ação ou entidade;
//   - Interface de repositório (`Repository`) para persistência e consulta no banco de dados;
//   - Interfaces de serviço (`Service`, `Leitor`, `Escritor`) para casos de uso.
//
// # Validações de Domínio
//
// Principais validações da entidade Log:
//
//   - ID não pode estar vazio;
//   - usuarioID não pode estar vazio;
//   - Ação deve ser válida, conforme definido em `Acao`;
//   - Entidade e detalhes não podem estar vazios.
//
// Todas as validações utilizam erros estruturados do pacote `domain` (`ErrosValidacao`).
//
// # Descrição Detalhada
//
// A entidade Log oferece:
//
//   - Novo():
//       Cria um novo log aplicando validações de domínio;
//
//   - ValidarAcao():
//       Verifica se a ação informada é válida;
//
//   - AcoesValidas():
//       Retorna todas as ações válidas para referência;
//
//   - Métodos de acesso:
//       ID(), UsuarioID(), Acao(), Entidade(), Detalhes(), CriadoEm();
//
//   - String():
//       Retorna uma representação textual do log.
//
// # Filtros de Consulta
//
// A estrutura `LogFiltro` permite consultas flexíveis:
//
//   - Paginação (Pagina, Limite);
//   - Intervalo de datas (DataInicio, DataFim);
//   - Filtro por usuário, ação ou entidade;
//   - Busca textual em detalhes.
//
// Métodos auxiliares:
//
//   - Pagina(), Limite(), Offset();
//   - Normalizar();
//   - SemLimite().
//
// # Repositório
//
// Interfaces definidas:
//
//   - Leitor:
//       BuscarPorID, Listar, ListarPorUsuario, ListarPorAcao, ListarPorEntidade;
//
//   - Escritor:
//       Criar;
//
//   - Service:
//       Agrega Leitor e Escritor.
//
// # Representação no Banco de Dados
//
// A estrutura LogDB representa o log no banco de dados.
//
// # Exemplo de Uso
//
//	l, err := log.Novo("id123", "usr001", log.Criar, "Chamado", "Criado novo chamado")
//	if err != nil {
//	    // tratar erro de validação
//	}
//
//	if err := repo.Criar(ctx, *l); err != nil {
//	    // tratar erro de persistência
//	}
package log
