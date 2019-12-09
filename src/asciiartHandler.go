package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"

	metrics "github.com/armon/go-metrics"
	"github.com/go-redis/redis"
)

func ascii(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())

	var m messageText

	switch r.Method {

	case "GET":
		var b []byte
		b = append([]byte(""), Environment...)
		w.Write(b)

	case "POST":
		b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err := json.Unmarshal(b, &m); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Print(err)
			}
		}

		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", AppCache, AppCachePort),
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		log.Print("Text: ", m.Text)

		hashStr, encodedStr := hash(m.Text)

		cached, err := client.Get(hashStr).Result()

		if err != nil {
			log.Print(err)
			sec, _ := time.ParseDuration(AppCacheExpire)
			client.Set(hashStr, encodedStr, sec)
			client.Set(fmt.Sprintf("%x", md5.Sum([]byte(m.Text))), hashStr, sec)

			log.Print("Hash:", hashStr)
			// message brocker placeholder
			// Create a unique subject name for replies.
			uniqueReplyTo := nats.NewInbox()

			// Listen for a single response
			sub, err := NC.SubscribeSync(uniqueReplyTo)
			if err != nil {
				log.Print(err)
			}
			// Send the request.
			// If processing is synchronous, use Request() which returns the response message.
			if err := NC.PublishRequest(AppDatastore+".hash", uniqueReplyTo, []byte(fmt.Sprintf(`{"hash":"%s"}`, hashStr))); err != nil {
				log.Print(err)
			}

			// Read the reply
			msg, err := sub.NextMsg(time.Second)
			var reply []byte
			if err != nil {
				log.Print(err)
				reply = []byte(fmt.Sprintf("{%s}", err))
			} else {
				reply = msg.Data
			}

			// Use the response
			log.Printf("Reply: %s", reply)
			w.Write(reply)
			//w.Write(rest("http://"+AppDatastore, fmt.Sprintf(`{"hash":"%s"}`, hashStr)))
		} else {
			decoded, err := hex.DecodeString(cached)
			if err != nil {
				log.Print(err)
				w.Write([]byte("undef"))
			} else {
				log.Print("Cached")
				w.Write([]byte(string(decoded)))
			}

		}
	}
}
