package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

func ascii(w http.ResponseWriter, r *http.Request) {

	var b []byte
	b = append([]byte(""), Environment...)

	w.Write(b)

}

func data(w http.ResponseWriter, r *http.Request) {

	var b []byte
	b = append([]byte(""), Environment...)

	w.Write(b)

}

func version(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(APIReg)
	if err != nil {
		log.Print(err)
	}

	REQ0 = REQ0 + 1

	w.Write([]byte(b))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	REQ0 = REQ0 + 1

	w.Write([]byte("Healthz: alive!"))
}

func readinez(w http.ResponseWriter, r *http.Request) {

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
