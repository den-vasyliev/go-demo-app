package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/CrowdSurge/banner"
)

func hash(decodedStr string) (string, string) {
	// defer metrics.MeasureSince([]string{"API"}, time.Now())
	//log.Print("DecodedStr: ", decodedStr)
	encodedStr := hex.EncodeToString([]byte(banner.PrintS(decodedStr)))
	//log.Print("EncodedStr: ", encodedStr)
	hashStr := fmt.Sprintf("%x", md5.Sum([]byte(encodedStr)))
	return hashStr, encodedStr
}
