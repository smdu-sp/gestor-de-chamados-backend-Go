// Package id fornece implementações para geração de identificadores
// únicos utilizados no domínio da aplicação.
//
// Este pacote é responsável por gerar IDs de forma consistente,
// previsível e desacoplada da infraestrutura, permitindo que o
// domínio dependa apenas de abstrações.
//
// A implementação atual utiliza UUID versão 7 (RFC 9562), que combina
// timestamp em milissegundos com aleatoriedade, garantindo ordenação
// temporal aproximada e unicidade mesmo sob alta concorrência.
//
// # Design e Funcionalidades
//
// 1. Geração de Identificadores:
//
//   - Gerador:
//       Implementa a interface dmn.GeradorID definida no domínio,
//       permitindo injeção de dependência e substituição da estratégia
//       de geração de IDs sem impacto no código de domínio.
//
//   - NovoID:
//       Retorna um identificador único no formato string.
//       Atualmente gera UUIDv7, mas pode ser alterado para outro formato
//       (ex.: ULID, KSUID, Snowflake) sem modificar o domínio.
//
// 2. Implementação de UUIDv7 (RFC 9562):
//
//   Arquivo uuidv7.go contém a implementação completa:
//
//   - GerarUUIDv7Bytes:
//       Gera um UUIDv7 monotônico e RFC-compliant em formato binário,
//       combinando timestamp (48 bits), sequência (12 bits) e entropia aleatória.
//
//   - GerarUUIDv7String:
//       Converte o UUIDv7 gerado para o formato string padrão (8-4-4-4-12).
//
//   - UUIDv7.String:
//       Realiza a serialização do UUID para string hexadecimal.
//
//   Características da implementação:
//     - Thread-safe via mutex.
//     - Monotonicidade garantida no mesmo milissegundo.
//     - Variante RFC 4122 correta.
//     - Uso de crypto/rand para entropia criptograficamente segura.
//
// # Justificativa da Localização do Pacote
//
//   - Este pacote pertence à camada de infraestrutura.
//   - Implementa uma interface do domínio (dmn.GeradorID).
//   - Evita que o domínio conheça detalhes de UUID, timestamps ou RNG.
//   - Facilita testes, mocks e troca futura do algoritmo de geração.
//
// Em uma arquitetura Hexagonal / Clean Architecture:
//
//   Domínio  ---> interface GeradorID
//                  ▲
//                  │
//           identity.Generator (implementação concreta)
//
// # Exemplo de Uso
//
//	O exemplo abaixo demonstra como utilizar o gerador de IDs
//
//	import (
//	    "fmt"
//
//	    dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/identity"
//	)
//
//	func main() {
//	    var gerador dmn.GeradorID = identity.Generator{}
//
//	    id, err := gerador.NovoID()
//	    if err != nil {
//	        panic(err)
//	    }
//
//	    fmt.Println(id)
//	}
//
// # Vantagens
//
//   - O domínio nunca deve instanciar diretamente UUIDs ou depender
//     de bibliotecas externas de geração de IDs.
//   - A troca do algoritmo de geração requer apenas uma nova implementação
//     da interface GeradorID.
//   - UUIDv7 é particularmente adequado para bancos de dados,
//     pois melhora localidade de índice em comparação ao UUIDv4.
//   - Este pacote não possui dependências de banco de dados ou I/O.
package id