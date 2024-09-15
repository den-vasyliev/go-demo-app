package main

import (
	"fmt"
	_ "image/jpeg"
	"log"
	"strconv"
	"time"

	"github.com/den-vasyliev/image2ascii/convert"
)

func ImgHandler(r *Req, i int) {
	sec, _ := time.ParseDuration(AppCacheExpire)

	img, err := CACHE.Get(fmt.Sprintf("%d", r.Token)).Result()

	if err != nil {
		log.Print("No image in cache", err)
	}

	if convertOptions, err := parseOptions(); err == nil {
		converter := convert.NewImageConverter()
		CACHE.Set(fmt.Sprintf("%d", r.Token), converter.ImageBuf2ASCIIString([]byte(img), convertOptions), sec)
		log.Print("Image converted")
	} else {
		log.Print("No options provided")
	}
	tokenStr := strconv.FormatUint(uint64(r.Token), 10)

	err = NC.Publish(r.Reply, []byte(tokenStr))

	if err != nil {
		log.Print(err)
	}
}
