package main

import (
	//"crypto/md5"
	"encoding/hex"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"time"

	"github.com/CrowdSurge/banner"
	"github.com/armon/go-metrics"
	//metrics "github.com/armon/go-metrics"
)

// ASCIIHandler msr broker
func ASCIIHandler(r *Req, i int) {
	//var t messageText

	//json.Unmarshal(m.Data, &t)
	hexDecodedStr, _ := hex.DecodeString(r.Hextr)

	hexEncodedStr := hex.EncodeToString([]byte(banner.PrintS(string(hexDecodedStr))))

	if err := EC.Publish("data.json.hash", &Req{Token: r.Token, Hextr: hexEncodedStr, Reply: r.Reply}); err != nil {
		log.Print(err)
	}
	REQ0 = REQ0 + 1
	//return []byte("0")
	/*
		cached, err := CACHE.Get(strconv.FormatUint(uint64(hashStr), 10)).Result()

		if *Cache == "false" {
			err = errors.New("Processing")
			cached = "636163686564"
		}

		if err != nil {
			log.Print(err)

			sec, _ := time.ParseDuration(AppCacheExpire)

			CACHE.Set(strconv.FormatUint(uint64(hashStr), 10), encodedStr, sec)

			msg, err := NC.Request(AppDatastore+".hash", []byte(fmt.Sprintf(`{"hash":"%s"}`, strconv.FormatUint(uint64(hashStr), 10))), 2*time.Second)
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

		decoded, err := hex.DecodeString(s)
		if err != nil {
			log.Print(err)
			return []byte("undef")
		}

		return []byte(string(decoded))
	*/
}

func ascii(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())

	switch r.Method {

	case "GET":
		var b []byte
		b = append([]byte(""), Environment...)

		w.Write(b)

	}
}
