package log

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Default returns a Logger using the default namespace.
func Default() *Logger {
	return &Logger{FileInfoDepth: 0, Namespace: DefaultNamespace}
}

// NewLogger returns a Logger using the given namespace. If namespace is empty,
// DefaultNamespace is used.
func NewLogger(namespace string) *Logger {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return &Logger{FileInfoDepth: 0, Namespace: namespace}
}

// SetInfoDepth sets the additional caller skip depth for file info.
func (l *Logger) SetInfoDepth(depth int) {
	l.FileInfoDepth = depth
}

// Trace prints out logs on trace level
func (l Logger) Trace(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Tracef is a formatted print for Trace
func (l Logger) Tracef(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Traceln prints out logs on trace level with newline
func (l Logger) Traceln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Debug prints out logs on debug level
func (l Logger) Debug(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Debugf is a formatted print for Debug
func (l Logger) Debugf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Debugln prints out logs on debug level with a newline
func (l Logger) Debugln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Info prints out logs on info level
func (l Logger) Info(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Infof is a formatted print for Info
func (l Logger) Infof(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Infoln prints out logs on info level with newline
func (l Logger) Infoln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Notice prints out logs on notice level
func (l Logger) Notice(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Noticef is a formatted print for Notice
func (l Logger) Noticef(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Noticeln prints out logs on notice level with newline
func (l Logger) Noticeln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Warn prints out logs on warn level
func (l Logger) Warn(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Warnf is a formatted print for Warn
func (l Logger) Warnf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Warnln prints out logs on warn level with a newline
func (l Logger) Warnln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Error prints out logs on error level
func (l Logger) Error(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Errorf is a formatted print for Error
func (l Logger) Errorf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Errorln prints out logs on error level with a new line
func (l Logger) Errorln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Panic prints out logs on panic level
func (l Logger) Panic(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: l.Namespace,
	}
	createLog(e)
	if len(args) > 0 {
		if err, ok := args[0].(error); ok {
			panic(err)
		}
	}
	Flush()
	panic(errors.New(output))
}

// Panicf is a formatted print for Panic
func (l Logger) Panicf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: l.Namespace,
	}
	createLog(e)
	if len(args) > 0 {
		if err, ok := args[0].(error); ok {
			panic(err)
		}
	}
	Flush()
	panic(errors.New(output))
}

// Panicln prints out logs on panic level with a newline
func (l Logger) Panicln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: l.Namespace,
	}
	createLog(e)
	if len(args) > 0 {
		if err, ok := args[0].(error); ok {
			panic(err)
		}
	}
	Flush()
	panic(errors.New(output))
}

// Fatal prints out logs on fatal level
func (l Logger) Fatal(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: l.Namespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Fatalf is a formatted print for Fatal
func (l Logger) Fatalf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: l.Namespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Fatalln prints fatal level with a new line
func (l Logger) Fatalln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: l.Namespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Print delegates to Info
func (l Logger) Print(args ...any) {
	l.Info(args...)
}

// Printf delegates to Infof
func (l Logger) Printf(format string, args ...any) {
	l.Infof(format, args...)
}

// Println delegates to Infoln
func (l Logger) Println(args ...any) {
	l.Infoln(args...)
}
