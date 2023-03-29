package main

import (
	"flag"
	"time"

	"github.com/taigrr/log-nats/log"
)

func generateLogs() {
	for {
		log.Info("This is an info log!")
		log.Trace("This is a trace log!")
		log.Debug("This is a debug log!")
		log.Warn("This is a warn log!")
		log.Error("This is an error log!")
		time.Sleep(2 * time.Second)
	}
}

func main() {
	defer log.Flush()
	flag.Parse()
	generateLogs()
}
