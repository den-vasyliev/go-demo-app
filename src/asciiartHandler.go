package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"

	metrics "github.com/armon/go-metrics"
)

// ASCIIHandler msr broker
func ASCIIHandler(m *nats.Msg, i int) []byte {
	defer metrics.MeasureSince([]string{"API"}, time.Now())

	var t messageText

	json.Unmarshal(m.Data, &t)

	//log.Print("Text: ", t.Text)

	hashStr, encodedStr := hash(t.Text)

	cached, err := CACHE.Get(hashStr).Result()

	if err != nil {
		log.Print("Processing")

		sec, _ := time.ParseDuration(AppCacheExpire)

		CACHE.Set(hashStr, encodedStr, sec)

		CACHE.Set(fmt.Sprintf("%x", md5.Sum([]byte(t.Text))), hashStr, sec)

		//log.Print("Hash:", hashStr)

		msg, err := NC.Request(AppDatastore+".hash", []byte(fmt.Sprintf(`{"hash":"%s"}`, hashStr)), 2*time.Second)
		if err != nil {
			log.Printf("ErrRequest: %e", err)
		}

		var reply []byte

		if err != nil {
			log.Printf("ErrReply: %e", err)
			reply = []byte(fmt.Sprintf("{Reply:%s}", err))
		} else {
			reply = msg.Data
		}

		// Use the response
		//log.Printf("Reply: %s", reply)

		return reply
	}

	decoded, err := hex.DecodeString(cached)
	if err != nil {
		log.Print(err)
		return []byte("undef")
	}

	log.Print("Cached")
	return []byte(string(decoded))

}

func ascii(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())

	//var m messageText

	switch r.Method {

	case "GET":
		var b []byte
		b = append([]byte(""), Environment...)
		w.Write(b)
		/*
			case "POST":
				b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
				if err := json.Unmarshal(b, &m); err != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(422) // unprocessable entity
					if err := json.NewEncoder(w).Encode(err); err != nil {
						log.Print(err)
					}
				}*/
	}
}
