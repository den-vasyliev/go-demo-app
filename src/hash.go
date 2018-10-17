package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/CrowdSurge/banner"
	"github.com/go-redis/redis"
)

func hash(decodedStr string) string {
	// defer metrics.MeasureSince([]string{"API"}, time.Now())
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	log.Print("DecodedStr: ", decodedStr)
	encodedStr := hex.EncodeToString([]byte(banner.PrintS(decodedStr)))
	log.Print("EncodedStr: ", encodedStr)
	hashStr := fmt.Sprintf("%x", md5.Sum([]byte(encodedStr)))
	client.Set(hashStr, encodedStr, 0)
	client.Set(fmt.Sprintf("%x", md5.Sum([]byte(decodedStr))), hashStr, 0)
	return hashStr
}
