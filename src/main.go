package main

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CrowdSurge/banner"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

// AppRole app
var AppRole = os.Getenv("APP_ROLE")

// AppPort app
var AppPort = os.Getenv("APP_PORT")

// AppBackend app
var AppBackend = os.Getenv("APP_BACKEND")

// AppDatastore app
var AppDatastore = os.Getenv("APP_DATASTORE")

// AppDb name
var AppDb = os.Getenv("APP_DB")

// Version app
var Version = "version"

// BuildInfo app
var BuildInfo = "commit"

// Revision app
var Revision = fmt.Sprintf("%s version: %s+%s", AppRole, Version, BuildInfo)

// NewFeature changes mock
var NewFeature = ""

type messageText struct {
	Text string `json:"Text"`
}

type messageToken struct {
	Hash string `json:"Hash"`
}

func main() {
	log.Print(Revision)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/readinez", readinessHandler)
	// router.Handle("/metrics", promhttp.Handler())

	switch AppRole {
	case "frontend":
		router.HandleFunc("/", frontendHandler)

	case "backend":
		router.HandleFunc("/", backendHandler)

	case "datastore":
		router.HandleFunc("/", datastoreHandler)

	}
	log.Fatal(http.ListenAndServe(":"+AppPort, router))
}

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
	// defer metrics.MeasureSince([]string{"API"}, time.Now())
	message := fmt.Sprintf(`{"text":"%s"}`, r.URL.Query()["message"])
	w.Write([]byte(fmt.Sprintf("%s", rest("http://"+AppBackend, message))))

}

func backendHandler(w http.ResponseWriter, r *http.Request) {
	// defer metrics.MeasureSince([]string{"API"}, time.Now())
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
		if NewFeature != "" {
			m.Text = NewFeature
		}
		log.Print("Text: ", m.Text)
		hashStr := fmt.Sprintf(`{"hash":"%s"}`, messageID(m.Text))
		log.Print("Hash:", hashStr)
		w.Write(rest("http://"+AppDatastore, hashStr))

	}
}

func messageID(decodedStr string) string {
	// defer metrics.MeasureSince([]string{"API"}, time.Now())
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	log.Print("DecodedStr: ", decodedStr)
	encodedStr := hex.EncodeToString([]byte(banner.PrintS(decodedStr)))
	log.Print("EncodedStr: ", encodedStr)
	hashStr := fmt.Sprintf("%x", md5.Sum([]byte(encodedStr)))
	client.Set(hashStr, encodedStr, 0)
	return hashStr
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

func dataStore(hash string) string {
	var Payload string

	db, err := sql.Open("mysql", AppDb)
	if err != nil {
		log.Print("Open db err: ")
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Print("Ping db err")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	hexStr, err := client.Get(hash).Result()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS demoTable (id INT NOT NULL AUTO_INCREMENT, token VARCHAR(100), text TEXT, PRIMARY KEY(id))")
	_, err = db.Exec("insert into demoTable values(null,?,?)", hash, hexStr)

	err = db.QueryRow("SELECT text FROM demoTable WHERE token = ?", hash).Scan(&Payload) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	decoded, err := hex.DecodeString(Payload)
	return string(decoded)
}

func rest(url string, jsonStr string) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Print("response Status:", resp.Status)
	log.Print("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func readiness(url string) string {
	c := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := c.Get(url)
	if err != nil {
		log.Print(err)
		return resp.Status
	}
	defer resp.Body.Close()

	log.Print("response Status:", resp.Status)
	return resp.Status
}
