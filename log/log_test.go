package log

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// startTestNATS starts an embedded NATS server and connects log-nats to it.
// Returns the connection and a cleanup function.
func startTestNATS(t *testing.T) (*nats.Conn, func()) {
	t.Helper()
	opts := &natsserver.Options{
		Host: "127.0.0.1",
		Port: -1, // random port
	}
	ns, err := natsserver.NewServer(opts)
	if err != nil {
		t.Fatalf("failed to create nats server: %v", err)
	}
	ns.Start()
	if !ns.ReadyForConnections(5 * time.Second) {
		t.Fatal("nats server not ready")
	}

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		ns.Shutdown()
		t.Fatalf("failed to connect to nats: %v", err)
	}
	SetDefaultConn(nc)

	return nc, func() {
		connMux.Lock()
		gNats = nil
		connMux.Unlock()
		nc.Close()
		ns.Shutdown()
	}
}

// subscribe subscribes to a NATS subject and sends decoded entries to the
// returned channel.
func subscribe(t *testing.T, nc *nats.Conn, subject string) (<-chan Entry, *nats.Subscription) {
	t.Helper()
	ch := make(chan Entry, 100)
	sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		var e Entry
		if err := json.Unmarshal(msg.Data, &e); err == nil {
			ch <- e
		}
	})
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}
	nc.Flush()
	return ch, sub
}

func receiveEntry(ch <-chan Entry, timeout time.Duration) (Entry, bool) {
	select {
	case e := <-ch:
		return e, true
	case <-time.After(timeout):
		return Entry{}, false
	}
}

func TestTrace(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.TRACE")
	defer sub.Unsubscribe()

	Trace("trace message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "TRACE" {
		t.Errorf("level = %q, want TRACE", e.Level)
	}
	if e.Output != "trace message" {
		t.Errorf("output = %q, want %q", e.Output, "trace message")
	}
	if e.Namespace != DefaultNamespace {
		t.Errorf("namespace = %q, want %q", e.Namespace, DefaultNamespace)
	}
}

func TestTracef(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.TRACE")
	defer sub.Unsubscribe()

	Tracef("trace %s %d", "msg", 1)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "trace msg 1" {
		t.Errorf("output = %q, want %q", e.Output, "trace msg 1")
	}
}

func TestDebug(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.DEBUG")
	defer sub.Unsubscribe()

	Debug("debug message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "DEBUG" {
		t.Errorf("level = %q, want DEBUG", e.Level)
	}
	if e.Output != "debug message" {
		t.Errorf("output = %q, want %q", e.Output, "debug message")
	}
}

func TestDebugf(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.DEBUG")
	defer sub.Unsubscribe()

	Debugf("hello %s %d", "world", 42)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "hello world 42" {
		t.Errorf("output = %q, want %q", e.Output, "hello world 42")
	}
}

func TestInfo(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.INFO")
	defer sub.Unsubscribe()

	Info("info message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
	if e.Output != "info message" {
		t.Errorf("output = %q, want %q", e.Output, "info message")
	}
}

func TestInfof(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.INFO")
	defer sub.Unsubscribe()

	Infof("count: %d", 99)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "count: 99" {
		t.Errorf("output = %q, want %q", e.Output, "count: 99")
	}
}

func TestNotice(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.NOTICE")
	defer sub.Unsubscribe()

	Notice("notice message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "NOTICE" {
		t.Errorf("level = %q, want NOTICE", e.Level)
	}
}

func TestWarn(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.WARN")
	defer sub.Unsubscribe()

	Warn("warning message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "WARN" {
		t.Errorf("level = %q, want WARN", e.Level)
	}
	if e.Output != "warning message" {
		t.Errorf("output = %q, want %q", e.Output, "warning message")
	}
}

func TestWarnf(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.WARN")
	defer sub.Unsubscribe()

	Warnf("warn %d", 1)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "warn 1" {
		t.Errorf("output = %q, want %q", e.Output, "warn 1")
	}
}

func TestError(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.ERROR")
	defer sub.Unsubscribe()

	Error("error message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "ERROR" {
		t.Errorf("level = %q, want ERROR", e.Level)
	}
	if e.Output != "error message" {
		t.Errorf("output = %q, want %q", e.Output, "error message")
	}
}

func TestErrorf(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.ERROR")
	defer sub.Unsubscribe()

	Errorf("err: %s", "something broke")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "err: something broke" {
		t.Errorf("output = %q, want %q", e.Output, "err: something broke")
	}
}

func TestPrint(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.INFO")
	defer sub.Unsubscribe()

	Print("print message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
}

func TestPrintf(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.default.INFO")
	defer sub.Unsubscribe()

	Printf("formatted %s", "print")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "formatted print" {
		t.Errorf("output = %q, want %q", e.Output, "formatted print")
	}
}

func TestPanic(t *testing.T) {
	_, cleanup := startTestNATS(t)
	defer cleanup()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
	}()

	Panic("panic message")
}

func TestPanicf(t *testing.T) {
	_, cleanup := startTestNATS(t)
	defer cleanup()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
	}()

	Panicf("panic %d", 42)
}

func TestPanicln(t *testing.T) {
	_, cleanup := startTestNATS(t)
	defer cleanup()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
	}()

	Panicln("panic line")
}

