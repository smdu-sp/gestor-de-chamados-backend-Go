// Package atendimento contém a lógica de domínio relacionada à atribuição
// de técnicos a chamados.
//
// Um atendimento representa a vinculação de um técnico responsável a um
// chamado, permitindo identificar e acompanhar quem está designado para
// cada chamado dentro do sistema.
//
// O pacote fornece:
//
//   - A entidade principal `Atendimento`, com validações e métodos de acesso;
//   - Estruturas de filtro (`Filtro`) para listagem de atendimentos,
//     com suporte a paginação;
//   - Interface de repositório (`Repository`) para persistência;
//   - Interfaces de serviço (`Service`, `Leitor`, `Escritor`) que definem
//     os casos de uso do domínio.
//
// # Validações de Domínio
//
// As principais validações aplicadas à entidade Atendimento incluem:
//
//   - O identificador do atendimento não pode ser vazio;
//   - O identificador do técnico atribuído (`AtribuidoID`) não pode ser vazio;
//   - O identificador do chamado (`ChamadoID`) não pode ser vazio.
//
// As validações utilizam erros estruturados do pacote `domain`,
// por meio do tipo `ErrosValidacao`.
//
// # Fluxo Típico de Uso
//
//   1. Criar um novo atendimento utilizando a função `Novo()`;
//   2. Persistir a entidade por meio de uma implementação de `Repository`;
//   3. Ler atendimentos utilizando a interface `Leitor` do serviço;
//   4. Atualizar campos mutáveis através de `AtualizarDados` e da
//      interface `Escritor` do serviço;
//   5. Aplicar filtros e paginação utilizando a estrutura `Filtro`.
//
// # Considerações de Implementação
//
//   - A entidade Atendimento possui campos imutáveis (`id`) e campos
//     mutáveis (`atribuidoID`, `chamadoID`, `atualizadoEm`);
//   - O pacote não possui dependência direta de infraestrutura ou HTTP;
//   - Permite consultas detalhadas através de filtros que combinam
//     `ChamadoID` e `AtribuidoID`.
//
// # Exemplo de Uso
//
//	a, err := atendimento.Novo(
//	    "id123",
//	    "tecnico456",
//	    "chamado789",
//	)
//	if err != nil {
//	    // tratar erro de validação
//	}
//
//	if err := repo.Criar(ctx, *a); err != nil {
//	    // tratar erro de persistência
//	}
package atendimento
