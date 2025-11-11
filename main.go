package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initLogging()

	// Capture termination signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("Watchdog terminated. Received: %v", sig)
			os.Exit(1)
		}
	}()

	plug := initPlug()

	modemwatchdog(plug)
}
