package logx

// Logger defines a structured, leveled logging interface.
//
// Implementations must be safe for concurrent use.
//
// Logger does not prescribe:
//   - output format
//   - log destination
//   - buffering or flushing behavior
//
// These concerns are left to concrete implementations.
type Logger interface {
	// Debug logs a message at DEBUG level.
	Debug(msg string, fields ...Field)

	// Info logs a message at Info level.
	Info(msg string, fields ...Field)

	// Warn logs a message at Warn level.
	Warn(msg string, fields ...Field)

	// Error logs a message at Error level.
	Error(msg string, fields ...Field)

	// Fatal logs a message at FATAL level.
	//
	// Implementations may choose to terminate the process,
	// but this behavior is not required by the interface.
	Fatal(msg string, fields ...Field)

	// Trace logs a message at TRACE level.
	//
	// TRACE is intended for very detailed, high-volume logging.
	Trace(msg string, fields ...Field)

	// With returns a new Logger with additional structured fields
	// attached to every log entry.
	With(fields ...Field) Logger

	// WithGroup returns a new Logger that groups all subsequent
	// fields under the specified name.
	WithGroup(name string) Logger
}
