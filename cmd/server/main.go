package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"telemetryguard/internal/console"
	"telemetryguard/internal/escalate"
	"telemetryguard/internal/notify"
	"telemetryguard/internal/cycle"
	"telemetryguard/internal/store"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	data := flag.String("data", "./data", "data directory")
	flag.Parse()

	state := store.NewState(*data)
	handler := console.NewAPI(state)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now().UTC()
			if _, err := cycle.RunCycle(state, now); err != nil {
				log.Printf("eval cycle: %v", err)
			}
			if _, err := escalate.Tick(state, now); err != nil {
				log.Printf("escalation tick: %v", err)
			}
			notify.RetryFailed(state, now)
		}
	}()

	log.Printf("TelemetryGuard listening on %s with data dir %s", *addr, *data)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal(err)
	}
}
