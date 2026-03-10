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

	"github.com/nats-io/nats.go"
)

var (
	gNats       *nats.Conn
	connMux     sync.Mutex
	subjectTmpl = "logging.{{.Namespace}}.{{.Level}}"
	cleanup     sync.Once
	namespaces  map[string]bool
	nsMux       sync.RWMutex
)

func init() {
	namespaces = make(map[string]bool)
}

// ConnectDefault connects to a NATS server and sets it as the default connection.
func ConnectDefault(url string, opts ...nats.Option) error {
	connMux.Lock()
	defer connMux.Unlock()
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return err
	}
	gNats = nc
	return nil
}

// SetDefaultConn sets a pre-existing NATS connection as the default.
func SetDefaultConn(nc *nats.Conn) {
	connMux.Lock()
	defer connMux.Unlock()
	gNats = nc
}

// SetSubjectTemplate sets the NATS subject template used when publishing log
// entries. The template receives an Entry and defaults to
// "logging.{{.Namespace}}.{{.Level}}".
func SetSubjectTemplate(tmpl string) {
	connMux.Lock()
	defer connMux.Unlock()
	subjectTmpl = tmpl
}

// Flush flushes and closes the NATS connection.
func Flush() {
	cleanup.Do(func() {
		connMux.Lock()
		if gNats != nil {
			gNats.Flush()
			gNats.Close()
		}
		connMux.Unlock()
	})
}

func createLog(e Entry) {
	// Track namespace
	nsMux.Lock()
	namespaces[e.Namespace] = true
	nsMux.Unlock()

	// Publish to NATS
	connMux.Lock()
	nc := gNats
	connMux.Unlock()
	if nc != nil {
		publishToNATS(nc, e)
	}
}

func publishToNATS(nc *nats.Conn, e Entry) {
	tmpl := template.New("")
	ut, err := tmpl.Parse(subjectTmpl)
	if err != nil {
		return
	}
	ut.Option("missingkey=zero")
	buf := bytes.NewBuffer(nil)
	if err := ut.Execute(buf, e); err != nil {
		return
	}
	b, _ := json.Marshal(e)
	nc.Publish(buf.String(), b)
}

// GetNamespaces returns a list of all namespaces that have been used.
func GetNamespaces() []string {
	nsMux.RLock()
	defer nsMux.RUnlock()

	result := make([]string, 0, len(namespaces))
	for ns := range namespaces {
		result = append(result, ns)
	}
	return result
}

// Broadcast sends an Entry through the logging system. This is the public
// entry point used by adapter packages (such as the slog handler) that
// construct entries themselves.
func Broadcast(e Entry) {
	if e.level == 0 && e.Level != "" && e.Level != "TRACE" {
		e.level = parseLevelString(e.Level)
	}
	createLog(e)
}

func parseLevelString(s string) Level {
	switch s {
	case "TRACE":
		return LTrace
	case "DEBUG":
		return LDebug
	case "INFO":
		return LInfo
	case "NOTICE":
		return LNotice
	case "WARN":
		return LWarn
	case "ERROR":
		return LError
	case "PANIC":
		return LPanic
	case "FATAL":
		return LFatal
	default:
		return LInfo
	}
}

// Trace prints out logs on trace level
func Trace(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Tracef is a formatted print for Trace
func Tracef(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Traceln prints out logs on trace level with newline
func Traceln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Debug prints out logs on debug level
func Debug(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Debugf is a formatted print for Debug
func Debugf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Debugln prints out logs on debug level with a newline
func Debugln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Info prints out logs on info level
func Info(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Infof is a formatted print for Info
func Infof(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Infoln prints out logs on info level with a newline
func Infoln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Notice prints out logs on notice level
func Notice(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Noticef is a formatted print for Notice
func Noticef(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Noticeln prints out logs on notice level with a newline
func Noticeln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Warn prints out logs on warn level
func Warn(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Warnf is a formatted print for Warn
func Warnf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Warnln is a newline print for Warn
func Warnln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Error prints out logs on error level
func Error(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Errorf is a formatted print for Error
func Errorf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Errorln prints out logs on error level with a newline
func Errorln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Panic prints out logs on panic level
func Panic(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: DefaultNamespace,
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
func Panicf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: DefaultNamespace,
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
func Panicln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: DefaultNamespace,
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
func Fatal(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Fatalf is a formatted print for Fatal
func Fatalf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Fatalln prints out logs on fatal level with a newline
func Fatalln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Print delegates to Info
func Print(args ...any) {
	Info(args...)
}

// Printf delegates to Infof
func Printf(format string, args ...any) {
	Infof(format, args...)
}

// Println delegates to Infoln
func Println(args ...any) {
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
