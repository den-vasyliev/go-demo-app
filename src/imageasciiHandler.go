package main

import (
	"bytes"
	_ "image/jpeg"
	"io"
	"log"
	"net/http"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/qeesung/image2ascii/convert"
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

func img(w http.ResponseWriter, r *http.Request) {
	defer metrics.MeasureSince([]string{"API"}, time.Now())
	switch r.Method {
	case "GET":
		var b []byte
		b = append([]byte(""), Environment...)
		w.Write(b)
	case "POST":
		var Buf bytes.Buffer

		f, _, _ := r.FormFile("image")
		defer f.Close()
		io.Copy(&Buf, f)
		b := Buf.Bytes()

		Buf.Reset()
		w.Header().Set("Content-Type", "text/plain")

		if convertOptions, err := parseOptions(); err == nil {
			converter := convert.NewImageConverter()

			w.Write([]byte(converter.ImageFile2ASCIIString(b, convertOptions)))
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
