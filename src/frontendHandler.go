package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"./image2ascii/convert"

	metrics "github.com/armon/go-metrics"
	"github.com/go-redis/redis"
)

var imageFilename string
var ratio float64
var fixedWidth int
var fixedHeight int
var fitScreen bool
var stretchedScreen bool
var colored bool
var reversed bool

var convertDefaultOptions = convert.DefaultOptions

func frontendHandler(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())
	switch r.Method {
	case "GET":

		log.Printf("Get Request: %s", r.URL.Query().Get("banner"))

		if r.URL.Query().Get("banner") == "" {
			log.Printf("No Banner Request - read index.html")
			dat, err := ioutil.ReadFile("/data/index.html")
			if err != nil {
				log.Printf("No found: index.html")
				w.Write([]byte(fmt.Sprintf("%s", "Not found")))
			}
			w.Write(dat)
		} else {

			message := fmt.Sprintf(`{"text":"%s"}`, r.URL.Query().Get("banner"))

			client := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", AppDbNoSQL, AppDbNoSQLPort),
				Password: "", // no password set
				DB:       0,  // use default DB
			})

			cacheItem, err := client.Get(fmt.Sprintf("%x", md5.Sum([]byte(message)))).Result()
			if err != nil {
				w.Write([]byte(fmt.Sprintf("%s", rest("http://"+AppBackend, message))))

			} else {
				hexStr, _ := client.Get(cacheItem).Result()
				decoded, _ := hex.DecodeString(hexStr)
				w.Write([]byte(decoded))
			}
			//w.Write([]byte(fmt.Sprintf("%s", rest("http://"+AppBackend, message))))
		}
	case "POST":
		var Buf bytes.Buffer
		//file, err := os.OpenFile("./downloaded.png", os.O_WRONLY|os.O_CREATE, 0666)
		//log.Print(err)
		f, _, _ := r.FormFile("image")
		defer f.Close()
		io.Copy(&Buf, f)
		b := Buf.Bytes()
		//io.Copy(file, f)

		Buf.Reset()
		w.Header().Set("Content-Type", "image/png")

		//img, _, err := image.Decode(bytes.NewReader(b))
		if convertOptions, err := parseOptions(); err == nil {
			converter := convert.NewImageConverter()
			img := converter.ImageFile2ASCIIString(b, convertOptions)
			//log.Print("Image: ", img)
			//w.Write(rest("http://"+AppDatastore, fmt.Sprintf(`{"hash":"%s"}`, img)))
			w.Header().Set("Content-Length", strconv.Itoa(len(img)))

			w.Write([]byte(converter.ImageFile2ASCIIString(b, convertOptions)))
			//w.Write([]byte("ok"))
		} else {
			log.Print("No opt")
		}

	}
}

func parseOptions() (*convert.Options, error) {

	// config  the options
	convertOptions := &convert.Options{
		Ratio:           ratio,
		FixedWidth:      fixedWidth,
		FixedHeight:     fixedHeight,
		FitScreen:       fitScreen,
		StretchedScreen: stretchedScreen,
		Colored:         colored,
		Reversed:        reversed,
	}
	return convertOptions, nil
}

func initOptions() {
	flag.StringVar(&imageFilename,
		"f",
		"",
		"Image filename to be convert")
	flag.Float64Var(&ratio,
		"r",
		convertDefaultOptions.Ratio,
		"Ratio to scale the image, ignored when use -w or -g")
	flag.IntVar(&fixedWidth,
		"w",
		convertDefaultOptions.FixedWidth,
		"Expected image width, -1 for image default width")
	flag.IntVar(&fixedHeight,
		"g",
		convertDefaultOptions.FixedHeight,
		"Expected image height, -1 for image default height")
	flag.BoolVar(&fitScreen,
		"s",
		convertDefaultOptions.FitScreen,
		"Fit the terminal screen, ignored when use -w, -g, -r")
	flag.BoolVar(&colored,
		"c",
		convertDefaultOptions.Colored,
		"Colored the ascii when output to the terminal")
	flag.BoolVar(&reversed,
		"i",
		convertDefaultOptions.Reversed,
		"Reversed the ascii when output to the terminal")
	flag.BoolVar(&stretchedScreen,
		"t",
		convertDefaultOptions.StretchedScreen,
		"Stretch the picture to overspread the screen")
}
