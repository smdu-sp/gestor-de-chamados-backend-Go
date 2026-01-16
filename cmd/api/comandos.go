package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/banner"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	cntr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/container"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

const (
	comandoMigrate = "migrate" // comando para executar migrations
	comandoServer  = "server"  // comando para iniciar o servidor HTTP
	comandoSeed    = "seed"    // comando para popular dados iniciais
	comandoHelp    = "help"    // comando para exibir ajuda
)

// executarComando é o dispatcher principal de comandos CLI.
func executarComando(args []string, config *cfg.Config, logger *slog.Logger) error {
	if len(args) < 2 {
		return fmt.Errorf("nenhum comando fornecido; use 'help' para ver os comandos disponíveis")
	}

	switch comando := args[1]; comando {
	case comandoMigrate:
		return executarMigrador(config, logger)

	case comandoServer:
		return executarServidor(config, logger)

	case comandoSeed:
		return executarSeed(config, logger)

	case comandoHelp:
		exibirAjuda()
		return nil

	default:
		return fmt.Errorf("comando desconhecido: %s, use 'help' para ver os comandos disponíveis", comando)
	}
}

// executarMigrador executa apenas as migrations do banco de dados.
func executarMigrador(config *cfg.Config, logger *slog.Logger) error {
	fmt.Println(banner.CarregarBanner(cfg.GoVersion))
	logger.Info("EXECUTANDO MIGRATIONS")

	// 1. Abrir conexão com o banco
	dbConn, err := mysql.AbrirConexao(config)
	if err != nil {
		return fmt.Errorf("executar migrador: %w", err)
	}
	defer dbConn.Close()

	// 2. Criar container apenas para acessar o migrador
	container, err := cntr.NovoContainer(config, dbConn)
	if err != nil {
		return fmt.Errorf("executar migrador: %w", err)
	}
	defer container.FecharRecursos(logger)

	// 3. Executar migrations
	ctx, cancel := context.WithTimeout(context.Background(), config.DBMigrationTimeout())
	defer cancel()

	if err := container.ExecutarMigracoes(ctx); err != nil {
		return fmt.Errorf("executar migrador: %w", err)
	}

	logger.Info("MIGRATIONS EXECUTADAS COM SUCESSO")
	return nil
}

// executarServidor inicia apenas o servidor HTTP (sem rodar migrations).
func executarServidor(config *cfg.Config, logger *slog.Logger) error {
	fmt.Println(banner.CarregarBanner(cfg.GoVersion))
	logger.Info("INICIANDO SERVIDOR")

	// Validação: verificar se migrations já foram executadas
	if err := verificarMigracoesAplicadas(config); err != nil {
		logger.Warn("atenção: migrations podem não ter sido executadas",
			"sugestão", "execute 'go run ./cmd/api migrate' antes de iniciar o servidor",
			"erro", err,
		)
		// Não retorna erro, apenas avisa
	}

	// Delega para a função rodarServidor (não roda migrations automaticamente)
	return rodarServidor(config, logger)
}

// verificarMigracoesAplicadas verifica se a tabela schema_migrations existe.
// Retorna erro se a tabela não existir, indicando que migrations não foram executadas.
func verificarMigracoesAplicadas(config *cfg.Config) error {
	dbConn, err := mysql.AbrirConexao(config)
	if err != nil {
		return fmt.Errorf("verificar migrações aplicadas: %w", err)
	}
	defer dbConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM schema_migrations LIMIT 1`
	var count int

	if err := dbConn.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return fmt.Errorf("tabela schema_migrations não encontrada: %w", err)
	}

	return nil
}

// executarSeed popula o banco de dados com dados iniciais (apenas em desenvolvimento).
func executarSeed(config *cfg.Config, logger *slog.Logger) error {
	fmt.Println(banner.CarregarBanner(cfg.GoVersion))
	logger.Info("EXECUTANDO SEED")

	dbConn, err := mysql.AbrirConexao(config)
	if err != nil {
		return fmt.Errorf("executar seed: %w", err)
	}
	defer dbConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := mysql.ExecutarSeed(ctx, dbConn, config); err != nil {
		return fmt.Errorf("executar seed: %w", err)
	}

	logger.Info("SEED EXECUTADA COM SUCESSO")
	return nil
}

// exibirAjuda exibe instruções de uso dos comandos CLI.
func exibirAjuda() {
	ajuda := `
USO: go run ./cmd/api [comando]

COMANDOS DISPONÍVEIS:
	migrate    Executa as migrations do banco de dados.
	server     Inicia o servidor HTTP.
	seed       Popula o banco de dados com dados iniciais (apenas em desenvolvimento).
	help       Exibe esta mensagem de ajuda.

EXEMPLOS:
	# Executar migrations
	go run ./cmd/api migrate

	# Iniciar o servidor HTTP
	go run ./cmd/api server

	# Popular dados iniciais (no ambiente de desenvolvimento)
	go run ./cmd/api seed

	# Exibir ajuda
	go run ./cmd/api help

`
	fmt.Fprint(os.Stderr, ajuda)
}
