package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

func SendRequest(url string, httpBody io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, httpBody)
	req.Close = true

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
