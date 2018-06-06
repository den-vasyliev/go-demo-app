package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/CrowdSurge/banner"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// AppName app
var AppName = "name"

// Version app
var Version = "version"

// BuildInfo app
var BuildInfo = "commit"

// Revision app
var Revision = fmt.Sprintf("%s version: %s+%s", AppName, Version, BuildInfo)

// AppPort app
var AppPort = os.Getenv("APP_PORT")

//AppDB name
var AppDb = "db/name"

type greetingsToken struct {
	Token string `json:"token"`
}

type greetingsText struct {
	Key string `json:"key"`
}

func main() {
	log.Print(Revision)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthzHandler)

	switch AppName {
	case "front":
		router.HandleFunc("/", frontHandler)

	case "service":
		router.HandleFunc("/", serviceHandler)

	case "data":
		router.HandleFunc("/", dataHandler)

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

func frontHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	log.Print(b)
	b = []byte("test")
	w.Write(rest("http://service", `{"token":"test"}`))

}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	var m greetingsToken
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
		log.Print(m.Token)
		key := fmt.Sprintf(`{"key":"%s"}`, greetingsID(m.Token))
		log.Print(key)
		w.Write(rest("http://data", key))

	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	var m greetingsText
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
		w.Write([]byte(greetingsDB(m.Key)))
	}
}

func greetingsID(token string) string {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Set("token", "token", 0).Err()
	val, err := client.Get(token).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func greetingsDB(key string) string {
	var text string
	db, err := sql.Open("mysql", AppDb)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	//stmtOut, err := db.Prepare("SELECT text FROM greetings WHERE token = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	//defer stmtOut.Close()
	//var squareNum int
	log.Print(key)
	log.Print(banner.PrintS(key))
	_, err = db.Exec("drop table IF EXISTS demoTable")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS demoTable (id INT NOT NULL AUTO_INCREMENT, token VARCHAR(100), text VARCHAR(100), PRIMARY KEY(id))")
	_, err = db.Exec("insert into demoTable values(1,?,?)", key, banner.PrintS(key))
	err = db.QueryRow("SELECT text FROM demoTable WHERE token = ?", key).Scan(&text) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return text
}

func rest(url string, jsonStr string) []byte {
	//	var jsonStr = []byte(`{"token":"devopscareerday"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
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
