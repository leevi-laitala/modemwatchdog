package main

import (
	"time"
	"fmt"
)

type MockPlug struct {}

var status bool = false

func (p MockPlug) checkStatus() (bool, error) {
	time.Sleep(3 * time.Second)
	return status, nil
}

func (p MockPlug) turnOn() error {
	status = true
	time.Sleep(3 * time.Second)
	return fmt.Errorf("cannot turn on")
}

func (p MockPlug) turnOff() error {
	status = false
	time.Sleep(3 * time.Second)
	return nil
}
