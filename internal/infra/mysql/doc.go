// Package mysql fornece implementações de repositórios e utilitários
// para interação com bancos de dados MySQL.
//
// Este pacote é responsável por criar conexões, configurar pools de conexão,
// gerenciar retries e fornecer implementações concretas das interfaces
// de repositório definidas nos domínios da aplicação.
//
// # Visão Geral
//
// 1. Conexão com MySQL:
//
//   - AbrirConexao: abre uma conexão com o banco de dados usando as
//     configurações fornecidas e testa a conexão com retry exponencial.
//   - FecharConexao: fecha a conexão, útil em testes ou scripts.
//   - configurarPool: ajusta parâmetros de pool, como máximo de conexões abertas,
//     ociosas e tempo de vida das conexões.
//   - montarDSN: monta a string de conexão (DSN) com charset, collation, timeouts e TLS.
//   - testarConexao: testa a conexão usando ping com retries e backoff exponencial.
//   - executarSeed: executa scripts SQL para inicialização do banco (apenas em ambiente de desenvolvimento).
//
// 2. Funções auxiliares para manipulação de linhas:
//
//   Arquivo rows.go fornece helpers genéricos:
//     - obterTotalRegistros: obtém o total de registros da última query usando SQL_CALC_FOUND_ROWS.
//     - scanRows: percorre linhas e aplica função de scan genérica para converter resultados em structs.
//
// 3. Implementação de Repositórios:
//
//   Cada entidade (usuário, chamado, categoria etc.) possui um repositório
//   que implementa a interface correspondente do domínio. Exemplo: UsuarioRepository.
//
//   - UsuarioRepository
//       - Métodos: BuscarPorID, BuscarPorLogin, Criar, Atualizar, Listar, ExistePorLogin
//       - Erros específicos: ErrUsuarioNaoEncontrado, ErrUsuarioJaExisteComEmail,
//         ErrUsuarioJaExisteComLogin
//       - scanUsuario: converte resultados de SQL para structs de domínio
//       - construirQueryListar: constrói consultas SQL com filtros, paginação e ordenação
//
//   Padrões comuns:
//     - Cada repositório recebe um *sql.DB injetado.
//     - Tratamento de erros específicos do MySQL (ex.: código 1062 para duplicidade).
//     - Uso de contexto (context.Context) em todas operações de banco.
//     - Uso de helpers genéricos de scan e total de registros.
//
// # Exemplo de Uso
//
//	import (
//	    "context"
//	    "database/sql"
//	    "fmt"
//	    "log/slog"
//
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
//	    "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
//	)
//
//	func main() {
//	    cfg := config.Carregar()
//	    db, err := mysql.AbrirConexao(cfg)
//	    if err != nil {
//	        slog.Error("erro ao abrir conexão", "err", err)
//	        return
//	    }
//	    defer mysql.FecharConexao(slog.Default(), db)
//
//	    repo := mysql.NovoUsuarioRepository(db)
//	    usuario, err := repo.BuscarPorID(context.Background(), "id-do-usuario")
//	    if err != nil {
//	        slog.Error("erro ao buscar usuário", "err", err)
//	    }
//	    fmt.Println(usuario)
//	}
//
// # Características
//
//   - Todo repositório usa prepared statements implicitamente via ExecContext/QueryContext.
//   - Padrão de erros sentinela facilita tratamento de erros específicos no domínio.
//   - Funções auxiliares (rows.go) promovem DRY e consistência no mapeamento de resultados.
//   - Este pacote é projetado para ser usado pelo container de dependências da aplicação,
//     que injeta *sql.DB em todos os repositórios.
//   - Configurações de conexão são carregadas do pacote config.
package mysql
