// Package config centraliza todo o carregamento, validação e acesso às
// configurações da aplicação Gestor de Chamados.
//
// Este pacote é responsável por processar variáveis de ambiente (diretas
// ou provenientes de um arquivo .env), aplicar valores padrão quando
// necessário, validar campos obrigatórios e fornecer acesso fortemente
// tipado às configurações da aplicação.
//
// Todas as configurações são carregadas uma única vez no bootstrap da
// aplicação e expostas apenas via métodos de leitura (getters),
// garantindo imutabilidade externa.
//
// # Estrutura principal
//
// O tipo Config agrega todas as subconfigurações do sistema:
//
//   - AppConfig:
//     Configurações gerais da aplicação (nome, versão, ambiente,
//     porta HTTP, nível de log e CORS).
//
//   - ServerConfig:
//     Parâmetros do servidor HTTP, incluindo timeouts,
//     desligamento gracioso e buffers internos de canais.
//
//   - DatabaseConfig:
//     Parâmetros de conexão, pooling e ciclo de vida do banco
//     de dados MySQL.
//
//   - AuthConfig:
//     Segredos e durações de tokens JWT de acesso e refresh.
//
//   - LDAPConfig:
//     Parâmetros opcionais de autenticação externa via LDAP
//     (OpenLDAP ou Active Directory).
//
// # Fluxo de carregamento e validação
//
// O método Carregar() é o ponto de entrada do pacote e segue o fluxo:
//
//   1. Carrega automaticamente um arquivo .env caso exista (via godotenv);
//   2. Inicializa cada subconfiguração a partir das variáveis de ambiente
//      ou de valores padrão definidos pelo sistema;
//   3. Converte valores tipados (int, time.Duration, etc);
//   4. Executa a validação de cada conjunto de configurações;
//   5. Agrega erros de validação no tipo especializado ConfigErros;
//   6. Retorna uma instância de *Config pronta para uso ou um erro estruturado.
//
// # Subconfigurações detalhadas
//
// ## App Config
//
// Definida em app.go, a estrutura AppConfig é responsável por:
//
//   - Nome e versão da aplicação;
//   - Ambiente de execução (local, development, production);
//   - Nível de log (DEBUG, INFO, WARN, ERROR);
//   - Porta do servidor HTTP;
//   - Configurações de CORS.
//
// As validações garantem que:
//
//   - o ambiente informado seja válido;
//   - a porta esteja no intervalo permitido (1 a 65535);
//   - CORS_ORIGIN seja obrigatório em ambiente de produção;
//   - o nível de log configurado seja reconhecido.
//
// ## Server Config
//
// Definida em server.go, a estrutura ServerConfig contém:
//
//   - ReadTimeout, WriteTimeout e IdleTimeout do servidor HTTP;
//   - Timeout para desligamento gracioso (graceful shutdown);
//   - Tamanho dos buffers de canais internos (sinais e erros).
//
// As validações garantem que:
//
//   - todos os timeouts sejam positivos;
//   - os buffers de canais sejam maiores ou iguais a 1;
//   - a configuração seja compatível com um shutdown seguro.
//
// ## Auth Config
//
// Definida em auth.go, a estrutura AuthConfig contém:
//
//   - Segredo de assinatura dos tokens JWT de acesso (TOKEN_SECRET);
//   - Segredo dos refresh tokens (REFRESH_TOKEN_SECRET);
//   - Durações de validade dos tokens (TOKEN_TTL e REFRESH_TOKEN_TTL).
//
// As validações garantem que:
//
//   - os segredos existam e possuam tamanho mínimo de segurança;
//   - as durações sejam válidas e positivas;
//   - limites máximos recomendados sejam respeitados
//     (até 48h para access token e até 30 dias para refresh token).
//
// ## Database Config
//
// Definida em database.go, a estrutura DatabaseConfig inclui:
//
//   - Host, porta, usuário, senha e nome do banco de dados;
//   - Limites de conexões (máximo, mínimo e idle);
//   - Ciclo de vida das conexões (lifetime e idle);
//   - Número máximo de tentativas de retry;
//   - Timeout máximo para migrações.
//
// As validações garantem consistência entre os valores configurados
// e evitam configurações conflitantes ou inseguras.
//
// ## LDAP Config
//
// Definida em ldap.go, a estrutura LDAPConfig é utilizada quando
// LDAP_BASE está definido. Nesse cenário, todos os campos necessários
// para autenticação LDAP tornam-se obrigatórios:
//
//   - Servidor LDAP;
//   - Base DN;
//   - Domínio;
//   - Usuário e senha de bind;
//   - Atributo de login.
//
// # Tratamento de erros
//
// O arquivo erros.go define o tipo ErrosConfig, responsável por:
//
//   - Agregar múltiplos erros de configuração;
//   - Fornecer mensagens amigáveis via Error();
//   - Permitir serialização estruturada (ex.: JSON),
//     facilitando o uso em logs e observabilidade.
//
// # Utilitários internos
//
// O arquivo config.go define utilitários internos como:
//
//   - getEnv(): obtém variáveis de ambiente com fallback;
//   - getEnvInt(): trata valores inteiros inválidos ou ausentes;
//   - getEnvDuration(): converte valores para time.Duration,
//     aplicando fallback seguro.
//
// # Acesso às configurações
//
// O tipo Config expõe apenas métodos getters,
// garantindo isolamento das estruturas internas e evitando mutações externas.
//
// Exemplo de uso:
//
//   cfg, err := config.Carregar()
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(cfg.AppPort())
//
// # Regras específicas para ambiente de produção
//
// Em ambiente "production", regras adicionais de segurança são aplicadas:
//
//   - CORS_ORIGIN não pode ser '*', localhost ou vazio;
//   - DB_HOST não pode apontar para localhost;
//   - Segredos de autenticação não podem ser valores default.
//
// Essas restrições ajudam a prevenir deploys inseguros ou mal configurados.
//
// Em resumo, o pacote config fornece uma base robusta, validada e segura
// para inicializar e operar todas as demais camadas da aplicação.
package config
