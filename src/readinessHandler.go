package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {

	flag.Parse()
	switch Role {

	case "api":
		w.Write([]byte("READY"))

	case "img", "ml5":
		w.Write([]byte("READY"))

	case "ascii":
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", AppCache, AppCachePort),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		probe, err := client.Ping().Result()
		log.Print(probe, err)
		if err != nil {
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}

	case "data":
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", AppCache, AppCachePort),
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
		db.SetConnMaxLifetime(time.Second * 20)

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
