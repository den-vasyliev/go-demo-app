package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"

	metrics "github.com/armon/go-metrics"
)

// ASCIIHandler msr broker
func ASCIIHandler(m *nats.Msg, i int) []byte {
	defer metrics.MeasureSince([]string{"API"}, time.Now())

	var t messageText

	json.Unmarshal(m.Data, &t)

	stmt, err := DB.Prepare("insert into demo values(null,?,?)")

	_, err = stmt.Exec(strconv.Itoa(i), strconv.Itoa(i))

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	log.Print("done: " + strconv.Itoa(i))
	return []byte(string("done: " + strconv.Itoa(i)))

	hashStr, encodedStr := hash(t.Text)

	cached, err := CACHE.Get(hashStr).Result()

	if *Cache == "false" {
		err = errors.New("Processing")
		cached = "636163686564"
	}

	if err != nil {
		log.Print(err)

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

		REQ0 = REQ0 + 1

		return reply
	}

	decoded, err := hex.DecodeString(cached)
	if err != nil {
		log.Print(err)
		return []byte("undef")
	}
	log.Print("done: " + strconv.Itoa(i))
	//return []byte(string("done: " + strconv.Itoa(i)))

	//log.Print(string(decoded))
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
