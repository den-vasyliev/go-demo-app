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
		_, err := client.Ping().Result()
		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		} else {

			w.Write([]byte("READY"))
		}

	case "data":
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", AppCache, AppCachePort),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		_, err := client.Set("readiness_probe", 0, 0).Result()
		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}
		log.Print(AppDb)
		db, err := sql.Open("mysql", AppDb)
		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}
		db.SetConnMaxLifetime(time.Second * 20)

		defer db.Close()
		err = db.Ping()

		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		} else {

			w.Write([]byte("READY"))

		}

	default:
		http.Error(w, "Not Ready", http.StatusServiceUnavailable)

	}

}
