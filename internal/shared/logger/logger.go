package logger

import (
	"log/slog"

	pkglogger "starter-boilerplate/pkg/logger"
)

type LoggerConfig struct {
	Format string `yaml:"format"` // json | console
	Level  string `yaml:"level"`  // debug | info | warn | error
}

func SetupLogger(cfg LoggerConfig) *pkglogger.Logger {
	log := pkglogger.SetupLogger(cfg.Format, cfg.Level)
	slog.SetDefault(log)
	return log
}
