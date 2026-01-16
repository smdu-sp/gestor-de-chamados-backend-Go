package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	rtr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/router"
	cntr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/container"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// servidor encapsula o servidor HTTP e suas dependências
type servidor struct {
	http   *http.Server
	logger *slog.Logger
	porta  string
	config *cfg.Config
}

// novoServer cria uma nova instância do servidor
func novoServidor(config *cfg.Config, handler http.Handler, logger *slog.Logger) *servidor {
	http := &http.Server{
		Addr:         ":" + config.AppPort(),
		Handler:      handler,
		ReadTimeout:  config.ServerReadTimeout(),
		WriteTimeout: config.ServerWriteTimeout(),
		IdleTimeout:  config.ServerIdleTimeout(),
	}

	return &servidor{
		http:   http,
		logger: logger,
		porta:  config.AppPort(),
		config: config,
	}
}

// rodarServidor inicia o servidor HTTP com as configurações e dependências necessárias.
func rodarServidor(config *cfg.Config, logger *slog.Logger) error {
	// Contexto base para o servidor
	ctx := context.Background()

	// 1. Logar informações iniciais da aplicação
	logger.Info("iniciando aplicação",
		"service", config.AppName(),
		"version", config.AppVersion(),
		"environment", config.AppEnvironment(),
		"port", config.AppPort(),
	)

	// 2. Abrir conexão com o banco de dados
	dbConn, err := mysql.AbrirConexao(config)
	if err != nil {
		return fmt.Errorf("rodar servidor: %w", err)
	}

	// 3. Criar container de dependências
	container, err := cntr.NovoContainer(config, dbConn)
	if err != nil {
		return fmt.Errorf("rodar servidor: %w", err)
	}
	defer container.FecharRecursos(logger)

	// 4. Configurar roteador HTTP
	router := rtr.NovoRouter(
		container.Handlers(),
		container.JWTService(),
		container.UsuarioService(),
		config,
	)

	// 5. Obter handler configurado
	handler := router.Configurar()

	// 6. Iniciar servidor HTTP
	servidor := novoServidor(config, handler, logger)
	return servidor.iniciar(ctx)
}

// iniciar inicia o servidor HTTP e aguarda sinais de interrupção para shutdown gracioso.
func (s *servidor) iniciar(ctx context.Context) error {
	// Canal para capturar erros do servidor
	servErros := make(chan error, s.config.ServerErrorChannelBufferSize())

	// Inicia o servidor em goroutine
	go func() {
		s.logger.Info(
			"iniciando servidor HTTP",
			"address", fmt.Sprintf("http://localhost:%s", s.porta),
		)

		s.logger.Info(
			"CORS configurado para origem",
			"origin", s.config.AppCORSOrigin(),
		)

		s.logger.Info(
			"documentação Swagger disponível em",
			"swagger", fmt.Sprintf("http://localhost:%s/swagger/index.html", s.porta),
		)

		s.logger.Info(
			"SERVIDOR INICIADO COM SUCESSO",
			"address", fmt.Sprintf("http://localhost:%s", s.porta),
		)

		// 'ListenAndServe' bloqueia até o servidor ser desligado ou ocorrer um erro
		if err := s.http.ListenAndServe(); err != nil {
			select {
			case servErros <- err:
			default:
			}
		}
	}()

	// Aguarda sinal de interrupção ou erro do servidor
	quit := make(chan os.Signal, s.config.ServerSignalChannelBufferSize())
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Bloqueia até receber um sinal ou erro
	select {
	case err := <-servErros:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("erro no servidor: %w", err)
		}

	case sig := <-quit:
		s.logger.Info(
			"sinal de interrupção recebido",
			"signal", sig.String(),
			"event", "shutdown_iniciada",
		)
	}

	return s.desligar(ctx)
}

// desligar realiza o graceful shutdown do servidor HTTP.
func (s *servidor) desligar(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.ServerShutdownTimeout())
	defer cancel()

	s.logger.Info(
		"iniciando shutdown do servidor",
		"timeout", s.config.ServerShutdownTimeout().String(),
	)

	// Tenta desligar o servidor graciosamente
	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("erro durante shutdown: %w", err)
	}

	s.logger.Info(
		"servidor desligado com sucesso",
		"event", "shutdown_complete",
		"service", s.config.AppName(),
	)

	return nil
}
