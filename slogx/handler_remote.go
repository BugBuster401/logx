package slogx

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"
)

// RemoteLogClient defines an interface for sending logs
// to an external system (e.g. HTTP, Kafka, gRPC).
type RemoteLogClient interface {
	Send(ctx context.Context, entry LogEntry) error
}

// LogEntry represents a structured log entry
// suitable for remote storage.
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]any
}

// RemoteHandler is a slog.Handler that:
//   - writes logs locally
//   - sends structured logs to a remote sink
type RemoteHandler struct {
	baseHandler slog.Handler
	client      RemoteLogClient
	groups      []string
	attrs       []slog.Attr
}

// NewRemoteHandler creates a slog.Handler
// that logs locally and forwards entries to a remote client.
//
// Local output is written to stdout by default.
func NewRemoteHandler(
	client RemoteLogClient,
	level slog.Leveler,
) slog.Handler {
	base := slog.NewTextHandler(
		os.Stdout,
		DefaultHandlerOptions(level),
	)

	return &RemoteHandler{
		baseHandler: base,
		client:      client,
	}
}

func (h *RemoteHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.baseHandler.Enabled(ctx, lvl)
}

func (h *RemoteHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &RemoteHandler{
		baseHandler: h.baseHandler.WithAttrs(attrs),
		client:      h.client,
		groups:      append([]string{}, h.groups...),
		attrs:       append(h.attrs, attrs...),
	}
}

func (h *RemoteHandler) WithGroup(name string) slog.Handler {
	return &RemoteHandler{
		baseHandler: h.baseHandler.WithGroup(name),
		client:      h.client,
		groups:      append(append([]string{}, h.groups...), name),
		attrs:       append([]slog.Attr{}, h.attrs...),
	}
}

func (h *RemoteHandler) Handle(ctx context.Context, r slog.Record) error {
	if err := h.baseHandler.Handle(ctx, r); err != nil {
		return err
	}

	fields := map[string]any{}

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	r.Attrs(func(a slog.Attr) bool {
		if len(h.groups) > 0 {
			a.Key = strings.Join(h.groups, ".") + "." + a.Key
		}
		fields[a.Key] = a.Value.Any()
		return true
	})

	entry := LogEntry{
		Timestamp: r.Time,
		Level:     levelToString(r.Level),
		Message:   r.Message,
		Fields:    fields,
	}

	return h.client.Send(ctx, entry)
}

func levelToString(lvl slog.Level) string {
	switch lvl {
	case LevelTrace:
		return "TRACE"
	case LevelFatal:
		return "FATAL"
	default:
		return lvl.String()
	}
}
