package logger

import "log/slog"

type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func Err(err error) slog.Attr {
	return slog.String("err", err.Error())
}
