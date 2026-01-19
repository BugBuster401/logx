package slogx

import (
	"io"
	"log/slog"
)

// NewTextHandler creates a slog.TextHandler
// with slogx custom level support.
//
// The returned handler is fully configured and safe to use
// with slogx.New.
func NewTextHandler(
	w io.Writer,
	level slog.Leveler,
) slog.Handler {
	return slog.NewTextHandler(
		w,
		DefaultHandlerOptions(level),
	)
}
