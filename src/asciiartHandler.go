package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/go-redis/redis"
)

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())

	var m messageText

	switch r.Method {

	case "GET":
		var b []byte
		b = append([]byte(""), Environment...)
		w.Write(b)

	case "POST":
		b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err := json.Unmarshal(b, &m); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", AppDbNoSQL, AppDbNoSQLPort),
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		log.Print("Text: ", m.Text)

		hashStr, encodedStr := hash(m.Text)

		client.Set(hashStr, encodedStr, 0)
		client.Set(fmt.Sprintf("%x", md5.Sum([]byte(m.Text))), hashStr, 0)

		log.Print("Hash:", hashStr)
		// message brocker placeholder
		w.Write(rest("http://"+AppDatastore, fmt.Sprintf(`{"hash":"%s"}`, hashStr)))

	}
}
