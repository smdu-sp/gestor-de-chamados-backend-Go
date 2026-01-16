// Package migrations fornece a infraestrutura para versionamento e execução
// controlada de migrações de banco de dados.
//
// Este pacote é responsável exclusivamente por aplicar alterações
// estruturais no schema do banco de dados, garantindo:
//
//   - Execução ordenada e determinística das migrações;
//   - Controle de versões aplicadas;
//   - Segurança transacional durante a aplicação das mudanças.
//
// O pacote migrations NÃO decide quando as migrações devem ser executadas.
// Essa decisão é responsabilidade da camada de orquestração
// (tipicamente o pacote main ou comandos CLI explícitos).
//
// # Funcionamento das Migrações
//
// As migrações são definidas como arquivos SQL versionados,
// embutidos no binário utilizando embed.FS.
//
// Cada arquivo representa uma migração atômica e é executado
// exatamente uma vez, com seu estado registrado na tabela
// de controle `schema_migrations`.
//
// A execução segue as etapas:
//
//   1. Garantir a existência da tabela schema_migrations;
//   2. Identificar migrações já aplicadas;
//   3. Identificar migrações disponíveis no filesystem;
//   4. Executar apenas as migrações pendentes, em ordem lexicográfica;
//   5. Registrar cada migração aplicada com sucesso.
//
// # Abstração do Migrador
//
// O pacote define a interface Migrador, que abstrai o mecanismo
// de execução de migrações:
//
//   type Migrador interface {
//       Migrar(ctx context.Context) error
//   }
//
// Essa abstração permite:
//
//   - Substituição do banco de dados (ex: PostgreSQL);
//   - Implementações alternativas (ex: dry-run, testes);
//   - Isolamento da lógica de migração da aplicação.
//
// # Decisões de Design
//
//   - Migrações são executadas dentro de transações;
//   - Falhas resultam em rollback completo da migração;
//   - Nenhuma migração é aplicada parcialmente;
//   - O pacote não executa AutoMigrate;
//   - Não há efeitos colaterais implícitos na inicialização da aplicação.
//
// Essas decisões seguem práticas profissionais para ambientes
// de produção, evitando alterações não intencionais no banco.
//
// # Funcionalidades Principais
//
//   - Embutir arquivos de migração SQL no binário;
//   - Executar migrações pendentes de forma segura;
//   - Manter o histórico de migrações aplicadas;
//   - Fornecer abstrações claras para execução de migrations.
//
// # Uso do Pacote
//
// Este pacote pertence à camada de infraestrutura da aplicação
// e deve ser utilizado apenas por componentes responsáveis
// pela orquestração do ciclo de vida do sistema.
package migrations
