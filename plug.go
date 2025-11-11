package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type SmartPlug interface {
	// true = power is on
	checkStatus() (bool, error)

	turnOn() error
	turnOff() error
}

var client = http.Client{
	Timeout: 5 * time.Second,
}

func makeReq(method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response failed with http code %d", resp.StatusCode)
	}

	return resp, nil
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
	if isOn {
		// Send poweroff signal
		_, err = retryFunc(powercycleRetries, powercycleRetryWait*time.Second, func() (struct{}, error) {
			return struct{}{}, plug.turnOff()
		})
		if err != nil {
			return fmt.Errorf("failed to power off: %v", err)
		}

		// Check status
		isOn, err = retryFunc(powercycleRetries, powercycleRetryWait*time.Second, plug.checkStatus)
		if err != nil || isOn {
			return fmt.Errorf("plug status cloud not be powered off: %v", err)
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
		return fmt.Errorf("failed to power on")
	}

	// Check status
	isOn, err = retryFunc(powercycleRetries, powercycleRetryWait*time.Second, plug.checkStatus)
	if err != nil || !isOn {
		return fmt.Errorf("failed to powercycle: Plug status cloud not be powered on")
	}

	log.Printf("Powercycle succeeded")
	return nil
}
