package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// API is api ref
var API = make(map[string]string)

// AppLicense app
var AppLicense = os.Getenv("APP_LICENSE")

// AppASCII app
var AppASCII = getEnv("APP_ASCII", "ascii")

// AppDatastore app
var AppDatastore = getEnv("APP_DATASTORE", "data")

// AppDbUser name
var AppDbUser = getEnv("AppDbPort", "root")

// AppDbPort name
var AppDbPort = getEnv("APP_DB_PORT", "3306")

// AppDbName name
var AppDbName = getEnv("APP_DB_NAME", "demo")

// AppDb name
var AppDb = getEnv("APP_DB", AppDbUser+"@tcp(127.0.0.1:"+AppDbPort+")/"+AppDbName)

// AppCache app
var AppCache = getEnv("APP_CACHE", "127.0.0.1")

// AppCachePort app
var AppCachePort = getEnv("APP_CACHE_PORT", "6379")

// AppCacheExpire app
var AppCacheExpire = getEnv("APP_CACHE_EXPIRE", "120s")

// Version app
var Version = "v3"

// Environment app
var Environment = ""

// APIReg is a api map
var APIReg = make(map[string]string)

// NC nats brocker
var NC *nats.Conn

type messageText struct {
	Text string `json:"Text"`
}

type messageToken struct {
	Hash string `json:"Hash"`
}

// Role application name
var Role = ""

func main() {

	API["ascii"] = "curl -XPOST --data '{text:TEXT}' HOST/ascii/"
	API["img"] = "curl -F 'image=@IMAGE' HOST/img/"
	API["ml5"] = "curl HOST/ml5/"
	API["data"] = "broker message api"

	initOptions()
	AppName := flag.String("name", "k8sdiy", "application name")
	AppRole := flag.String("role", "api", "app role: api data ascii img ml5")
	AppPort := flag.String("port", "8080", "application port")
	AppPath := flag.String("path", "/static/", "path to serve static files")
	AppDir := flag.String("dir", "./ml5", "the directory of static files to host")
	ModelsPath := flag.String("mpath", "/models/", "path to serve models files")
	ModelsDir := flag.String("mdir", "./ml5/models", "the directory of models files to host")

	var urls = flag.String("server", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var showTime = flag.Bool("timestamp", false, "Display timestamps")
	var queueName = flag.String("queGroupName", "K8S-NATS-Q", "Queue Group Name")
	var showHelp = flag.Bool("help", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage

	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	// Environment app
	Role = *AppRole

	Environment = fmt.Sprintf("%s-%s:%s", *AppName, Role, Version)

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Responder")}
	opts = setupConnOptions(opts)

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Connect to NATS
	var err error

	NC, err = nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	if err := NC.LastError(); err != nil {
		log.Fatal(err)
	}

	subj, i := *AppRole+".*", 0

	NC.QueueSubscribe(subj, *queueName+*AppRole, func(msg *nats.Msg) {
		i++
		//log
		printMsg(msg, i)

		if *AppRole == "api" {

			APIReg[msg.Subject] = string(msg.Data)

		} else if *AppRole == "ascii" {

			msg.Respond(ASCIIHandler(msg, i))

		} else {

			msg.Respond(DataHandler(msg, i))

		}
	})
	NC.Flush()

	log.Printf("Listening on [%s]: %s", subj, Environment)

	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/readinez", readinessHandler)
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/perf", perfHandler)

	switch *AppRole {

	case "api":

		router.HandleFunc("/", api)

	case "ascii":

		if err := NC.Publish("api."+Environment, []byte(API["ascii"])); err != nil {
			log.Fatal(err)
		}

		router.HandleFunc("/", ascii)

	case "img":
		if err := NC.Publish("api."+Environment, []byte(API["img"])); err != nil {
			log.Fatal(err)
		}
		router.HandleFunc("/", img)

	case "ml5":
		if err := NC.Publish("api."+Environment, []byte(API["img"])); err != nil {
			log.Fatal(err)
		}

		router.PathPrefix(*AppPath).Handler(http.StripPrefix(*AppPath, http.FileServer(http.Dir(*AppDir))))
		router.PathPrefix(*ModelsPath).Handler(http.StripPrefix(*ModelsPath, http.FileServer(http.Dir(*ModelsDir))))

		router.HandleFunc("/", ml5)

	case "data":

		if err := NC.Publish("api."+Environment, []byte(API["data"])); err != nil {
			log.Fatal(err)
		}
		router.HandleFunc("/", dataHandler)
	}
	log.Fatal(http.ListenAndServe(":"+*AppPort, router))

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	NC.Drain()
	log.Fatalf("Exiting")

}

func usage() {
	log.Printf("Usage: app [-name name] [-role role] [-port port] \n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
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

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	opts = append(opts, nats.ErrorHandler(natsErrHandler))
	return opts
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

// getEnv get key environment variable if exist otherwise return defalutValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func natsErrHandler(NC *nats.Conn, sub *nats.Subscription, natsErr error) {
	fmt.Printf("error: %v\n", natsErr)
	if natsErr == nats.ErrSlowConsumer {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			fmt.Printf("couldn't get pending messages: %v", err)
			return
		}
		fmt.Printf("Falling behind with %d pending messages on subject %q.\n",
			pendingMsgs, sub.Subject)
		// Log error, notify operations...
	}
	// check for other errors
}
