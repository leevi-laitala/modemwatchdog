package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	pingInterval = 10

	powercyclePeriod = 5
	powercycleRetries = 10
	powercycleReqRestartInterval = 3

	syslogAddr = "logger.internal.bebanen.com"
	syslogPort = 514
	syslogProtocol = "tcp"
)

var urls = []string{
	"http://127.0.0.1:5000/gen_204",
}

//var urls = []string{
//	"https://www.google.com/generate_204",
//	"https://1.1.1.1/cdn-cgi/trace",
//}

func main() {
	initLogging()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("Watchdog terminated. Received: %v", sig)
			os.Exit(1)
		}
	}()

	modemwatchdog(PythonPlug{
		apiUrl: "http://127.0.0.1:5000/api",
		apiTurnOff: "/turnoff",
		apiTurnOn: "/turnon",
		apiStatus: "/status",
	})
}
