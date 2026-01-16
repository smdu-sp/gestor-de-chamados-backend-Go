// Package domain define tipos e utilitários fundamentais do domínio da aplicação.
//
// Este pacote centraliza estruturas genéricas e reutilizáveis utilizadas
// por todas as entidades e regras de negócio, mantendo conceitos que são
// independentes de infraestrutura, frameworks ou detalhes de persistência.
//
// O pacote domain fornece:
//
//   1. ErrosValidacao:
//       Estrutura para acumular e representar erros de validação de dados
//       do domínio de forma estruturada e serializável.
//
//   2. Paginacao:
//       Estrutura padrão para paginação de listas e consultas.
//
//   3. GeradorID:
//       Interface para geração de identificadores únicos de entidades
//       de forma desacoplada da estratégia de geração.
//
// ====================================================================
// ERROS DE VALIDAÇÃO
// ====================================================================
//
// O arquivo `erros.go` define o tipo ErrosValidacao, utilizado para
// acumular erros de validação de entidades ou comandos de domínio.
//
// Funcionalidades principais:
//
//   - Add(campo, mensagem):
//       Adiciona um erro associado a um campo específico.
//
//   - HaErros():
//       Indica se há erros de validação acumulados.
//
//   - Error():
//       Retorna uma representação textual amigável para logs.
//
//   - Erros() / MarshalJSON():
//       Retorna a estrutura de erros serializável em JSON.
//
// Esse modelo permite que validações complexas retornem múltiplos erros
// de forma organizada e previsível.
//
// Exemplo de uso:
//
//	ev := domain.NovoErrosValidacao()
//
//	if usuario.Nome == "" {
//	    ev.Add("nome", "é obrigatório")
//	}
//
//	if ev.HaErros() {
//	    return ev
//	}
//
// ====================================================================
// PAGINAÇÃO
// ====================================================================
//
// O arquivo `paginacao.go` define o tipo Paginacao, uma estrutura genérica
// para lidar com paginação de resultados.
//
// Campos:
//
//   - Pagina:
//       Número da página (inicia em 1).
//
//   - Limite:
//       Quantidade máxima de registros por página.
//
// Métodos principais:
//
//   - NovoPaginacao(pagina, limite int):
//       Cria uma nova instância de paginação.
//
//   - Normalizar():
//       Aplica valores padrão caso página ou limite sejam inválidos.
//
//   - Offset():
//       Retorna o offset para queries SQL
//       (OFFSET = (Pagina - 1) * Limite).
//
//   - SemLimite():
//       Define limites altos para retornar todos os registros.
//
// Exemplo de uso:
//
//	p := domain.NovoPaginacao(0, 200)
//	p.Normalizar() // Pagina = 1, Limite = 10
//
//	offset := p.Offset()
//
// ====================================================================
// GERADOR DE ID
// ====================================================================
//
// O arquivo `id.go` define a interface GeradorID, responsável por abstrair
// a geração de identificadores únicos para entidades do domínio.
//
// A interface permite que diferentes estratégias de geração de IDs
// (UUID, ULID, IDs sequenciais, etc.) sejam utilizadas sem acoplamento
// direto ao domínio.
//
// Interface:
//
//   - NovoID() (string, error):
//       Retorna um novo identificador único.
//
// O domínio depende apenas da interface, enquanto as implementações
// concretas devem residir em camadas de infraestrutura.
//
// ====================================================================
// RESUMO
// ====================================================================
//
// O pacote domain provê utilitários essenciais para:
//
//   - Validar dados de entidades e comandos de domínio de forma estruturada;
//   - Garantir paginação consistente em listas e consultas;
//   - Gerar identificadores únicos de forma desacoplada.
//
// Ele é totalmente independente de bancos de dados, protocolos de rede
// ou frameworks HTTP, podendo ser utilizado em qualquer camada do domínio
// ou dos serviços da aplicação.
package domain
