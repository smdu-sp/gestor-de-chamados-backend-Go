package main

import (
	"log/slog"
	"os"
	"strings"
	"time"

	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
)

// configurarLogger configura o logger estruturado em formato JSON para ELK/Kibana.
func configurarLogger(config *cfg.Config) *slog.Logger {

	opcoes := &slog.HandlerOptions{
		Level:       mapearNivel(config.AppLogLevel()),
		ReplaceAttr: formatarParaELK,
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, opcoes)

	return slog.New(jsonHandler)
}

// mapearNivel converte uma string de nível de log para o tipo slog.Level.
func mapearNivel(nivel string) slog.Level {
	switch strings.ToUpper(nivel) {
	case "DEBUG":
		return slog.LevelDebug

	case "WARN":
		return slog.LevelWarn

	case "ERROR":
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}

// formatarAtributosParaELK formata os atributos para compatibilidade com ELK.
func formatarParaELK(grupos []string, atributos slog.Attr) slog.Attr {

	// 1- Mapeia as chaves padrão para o formato ELK
	switch atributos.Key {

	// Renomeia chaves padrão
	case slog.MessageKey:
		atributos.Key = "message"

	// Formata timestamp para RFC3339 UTC
	case slog.TimeKey:
		atributos.Key = "@timestamp"

		// força formato timestamp ECS: RFC3339 + UTC
		if t, ok := atributos.Value.Any().(time.Time); ok {
			atributos.Value = slog.StringValue(t.UTC().Format(time.RFC3339Nano))
		}

	// Formata nível de log para string
	case slog.LevelKey:
		atributos.Key = "log.level"

		// evitar panic: slog pode entregar string OU Level
		switch v := atributos.Value.Any().(type) {
		case slog.Level:
			atributos.Value = slog.StringValue(v.String())

		case string:
			atributos.Value = slog.StringValue(v)

		default:
			atributos.Value = slog.StringValue("info")
		}
	}

	return atributos
}
