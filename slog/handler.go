// Package slog provides a [log/slog.Handler] that routes structured log
// records into the log-nats system, publishing entries to NATS subjects
// and stderr via the normal broadcast path.
package slog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/taigrr/log-nats/v2/log"
)

// Handler implements [slog.Handler] by converting each [slog.Record] into a
// log-nats [log.Entry] and feeding it through the normal broadcast path.
//
// Attributes accumulated via [Handler.WithAttrs] are prepended to the log
// message as key=value pairs. Groups set via [Handler.WithGroup] prefix
// attribute keys with "group.".
type Handler struct {
	namespace string
	level     slog.Level
	attrs     []slog.Attr
	groups    []string
}

// ensure interface compliance at compile time.
var _ slog.Handler = (*Handler)(nil)

// Option configures a [Handler].
type Option func(*Handler)

// WithNamespace sets the log-nats namespace for entries produced by this
// handler. If empty, [log.DefaultNamespace] is used.
func WithNamespace(ns string) Option {
	return func(h *Handler) {
		h.namespace = ns
	}
}

// WithLevel sets the minimum slog level the handler will accept.
func WithLevel(l slog.Level) Option {
	return func(h *Handler) {
		h.level = l
	}
}

// NewHandler returns a new [Handler] that writes to the log-nats broadcast
// system. Options may be used to set the namespace and minimum level.
func NewHandler(opts ...Option) *Handler {
	h := &Handler{
		namespace: log.DefaultNamespace,
		level:     slog.LevelDebug,
	}
	for _, o := range opts {
		o(h)
	}
	return h
}

// Enabled reports whether the handler is configured to process records at
// the given level.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle converts r into a log-nats Entry and broadcasts it.
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder
	b.WriteString(r.Message)

	// Append pre-collected attrs.
	for _, a := range h.attrs {
		writeAttr(&b, h.groups, a)
	}

	// Append record-level attrs.
	r.Attrs(func(a slog.Attr) bool {
		writeAttr(&b, h.groups, a)
		return true
	})

	file := "???"
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			short := f.File
			if idx := strings.LastIndex(short, "/"); idx >= 0 {
				short = short[idx+1:]
			}
			file = fmt.Sprintf("%s:%d", short, f.Line)
		}
	}

	e := log.Entry{
		Timestamp: r.Time,
		Output:    b.String(),
		File:      file,
		Level:     slogLevelToString(r.Level),
		Namespace: h.namespace,
	}
	log.Broadcast(e)
	return nil
}

// WithAttrs returns a new Handler with the given attributes appended.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := h.clone()
	h2.attrs = append(h2.attrs, attrs...)
	return h2
}

// WithGroup returns a new Handler where subsequent attributes are nested
// under the given group name.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *Handler) clone() *Handler {
	h2 := &Handler{
		namespace: h.namespace,
		level:     h.level,
		attrs:     make([]slog.Attr, len(h.attrs)),
		groups:    make([]string, len(h.groups)),
	}
	copy(h2.attrs, h.attrs)
	copy(h2.groups, h.groups)
	return h2
}

func writeAttr(b *strings.Builder, groups []string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	b.WriteByte(' ')
	for _, g := range groups {
		b.WriteString(g)
		b.WriteByte('.')
	}
	b.WriteString(a.Key)
	b.WriteByte('=')
	b.WriteString(a.Value.String())
}

func slogLevelToString(l slog.Level) string {
	switch {
	case l >= slog.LevelError:
		return "ERROR"
	case l >= slog.LevelWarn:
		return "WARN"
	case l >= slog.LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}
