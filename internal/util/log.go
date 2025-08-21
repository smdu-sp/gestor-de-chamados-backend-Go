package util

import (
	"log/slog"
	"os"
)

// Logger cria e retorna uma instância de *slog.Logger
// Configurado para saída JSON no stdout e nível de log INFO
func Logger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout, // saída dos logs no terminal
			&slog.HandlerOptions{
				Level: slog.LevelInfo, // nível mínimo de log que será registrado
			},
		),
	)
}
