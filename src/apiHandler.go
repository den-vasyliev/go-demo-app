package main

import (
	"encoding/json"
	_ "image/jpeg"
	"log"
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

func api(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())

	b, err := json.Marshal(APIReg)
	if err != nil {
		log.Print(err)
	}
	w.Write([]byte(b))
}
