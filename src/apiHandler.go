package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash/fnv"
	_ "image/jpeg"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	//metrics "github.com/armon/go-metrics"
)

func perfHandler(w http.ResponseWriter, r *http.Request) {

	REQ0 = REQ0 + 1

	var reply []byte

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Print(err)
	}
	q := u.Query()

	h := fnv.New32a()

	h.Write([]byte(q.Get("text")))

	tokenStr := strconv.FormatUint(uint64(h.Sum32()), 10)

	token := h.Sum32()
	hexEncodedStr := hex.EncodeToString([]byte(q.Get("text")))

	cached, err := CACHE.Get(tokenStr).Result()

	if *Cache == "false" {
		err = errors.New("NoCache")
		cached = "636163686564"
	}

	if err == nil {
		reply, err = hex.DecodeString(cached)
		w.Write(reply)
	} else {
		// Create a unique subject name for replies.
		uniqueReplyTo := nats.NewInbox()

		// Listen for a single response
		sub, err := NC.SubscribeSync(uniqueReplyTo)
		if err != nil {
			log.Print(err)
		}

		// Send the request.
		// If processing is synchronous, use Request() which returns the response message.
		if err := EC.Publish("ascii.json.banner", &Req{Token: token, Hextr: hexEncodedStr, Reply: uniqueReplyTo}); err != nil {
			log.Print(err)
		}

		// Read the reply
		msg, err := sub.NextMsg(2 * time.Second)
		if err != nil {
			log.Print(err)
		}

		cached, err := CACHE.Get(string(msg.Data)).Result()

		reply, _ = hex.DecodeString(cached)
		w.Write(reply)

	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	var b []byte
	REQ0 = REQ0 + 1

	b = append([]byte(""), Environment...)
	w.Write(b)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	REQ0 = REQ0 + 1

	w.Write([]byte("Healthz: alive!"))
}

func api(w http.ResponseWriter, r *http.Request) {

	b, err := json.Marshal(APIReg)
	if err != nil {
		log.Print(err)
	}

	REQ0 = REQ0 + 1

	w.Write([]byte(b))
}
