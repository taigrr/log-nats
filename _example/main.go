package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/taigrr/log-nats/v2/log"
)

func main() {
	url := flag.String("url", "nats://127.0.0.1:4222", "NATS server URL")
	flag.Parse()

	if err := log.ConnectDefault(*url); err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to NATS: %v\n", err)
		os.Exit(1)
	}
	defer log.Flush()

	// Package-level functions use the "default" namespace.
	log.Info("application started")
	log.Debugf("connected to %s", *url)

	// Create a namespaced logger for a subsystem.
	apiLog := log.NewLogger("api")
	apiLog.Infof("listening on :8080")

	for i := 0; ; i++ {
		log.Tracef("tick %d", i)
		apiLog.Warnf("request %d slow", i)
		time.Sleep(2 * time.Second)
	}
}
