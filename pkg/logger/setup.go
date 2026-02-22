package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.Logger

func SetupLogger(format, level, stacktraceLevel string) *Logger {
	var cfg zap.Config

	if format == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if lvl, err := zapcore.ParseLevel(level); err == nil {
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	}

	var opts []zap.Option
	switch stacktraceLevel {
	case "off":
		cfg.DisableStacktrace = true
	case "":
		// keep defaults
	default:
		if lvl, err := zapcore.ParseLevel(stacktraceLevel); err == nil {
			opts = append(opts, zap.AddStacktrace(lvl))
		}
	}

	logger, err := cfg.Build(opts...)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	return logger
}

func NewNop() *Logger {
	return zap.NewNop()
}
