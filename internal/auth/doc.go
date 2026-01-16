// Package auth implementa a camada de autenticação e autorização da
// aplicação Gestor de Chamados.
//
// O pacote auth centraliza todos os mecanismos relacionados a:
//   - Autenticação interna (usuário/senha)
//   - Autenticação externa (LDAP ou outros provedores)
//   - Geração e validação de tokens JWT (access e refresh)
//   - Extração e manipulação de claims no contexto da requisição
//   - Autorização baseada em permissões
//
// O pacote é utilizado tanto pela camada de handlers HTTP quanto por
// middlewares de autenticação e autorização.
//
// # Autenticação Interna
//
// A autenticação interna é orquestrada pelo tipo AuthInterno, que atua
// como o caso de uso central de autenticação da aplicação.
//
// O fluxo de login contempla:
//  1. Validação do payload de entrada (LoginReq)
//  2. Autenticação via LDAP ou credenciais internas
//  3. Criação do usuário local, se necessário
//  4. Geração de tokens JWT (access e refresh)
//
// O fluxo de refresh contempla:
//  1. Validação do refresh token
//  2. Geração de um novo par de tokens
//
// O fluxo "Me" contempla:
//  1. Extração do usuário autenticado a partir das claims do token
//  2. Retorno dos dados completos do usuário
//
// # Autenticação Externa via LDAP
//
// O arquivo `ldap.go` define o serviço LDAPService, responsável por
// autenticar usuários contra servidores externos (AD, OpenLDAP, etc.).
//
// Funcionalidades principais:
//   - Bind: valida usuário e senha no servidor LDAP
//   - Pesquisar: realiza consultas genéricas no diretório
//   - PesquisarPorLogin: busca dados por login, nome ou e-mail
//
// O AuthExternoService é exposto como interface para permitir a
// substituição ou extensão para outros provedores de autenticação.
//
// # Tokens JWT
//
// O pacote implementa um serviço de tokens JWT responsável por:
//   - Gerar tokens de acesso (access token)
//   - Gerar tokens de refresh
//   - Validar tokens de acesso
//   - Validar tokens de refresh
//
// As claims personalizadas do token são definidas no tipo Claims,
// incluindo:
//   - ID do usuário
//   - Login
//   - Nome
//   - E-mail
//   - Permissão
//
// As claims são armazenadas no contexto da requisição e podem ser
// extraídas através de helpers:
//
//   - ClaimsDoRequest
//   - ClaimsdoContexto
//
// # Usuário Autenticado
//
// O tipo UsuarioAutenticado representa o usuário autenticado a partir
// das claims presentes no token JWT.
//
// Ele fornece acesso tipado aos dados do usuário e métodos auxiliares
// de verificação de permissões:
//
//   - EhTEC()
//   - EhADM()
//   - EhDEV()
//
// Helpers adicionais permitem extrair o usuário autenticado diretamente
// do contexto:
//
//   - UsuarioAutenticadoDoContexto
//   - UsuarioAutenticadoDoRequest
//
// # Autorização
//
// O pacote fornece mecanismos de autorização baseados nas permissões
// contidas nas claims do token JWT.
//
// Middlewares podem utilizar:
//   - Claims
//   - UsuarioAutenticado
//   - Métodos EhTEC / EhADM / EhDEV
//
// para restringir o acesso a rotas específicas.
//
// # Handlers HTTP
//
// O pacote expõe handlers HTTP responsáveis por:
//   - Login
//   - Refresh
//   - Me
//
// Os handlers utilizam os serviços internos de autenticação e respondem
// sempre com JSON padronizado.
//
// # DTOs
//
// O arquivo `dto.go` define as estruturas de transporte de dados, incluindo:
//   - LoginReq
//   - RefreshReq
//   - TokenResp
//   - UsuarioResp
//
// Funções auxiliares realizam a conversão entre modelos de domínio e
// estruturas de resposta.
//
// # Interfaces Principais
//
// O pacote define as interfaces centrais:
//
//   - AuthInternoService:
//     Login, Refresh e Me
//
//   - AuthExternoService:
//     Abstração para provedores externos de autenticação (ex: LDAP)
//
//   - JWTService:
//     Abstração para geração e validação de tokens
//
// # Fluxo de Autenticação e Autorização
//
//  1. Handler recebe a requisição HTTP
//  2. Payload é decodificado e validado
//  3. AuthInterno coordena autenticação interna ou externa
//  4. Tokens JWT são gerados e retornados
//  5. Claims são anexadas ao contexto da requisição
//  6. Middlewares extraem claims ou UsuarioAutenticado
//  7. Permissões são verificadas
//
// O pacote auth fornece uma implementação robusta, extensível e
// centralizada de autenticação e autorização para toda a aplicação.
package auth
