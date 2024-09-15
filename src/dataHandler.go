package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strconv"
	"time"
	//metrics "github.com/armon/go-metrics"
)

// DataHandler export broker msg func
func DataHandler(r *Req, i int) {
	var err error
	REQ0 = REQ0 + 1
	var Payload string = r.Hextr

	tokenStr := strconv.FormatUint(uint64(r.Token), 10)

	_, err = STMTIns.Exec(r.Token, r.Hextr)
	if err != nil {
		log.Print("InsertErr: %s", err)
	}

	sec, _ := time.ParseDuration(AppCacheExpire)

	err = CACHE.Set(tokenStr, Payload, sec).Err()

	if err != nil {
		log.Print(err)
	}

	err = NC.Publish(r.Reply, []byte(tokenStr))

	if err != nil {
		log.Print(err)
	}

}
