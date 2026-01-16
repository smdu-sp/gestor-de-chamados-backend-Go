// Package acompanhamento contém a lógica de domínio relacionada aos
// acompanhamentos de chamados.
//
// Um acompanhamento representa um comentário ou atualização registrada
// em um chamado dentro do sistema. Ele é sempre associado a um chamado,
// a um usuário e possui um remetente, representado pelo tipo de permissão
// do usuário que realizou a ação.
//
// O pacote fornece:
//
//   - A entidade principal `Acompanhamento`, com validações e métodos
//     de acesso aos seus dados;
//   - Estruturas de filtro (`Filtro`) para listagem de acompanhamentos,
//     com suporte a paginação;
//   - Interface de repositório (`Repository`) para persistência;
//   - Interfaces de serviço (`Service`, `Leitor`, `Escritor`) que definem
//     os casos de uso do domínio.
//
// # Validações de Domínio
//
// As principais validações aplicadas à entidade Acompanhamento incluem:
//
//   - O conteúdo deve possuir entre 1 e 1000 caracteres;
//   - O remetente deve ser um tipo de permissão válido
//     (ex: usr.PermTEC ou usr.PermUSR);
//   - Os identificadores do acompanhamento, chamado e usuário
//     não podem ser vazios.
//
// As validações utilizam erros estruturados do pacote `domain`,
// por meio do tipo `ErrosValidacao`.
//
// # Fluxo de Uso
//
//   1. Criar um novo acompanhamento utilizando a função `Novo()`;
//   2. Persistir a entidade por meio de uma implementação de `Repository`;
//   3. Ler acompanhamentos utilizando a interface `Leitor` do serviço;
//   4. Atualizar campos mutáveis através de `AtualizarDados` e da
//      interface `Escritor` do serviço;
//   5. Aplicar filtros e paginação utilizando a estrutura `Filtro`.
//
// # Considerações de Implementação
//
//   - A entidade Acompanhamento possui campos imutáveis
//     (`id`, `chamadoID`, `usuarioID`) e campos mutáveis
//     (`conteudo`, `atualizadoEm`);
//   - Este pacote depende do pacote `usuario` para validação
//     das permissões de remetente;
//   - Não possui dependência direta de infraestrutura ou HTTP.
//
// # Exemplo de Uso
//
//	a, err := acompanhamento.Novo(
//	    "id123",
//	    "chamado456",
//	    "usuario789",
//	    "Comentário de teste",
//	    usr.PermTEC,
//	)
//	if err != nil {
//	    // tratar erro de validação
//	}
//
//	if err := repo.Criar(ctx, *a); err != nil {
//	    // tratar erro de persistência
//	}
package acompanhamento
