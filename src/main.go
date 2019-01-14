package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// AppNmae app
var AppName = os.Getenv("APP_Name")

// AppRole app
var AppRole = os.Getenv("APP_ROLE")

// AppPort app
var AppPort = os.Getenv("APP_PORT")

// AppLicense app
var AppLicense = os.Getenv("APP_LICENSE")

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
var Revision = fmt.Sprintf("%s %s version: %s+%s", AppName, AppRole, Version, BuildInfo)

// NewFeature changes mock
var NewFeature = ""

// AppDbNoSql app
var AppDbNoSql = os.Getenv("APP_DB_NO_SQL")

// AppDbNoSql app
var AppDbNoSqlPort = os.Getenv("APP_DB_NO_SQL_PORT")

type messageText struct {
	Text string `json:"Text"`
}

type messageToken struct {
	Hash string `json:"Hash"`
}

func main() {
	log.Print(Revision)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/readinez", readinessHandler)
	router.Handle("/metrics", promhttp.Handler())

	switch AppRole {
	case "frontend":
		router.HandleFunc("/", frontendHandler)

	case "backend":
		router.HandleFunc("/", backendHandler)

	case "datastore":
		router.HandleFunc("/", datastoreHandler)

	}
	log.Fatal(http.ListenAndServe(":"+AppPort, router))
}
