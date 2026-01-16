// Package container centraliza a criação, configuração e gerenciamento
// de todas as dependências da aplicação.
//
// Este pacote atua como o Composition Root do sistema, sendo responsável
// por instanciar, configurar e conectar explicitamente as camadas de:
//
//   - Infraestrutura (repositórios, banco de dados, autenticação externa)
//   - Domínio / Casos de uso (services)
//   - Apresentação (handlers HTTP)
//
// A implementação segue os princípios de Dependency Injection explícita e SOLID,
// evitando o Service Locator anti-pattern e promovendo:
//
//   - Baixo acoplamento entre camadas
//   - Facilidade de testes e substituição de dependências
//   - Clareza e previsibilidade na inicialização da aplicação
//
// # Responsabilidades
//
// O Container é responsável por:
//
//   - Criar e configurar repositórios de infraestrutura (ex.: MySQL).
//   - Criar e orquestrar serviços de domínio e aplicação.
//   - Criar e configurar serviços de autenticação (JWT, LDAP e autenticação interna).
//   - Criar e expor handlers HTTP com todas as dependências resolvidas.
//   - Gerenciar o ciclo de vida de recursos compartilhados
//     (ex.: conexões com banco de dados).
//
// # Princípios de Design
//
// Container:
//
//   - Atua como uma fachada para acesso controlado às dependências.
//   - Expõe apenas os serviços necessários para a camada de apresentação
//     e roteamento.
//   - Mantém dependências internas encapsuladas, evitando acesso direto
//     às implementações concretas.
//
// # Dependências Expostas
//
//   - Handlers():
//       Retorna a estrutura de handlers HTTP utilizada pelo roteador.
//
//   - JWTService():
//       Retorna o serviço responsável pela geração e validação de tokens JWT,
//       utilizado principalmente por middlewares de autenticação.
//
//   - UsuarioService():
//       Retorna o serviço de usuários, utilizado por middlewares e fluxos
//       de autenticação.
//
// # Estrutura Interna
//
//   - repositories:
//       Agrupa todos os repositórios de infraestrutura (MySQL).
//
//   - services:
//       Agrupa os serviços de domínio e casos de uso.
//
//   - authServices:
//       Agrupa os serviços de autenticação
//       (JWT, LDAP e autenticação interna).
//
// # Fluxo de Inicialização
//
//   1. Validação das dependências básicas
//      (configuração e conexão com banco de dados).
//   2. Criação dos repositórios (infraestrutura).
//   3. Criação dos serviços de domínio (casos de uso).
//   4. Criação dos serviços de autenticação.
//   5. Criação dos handlers HTTP.
//   6. Exposição controlada das dependências necessárias.
//
// # Gerenciamento de Recursos
//
//   - O método FecharRecursos é responsável por encerrar recursos
//     compartilhados, como a conexão com o banco de dados.
//   - Deve ser chamado explicitamente durante o shutdown da aplicação.
//
// # Exemplo de Uso
//
//	cfg := carregarConfig()
//	db := conectarBanco(cfg)
//
//	container, err := container.NovoContainer(cfg, db)
//	if err != nil {
//	    // tratar erro
//	}
//
//	router := router.Novo(
//	    container.Handlers(),
//	    container.JWTService(),
//	    container.UsuarioService(),
//	)
//
//	defer container.FecharRecursos(logger)
//
// # Observações Importantes
//
//   - Este pacote não deve conter regras de negócio.
//   - Não deve ser utilizado diretamente dentro de serviços ou repositórios.
//   - Toda nova dependência global da aplicação deve ser registrada aqui,
//     garantindo um único ponto de composição do sistema.
package container
