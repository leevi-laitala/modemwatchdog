package main

import (
	"fmt"
	"log"
	"time"
)

func ping() bool {
	// Send head request to each url
	for _, url := range pingUrls {
		resp, err := client.Head(url)

		if err == nil && resp.StatusCode < 400 {
			resp.Body.Close()
			return true
		}
	}

	// If connection to all urls fail, internet access is probably down
	return false
}

func waitModemBoot() error {
	deadline := time.After(modemBootTimeDeadline * time.Second)
	ticker := time.NewTicker(pingInterval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			return fmt.Errorf("waiting for internet connectivity exceeded %d second timeout", modemBootTimeDeadline)
		case <-ticker.C:
			if ping() {
				log.Println("internet connectivity restored")
				return nil
			}
		}
	}
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

			err = waitModemBoot()
			if err != nil {
				log.Printf("%v", err)
			}
		}

		time.Sleep(pingInterval * time.Second)
	}
}
