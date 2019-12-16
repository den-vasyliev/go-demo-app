package main

import (
	//"crypto/md5"
	"encoding/hex"
	"hash/fnv"

	"github.com/CrowdSurge/banner"
)

func hash(decodedStr string) (uint32, string) {
	///defer metrics.MeasureSince([]string{"API"}, time.Now())
	//log.Print("DecodedStr: ", decodedStr)
	encodedStr := hex.EncodeToString([]byte(banner.PrintS(decodedStr)))
	//log.Print("EncodedStr: ", encodedStr)
	h := fnv.New32a()
	h.Write([]byte(encodedStr))
	//hashStr := fmt.Sprintf("%x", md5.Sum([]byte(encodedStr)))
	return h.Sum32(), encodedStr
}
