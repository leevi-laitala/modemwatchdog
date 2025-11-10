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

func powercycleModem(plug SmartPlug) error {
	var (
		err  error
		isOn bool
	)

	log.Println("Attempting powercycle")

	isOn, err = retryFunc(powercycleRetries, powercycleReqRestartInterval*time.Second, plug.checkStatus)
	if err != nil {
		log.Printf("Cannot check plug status. Assuming power on...")
	}

	if !isOn {
		_, err = retryFunc(powercycleRetries, powercycleReqRestartInterval*time.Second, func() (struct{}, error) {
			return struct{}{}, plug.turnOff()
		})
		if err != nil {
			return fmt.Errorf("Failed to power off: %v", err)
		}

		isOn, err = retryFunc(powercycleRetries, powercycleReqRestartInterval*time.Second, plug.checkStatus)
		if err != nil || isOn {
			return fmt.Errorf("Plug status cloud not be powered off: %v", err)
		}
	}

	log.Printf("Plug is powered off")

	log.Printf("Waiting for %d seconds before turning power back on", powercyclePeriod)
	time.Sleep(powercyclePeriod * time.Second)

	_, err = retryFunc(powercycleRetries, powercyclePeriod*time.Second, func() (struct{}, error) {
		return struct{}{}, plug.turnOn()
	})
	if err != nil {
		return fmt.Errorf("Failed to power on")
	}

	isOn, err = retryFunc(powercycleRetries, powercyclePeriod*time.Second, plug.checkStatus)
	if err != nil || !isOn {
		return fmt.Errorf("Failed to powercycle: Plug status cloud not be powered on")
	}

	log.Printf("Powercycle succeeded")
	return nil
}
