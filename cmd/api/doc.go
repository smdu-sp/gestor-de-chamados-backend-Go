// Package main é o ponto de entrada do binário da aplicação Gestor de Chamados.
//
// Este pacote é responsável exclusivamente pelo bootstrap da aplicação,
// atuando como camada de orquestração da infraestrutura e do ciclo de vida
// do processo.
//
// O pacote main NÃO contém regras de negócio, lógica de domínio ou casos de uso.
//
// # Seu papel
//
//   - Interpretar comandos CLI;
//
//   - Inicializar dependências técnicas;
//
//   - Coordenar a execução controlada da aplicação.
//
// # Modelo de execução por comandos
//
// A aplicação segue um modelo explícito de execução via CLI,
// onde cada responsabilidade operacional é acionada por um comando:
//   - migrate → executa migrations do banco de dados
//   - serve   → inicia o servidor HTTP
//   - seed    → popula dados iniciais no banco de dados (apenas para desenvolvimento)
//	 - help    → exibe ajuda sobre os comandos disponíveis
//
// Exemplos:
//
//	go run ./cmd/api migrate
//	go run ./cmd/api serve
//
// Em produção:
//
//	./api migrate
//	./api serve
//
// Esse modelo evita efeitos colaterais implícitos, garantindo que
// alterações estruturais no banco de dados nunca ocorram
// automaticamente durante a inicialização do servidor.
//
// # Responsabilidades do pacote
//
// O pacote main é responsável por:
//   - Carregamento e validação das configurações da aplicação;
//   - Configuração e registro do logger estruturado (slog);
//   - Exibição do banner e informações iniciais da aplicação;
//   - Dispatch de comandos CLI (migrate, serve);
//   - Inicialização explícita da conexão com o banco de dados;
//   - Criação do container de injeção de dependências;
//   - Inicialização e gerenciamento do servidor HTTP;
//   - Gerenciamento de sinais do sistema e shutdown gracioso.
//
// # Organização do fluxo de execução
//
// A função main atua apenas como ponto de entrada do binário,
// delegando responsabilidades para funções especializadas:
//   - executarComando → interpreta e despacha comandos CLI;
//   - executarMigrador → executa apenas as migrations;
//   - rodarServidor → inicia o servidor HTTP;
//   - servidor.iniciar → gerencia execução e graceful shutdown.
//
// Cada comando inicializa apenas as dependências necessárias
// para sua responsabilidade específica.
//
// #	Documentação da API
//
// O pacote importa o módulo Swagger apenas quando o comando `serve`
// é executado, permitindo a exposição automática da documentação
// dos endpoints HTTP.
//
// # Tratamento de erros
//
// Falhas críticas durante a inicialização ou execução de comandos
// resultam em:
//   - Log estruturado com contexto;
//   - Exibição clara do erro no stderr;
//   - Encerramento controlado do processo.
//
// Esse comportamento garante previsibilidade operacional
// e facilidade de diagnóstico em ambientes de produção.
package main
