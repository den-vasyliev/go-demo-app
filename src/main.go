package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// AppName app
var AppName = os.Getenv("APP_Name")

/**
// AppRole app
var AppRole = os.Getenv("APP_ROLE")

// AppPort app
var AppPort = os.Getenv("APP_PORT")

// AppLicense app
var AppLicense = os.Getenv("APP_LICENSE")
**/

// AppBackend app
var AppBackend = os.Getenv("APP_BACKEND")

// AppDatastore app
var AppDatastore = os.Getenv("APP_DATASTORE")

// AppDb name
var AppDb = os.Getenv("APP_DB")

// Version app
var Version = "version"

// BuildInfo app
var BuildInfo = "commit"

// Revision app
var Revision = fmt.Sprintf("%s %s version: %s+%s", AppName, Version, BuildInfo)

// NewFeature changes mock
var NewFeature = ""

// AppDbNoSQL app
var AppDbNoSQL = os.Getenv("APP_DB_NO_SQL")

// AppDbNoSQLPort app
var AppDbNoSQLPort = os.Getenv("APP_DB_NO_SQL_PORT")

type messageText struct {
	Text string `json:"Text"`
}

type messageToken struct {
	Hash string `json:"Hash"`
}

func main() {
	AppRole := flag.String("role", "neuart", "application role")
	AppPort := flag.String("port", "8080", "application port")

	flag.Parse()

	log.Print(Revision)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/readinez", readinessHandler)
	router.Handle("/metrics", promhttp.Handler())

	switch *AppRole {
	case "frontend":
		initOptions()
		router.HandleFunc("/", frontendHandler)

	case "backend":
		router.HandleFunc("/", backendHandler)

	case "datastore":
		router.HandleFunc("/", datastoreHandler)

	case "neuart":
		path := flag.String("path", "/s/", "path to serve static files")
		directory := flag.String("dir", "./art", "the directory of static file to host")
		flag.Parse()

		router.PathPrefix(*path).Handler(http.StripPrefix(*path, http.FileServer(http.Dir(*directory))))

		router.HandleFunc("/", uploadHandler)

	}
	log.Fatal(http.ListenAndServe(":"+*AppPort, router))
}

