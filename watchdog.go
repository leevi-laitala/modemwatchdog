package main

import (
	"net/http"
	"time"
	"log"
)

func ping() bool {
	// Send head request to each url
	for _, url := range pingUrls {
		resp, err := http.Head(url)

		if err == nil && resp.StatusCode < 400 {
			resp.Body.Close()
			return true
		}
	}

	// If connection to all urls fail, internet access is probably down
	return false
}

func modemwatchdog(plug SmartPlug) {
	log.Println("Modem watchdog started")

	for {
		if !ping() {
			// Check again after a short period
			time.Sleep(pingInterval * time.Second)
			if ping() {
				continue
			}

			log.Println("No internet connection")

			err := powercycleModem(plug)
			if err != nil {
				log.Printf("Powercycle failed: %v", err)
			}
		}

		time.Sleep(pingInterval * time.Second)
	}
}
