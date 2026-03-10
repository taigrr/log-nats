package slog

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/taigrr/log-nats/v2/log"
)

func startTestNATS(t *testing.T) (*nats.Conn, func()) {
	t.Helper()
	opts := &natsserver.Options{Host: "127.0.0.1", Port: -1}
	ns, err := natsserver.NewServer(opts)
	if err != nil {
		t.Fatalf("nats server: %v", err)
	}
	ns.Start()
	if !ns.ReadyForConnections(5 * time.Second) {
		t.Fatal("nats not ready")
	}
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		ns.Shutdown()
		t.Fatalf("nats connect: %v", err)
	}
	log.SetDefaultConn(nc)
	return nc, func() {
		log.SetDefaultConn(nil)
		nc.Close()
		ns.Shutdown()
	}
}

func receiveEntry(ch <-chan log.Entry, timeout time.Duration) (log.Entry, bool) {
	select {
	case e := <-ch:
		return e, true
	case <-time.After(timeout):
		return log.Entry{}, false
	}
}

func subscribeCh(t *testing.T, nc *nats.Conn, subject string) <-chan log.Entry {
	t.Helper()
	ch := make(chan log.Entry, 100)
	sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		var e log.Entry
		if err := json.Unmarshal(msg.Data, &e); err == nil {
			ch <- e
		}
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	nc.Flush()
	t.Cleanup(func() { sub.Unsubscribe() })
	return ch
}

func TestHandler_Enabled(t *testing.T) {
	h := NewHandler(WithLevel(slog.LevelWarn))
	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected Info to be disabled when min level is Warn")
	}
	if !h.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("expected Warn to be enabled")
	}
	if !h.Enabled(context.Background(), slog.LevelError) {
		t.Error("expected Error to be enabled")
	}
}

func TestHandler_Handle(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch := subscribeCh(t, nc, "logging.test-ns.INFO")

	h := NewHandler(WithNamespace("test-ns"))
	logger := slog.New(h)
	logger.Info("hello world", "key", "value")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out waiting for log entry")
	}
	if e.Namespace != "test-ns" {
		t.Errorf("namespace = %q, want %q", e.Namespace, "test-ns")
	}
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
	if e.Output == "" {
		t.Error("output should not be empty")
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch := subscribeCh(t, nc, "logging.default.INFO")

	h := NewHandler()
	h2 := h.WithAttrs([]slog.Attr{slog.String("service", "api")})
	logger := slog.New(h2)
	logger.Info("request")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "request service=api" {
		t.Errorf("output = %q, want %q", e.Output, "request service=api")
	}
}

func TestHandler_WithGroup(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch := subscribeCh(t, nc, "logging.default.INFO")

	h := NewHandler()
	h2 := h.WithGroup("http").WithAttrs([]slog.Attr{slog.Int("status", 200)})
	logger := slog.New(h2)
	logger.Info("done")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "done http.status=200" {
		t.Errorf("output = %q, want %q", e.Output, "done http.status=200")
	}
}

func TestSlogLevelMapping(t *testing.T) {
	tests := []struct {
		level slog.Level
		want  string
	}{
		{slog.LevelDebug, "DEBUG"},
		{slog.LevelInfo, "INFO"},
		{slog.LevelWarn, "WARN"},
		{slog.LevelError, "ERROR"},
	}
	for _, tt := range tests {
		got := slogLevelToString(tt.level)
		if got != tt.want {
			t.Errorf("slogLevelToString(%v) = %q, want %q", tt.level, got, tt.want)
		}
	}
}
