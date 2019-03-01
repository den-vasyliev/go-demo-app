package main

import (
	"fmt"
	_ "image/jpeg"
	"io/ioutil"
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

func apiHandler(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())
	log.Printf("Read index.html")
	dat, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Printf("No found: index.html")
		w.Write([]byte(fmt.Sprintf("%s", "Not found")))
	}
	w.Write(dat)

}
