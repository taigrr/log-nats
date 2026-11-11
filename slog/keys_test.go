package slog

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestWriteAttrQuotesKeys(t *testing.T) {
	tests := []struct {
		name   string
		groups []string
		attr   slog.Attr
		want   string
	}{
		{name: "space", attr: slog.String("user name", "tai"), want: ` "user name"=tai`},
		{name: "equals", attr: slog.String("user=admin", "tai"), want: ` "user=admin"=tai`},
		{name: "newline", attr: slog.String("user\nforged", "tai"), want: ` "user\nforged"=tai`},
		{name: "quote", attr: slog.String(`user"name`, "tai"), want: ` "user\"name"=tai`},
		{name: "group", groups: []string{"request context"}, attr: slog.String("id", "abc"), want: ` "request context.id"=abc`},
		{name: "nested group", attr: slog.Group("http\nrequest", slog.Int("status", 200)), want: ` "http\nrequest.status"=200`},
		{name: "plain", groups: []string{"request", "http"}, attr: slog.Int("status", 200), want: ` request.http.status=200`},
		{name: "inline group", attr: slog.Group("", slog.String("user name", "tai")), want: ` "user name"=tai`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder
			writeAttr(&output, tt.groups, tt.attr)
			if got := output.String(); got != tt.want {
				t.Errorf("writeAttr() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHandlerQuotesGroupedKeys(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()
	entries := subscribeCh(t, nc, "logging.default.INFO")
	logger := slog.New(NewHandler()).WithGroup("request context").With("user\nname", "tai")
	logger.Info("done", slog.Group("http", slog.Int("status", 200)))
	entry, ok := receiveEntry(entries, time.Second)
	if !ok {
		t.Fatal("timed out waiting for log entry")
	}
	want := `done "request context.user\nname"=tai "request context.http.status"=200`
	if entry.Output != want {
		t.Errorf("output = %q, want %q", entry.Output, want)
	}
}
