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
	w.Write([]byte("Welcome to k8s-art!\n\nHave a fun with:\n1. /ascii/ POST --data '{\"text\":\"<TEXT>\"}'\n2. /img/ POST <IMG>\n3. /ml5/"))

}
