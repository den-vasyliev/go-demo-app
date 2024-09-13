package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	//"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/armon/go-metrics"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"
)

// Req Define the object
type Req struct {
	Token uint32
	Hextr string
	Reply string
	Db    string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	INM = metrics.NewInmemSink(10*time.Second, time.Minute)
	sig := metrics.DefaultInmemSignal(INM)

	defer sig.Stop()

	API["ascii"] = "curl -XPOST --data '{text:TEXT}' HOST/ascii/"
	API["img"] = "curl -F 'image=@IMAGE' HOST/img/"
	API["ml5"] = "curl HOST/ml5/"
	API["data"] = "broker message api"

	initOptions()

	AppName := flag.String("name", "k8sdiy", "application name")
	AppRole := flag.String("role", "api", "app role: api data ascii img ml5 iot")
	AppPort := flag.String("port", "8080", "application port")
	//AppPath := flag.String("path", "/static/", "path to serve static files")
	//AppDir := flag.String("dir", "./ml5", "the directory of static files to host")
	//ModelsPath := flag.String("mpath", "/models/", "path to serve models files")
	//ModelsDir := flag.String("mdir", "./ml5/models", "the directory of models files to host")
	Cache = flag.String("cache", "true", "cache enable")
	Wait = flag.String("wait", "2s", "wait timeout")

	Urls = flag.String("server", nats.DefaultURL, "The nats server URLs (separated by comma)")
	//var userCreds = flag.String("creds", "", "User Credentials File")
	var showTime = flag.Bool("timestamp", false, "Display timestamps")
	//var queueName = flag.String("queGroupName", "K8S-NATS-Q", "Queue Group Name")
	var showHelp = flag.Bool("help", false, "Show help message")
	//var err error

	log.SetFlags(0)
	flag.Usage = usage

	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	// Environment app
	Role = *AppRole

	metrics.NewGlobal(metrics.DefaultConfig(Role), INM)

	// Perf

	REQ0 = 0.0
	REQ1 = 0.0
	//t0 := time.Now()

	go func() { // Daniel told me to write this handler this way.
		for {
			select {
			case <-time.After(time.Second * 1):
				//	ts := time.Since(t0)
				//	log.Println("[", Role, "] time: ", ts, " requests: ", REQ0, " rps: ", (REQ0-REQ1)/1, " throughput:", float64(REQ0)/ts.Seconds())
				REQ1 = REQ0
			}
		}
	}()

	Environment = fmt.Sprintf("%s-%s:%s", *AppName, Role, Version)
	log.Printf("Listening on [%s]: %s port: %s", *AppRole, Environment, *AppPort)

	// Connect Options.

	subj, subjJSON, i := Role+".*", Role+".json.*", 0
	//opts := []nats.Option{nats.Name(*AppRole + " on " + subj)}
	//opts = setupConnOptions(opts)

	// Connect to NATS
	var err error
	NC, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	if err := NC.LastError(); err != nil {
		log.Fatalf("Last error from NATS client: %v", err)
	}
	defer NC.Close()

	EC, err = nats.NewEncodedConn(NC, nats.JSON_ENCODER)
	defer EC.Close()

	// Subscribe
	if _, err = EC.Subscribe(subjJSON, func(r *Req) {

		REQ0 = REQ0 + 1

		i++
		if *AppRole == "ascii" {

			go AsciiHandler(r, i)

		} else if *AppRole == "data" {

			go DataHandler(r, i)
		}
	}); err != nil {
		log.Fatalf("Last error from NATS client: %v", err)
	}


	cache()
	router := func(ctx *fasthttp.RequestCtx) {
		switch *AppRole {
		case "api":
			api(ctx)
		case "img":
			img(ctx)
		default:
			log.Print(subj)
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.Write([]byte("200 - OK"))
		}
	}
	log.Fatal(fasthttp.ListenAndServe(":"+*AppPort, router))

	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	// the corresponding fasthttp code
	

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	//NC.Drain()
	log.Fatalf("Exiting")

}

func cache() {
	// Connect to cache
	CACHE = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", AppCache, AppCachePort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := CACHE.Ping().Result()
	if err != nil {
		log.Print(err)
	}
}



func mq() {
	subj, subjJSON, i := Role+".*", Role+".json.*", 0
	//opts := []nats.Option{nats.Name(*AppRole + " on " + subj)}
	//opts = setupConnOptions(opts)

	// Connect to NATS
	var err error
	NC, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	if err := NC.LastError(); err != nil {
		log.Fatalf("Last error from NATS client: %v", err)
	}
	defer NC.Close()

	EC, err = nats.NewEncodedConn(NC, nats.JSON_ENCODER)
	defer EC.Close()

	// Subscribe
	if _, err = EC.Subscribe(subjJSON, func(r *Req) {

		REQ0 = REQ0 + 1

		i++

		router := func(ctx *fasthttp.RequestCtx) {
			switch *AppRole {
			case "api":
				api(ctx)
			case "img":
				img(ctx)
			default:
				log.Print(subj)
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.Write([]byte("200 - OK"))
			}
		}
	
		log.Fatal(fasthttp.ListenAndServe(":"+*AppPort, router))

	}); err != nil {
		log.Fatalf("Last error from NATS client: %v", err)
	}


	//`log.Printf("Listening on [%s]: %s port: %s", subj, Environment, *AppPort)

}

// subscribe sets up a subscriber that listens for messages on the specified subject.
func subscribe(nc *nats.Conn, subject string) {
	// Subscribe to subject
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {
		log.Printf("Received a message: %s", string(m.Data))
		// Optionally acknowledge the message (if needed)
	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}

	// Keep the subscriber running
	select {} // Block forever
}

func db() {
	// Connect to db

	DB, err := sql.Open("mysql", AppDb)
	//	DB.SetMaxIdleConns(10000)
	if err != nil {
		log.Print(err)
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		log.Print(err) // proper error handling instead of panic in your app
	}

	STMTIns, err = DB.Prepare("insert into demo values(null,?,?)")
	STMTSel, err = DB.Prepare("SELECT text FROM demo WHERE token = ? limit 1")
	defer STMTIns.Close()
	defer STMTSel.Close()
}