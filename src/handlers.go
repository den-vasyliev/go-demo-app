package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

	switch AppRole {

	case "frontend":
		w.Write([]byte("OK"))

	case "backend":
		client := redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
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
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		probe, err := client.Set("readiness_probe", 0, 0).Result()
		log.Print(probe)
		if err != nil {
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}

		db, err := sql.Open("mysql", "")
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

func frontendHandler(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())
	message := fmt.Sprintf(`{"text":"%s"}`, r.URL.Query()["message"])
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cacheItem, err := client.Get(fmt.Sprintf("%x", md5.Sum([]byte(message)))).Result()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", rest("http://"+AppBackend, message))))
		//hashStr := fmt.Sprintf(`{"hash":"%s"}`, hash(message))
		//log.Print("Hash:", hashStr)
		//w.Write(rest("http://"+AppBackend, hashStr))

	} else {
		hexStr, _ := client.Get(cacheItem).Result()
		decoded, _ := hex.DecodeString(hexStr)
		w.Write([]byte(decoded))
	}
	//w.Write([]byte(fmt.Sprintf("%s", rest("http://"+AppBackend, message))))

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

		log.Print("Text: ", m.Text)

		client := redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		cacheItem, err := client.Get(fmt.Sprintf("%x", md5.Sum([]byte(m.Text)))).Result()
		if err != nil {

			hashStr := fmt.Sprintf(`{"hash":"%s"}`, hash(m.Text))
			log.Print("Hash:", hashStr)
			w.Write(rest("http://"+AppDatastore, hashStr))

		} else {
			hexStr, _ := client.Get(cacheItem).Result()
			decoded, _ := hex.DecodeString(hexStr)
			w.Write([]byte(decoded))
		}
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
