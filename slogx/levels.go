package slogx

import "log/slog"

// Custom log levels supported by slogx.
//
// TRACE is more verbose than DEBUG.
// FATAL is more severe than ERROR.
//
// These levels are intentionally placed outside
// the standard slog level range.
const (
	LevelTrace slog.Level = slog.LevelDebug - 4
	LevelFatal slog.Level = slog.LevelError + 4
)

// ParseLevel converts a string representation of a level
// into a slog.Leveler.
//
// This helper is intended for application configuration
// and is optional to use.
func ParseLevel(level string) slog.Leveler {
	switch level {
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "trace":
		return LevelTrace
	case "fatal":
		return LevelFatal
	default:
		return slog.LevelInfo
	}
}
