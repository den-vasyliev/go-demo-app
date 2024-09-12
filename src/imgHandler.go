package main

import (
	"bytes"
	_ "image/jpeg"
	"io"
	"log"

	"github.com/den-vasyliev/image2ascii/convert"
	"github.com/valyala/fasthttp"
)

func img(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case "GET":
		if Environment == nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte("Environment variable not set"))
			return
		}
		ctx.Write(Environment)
	case "POST":
		var buf bytes.Buffer

		// Attempt to retrieve the file from the form
		file, err := ctx.FormFile("image")
		if err != nil {
			log.Printf("Error retrieving file: %v", err)
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Write([]byte("Failed to retrieve file"))
			return
		}

		ff, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		defer ff.Close() // Ensure the file gets closed

		// Copy the file data to the buffer
		if _, err := io.Copy(&buf, ff); err != nil {
			log.Printf("Error copying file: %v", err)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		if convertOptions, err := parseOptions(); err == nil {
			converter := convert.NewImageConverter()
			ctx.Write([]byte(converter.ImageBuf2ASCIIString(buf.Bytes(), convertOptions)))
		} else {
			log.Print("No options provided")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Write([]byte("No options provided"))
		}
	}
}
