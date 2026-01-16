package slogx

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/BugBuster401/logx"
)

const (
	LevelTrace slog.Level = slog.LevelDebug - 4
	LevelFatal slog.Level = slog.LevelError + 4
)

type SlogLogger struct {
	logger  *slog.Logger
	logFile *os.File
}

func New(env, level string, remoteClient RemoteLogClient) (logx.Logger, error) {
	var h slog.Handler
	var l slog.Leveler
	var logFile *os.File

	switch strings.ToLower(level) {
	case "error":
		l = slog.LevelError
	case "warn":
		l = slog.LevelWarn
	case "info":
		l = slog.LevelInfo
	case "debug":
		l = slog.LevelDebug
	case "trace":
		l = LevelTrace
	case "fatal":
		l = LevelFatal
	default:
		l = slog.LevelInfo
	}

	opt := &slog.HandlerOptions{
		Level: l,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				switch {
				case level == LevelTrace:
					return slog.String(slog.LevelKey, "TRACE")
				case level == LevelFatal:
					return slog.String(slog.LevelKey, "FATAL")
				default:
					return a
				}
			}
			return a
		},
	}

	if env == "local" {
		h = slog.NewTextHandler(os.Stdout, opt)
	} else {
		h = NewRemoteHandler(remoteClient, opt)
	}

	logger := slog.New(h)

	slog.SetDefault(logger)

	return &SlogLogger{
		logger:  logger,
		logFile: logFile,
	}, nil
}

func Noop() logx.Logger {
	return &SlogLogger{
		logger: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
	}
}

func (l *SlogLogger) Debug(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, toSlogAttrs(fields)...)
}

func (l *SlogLogger) Info(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), slog.LevelInfo, msg, toSlogAttrs(fields)...)
}

func (l *SlogLogger) Warn(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg, toSlogAttrs(fields)...)
}

func (l *SlogLogger) Error(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, toSlogAttrs(fields)...)
}

func (l *SlogLogger) Fatal(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), LevelFatal, msg, toSlogAttrs(fields)...)
}

func (l *SlogLogger) Trace(msg string, fields ...logx.Field) {
	l.logger.LogAttrs(context.Background(), LevelTrace, msg, toSlogAttrs(fields)...)
}

func toSlogAttrs(fields []logx.Field) []slog.Attr {
	attrs := make([]slog.Attr, len(fields))
	for i, f := range fields {
		attrs[i] = slog.Any(f.Key, f.Value)
	}
	return attrs
}

func (l *SlogLogger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

func (l *SlogLogger) With(fields ...logx.Field) logx.Logger {
	args := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return &SlogLogger{logger: l.logger.With(args...)}
}

func (l *SlogLogger) WithGroup(name string) logx.Logger {
	return &SlogLogger{logger: l.logger.WithGroup(name)}
}
