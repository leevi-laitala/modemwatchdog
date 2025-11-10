package main

import (
	"net/http"
	"time"
	"log"
)

func ping() bool {
	for _, url := range urls {
		resp, err := http.Head(url)

		if err == nil && resp.StatusCode < 400 {
			resp.Body.Close()
			return true
		}

	}

	return false
}

func modemwatchdog(plug SmartPlug) {
	log.Println("Modem watchdog started")

	for {
		netconnection := ping()

		if !netconnection {
			log.Println("No network connection")

			err := powercycleModem(plug)
			if err != nil {
				log.Printf("Powercycle failed: %v", err)
			}
		}

		time.Sleep(pingInterval * time.Second)
	}
}
