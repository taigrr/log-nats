# log-nats

`log-nats` is an opinionated logging library that publishes structured log
entries to a [NATS](https://nats.io) message bus. It ships with both a
package-level API (similar to the standard `log` package) and a
[`log/slog`](https://pkg.go.dev/log/slog) handler.

## Installation

```bash
go get github.com/taigrr/log-nats/v2
```

## Quick Start

```go
package main

import (
	"github.com/taigrr/log-nats/v2/log"
)

func main() {
	// Connect to NATS and set as the default connection.
	if err := log.ConnectDefault("nats://127.0.0.1:4222"); err != nil {
		log.Fatal(err)
	}
	defer log.Flush()

	log.Info("hello world")
	log.Debugf("version: %s", "2.0.0")
}
```

## Namespaced Loggers

Create dedicated loggers per subsystem. Each logger publishes to its own
NATS subject (by default `logging.<namespace>.<LEVEL>`).

```go
apiLog := log.NewLogger("api")
apiLog.Info("request handled")

dbLog := log.NewLogger("database")
dbLog.Warnf("slow query: %dms", 230)
```

## slog Handler

Use the `slog` sub-package to integrate with Go's structured logging:

```go
package main

import (
	"log/slog"

	"github.com/taigrr/log-nats/v2/log"
	logslog "github.com/taigrr/log-nats/v2/slog"
)

func main() {
	if err := log.ConnectDefault("nats://127.0.0.1:4222"); err != nil {
		log.Fatal(err)
	}
	defer log.Flush()

	handler := logslog.NewHandler(
		logslog.WithNamespace("myapp"),
		logslog.WithLevel(slog.LevelInfo),
	)
	logger := slog.New(handler)
	logger.Info("request", "method", "GET", "path", "/health")
}
```

## NATS Subjects

Log entries are published to subjects derived from a configurable Go
template. The default is:

```
logging.{{.Namespace}}.{{.Level}}
```

Override it with `log.SetSubjectTemplate`:

```go
log.SetSubjectTemplate("logs.{{.Level}}.{{.Namespace}}")
```

Subscribe to all logs with a wildcard: `logging.>`

## Log Levels

`TRACE` · `DEBUG` · `INFO` · `NOTICE` · `WARN` · `ERROR` · `PANIC` · `FATAL`

`Panic` and `Fatal` behave like their `log` package counterparts — `Panic`
calls `panic()` and `Fatal` calls `os.Exit(1)` after publishing.

## License

[0BSD](LICENSE)
