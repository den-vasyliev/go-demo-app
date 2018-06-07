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

type greetingsText struct {
	Text string `json:"Text"`
}

type greetingsToken struct {
	Hash string `json:"Hash"`
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
	w.Write([]byte(fmt.Sprintf("%s", rest("http://service", `{"text":"devops career day"}`))))

}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Print("Text: ", m.Text)
		hashStr := fmt.Sprintf(`{"hash":"%s"}`, greetingsID(m.Text))
		log.Print("Hash:", hashStr)
		w.Write(rest("http://data", hashStr))

	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte(greetingsDB(m.Hash)))
	}
}

func greetingsID(decodedStr string) string {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	decodedStr = "New Feature"
	log.Print("DecodedStr: ", decodedStr)
	encodedStr := hex.EncodeToString([]byte(banner.PrintS(decodedStr)))
	log.Print("EncodedStr: ", encodedStr)
	hashStr := fmt.Sprintf("%x", md5.Sum([]byte(encodedStr)))
	client.Set(hashStr, encodedStr, 0)
	return hashStr
}

func greetingsDB(hash string) string {
	var Payload string
	db, err := sql.Open("mysql", AppDb)
	if err != nil {
		log.Print("Open db err: ")
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Print("Ping db err: ")
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
