package logger

import (
	"log/slog"
	"os"
)

type jsonLogger struct {
	logger *slog.Logger
}

func NewJsonLogger() Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &jsonLogger{logger: slog.New(handler)}
}

func (l *jsonLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *jsonLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *jsonLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
