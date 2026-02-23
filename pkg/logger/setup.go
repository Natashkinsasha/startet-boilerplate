package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

type Logger = slog.Logger

func SetupLogger(format, level string) *Logger {
	lvl := parseLevel(level)
	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func NewNop() *Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
