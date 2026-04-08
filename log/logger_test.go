package log

import (
	"testing"
	"time"
)

func TestLoggerTrace(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-trace.TRACE")
	defer sub.Unsubscribe()

	l := NewLogger("logger-trace")
	l.Trace("trace message")

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
	if e.Namespace != "logger-trace" {
		t.Errorf("namespace = %q, want %q", e.Namespace, "logger-trace")
	}
}

func TestLoggerTracef(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-tracef.TRACE")
	defer sub.Unsubscribe()

	l := NewLogger("logger-tracef")
	l.Tracef("trace %s %d", "msg", 1)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "trace msg 1" {
		t.Errorf("output = %q, want %q", e.Output, "trace msg 1")
	}
}

func TestLoggerDebug(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-debug.DEBUG")
	defer sub.Unsubscribe()

	l := NewLogger("logger-debug")
	l.Debug("debug message")

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

func TestLoggerDebugf(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-debugf.DEBUG")
	defer sub.Unsubscribe()

	l := NewLogger("logger-debugf")
	l.Debugf("debug %d", 42)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "debug 42" {
		t.Errorf("output = %q, want %q", e.Output, "debug 42")
	}
}

func TestLoggerDebugln(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-debugln.DEBUG")
	defer sub.Unsubscribe()

	l := NewLogger("logger-debugln")
	l.Debugln("debugln message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "debugln message\n" {
		t.Errorf("output = %q, want %q", e.Output, "debugln message\n")
	}
}

func TestLoggerInfo(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-info.INFO")
	defer sub.Unsubscribe()

	l := NewLogger("logger-info")
	l.Info("info message")

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

func TestLoggerInfof(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-infof.INFO")
	defer sub.Unsubscribe()

	l := NewLogger("logger-infof")
	l.Infof("count: %d", 99)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "count: 99" {
		t.Errorf("output = %q, want %q", e.Output, "count: 99")
	}
}

func TestLoggerNotice(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-notice.NOTICE")
	defer sub.Unsubscribe()

	l := NewLogger("logger-notice")
	l.Notice("notice message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "NOTICE" {
		t.Errorf("level = %q, want NOTICE", e.Level)
	}
}

func TestLoggerWarn(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-warn.WARN")
	defer sub.Unsubscribe()

	l := NewLogger("logger-warn")
	l.Warn("warn message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "WARN" {
		t.Errorf("level = %q, want WARN", e.Level)
	}
	if e.Output != "warn message" {
		t.Errorf("output = %q, want %q", e.Output, "warn message")
	}
}

func TestLoggerError(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-error.ERROR")
	defer sub.Unsubscribe()

	l := NewLogger("logger-error")
	l.Error("error message")

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

func TestLoggerPanic(t *testing.T) {
	_, cleanup := startTestNATS(t)
	defer cleanup()

	l := NewLogger("logger-panic")

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
	}()

	l.Panic("panic message")
}

func TestLoggerPrint(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-print.INFO")
	defer sub.Unsubscribe()

	l := NewLogger("logger-print")
	l.Print("print message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
	if e.Output != "print message" {
		t.Errorf("output = %q, want %q", e.Output, "print message")
	}
}

func TestLoggerPrintf(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-printf.INFO")
	defer sub.Unsubscribe()

	l := NewLogger("logger-printf")
	l.Printf("formatted %s", "print")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "formatted print" {
		t.Errorf("output = %q, want %q", e.Output, "formatted print")
	}
}

func TestLoggerSetInfoDepth(t *testing.T) {
	l := NewLogger("depth-test")
	l.SetInfoDepth(3)
	if l.FileInfoDepth != 3 {
		t.Errorf("FileInfoDepth = %d, want 3", l.FileInfoDepth)
	}
}

func TestLoggerTraceln(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-traceln.TRACE")
	defer sub.Unsubscribe()

	l := NewLogger("logger-traceln")
	l.Traceln("traceln message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "traceln message\n" {
		t.Errorf("output = %q, want %q", e.Output, "traceln message\n")
	}
}

func TestLoggerInfoln(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-infoln.INFO")
	defer sub.Unsubscribe()

	l := NewLogger("logger-infoln")
	l.Infoln("infoln message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "infoln message\n" {
		t.Errorf("output = %q, want %q", e.Output, "infoln message\n")
	}
}

func TestLoggerNoticef(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-noticef.NOTICE")
	defer sub.Unsubscribe()

	l := NewLogger("logger-noticef")
	l.Noticef("notice %d", 42)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "notice 42" {
		t.Errorf("output = %q, want %q", e.Output, "notice 42")
	}
}

func TestLoggerNoticeln(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-noticeln.NOTICE")
	defer sub.Unsubscribe()

	l := NewLogger("logger-noticeln")
	l.Noticeln("noticeln msg")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "noticeln msg\n" {
		t.Errorf("output = %q, want %q", e.Output, "noticeln msg\n")
	}
}

func TestLoggerWarnf(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-warnf.WARN")
	defer sub.Unsubscribe()

	l := NewLogger("logger-warnf")
	l.Warnf("warn %d", 1)

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "warn 1" {
		t.Errorf("output = %q, want %q", e.Output, "warn 1")
	}
}

func TestLoggerWarnln(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-warnln.WARN")
	defer sub.Unsubscribe()

	l := NewLogger("logger-warnln")
	l.Warnln("warnln message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "warnln message\n" {
		t.Errorf("output = %q, want %q", e.Output, "warnln message\n")
	}
}

func TestLoggerErrorf(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-errorf.ERROR")
	defer sub.Unsubscribe()

	l := NewLogger("logger-errorf")
	l.Errorf("err: %s", "broke")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "err: broke" {
		t.Errorf("output = %q, want %q", e.Output, "err: broke")
	}
}

func TestLoggerErrorln(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-errorln.ERROR")
	defer sub.Unsubscribe()

	l := NewLogger("logger-errorln")
	l.Errorln("errorln message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "errorln message\n" {
		t.Errorf("output = %q, want %q", e.Output, "errorln message\n")
	}
}

func TestLoggerPanicf(t *testing.T) {
	_, cleanup := startTestNATS(t)
	defer cleanup()

	l := NewLogger("logger-panicf")
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
	}()

	l.Panicf("panic %d", 42)
}

func TestLoggerPanicln(t *testing.T) {
	_, cleanup := startTestNATS(t)
	defer cleanup()

	l := NewLogger("logger-panicln")
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
	}()

	l.Panicln("panic line")
}

func TestLoggerPrintln(t *testing.T) {
	nc, cleanup := startTestNATS(t)
	defer cleanup()

	ch, sub := subscribe(t, nc, "logging.logger-println.INFO")
	defer sub.Unsubscribe()

	l := NewLogger("logger-println")
	l.Println("println message")

	e, ok := receiveEntry(ch, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "println message\n" {
		t.Errorf("output = %q, want %q", e.Output, "println message\n")
	}
}
