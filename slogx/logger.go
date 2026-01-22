package slogx

import (
	"io"
	"log/slog"

	"github.com/BugBuster401/logx"
)

// Logger is a slog-based implementation of logx.Logger.
type Logger struct {
	logger *slog.Logger
}

// New creates a logx.Logger using the provided slog.Handler.
//
// The handler must be fully configured before being passed in.
// slogx provides helper constructors for this purpose.
func New(handler slog.Handler) logx.Logger {
	return &Logger{
		logger: slog.New(handler),
	}
}

// Noop returns a null Logger, which ignores all calls.
// Ideal for test scenarios where logging is not required.
func Noop() logx.Logger {
	return &Logger{
		logger: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
	}
}
