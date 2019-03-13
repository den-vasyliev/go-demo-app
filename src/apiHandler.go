package main

import (
	_ "image/jpeg"
	"net/http"
	"time"

	metrics "github.com/armon/go-metrics"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	var b []byte
	b = append([]byte(""), Environment...)
	w.Write(b)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Healthz: alive!"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())
	switch Version {
	case "v2":
		w.Write([]byte("Welcome to k8s-art version:2 !\n\nHave a fun with:\n1. curl -XPOST --data '{\"text\":\"<TEXT>\"}' '<HOST>:<PORT>/ascii/'\n2. curl -F 'image=@<IMAGE>' '<HOST>:<PORT>/img/'\n"))
	case "v3":
		w.Write([]byte("Welcome to k8s-art version:3 !\n\nHave a fun with:\n1. curl -XPOST --data '{\"text\":\"<TEXT>\"}' '<HOST>:<PORT>/ascii/'\n2. curl -F 'image=@<IMAGE>' '<HOST>:<PORT>/img/'\n3. curl <HOST>:<PORT>/ml5/"))
	default:
		w.Write([]byte("Welcome to k8s-art version:1 !\n\nHave a fun with:\n1. curl -XPOST --data '{\"text\":\"<TEXT>\"}' '<HOST>:<PORT>/ascii/'\n"))

	}
}
