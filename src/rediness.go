package main

import (
	"log"
	"net/http"
	"time"
)

func readiness(url string) string {
	c := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := c.Get(url)
	if err != nil {
		log.Print(err)
		return resp.Status
	}
	defer resp.Body.Close()

	log.Print("response Status:", resp.Status)
	return resp.Status
}
