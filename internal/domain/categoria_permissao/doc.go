// Package categoria_permissao contém a lógica de domínio relacionada às permissões
// de usuários sobre categorias de chamados.
//
// Uma CategoriaPermissao representa a permissão que um usuário possui sobre uma
// determinada categoria. Cada permissão define as ações que o usuário pode executar
// dentro daquela categoria.
//
// O pacote fornece:
//
//   - A entidade principal `CategoriaPermissao`, com validações e métodos de acesso;
//   - Estruturas de filtro (`Filtro`) para listagem de permissões, com suporte a
//     paginação e filtragem por categoria, usuário ou nível de permissão;
//   - Interface de repositório (`Repository`) para persistência e consulta;
//   - Interfaces de serviço (`Service`, `Leitor`, `Escritor`) para os casos de uso.
//
// # Validações de Domínio
//
// Principais validações da entidade CategoriaPermissao:
//
//   - CategoriaID e UsuarioID não podem ser vazios;
//   - Permissao deve ser válida, conforme definido no pacote `usuario`.
//
// Todas as validações utilizam erros estruturados do pacote `domain` (`ErrosValidacao`).
//
// # Estrutura da Entidade CategoriaPermissao
//
// A entidade CategoriaPermissao oferece:
//
//   - Novo():
//       Cria uma nova permissão aplicando as validações de domínio;
//
//   - AtualizarDados():
//       Atualiza o nível de permissão e a data de atualização;
//
//   - Métodos de acesso:
//       CategoriaID(), UsuarioID(), Permissao(), CriadoEm(), AtualizadoEm();
//
//   - String():
//       Retorna uma representação textual da permissão.
//
// # Filtros de Consulta
//
// A estrutura `Filtro` permite consultas flexíveis:
//
//   - Paginação;
//   - Filtragem por CategoriaID, UsuarioID ou nível de Permissao;
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
//       BuscarPorCategoria, BuscarPorUsuario, Listar;
//
//   - Escritor:
//       Criar, Atualizar;
//
//   - Service:
//       Agrega Leitor e Escritor.
//
// # Estrutura para Persistência
//
// A estrutura CategoriaPermissaoDB representa a permissão no banco de dados.
//
// # Exemplo de Uso
//
//	cp, err := categoria_permissao.Novo("ctg123", "usr456", usuario.PermTEC)
//	if err != nil {
//	    // tratar erro de validação
//	}
//
//	if err := repo.Criar(ctx, *cp); err != nil {
//	    // tratar erro de persistência
//	}
package categoria_permissao
