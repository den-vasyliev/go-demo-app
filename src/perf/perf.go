package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/nats-io/nats.go"

	metrics "github.com/armon/go-metrics"
	_ "github.com/go-sql-driver/mysql"
)

// STMT prepare
var STMT *sql.Stmt

// BumpMySQL some finc
func BumpMySQL(db *sql.DB, buf chan int, qq chan int) {
	defer metrics.MeasureSince([]string{"DB"}, time.Now())

	//var k int
	// Open doesn't open a connection. Validate DSN data:

	for {
		select {
		case i := <-buf:

			//var Payload string
			q := <-qq
			//q1 := fmt.Sprintf("insert into demo values(null,%s,'bencha')", strconv.Itoa(i))
			//	q2 := fmt.Sprintf("SELECT text FROM demo WHERE token = %s limit 1", strconv.Itoa(i))
			_, err := STMT.Exec(strconv.Itoa(i), strconv.Itoa(q))

			//_, err = db.Exec(q1)

			if err != nil {
				log.Print(err)
			}

			//	err = db.QueryRow(q2).Scan(&Payload)

			//	if err != nil {
			//		log.Print(err)
			//	}

		}
	}

}

var threads = flag.Int("threads", 8, "number of user threads")

func perf() {

	flag.Parse()
	log.Println("User threads:", *threads)

	runtime.GOMAXPROCS(runtime.NumCPU())

	var buffer = make(chan int, 1000) // 1000 just looks big enough for this case
	var qq = make(chan int, 50000)    // 1000 just looks big enough for this case

	maxtime := 300 * time.Second // time to run benchmark

	db, err := sql.Open("mysql", "root:@tcp(db)/demo")

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(1000)
	db.SetConnMaxLifetime(time.Second)

	opts := []nats.Option{nats.Name("bench")}
	opts = setupConnOptions(opts)

	// Connect to NATS

	nc, err := nats.Connect("nats://perf-nats-cluster", opts...)
	if err != nil {
		log.Fatal(err)
	}

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe
	req0 := 0
	req1 := 0
	k := 0

	STMT, err = db.Prepare("insert into demo values(null,?,?)")

	wg := sync.WaitGroup{}
	wg.Add(1000)

	if _, err := nc.Subscribe("bench", func(msg *nats.Msg) {
		k++
		//printMsg(msg, k)

		buffer <- rand.Intn(1000000) + 1

		go BumpMySQL(db, buffer, qq)

		req0 = req0 + 1

		qq <- req0

	}); err != nil {

		log.Print(err)
	}

	// Wait for a message to come in
	wg.Wait()
	/*
		for j := 0; j < *threads; j++ {
			go BumpMySQL(db, buffer, j)
		}
	*/

	t0 := time.Now()

	go func() { // Daniel told me to write this handler this way.
		timer := time.NewTimer(maxtime)
		for {
			select {
			case <-time.After(time.Second * 1):
				ts := time.Since(t0)
				log.Println("time: ", ts, " requests: ", req0, " rps: ", (req0-req1)/1, " throughput:", float64(req0)/ts.Seconds())
				req1 = req0
			case <-timer.C:
				log.Println("Finish!")
				ts := time.Since(t0)
				log.Println("Final result: ", ts, " requests: ", req0, " throughput:", float64(req0)/ts.Seconds())
				os.Exit(1) // this is not a quite nice way to exit, but I do not care.
			}
		}
	}()

	/*	for {

			buffer <- rand.Intn(1000000) + 1
			req0 = req0 + 1
		}
	*/
	//time.Sleep(100000 * time.Millisecond)

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
func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}
