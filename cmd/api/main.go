package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	_ "github.com/smdu-sp/gestor-de-chamados-backend-Go/docs"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
)

// @title Gestor de Chamados API
// @version 1.0
// @description API para gestão de chamados
// @host localhost:8080
// @BasePath /
func main() {
	// 1. Carregar configuração
	config, err := cfg.Carregar()
	if err != nil {
		tratarErroConfiguracao(err)
	}

	// 2. Configurar logger
	logger := configurarLogger(config)
	slog.SetDefault(logger)

	// 3. Executar comando especificado via CLI
	if err := executarComando(os.Args, config, logger); err != nil {
		logger.Error(
			"falha ao executar comando CLI",
			"error", err,
			"stack", string(debug.Stack()),
		)

		fmt.Fprintf(os.Stderr, "\nErro: %v\n\n", err)
		exibirAjuda()
		os.Exit(1)
	}
}

// tratarErroConfiguracao trata erros relacionados ao carregamento e validação da configuração.
func tratarErroConfiguracao(err error) {
	if erros, ok := err.(*cfg.ErrosConfig); ok {
		fmt.Fprintln(os.Stderr, erros.Error())

		l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		l.Error(
			"falha na validação da configuração",
			"errors", erros.Erros(),
		)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Erro ao carregar configuração: %v\n", err)
	os.Exit(1)
}
