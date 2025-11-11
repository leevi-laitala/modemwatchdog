package main

import (
	"encoding/json"
	"io"
)

type ShellyPlug struct {
	apiUrl     string
	apiTurnOff string
	apiTurnOn  string
	apiStatus  string
}

func (p ShellyPlug) checkStatus() (bool, error) {
	resp, err := makeReq("GET", p.apiUrl+p.apiStatus)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var data map[string]any
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return false, err
	}

	return data["output"] == true, nil
}

func (p ShellyPlug) turnOn() error {
	resp, err := makeReq("GET", p.apiUrl+p.apiTurnOn)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data map[string]any
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return err
	}

	return nil
}

func (p ShellyPlug) turnOff() error {
	resp, err := makeReq("GET", p.apiUrl+p.apiTurnOff)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data map[string]any
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return err
	}

	return nil
}
