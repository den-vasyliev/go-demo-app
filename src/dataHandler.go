package main

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	var m messageToken
	switch r.Method {
	case "GET":
		log.Printf("Get GET Request!")

		w.Write([]byte(fmt.Sprintf("%s", Environment)))

	case "POST":
		b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err := json.Unmarshal(b, &m); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
		w.Write([]byte(dataStore(m.Hash)))
	}
}

//DataHandler export brocker msg func
func DataHandler(m *nats.Msg, i int) []byte {

	var t messageToken
	json.Unmarshal(m.Data, &t)

	return []byte(dataStore(string(t.Hash)))
}
