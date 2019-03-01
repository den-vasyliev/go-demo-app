package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	AppRole := flag.String("role", "neuart", "application role")
	AppLicense := flag.String("lic", "122345", "application license")

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
