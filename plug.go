package main

import (
	"log"
	"time"
	"fmt"
)

type SmartPlug interface {
	// true = power is on
	checkStatus() (bool, error)

	turnOn() error
	turnOff() error
}

func retryFunc[T any](maxRetries int, interval time.Duration, f func() (T, error)) (T, error) {
	var (
		err error
		res T
	)

	for retry := 1; retry < maxRetries+1; retry++ {
		res, err = f()
		if err == nil {
			break
		}

		log.Printf("Retrying... Attempt %d out of %d attempts. Error: %v", retry, maxRetries, err)

		time.Sleep(interval)
	}

	return res, err
}

// Turn off power, wait for some time and turn power back on
func powercycleModem(plug SmartPlug) error {
	var (
		err  error
		isOn bool
	)

	log.Println("Attempting powercycle")

	// Check initial status
	isOn, err = retryFunc(powercycleRetries, powercycleRetryWait*time.Second, plug.checkStatus)
	if err != nil {
		log.Printf("Cannot check plug status. Assuming power on...")
	}

	// Skip turning off if already off. This may be due to failed powercycle
	if !isOn {
		// Send poweroff signal
		_, err = retryFunc(powercycleRetries, powercycleRetryWait*time.Second, func() (struct{}, error) {
			return struct{}{}, plug.turnOff()
		})
		if err != nil {
			return fmt.Errorf("Failed to power off: %v", err)
		}

		// Check status
		isOn, err = retryFunc(powercycleRetries, powercycleRetryWait*time.Second, plug.checkStatus)
		if err != nil || isOn {
			return fmt.Errorf("Plug status cloud not be powered off: %v", err)
		}
	}

	log.Printf("Plug is powered off")

	log.Printf("Waiting for %d seconds before turning power back on", powercycleDuration)
	time.Sleep(powercycleDuration * time.Second)

	// Send poweron signal
	_, err = retryFunc(powercycleRetries, powercycleRetryWait*time.Second, func() (struct{}, error) {
		return struct{}{}, plug.turnOn()
	})
	if err != nil {
		return fmt.Errorf("Failed to power on")
	}

	// Check status
	isOn, err = retryFunc(powercycleRetries, powercycleRetryWait*time.Second, plug.checkStatus)
	if err != nil || !isOn {
		return fmt.Errorf("Failed to powercycle: Plug status cloud not be powered on")
	}

	log.Printf("Powercycle succeeded")
	return nil
}
