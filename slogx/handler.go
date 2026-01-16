package slogx

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"
)

// Interface for sending logs to external systems
type RemoteLogClient interface {
	Send(ctx context.Context, entry LogEntry) error
}

// slog.Handler, which sends logs to a remote source
type RemoteHandler struct {
	baseHandler slog.Handler    // basic slog.Handler for local output
	client      RemoteLogClient // client for remotely sending logs
	groups      []string        // current groups
	attrs       []slog.Attr     // accumulated attributes
}

// LogEntry describes a generic log entry,
// which can be sent to any remote storage.
type LogEntry struct {
	Timestamp time.Time      // time of event
	Level     string         // logging level
	Message   string         // message text
	Fields    map[string]any // additional fields
}

// Creates a new handler with any client
func NewRemoteHandler(client RemoteLogClient, opt *slog.HandlerOptions) *RemoteHandler {
	return &RemoteHandler{
		baseHandler: slog.NewTextHandler(os.Stdout, opt),
		client:      client,
		groups:      []string{},
		attrs:       []slog.Attr{},
	}
}

func (h *RemoteHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.baseHandler.Enabled(ctx, lvl)
}

func (h *RemoteHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Prefix keys if there are groups
	prefixedAttrs := make([]slog.Attr, len(attrs))
	copy(prefixedAttrs, attrs)

	if len(h.groups) > 0 {
		prefix := strings.Join(h.groups, ".") + "."
		for i := range prefixedAttrs {
			prefixedAttrs[i].Key = prefix + prefixedAttrs[i].Key
		}
	}

	return &RemoteHandler{
		baseHandler: h.baseHandler.WithAttrs(attrs),
		client:      h.client,
		groups:      append([]string{}, h.groups...),
		attrs:       append(h.attrs, prefixedAttrs...),
	}
}

func (h *RemoteHandler) WithGroup(name string) slog.Handler {
	return &RemoteHandler{
		baseHandler: h.baseHandler.WithGroup(name),
		client:      h.client,
		attrs:       append([]slog.Attr{}, h.attrs...),
		groups:      append(append([]string{}, h.groups...), name),
	}
}

func (h *RemoteHandler) Handle(ctx context.Context, record slog.Record) error {
	// Call the slog handler
	if err := h.baseHandler.Handle(ctx, record); err != nil {
		return err
	}

	// Collecting attributes
	fields := make(map[string]any)
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	record.Attrs(func(a slog.Attr) bool {
		if len(h.groups) > 0 {
			a.Key = strings.Join(h.groups, ".") + "." + a.Key
		}

		fields[a.Key] = a.Value.Any()
		return true
	})

	// Generating a log entry
	entry := LogEntry{
		Timestamp: record.Time,
		Level:     levelToString(record.Level),
		Message:   record.Message,
		Fields:    fields,
	}

	return h.client.Send(ctx, entry)
}

// A custom method for converting levels to strings to support your own levels
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
