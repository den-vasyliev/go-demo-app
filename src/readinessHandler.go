package main

import (
	"flag"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {

	flag.Parse()
	switch Role {

	case "api":
		w.Write([]byte("READY"))

	case "img", "ml5":
		w.Write([]byte("READY"))

	case "ascii":
		_, err := CACHE.Ping().Result()
		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		} else {

			w.Write([]byte("READY"))
		}

	case "data":
		_, err := CACHE.Set("readiness_probe", 0, 0).Result()
		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}
		log.Print(AppDb)

		err = DB.Ping()

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
