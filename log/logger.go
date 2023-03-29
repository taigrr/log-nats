package log

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	gNats       *nats.Conn
	connNats    sync.Once
	natsConnURL string
)

func init() {
	natsConnURL = nats.DefaultURL
}

func ConnectDefault(url string, opts ...nats.Option) error {
	connNats.Do(func() {})
	nc, err := nats.Connect(url, opts...)
	gNats = nc
	return err
}

func SetDefaultConn(nc *nats.Conn) {
	connNats.Do(func() {})
	gNats = nc
}

func Default() *Logger {
	connNats.Do(func() {
		var err error
		gNats, err = nats.Connect(natsConnURL)
		if err != nil {
			gNats = nil
		}
	})
	logger := Logger{FileInfoDepth: 0, LogLevel: LWarn, initialized: true, nc: gNats, SubTemplate: "global.logging.{{.Level}}"}
	if gNats == nil {
		logger.initialized = false
	}
	return &logger
}

func (l *Logger) SetNC(nc *nats.Conn) {
	l.nc = nc
}

func (l *Logger) SetInfoDepth(depth int) {
	l.FileInfoDepth = depth
}

// Trace prints out logs on trace level
func (l Logger) Trace(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
	}
	l.createLog(e)
}

// Formatted print for Trace
func (l Logger) Tracef(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
	}
	l.createLog(e)
}

// Trace prints out logs on trace level with newline
func (l Logger) Traceln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
	}
	l.createLog(e)
}

// Debug prints out logs on debug level
func (l Logger) Debug(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "DEBUG",
		level:     LDebug,
	}
	l.createLog(e)
}

// Formatted print for Debug
func (l Logger) Debugf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "DEBUG",
		level:     LDebug,
	}
	l.createLog(e)
}

// Info prints out logs on info level
func (l Logger) Info(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
	}
	l.createLog(e)
}

// Formatted print for Info
func (l Logger) Infof(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
	}
	l.createLog(e)
}

// Info prints out logs on info level with newline
func (l Logger) Infoln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
	}
	l.createLog(e)
}

// Notice prints out logs on notice level
func (l Logger) Notice(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
	}
	l.createLog(e)
}

// Formatted print for Notice
func (l Logger) Noticef(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
	}
	l.createLog(e)
}

// Notice prints out logs on notice level with newline
func (l Logger) Noticeln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
	}
	l.createLog(e)
}

// Warn prints out logs on warn level
func (l Logger) Warn(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
	}
	l.createLog(e)
}

// Formatted print for Warn
func (l Logger) Warnf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
	}
	l.createLog(e)
}

// Warn prints out logs on warn level with a newline
func (l Logger) Warnln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
	}
	l.createLog(e)
}

// Error prints out logs on error level
func (l Logger) Error(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
	}
	l.createLog(e)
}

// Formatted print for error
func (l Logger) Errorf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
	}
	l.createLog(e)
}

// Error prints out logs on error level with a new line
func (l Logger) Errorln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
	}
	l.createLog(e)
}

// Panic prints out logs on panic level
func (l Logger) Panic(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
	}
	l.createLog(e)
	if len(args) >= 0 {
		switch args[0].(type) {
		case error:
			panic(args[0])
		default:
			// falls through to default below
		}
	}
	Flush()
	panic(errors.New(output))
}

// Formatted print for panic
func (l Logger) Panicf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
	}
	l.createLog(e)
	if len(args) >= 0 {
		switch args[0].(type) {
		case error:
			panic(args[0])
		default:
			// falls through to default below
		}
	}
	Flush()
	panic(errors.New(output))
}

// Panic prints out logs on panic level with a newline
func (l Logger) Panicln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
	}
	l.createLog(e)
	if len(args) >= 0 {
		switch args[0].(type) {
		case error:
			panic(args[0])
		default:
			// falls through to default below
		}
	}
	Flush()
	panic(errors.New(output))
}

// Fatal prints out logs on fatal level
func (l Logger) Fatal(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
	}
	l.createLog(e)
	Flush()
	os.Exit(1)
}

// Formatted print for fatal
func (l Logger) Fatalf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
	}
	l.createLog(e)
	Flush()
	os.Exit(1)
}

// Fatal prints fatal level with a new line
func (l Logger) Fatalln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
	}
	l.createLog(e)
	Flush()
	os.Exit(1)
}

// Handles print to info
func (l Logger) Print(args ...interface{}) {
	l.Info(args...)
}

// Handles formatted print to info
func (l Logger) Printf(format string, args ...interface{}) {
	l.Infof(format, args...)
}

// Handles print to info with new line
func (l Logger) Println(args ...interface{}) {
	l.Infoln(args...)
}
