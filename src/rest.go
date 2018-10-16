package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func rest(url string, jsonStr string) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Print("response Status:", resp.Status)
	log.Print("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