func TestNamespaceTracking(t *testing.T) {
	_, cleanup := startTestNATS(t)
	defer cleanup()

	l := NewLogger("test-ns-tracking")
	l.Info("register namespace")

	nss := GetNamespaces()
	found := false
	for _, ns := range nss {
		if ns == "test-ns-tracking" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'test-ns-tracking' in GetNamespaces(), got %v", nss)
	}
}

func TestNamespacedLogger(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.api.INFO")
	defer sub.Unsubscribe()

	apiLogger := NewLogger("api")
	apiLogger.Info("api message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Namespace != "api" {
		t.Errorf("namespace = %q, want 'api'", e.Namespace)
	}
	if e.Output != "api message" {
		t.Errorf("output = %q, want %q", e.Output, "api message")
	}
}

func TestNamespaceIsolation(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.api.INFO")
	defer sub.Unsubscribe()

	dbLogger := NewLogger("database")
	apiLogger := NewLogger("api")

	dbLogger.Info("db message")
	apiLogger.Info("api message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "api message" {
		t.Errorf("output = %q, want 'api message'", e.Output)
	}
}

func TestWildcardSubscription(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.>")
	defer sub.Unsubscribe()

	apiLogger := NewLogger("api")
	dbLogger := NewLogger("database")

	apiLogger.Info("api msg")
	dbLogger.Warn("db msg")

	e1, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out on first entry")
	}
	e2, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out on second entry")
	}

	outputs := map[string]bool{e1.Output: true, e2.Output: true}
	if !outputs["api msg"] {
		t.Error("missing 'api msg'")
	}
	if !outputs["db msg"] {
		t.Error("missing 'db msg'")
	}
}

func TestBroadcast(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.broadcast-ns.WARN")
	defer sub.Unsubscribe()

	e := Entry{
		Timestamp: time.Now(),
		Output:    "broadcast test",
		File:      "test.go:1",
		Level:     "WARN",
		Namespace: "broadcast-ns",
	}
	Broadcast(e)

	got, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if got.Output != "broadcast test" {
		t.Errorf("output = %q, want %q", got.Output, "broadcast test")
	}
	if got.Level != "WARN" {
		t.Errorf("level = %q, want WARN", got.Level)
	}
}

func TestNewLoggerEmptyNamespace(t *testing.T) {
	l := NewLogger("")
	if l.Namespace != DefaultNamespace {
		t.Errorf("namespace = %q, want %q", l.Namespace, DefaultNamespace)
	}
}

func TestDefaultLogger(t *testing.T) {
	l := Default()
	if l.Namespace != DefaultNamespace {
		t.Errorf("namespace = %q, want %q", l.Namespace, DefaultNamespace)
	}
	if l.FileInfoDepth != 0 {
		t.Errorf("FileInfoDepth = %d, want 0", l.FileInfoDepth)
	}
}

func TestFileInfo(t *testing.T) {
	fi := fileInfo(1)
	if fi == "" || fi == "<???>:1" {
		t.Errorf("fileInfo returned unexpected value: %q", fi)
	}
}

func TestLogLevelConstants(t *testing.T) {
	if LTrace >= LDebug {
		t.Error("LTrace should be less than LDebug")
	}
	if LDebug >= LInfo {
		t.Error("LDebug should be less than LInfo")
	}
	if LInfo >= LNotice {
		t.Error("LInfo should be less than LNotice")
	}
	if LNotice >= LWarn {
		t.Error("LNotice should be less than LWarn")
	}
	if LWarn >= LError {
		t.Error("LWarn should be less than LError")
	}
	if LError >= LPanic {
		t.Error("LError should be less than LPanic")
	}
	if LPanic >= LFatal {
		t.Error("LPanic should be less than LFatal")
	}
}

func TestParseLevelString(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"TRACE", LTrace},
		{"DEBUG", LDebug},
		{"INFO", LInfo},
		{"NOTICE", LNotice},
		{"WARN", LWarn},
		{"ERROR", LError},
		{"PANIC", LPanic},
		{"FATAL", LFatal},
		{"UNKNOWN", LInfo},
		{"", LInfo},
	}
	for _, tt := range tests {
		got := parseLevelString(tt.input)
		if got != tt.want {
			t.Errorf("parseLevelString(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestNoNATSDoesNotPanic(t *testing.T) {
	connMux.Lock()
	old := gNats
	gNats = nil
	connMux.Unlock()

	defer func() {
		connMux.Lock()
		gNats = old
		connMux.Unlock()
	}()

	Info("no nats message")
	Warnf("warn %d", 1)
	l := NewLogger("test")
	l.Error("logger error")
}

func TestConcurrentLogging(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.>")
	defer sub.Unsubscribe()

	var wg sync.WaitGroup
	const goroutines = 10
	const msgsPerRoutine = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			l := NewLogger("concurrent")
			for j := 0; j < msgsPerRoutine; j++ {
				l.Infof("goroutine %d msg %d", id, j)
			}
		}(i)
	}

	wg.Wait()
	nc.Flush()

	received := 0
	deadline := time.After(3 * time.Second)
	for {
		select {
		case <-ch:
			received++
			if received >= goroutines*msgsPerRoutine {
				return
			}
		case <-deadline:
			t.Logf("received %d/%d messages before deadline", received, goroutines*msgsPerRoutine)
			if received < goroutines*msgsPerRoutine/2 {
				t.Errorf("too few messages received: %d", received)
			}
			return
		}
	}
}
