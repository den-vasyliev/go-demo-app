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
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/prometheus"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
)

// AppName app
//var AppName = "name"
var AppName = os.Getenv("APP_NAME")

// Version app
var Version = "version"

// BuildInfo app
var BuildInfo = "commit"

// Revision app
var Revision = fmt.Sprintf("%s version: %s+%s", AppName, Version, BuildInfo)

// AppPort app
var AppPort = os.Getenv("APP_PORT")

// AppDb name
var AppDb = os.Getenv("APP_DB_SERVER")

// AppRqm name
var AppRmq = os.Getenv("APP_RMQ_SERVER")

// AppRedis name
var AppRedis = os.Getenv("APP_REDIS_SERVER") + ":" + os.Getenv("APP_REDIS_PORT")

type apiStruct struct {
	APIVersion string `json:"apiVersion"`
}

func main() {
	log.Print(Revision)
	sink, _ := prometheus.NewPrometheusSink()
	metrics.NewGlobal(metrics.DefaultConfig("API"), sink)

	conn, err := amqp.Dial("amqp://guest:guest@" + AppRmq + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	
	q, err := ch.QueueDeclare(
	  "hello", // name
	  false,   // durable
	  false,   // delete when usused
	  false,   // exclusive
	  false,   // no-wait
	  nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	  failOnError(err, "Failed to register a consumer")
	  
	  forever := make(chan bool)
	  
	  go func() {
		for d := range msgs {
		  log.Printf("Received a message: %s", d.Body)
		}
	  }()
	  
	  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	  <-forever



	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/readinez", readinessHandler)
	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/", appHandler)

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

}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	var b []byte
	b = append([]byte(""), Revision...)
	w.Write(b)
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

func appHandler(w http.ResponseWriter, r *http.Request) {
	metrics.IncrCounter([]string{"requestCounter"}, 1)
	var m apiStruct
	switch r.Method {
	case "GET":
		log.Printf("Get GET Request!")
		u, err := url.Parse(r.RequestURI)
		if err != nil {
			log.Fatal(err)
		}
		q := u.Query()

		r, err := regexp.Compile(`^\*(?P<USSD_CODE>\d{3})[\*|\#](?P<USSD_DEST>\D{0,}\d{0,}).?(?P<USSD_EXT>.{0,}).?`)
		if err != nil {
			log.Print(err)
		}
		r2 := r.FindAllStringSubmatch(q.Get("calldestination"), -1)
		fmt.Println(r2[0][1])

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		  )
		  failOnError(err, "Failed to register a consumer")
		  
		  forever := make(chan bool)
		  go func() {
			for d := range msgs {
			  log.Printf("Received a message: %s", d.Body)
			}
		  }()
		  
		  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		  <-forever

		w.Write([]byte("Please use POST"))

	case "POST":
		b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err := json.Unmarshal(b, &m); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Print(err, len(m.APIVersion))
			}
		}
		if len(m.APIVersion) > 4 {
			enc := json.NewEncoder(w)
			enc.Encode(API(m.APIVersion))
		} else {
			w.Write([]byte("Incorrect Api Version"))
		}
		//w.Write([]byte(getRate(m.Msisdn)))
	}
}

// API func
func API(verson string) apiStruct {

	defer metrics.MeasureSince([]string{"API"}, time.Now())

	var apiResult apiStruct

	db, err := sql.Open("sqlite3", AppDb)
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
		Addr:     AppRedis, //Addr:     "redis:6379",
		Password: "",       // no password set
		DB:       0,        // use default DB
	})
	client.Ping()
	return apiResult
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
