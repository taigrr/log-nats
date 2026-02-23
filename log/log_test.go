package log

import (
	"testing"
)

func TestFileInfo(t *testing.T) {
	info := fileInfo(1)
	if info == "" {
		t.Error("fileInfo returned empty string")
	}
	// Should contain filename:line format
	if len(info) < 3 {
		t.Errorf("fileInfo too short: %s", info)
	}
}

func TestLogLevelFiltering(t *testing.T) {
	// createLog should not panic with nil connection when level is filtered out
	l := &Logger{
		LogLevel:    LWarn,
		initialized: true,
		nc:          nil,
		SubTemplate: "test.{{.Level}}",
	}

	// This should be filtered out (Trace < Warn) and return without accessing nc
	e := Entry{
		Msg:   "test",
		Level: "TRACE",
		level: LTrace,
	}
	// Should not panic
	l.createLog(e)
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

func TestNilConnectionSafety(t *testing.T) {
	l := &Logger{
		LogLevel:    LTrace,
		initialized: true,
		nc:          nil,
		SubTemplate: "test.{{.Level}}",
	}

	// Should not panic even when level passes filter but nc is nil
	e := Entry{
		Msg:   "test",
		Level: "ERROR",
		level: LError,
	}
	l.createLog(e)
}
