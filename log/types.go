package log

import (
	"time"

	"github.com/nats-io/nats.go"
)

type (
	LogWriter chan Entry
	Level     int
)

const (
	LTrace Level = iota
	LDebug
	LInfo
	LNotice
	LWarn
	LError
	LPanic
	LFatal
)

type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Msg       string    `json:"msg"`
	File      string    `json:"file"`
	Level     string    `json:"level"`
	level     Level
}

type Logger struct {
	LogLevel      Level `json:"level"`
	nc            *nats.Conn
	initialized   bool
	FileInfoDepth int
	SubTemplate   string
}
