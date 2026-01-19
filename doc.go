// Package logx defines a minimal, implementation-agnostic
// logging interface.
//
// The package provides:
//   - a common Logger interface
//   - a small set of helper constructors for structured fields
//
// logx is intentionally minimal and does not:
//   - define logging backends
//   - depend on any specific logging library
//   - manage global state
//
// Implementations of Logger are expected to provide
// structured, leveled logging.
//
// Typical usage:
//
//	logger.Info(
//		"user logged in",
//		logx.String("user_id", id),
//		logx.Duration("latency", d),
//	)
package logx
