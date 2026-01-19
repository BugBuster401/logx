// Package slogx provides a logx.Logger implementation
// based on the standard library log/slog package.
//
// slogx adapts slog.Logger to the logx.Logger interface,
// providing a thin compatibility layer without adding
// additional behavior or abstractions.
//
// The package is intended to:
//   - integrate log/slog into applications using logx
//   - keep logging infrastructure isolated from business code
//
// slogx does not:
//   - replace slog
//   - introduce its own logging semantics
//   - manage global loggers
//
// Configuration of log levels, handlers, and output
// is delegated to slog.
package slogx
