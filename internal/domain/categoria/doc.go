// Package categoria contém a lógica de domínio relacionada às categorias
// de chamados do sistema.
//
// Uma Categoria representa uma classificação de chamados, possuindo
// nome, status (ativo/inativo) e timestamps de criação e atualização.
//
// O pacote fornece:
//
//   - A entidade principal `Categoria`, com validações e métodos
//     para criação, atualização e ativação/desativação;
//   - Estruturas de filtro (`Filtro`) para listagem de categorias,
//     com suporte a paginação, busca textual e filtragem por status;
//   - Interface de repositório (`Repository`) para persistência e consulta;
//   - Interfaces de serviço (`Service`, `Leitor`, `Escritor`) que definem
//     os casos de uso do domínio.
//
// # Validações de Domínio
//
// As principais validações aplicadas à entidade Categoria incluem:
//
//   - O identificador da categoria não pode ser vazio;
//   - O nome deve possuir entre 3 e 100 caracteres.
//
// As validações utilizam erros estruturados do pacote `domain`,
// por meio do tipo `ErrosValidacao`.
//
// # Estrutura da Entidade Categoria
//
// A entidade Categoria expõe as seguintes operações:
//
//   - Novo():
//       Cria uma nova categoria aplicando as validações de domínio.
//
//   - CarregarDoBD():
//       Cria uma instância da categoria a partir de dados persistidos.
//
//   - AtualizarDados():
//       Permite atualizar o nome e o status da categoria.
//
//   - Ativar() / Desativar():
//       Alteram explicitamente o status da categoria.
//
//   - Métodos de acesso:
//       ID(), Nome(), Status(), CriadoEm(), AtualizadoEm().
//
//   - String():
//       Retorna uma representação textual da categoria.
//
// # Filtros de Consulta
//
// A estrutura `Filtro` permite a construção de consultas flexíveis,
// suportando:
//
//   - Paginação;
//   - Busca textual por nome;
//   - Filtragem por status (ativo/inativo).
//
// Métodos auxiliares incluem:
//
//   - Pagina(), Limite(), Offset();
//   - Normalizar();
//   - SemLimite().
//
// # Interfaces de Repositório e Serviço
//
// O pacote define as seguintes interfaces:
//
//   - Leitor:
//       BuscarPorID, BuscarPorNome, Listar;
//
//   - Escritor:
//       Criar, Atualizar, Ativar, Desativar;
//
//   - Service:
//       Agrega as interfaces Leitor e Escritor.
//
// # Estrutura para Persistência
//
// A estrutura CategoriaDB representa o formato da categoria
// utilizada na camada de persistência.
//
// # Exemplo de Uso
//
//	filtro := categoria.NovoFiltro(
//	    domain.Paginacao{Pagina: 1, Limite: 10},
//	    nil,
//	    nil,
//	)
//
//	c, err := categoria.Novo("c123", "Suporte")
//	if err != nil {
//	    // tratar erro de validação
//	}
//
//	if err := repo.Criar(ctx, *c); err != nil {
//	    // tratar erro de persistência
//	}
package categoria
