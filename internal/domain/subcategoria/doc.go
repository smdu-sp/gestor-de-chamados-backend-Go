// Package subcategoria gerencia subcategorias de chamados no sistema.
//
// Uma Subcategoria representa uma classificação secundária de chamados vinculada a uma categoria.
// Cada subcategoria possui um nome, status (ativa ou inativa) e referência à categoria pai.
//
// O pacote fornece:
//
//   - A entidade principal `Subcategoria` com validação de campos e campos de auditoria (criadoEm, atualizadoEm);
//   - Estruturas de filtro (`Filtro`) para listagem de subcategorias, com suporte a paginação, busca textual e filtragem por status;
//   - Interface de repositório (`Repository`) para persistência e consulta no banco de dados;
//   - Interfaces de serviço (`Service`, `Leitor`, `Escritor`) para casos de uso.
//
// # Validações de Domínio
//
// Principais validações da entidade Subcategoria:
//
//   - ID da subcategoria não pode estar vazio;
//   - Nome deve ter entre 3 e 100 caracteres;
//   - CategoriaID não pode estar vazio.
//
// Todas as validações utilizam erros estruturados do pacote `domain` (`ErrosValidacao`).
//
// # Descrição das Funcionalidades
//
// A entidade Subcategoria oferece:
//
//   - Novo():
//       Cria uma nova subcategoria aplicando validações de domínio;
//
//   - CarregarDoBD():
//       Cria instâncias a partir de dados do banco;
//
//   - AtualizarDados():
//       Permite atualizar nome, categoria e status;
//
//   - Ativar() / Desativar():
//       Alteram o status da subcategoria;
//
//   - Métodos de acesso:
//       ID(), Nome(), CategoriaID(), Status(), CriadoEm(), AtualizadoEm();
//
//   - String():
//       Retorna uma representação textual da subcategoria.
//
// # Filtros de Consulta
//
// A estrutura `Filtro` permite consultas flexíveis:
//
//   - Paginação (Pagina, Limite);
//   - Busca textual por nome;
//   - Filtragem por status (ativa/inativa);
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
//       BuscarPorID, BuscarPorNome, Listar;
//
//   - Escritor:
//       Criar, Atualizar, Ativar, Desativar;
//
//   - Service:
//       Agrega Leitor e Escritor.
//
// # Estrutura de Banco de Dados
//
// A estrutura SubcategoriaDB representa a subcategoria no banco de dados.
//
// # Exemplo de Uso
//
//	sc, err := subcategoria.Novo("sc123", "Suporte Técnico", "cat001")
//	if err != nil {
//	    // tratar erro de validação
//	}
//
//	if err := repo.Criar(ctx, *sc); err != nil {
//	    // tratar erro de persistência
//	}
package subcategoria
