package main

import (
	"fmt"
	"log"
	"log/syslog"
	"io"
	"os"
)

func initLogging() {
	addr := fmt.Sprintf("%s:%d", syslogAddr, syslogPort)

	writer, err := syslog.Dial(syslogProtocol, addr, syslog.LOG_INFO | syslog.LOG_DAEMON, "modemwatchdog")
	if err != nil {
		log.Fatalf("Failed to connect to syslog: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, writer)

	log.SetOutput(mw)
}
