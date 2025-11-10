package main

import (
	"fmt"
	"log"
	"log/syslog"
	"io"
	"os"
)

func initLogging() {
	// Use stdout by default
	var writer io.Writer = os.Stdout

	// Use syslog if server is defined
	if logSyslogAddr != "" {
		var err error

		addr := fmt.Sprintf("%s:%d", logSyslogAddr, logSyslogPort)

		writer, err = syslog.Dial(logSyslogProtocol, addr, syslog.LOG_INFO | syslog.LOG_DAEMON, "modemwatchdog")
		if err != nil {
			log.Fatalf("Failed to connect to syslog: %v", err)
		}
	}

	log.SetOutput(writer)
}
