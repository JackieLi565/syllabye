package logger

import (
	"log/slog"
	"os"
)

type testLogger struct {
	logger *slog.Logger
}

func NewTextLogger() Logger {
	handler := slog.NewTextHandler(os.Stdout, nil)
	return &testLogger{logger: slog.New(handler)}
}

func (l *testLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *testLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *testLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
