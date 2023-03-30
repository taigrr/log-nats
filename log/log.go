package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"
)

var cleanup sync.Once

func Flush() {
	cleanup.Do(func() {
		gNats.Flush()
		gNats.Close()
	})
}

func (l *Logger) SetSubscriptionTemplate(template string) {
	l.SubTemplate = template
}

func (l *Logger) createLog(e Entry) {
	temp := template.New("")
	if l.SubTemplate == "" {
		l.SubTemplate = Default().SubTemplate
	}
	ut, err := temp.Parse(l.SubTemplate)
	if err != nil {
		panic(err)
	}
	ut.Option("missingkey=zero")
	buf := bytes.NewBuffer([]byte{})
	ut.Execute(buf, e)
	b, _ := json.Marshal(e)

	l.nc.Publish(buf.String(), []byte(string(b)))
}

// SetLogLevel set log level of logger
func (l *Logger) SetLogLevel(level Level) {
	if !l.initialized {
		panic(errors.New("cannot set level for uninitialized logger, use Default to initialize instead"))
	}
	l.LogLevel = level
}

// Trace prints out logs on trace level
func Trace(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "TRACE",
		level:     LTrace,
	}
	Default().createLog(e)
}

// Formatted print for Trace
func Tracef(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "TRACE",
		level:     LTrace,
	}
	Default().createLog(e)
}

// Trace prints out logs on trace level with newline
func Traceln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "TRACE",
		level:     LTrace,
	}
	Default().createLog(e)
}

// Debug prints out logs on debug level
func Debug(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
	}
	Default().createLog(e)
}

// Formatted print for Debug
func Debugf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
	}
	Default().createLog(e)
}

// Debug prints out logs on debug level with a newline
func Debugln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
	}
	Default().createLog(e)
}

// Info prints out logs on info level
func Info(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
	}
	Default().createLog(e)
}

// Formatted print for Info
func Infof(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
	}
	Default().createLog(e)
}

// Info prints out logs on info level with a newline
func Infoln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
	}
	Default().createLog(e)
}

// Info prints out logs on info level
func Notice(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
	}
	Default().createLog(e)
}

// Formatted print for Info
func Noticef(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
	}
	Default().createLog(e)
}

// Info prints out logs on info level with a newline
func Noticeln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
	}
	Default().createLog(e)
}

// Warn prints out logs on warn level
func Warn(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
	}
	Default().createLog(e)
}

// Formatted print for Warn
func Warnf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
	}
	Default().createLog(e)
}

// Newline print for Warn
func Warnln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
	}
	Default().createLog(e)
}

// Error prints out logs on error level
func Error(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
	}
	Default().createLog(e)
}

// Formatted print for error
func Errorf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
	}
	Default().createLog(e)
}

// Error prints out logs on error level with a newline
func Errorln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
	}
	Default().createLog(e)
}

// Panic prints out logs on panic level
func Panic(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
	}
	Default().createLog(e)
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
func Panicf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
	}
	Default().createLog(e)
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

func Panicln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
	}
	Default().createLog(e)
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
func Fatal(args ...interface{}) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
	}
	Default().createLog(e)
	Flush()
	os.Exit(1)
}

// Formatted print for fatal
func Fatalf(format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
	}
	Default().createLog(e)
	Flush()
	os.Exit(1)
}

func Fatalln(args ...interface{}) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Msg:       output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
	}
	Default().createLog(e)
	Flush()
	os.Exit(1)
}

func Print(args ...interface{}) {
	Info(args...)
}

func Printf(format string, args ...interface{}) {
	Infof(format, args...)
}

func Println(args ...interface{}) {
	Infoln(args...)
}

// fileInfo for getting which line in which file
func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
