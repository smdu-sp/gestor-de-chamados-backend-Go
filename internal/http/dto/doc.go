// Package dto contém os objetos de transferência de dados (DTOs) utilizados no sistema.
//
// DTOs (Data Transfer Objects) são usados para estruturar os dados que entram ou saem da aplicação,
// especialmente em camadas de API, facilitando a conversão entre entidades do domínio e formatos
// de comunicação, como JSON.
//
// # Estrutura do Pacote
//
// Para cada entidade do domínio (ex.: Chamado, Usuario, Categoria, Subcategoria, etc.), existem:
//
//   1. Arquivo DTO principal:
//      - Define as estruturas de request e response (Req / Resp).
//      - Contém funções auxiliares para converter entre DTO e parâmetros de criação/atualização da entidade.
//
//   2. Map.go:
//      - Função genérica MapearSlice, que facilita a conversão de slices de entidades em slices de DTOs.
//
//   3. Paginacao.go:
//      - Estrutura RespPaginada[T] para respostas paginadas genéricas.
//
// # Convenções de Nomenclatura
//
// Os DTOs seguem nomenclatura consistente:
//
//   - [Entidade]Req: para payloads de criação ou atualização.
//   - [Entidade]Resp: para payloads de resposta.
//   - [Entidade]Paginado: para respostas paginadas.
//
// # Exemplos de Uso
//
// Exemplo de conversão de DTOs:
//
//	// Convertendo um request HTTP em parâmetros de criação de chamado
//	var req dto.CriarChamadoReq
//	params := req.ParaCriarParams()
//
//	// Convertendo uma entidade Chamado em DTO de response
//	chamadoResp := dto.ParaChamadoResp(chamado)
//
//	// Convertendo uma lista de chamados para resposta paginada
//	chamadosResp := dto.ParaChamadosResp(listaChamados)
//
// # Vantagens do Uso de DTOs
//
//   - Separa o modelo de domínio da camada de transporte (API/JSON);
//   - Facilita validações e transformações antes de persistir dados;
//   - Permite reutilização de estruturas de forma consistente entre diferentes endpoints;
//   - Suporta facilmente paginação genérica e respostas padronizadas.
//
// # Resumo
//
// O pacote dto fornece uma camada intermediária entre a API e o domínio, garantindo que
// dados recebidos ou enviados estejam estruturados, validados e consistentes, além de
// permitir fácil conversão para entidades do domínio e vice-versa.
package dto
