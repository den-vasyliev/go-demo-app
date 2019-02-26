package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"flag"
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

func versionHandler(w http.ResponseWriter, r *http.Request) {
	var b []byte
	b = append([]byte(""), Revision...)
	w.Write(b)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Healthz: alive!"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	AppRole := flag.String("r", "neuart", "application role")
	AppLicense := flag.String("l", "122345", "application license")

	flag.Parse()
	switch *AppRole {

	case "frontend":
		if *AppLicense != "" {
			w.Write([]byte("OK"))
		} else {
			http.Error(w, "No License", http.StatusServiceUnavailable)

		}

	case "backend":
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", AppDbNoSQL, AppDbNoSQLPort),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		probe, err := client.Ping().Result()
		log.Print(probe, err)
		if err != nil {
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}

	case "datastore":
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", AppDbNoSQL, AppDbNoSQLPort),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		probe, err := client.Set("readiness_probe", 0, 0).Result()
		log.Print(probe)
		if err != nil {
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}

		db, err := sql.Open("mysql", AppDb)
		if err != nil {
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}
		defer db.Close()
		err = db.Ping()

		if err != nil {
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}

		w.Write([]byte("200"))

	default:
		http.Error(w, "Not Ready", http.StatusServiceUnavailable)

	}

}

func backendHandler(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())

	var m messageText

	switch r.Method {
	case "GET":
		log.Printf("Get GET Request!")
		w.Write([]byte("Please use POST"))

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

func datastoreHandler(w http.ResponseWriter, r *http.Request) {
	var m messageToken
	switch r.Method {
	case "GET":
		log.Printf("Get GET Request!")
		w.Write([]byte("Please use POST"))

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

