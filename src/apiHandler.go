package main

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	"log"
	"net/http"
	"net/url"
	"time"

	metrics "github.com/armon/go-metrics"
)

func perfHandler(w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Print(err)
	}
	q := u.Query()

	msg, err := NC.Request("ascii.text", []byte(fmt.Sprintf(`{"text":"%s"}`, q.Get("text"))), 2*time.Second) // Read the reply

	var reply []byte
	if err != nil {
		log.Printf("Message: %e", err)
		reply = []byte(fmt.Sprintf("{%s}", err))
	} else {
		reply = msg.Data
	}

	//log.Printf("[api-reply]: %s", reply)
	REQ0 = REQ0 + 1

	w.Write(reply)
}

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

	INM.SetGauge([]string{"foo"}, 42)

	b, err := json.Marshal(APIReg)
	if err != nil {
		log.Print(err)
	}

	REQ0 = REQ0 + 1

	w.Write([]byte(b))
}
