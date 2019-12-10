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
	// Send the request.
	// If processing is synchronous, use Request() which returns the response message.
	log.Printf("pub: %s", q.Get("text"))

	msg, err := NC.Request("ascii.text", []byte(fmt.Sprintf(`{"text":"%s"}`, q.Get("text"))), 2*time.Second) // Read the reply

	var reply []byte
	if err != nil {
		log.Printf("Message: %e", err)
		reply = []byte(fmt.Sprintf("{%s}", err))
	} else {
		reply = msg.Data
	}

	log.Printf("[api-reply]: %s", reply)

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

	b, err := json.Marshal(APIReg)
	if err != nil {
		log.Print(err)
	}
	w.Write([]byte(b))
}
