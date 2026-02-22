package logger

import (
	pkglogger "starter-boilerplate/pkg/logger"

	"go.uber.org/zap"
)

type LoggerConfig struct {
	Format          string `yaml:"format"`           // json | console
	Level           string `yaml:"level"`            // debug | info | warn | error
	StacktraceLevel string `yaml:"stacktrace_level"` // debug | info | warn | error | off
}

func SetupLogger(cfg LoggerConfig) *pkglogger.Logger {
	log := pkglogger.SetupLogger(cfg.Format, cfg.Level, cfg.StacktraceLevel)
	zap.ReplaceGlobals(log)
	return log
}
