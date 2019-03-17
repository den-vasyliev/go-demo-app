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

// AppLicense app
var AppLicense = os.Getenv("APP_LICENSE")

// AppBackend app
var AppBackend = os.Getenv("APP_BACKEND")

// AppDatastore app
var AppDatastore = os.Getenv("APP_DATASTORE")

// AppDb name
var AppDb = os.Getenv("APP_DB")

// AppDbNoSQL app
var AppDbNoSQL = os.Getenv("APP_DB_NO_SQL")

// AppDbNoSQLPort app
var AppDbNoSQLPort = os.Getenv("APP_DB_NO_SQL_PORT")

// Version app
var Version = "v3"

// Environment app
var Environment = ""

type messageText struct {
	Text string `json:"Text"`
}

type messageToken struct {
	Hash string `json:"Hash"`
}

// Role application name
var Role = ""

func main() {

	initOptions()
	AppName := flag.String("name", "k8s:art", "application name")
	AppRole := flag.String("role", "api", "app role: api data ascii img ml5")
	AppPort := flag.String("port", "8080", "application port")
	AppPath := flag.String("path", "/static/", "path to serve static files")
	AppDir := flag.String("dir", "./ml5", "the directory of static files to host")
	ModelsPath := flag.String("mpath", "/models/", "path to serve models files")
	ModelsDir := flag.String("mdir", "./ml5/models", "the directory of models files to host")

	flag.Parse()
	// Environment app
	Environment = fmt.Sprintf("%s version:%s role:%s port:%s", *AppName, Version, *AppRole, *AppPort)
	Role = *AppRole

	log.Print(Environment)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/readinez", readinessHandler)
	router.Handle("/metrics", promhttp.Handler())

	switch *AppRole {

	case "api":
		router.HandleFunc("/", apiHandler)

	case "data":
		router.HandleFunc("/", dataHandler)

	case "ascii":
		router.HandleFunc("/", asciiHandler)

	case "img":
		router.HandleFunc("/", imgHandler)

	case "ml5":

		router.PathPrefix(*AppPath).Handler(http.StripPrefix(*AppPath, http.FileServer(http.Dir(*AppDir))))
		router.PathPrefix(*ModelsPath).Handler(http.StripPrefix(*ModelsPath, http.FileServer(http.Dir(*ModelsDir))))

		router.HandleFunc("/", ml5Handler)

	}
	log.Fatal(http.ListenAndServe(":"+*AppPort, router))
}

func initOptions() {

	flag.StringVar(&imageFilename,
		"f",
		"",
		"Image filename to be convert")
	flag.Float64Var(&ratio,
		"r",
		convertDefaultOptions.Ratio,
		"Ratio to scale the image, ignored when use -w or -g")
	flag.IntVar(&fixedWidth,
		"w",
		convertDefaultOptions.FixedWidth,
		"Expected image width, -1 for image default width")
	flag.IntVar(&fixedHeight,
		"g",
		convertDefaultOptions.FixedHeight,
		"Expected image height, -1 for image default height")
	flag.BoolVar(&fitScreen,
		"s",
		convertDefaultOptions.FitScreen,
		"Fit the terminal screen, ignored when use -w, -g, -r")
	flag.BoolVar(&colored,
		"c",
		convertDefaultOptions.Colored,
		"Colored the ascii when output to the terminal")
	flag.BoolVar(&reversed,
		"i",
		convertDefaultOptions.Reversed,
		"Reversed the ascii when output to the terminal")
	flag.BoolVar(&stretchedScreen,
		"t",
		convertDefaultOptions.StretchedScreen,
		"Stretch the picture to overspread the screen")
}
