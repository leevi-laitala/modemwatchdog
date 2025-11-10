package main

import (
	"io"
	"fmt"
	"net/http"
	"time"
	"encoding/json"
)

type PythonPlug struct {
	apiUrl string
	apiTurnOff string
	apiTurnOn string
	apiStatus string
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

func (p PythonPlug) checkStatus() (bool, error) {
	resp, err := makeReq("GET", p.apiUrl + p.apiStatus)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

    var data map[string]interface{}
    err = json.Unmarshal([]byte(body), &data)
	if err != nil || data["success"] != nil {
		return false, err
	}

	return data["currentStatus"] == "on", nil
}

func (p PythonPlug) turnOn() error {
	resp, err := makeReq("POST", p.apiUrl + p.apiTurnOn)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

    var data map[string]interface{}
    err = json.Unmarshal([]byte(body), &data)
	if err != nil || data["success"] != nil {
		return err
	}

	return nil
}

func (p PythonPlug) turnOff() error {
	resp, err := makeReq("POST", p.apiUrl + p.apiTurnOff)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

    var data map[string]interface{}
    err = json.Unmarshal([]byte(body), &data)
	if err != nil || data["success"] != nil {
		return err
	}

	return nil
}
