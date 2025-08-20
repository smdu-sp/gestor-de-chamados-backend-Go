package util

import (
	"log/slog"
	"os"
)

func Logger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo}))
}
