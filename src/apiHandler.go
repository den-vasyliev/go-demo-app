package main

import (
	_ "image/jpeg"
	"net/http"
	"time"

	metrics "github.com/armon/go-metrics"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	var b []byte
	b = append([]byte(""), Environment...)
	w.Write(b)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Healthz: alive!"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())
	w.Write("Welcome to k8s-art!\nHave a fun with /ascii/?banner=<TEXT> /img/ POST <IMG> /ml5/")

}
