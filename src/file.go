package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func file() {
	router := mux.NewRouter()
	path := flag.String("p", "/static/", "path to serve static files")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	router.PathPrefix(*path).Handler(http.StripPrefix(*path, http.FileServer(http.Dir(*directory))))
	//http.Handle("/static/", http.StripPrefix(strings.TrimRight(path, "/"), http.FileServer(http.Dir(*directory))))
	log.Fatal(http.ListenAndServe(":8880", router))

}
