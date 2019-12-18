package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	//metrics "github.com/armon/go-metrics"
)

//DataHandler export broker msg func
func DataHandler(r *Req, i int) {

	REQ0 = REQ0 + 1
	var err error
	var Payload string
	tokenStr := strconv.FormatUint(uint64(r.Token), 10)
	//json.Unmarshal(m.Data, &t)

	//	hexStr, err := CACHE.Get(t.Hash).Result()
	//	if err != nil {
	//		log.Print(err)
	//	}

	_, err = STMTIns.Exec(r.Token, r.Hextr)

	if err != nil {
		log.Print(err)
	}

	// additional iteration
	err = STMTSel.QueryRow(r.Token).Scan(&Payload) // WHERE number = 13

	if err != nil {
		log.Printf("QueryRowErr: %s", err) // proper error handling instead of panic in your app
	}

	//decoded, err = hex.DecodeString(Payload)

	//hex.Decode(decoded, []byte(Payload))
	sec, _ := time.ParseDuration(AppCacheExpire)

	err = CACHE.Set(tokenStr, Payload, sec).Err()

	if err != nil {
		log.Print(err)
	}

	NC.Publish(r.Reply, []byte(tokenStr))

	//return []byte()

}

func dataHandler(w http.ResponseWriter, r *http.Request) {

	var Payload string

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Print(err)
	}
	q := u.Query()

	stmt, err := DB.Prepare("insert into demo values(null,?,?)")

	_, err = stmt.Exec(q.Get("key"), q.Get("val"))

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	stmt, err = DB.Prepare("SELECT text FROM demo WHERE token = ? limit 1")

	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	// additional iteration
	_ = stmt.QueryRow(q.Get("key")).Scan(&Payload) // WHERE number = 13

	REQ0 = REQ0 + 1

	w.Write([]byte(fmt.Sprintf("%s", Payload)))
	/*
		case "POST":
			b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
			if err := json.Unmarshal(b, &m); err != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(422) // unprocessable entity
				if err := json.NewEncoder(w).Encode(err); err != nil {
					panic(err)
				}
			}
			w.Write([]byte(dataStore(m.Hash)))*/

}
