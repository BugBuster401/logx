package slogx

import (
	"io"
	"log/slog"
)

// NewJSONHandler creates a slog.JSONHandler
// with slogx custom level support.
//
// The returned handler is fully configured and safe to use
// with slogx.New.
func NewJSONHandler(
	w io.Writer,
	level slog.Leveler,
) slog.Handler {
	return slog.NewJSONHandler(
		w,
		DefaultHandlerOptions(level),
	)
}
