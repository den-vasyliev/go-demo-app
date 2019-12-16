package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/nats-io/nats.go"
)

//DataHandler export broker msg func
func DataHandler(m *nats.Msg, i int) []byte {
	defer metrics.MeasureSince([]string{"DB"}, time.Now())

	var err error
	var t messageToken
	var Payload string
	var decoded []byte

	log.Print("Processing")

	json.Unmarshal(m.Data, &t)

	hexStr, err := CACHE.Get(t.Hash).Result()
	if err != nil {
		log.Print(err)
	}

	stmt, err := DB.Prepare("insert into demo values(null,?,?)")

	_, err = stmt.Exec(t.Hash, hexStr)

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	stmt, err = DB.Prepare("SELECT text FROM demo WHERE token = ? limit 1")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	// additional iteration
	err = stmt.QueryRow(t.Hash).Scan(&Payload) // WHERE number = 13

	stmt, err = DB.Prepare("DELETE from demo where token = ?")

	_, err = stmt.Exec(t.Hash)

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	REQ0 = REQ0 + 1

	if err != nil {
		log.Printf("QueryRowErr: %s", err) // proper error handling instead of panic in your app
	}

	decoded, err = hex.DecodeString(Payload)
	if err != nil {
		log.Printf("DecodeStringErr:%s", err)
	}
	return []byte(string(decoded))

}

func dataHandler(w http.ResponseWriter, r *http.Request) {

	var Payload string

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Print(err)
	}
	q := u.Query()

	stmt, err := DB.Prepare("insert into demo values(null,?,?)")

	_, err = stmt.Exec(q.Get("key"), q.Get("val"))

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	stmt, err = DB.Prepare("SELECT text FROM demo WHERE token = ? limit 1")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	// additional iteration
	_ = stmt.QueryRow(q.Get("key")).Scan(&Payload) // WHERE number = 13

	REQ0 = REQ0 + 1

	w.Write([]byte(fmt.Sprintf("%s", Payload)))
	/*
		case "POST":
			b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
			if err := json.Unmarshal(b, &m); err != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(422) // unprocessable entity
				if err := json.NewEncoder(w).Encode(err); err != nil {
					panic(err)
				}
			}
			w.Write([]byte(dataStore(m.Hash)))*/

}
