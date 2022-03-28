package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func SendRequest(url string, httpBody io.Reader) ([]byte, error) {
	log.Println("Sending request to", url)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://"+url, httpBody)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
