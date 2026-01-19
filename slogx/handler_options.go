package slogx

import "log/slog"

// DefaultHandlerOptions returns slog.HandlerOptions
// with built-in support for slogx custom levels.
//
// This function is a core invariant of the package
// and should always be used when constructing handlers.
func DefaultHandlerOptions(level slog.Leveler) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: replaceLevelAttr,
	}
}

// replaceLevelAttr ensures that custom slogx levels
// are rendered as human-readable strings.
func replaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key != slog.LevelKey {
		return a
	}

	lvl, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a
	}

	switch lvl {
	case LevelTrace:
		return slog.String(slog.LevelKey, "TRACE")
	case LevelFatal:
		return slog.String(slog.LevelKey, "FATAL")
	default:
		return a
	}
}
