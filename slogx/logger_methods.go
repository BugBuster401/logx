package slogx

import (
	"context"
	"log/slog"

	"github.com/BugBuster401/logx"
)

func (l *Logger) Debug(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, toAttrs(fields)...)
}

func (l *Logger) Info(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), slog.LevelInfo, msg, toAttrs(fields)...)
}

func (l *Logger) Warn(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg, toAttrs(fields)...)
}

func (l *Logger) Error(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, toAttrs(fields)...)
}

func (l *Logger) Trace(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), LevelTrace, msg, toAttrs(fields)...)
}

func (l *Logger) Fatal(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), LevelFatal, msg, toAttrs(fields)...)
}

func (l *Logger) With(fields ...logx.Field) logx.Logger {
	args := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return &Logger{logger: l.logger.With(args...)}
}

func (l *Logger) WithGroup(name string) logx.Logger {
	return &Logger{logger: l.logger.WithGroup(name)}
}

func toAttrs(fields []logx.Field) []slog.Attr {
	attrs := make([]slog.Attr, len(fields))
	for i, f := range fields {
		attrs[i] = slog.Any(f.Key, f.Value)
	}
	return attrs
}
