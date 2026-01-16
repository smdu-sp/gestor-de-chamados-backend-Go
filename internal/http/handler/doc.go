// Package handler fornece os handlers HTTP para cada entidade do sistema.
//
// Os handlers encapsulam a lógica de recepção de requisições HTTP, validação de dados,
// decodificação de payloads JSON, tratamento de erros, logging estruturado e envio de
// respostas padronizadas em JSON.
//
// # Responsabilidades
//
// Para cada entidade (ex.: Usuario, Chamado, Categoria, Subcategoria, etc.), o pacote:
//
//   - Recebe requisições HTTP e extrai parâmetros da URL, query ou body;
//   - Valida dados de entrada e aplica regras de negócio básicas;
//   - Converte DTOs para parâmetros de serviço;
//   - Invoca serviços (Service, Leitor, Escritor) para processar as operações;
//   - Trata erros estruturados e envia respostas JSON padronizadas;
//   - Registra logs estruturados para auditoria e troubleshooting.
//
// # Características Principais
//
//   - Context com timeout para cada requisição;
//   - Tratamento detalhado de erros (validação, conflito, timeout, cancelamento, erros internos);
//   - Logging estruturado das operações críticas;
//   - Respostas JSON padronizadas, incluindo erros conforme RFC 7807;
//   - Reutilização de utilitários para decodificação e construção de respostas.
//
// # Estrutura dos Handlers
//
// Cada handler depende de:
//   - Serviços de domínio correspondentes à entidade (ex.: usr.Service);
//   - Serviços auxiliares como AuthInternoService, AuthExternoService e log.Service;
//   - DTOs para decodificação e resposta;
//   - Pacotes de utilitários para parsing de query params e respostas HTTP.
//
// # Utilitários Disponíveis
//
//   - DecodificarJSON: decodifica body de requisição de forma segura;
//   - ResponderJSONStatusOK / ResponderJSONCreated / ResponderJSONErro: envia respostas HTTP padronizadas;
//   - Query parser: auxilia na extração e conversão de parâmetros de query.
//
// # Exemplo de Uso
//
//	import (
//	    "net/http"
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/handler"
//	)
//
//	func main() {
//	    usrHandler := handler.NovoUsuarioHandler(usuarioSvc, authSvc, ldapSvc, logSvc)
//	    http.HandleFunc("/usuarios", usrHandler.Listar)
//	    http.ListenAndServe(":8080", nil)
//	}
//
// # Resumo
//
// O pacote handler oferece uma camada de interface HTTP completa, organizada por entidade,
// garantindo tratamento seguro e padronizado de requisições, erros, logs e respostas JSON,
// separando claramente a camada de transporte da lógica de negócio.
package handler
