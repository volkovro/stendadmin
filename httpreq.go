package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PostReq - Use for POST requests. Returns the request body []byte
func PostReq(data map[string]interface{}, url string, headers ...map[string]string) (response []byte, err error) {

	reqBody, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if len(headers) != 0 {
		for i := 0; i < len(headers); i++ {
			for s, e := range headers[i] {
				req.Header.Add(s, e)
			}
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode >= 300 {
		err := fmt.Errorf("HTTP error:\n       url: %s\n    status: %d", url, resp.StatusCode)
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return result, err
}

// GetReq - Use for GET request. Returns the request body []byte
func GetReq(data map[string]interface{}, url string, headers ...map[string]string) (response []byte, err error) {

	reqBody, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if len(headers) != 0 {
		for i := 0; i < len(headers); i++ {
			for s, e := range headers[i] {
				req.Header.Add(s, e)
			}
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode >= 300 {
		err := fmt.Errorf("HTTP error:\n       url: %s\n    status: %d", url, resp.StatusCode)
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return result, err
}

// DelReq - Use for DELETE request. Returns the request body []byte
func DelReq(data map[string]interface{}, url string, headers ...map[string]string) (response []byte, err error) {

	reqBody, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if len(headers) != 0 {
		for i := 0; i < len(headers); i++ {
			for s, e := range headers[i] {
				req.Header.Add(s, e)
			}
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode >= 300 {
		err := fmt.Errorf("HTTP error:\n       url: %s\n    status: %d", url, resp.StatusCode)
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return result, err
}
