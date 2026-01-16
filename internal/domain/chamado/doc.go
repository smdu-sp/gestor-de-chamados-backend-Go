// Package chamado contém a lógica de domínio relacionada aos chamados do sistema.
//
// Um Chamado representa uma solicitação registrada no sistema, podendo ser criada por
// um usuário e atribuída a categorias e subcategorias específicas. Cada chamado possui
// status, solução opcional e campos de auditoria.
//
// O pacote fornece:
//
//   - A entidade principal `Chamado` com validações e métodos de acesso;
//   - Manipulação de status via `StatusChamado` (ABERTO, ATRIBUIDO, RESOLVIDO, REJEITADO, FECHADO);
//   - Estruturas de filtro (`Filtro`) para listagem de chamados, com suporte a paginação
//     e pesquisa por título, descrição, categoria, subcategoria ou criador;
//   - Interface de repositório (`Repository`) para persistência e consulta no banco de dados;
//   - Interfaces de serviço (`Service`, `Leitor`, `Escritor`) para os casos de uso.
//
// # Validações de Domínio
//
// Principais validações da entidade Chamado:
//
//   - Título: entre 10 e 255 caracteres;
//   - Descrição: entre 10 e 1000 caracteres;
//   - IDs de categoria, subcategoria e criador não podem ser vazios;
//   - Status deve ser um valor válido definido em `StatusChamado`.
//
// Todas as validações utilizam erros estruturados do pacote `domain` (`ErrosValidacao`).
//
// # Funções e Métodos
//
// A entidade Chamado oferece:
//
//   - Novo():
//       Cria um novo chamado aplicando validações de domínio;
//
//   - AtualizarDados():
//       Atualiza título, descrição, categoria e subcategoria;
//
//   - AtualizarStatus():
//       Altera o status do chamado (ABERTO, ATRIBUIDO, RESOLVIDO, REJEITADO, FECHADO);
//
//   - PodeSerAtendido():
//       Indica se o chamado está em estado elegível para atendimento;
//
//   - Métodos de acesso:
//       ID(), CategoriaID(), SubcategoriaID(), CriadorID(), Titulo(), Descricao(),
//       Status(), Solucao(), CriadoEm(), AtualizadoEm(), SolucionadoEm(), FechadoEm();
//
//   - String():
//       Retorna uma representação textual do chamado.
//
// # Filtros de Consulta
//
// A estrutura `Filtro` permite consultas flexíveis:
//
//   - Paginação (Pagina, Limite);
//   - Pesquisa por título ou descrição;
//   - Filtragem por CategoriaID, SubcategoriaID, CriadorID ou Status;
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
//       BuscarPorID, Listar, ListarPorCategoria, ListarPorCriador;
//
//   - Escritor:
//       Criar, AtualizarDados, AtualizarStatus, Arquivar, Desarquivar;
//
//   - Service:
//       Agrega Leitor e Escritor.
//
// # Estrutura de Persistência
//
// A estrutura ChamadoDB representa o chamado no banco de dados.
//
// # Exemplo de Uso
//
//	ch, err := chamado.Novo("id123", "ctg001", "sub001", "usr001", "Título do chamado", "Descrição do chamado")
//	if err != nil {
//	    // tratar erro de validação
//	}
//
//	if err := repo.Criar(ctx, *ch); err != nil {
//	    // tratar erro de persistência
//	}
package chamado
