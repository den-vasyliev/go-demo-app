package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"
)

// Define a struct for the incoming JSON request
type RequestBody struct {
	Text string `json:"text"`
}

// Helper function to handle errors
func handleError(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	ctx.SetStatusCode(statusCode)
	ctx.Write([]byte(message))
	log.Print(message)
}

// API handler processing POST requests
func api(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != http.MethodPost {
		handleGetRequest(ctx)
		return
	}

	contentType := string(ctx.Request.Header.Peek("Content-Type"))

	switch {
	case contentType == "application/json":
		handleJSONRequest(ctx)
	case strings.Contains(contentType, "multipart/form-data"):
		handleMultipartRequest(ctx)
	default:
		handleError(ctx, http.StatusUnsupportedMediaType, "Unsupported Media Type: "+contentType)
	}
}

func handleGetRequest(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte(Version))
}

func handleJSONRequest(ctx *fasthttp.RequestCtx) {
	var reqBody RequestBody
	if err := json.Unmarshal(ctx.PostBody(), &reqBody); err != nil {
		handleError(ctx, http.StatusBadRequest, "Error parsing JSON: "+err.Error())
		return
	}

	h := fnv.New32a()
	h.Write([]byte(reqBody.Text))
	token := h.Sum32()
	tokenStr := strconv.FormatUint(uint64(token), 10)
	cached, err := CACHE.Get(tokenStr).Result()
	if err == nil {
		if reply, err := hex.DecodeString(cached); err == nil {
			ctx.Write(reply)
			return
		}
	}

	uniqueReplyTo := nats.NewInbox()
	subscribeAndPublish(ctx, uniqueReplyTo, "ascii.json.banner", &Req{Token: token, Hextr: hex.EncodeToString([]byte(reqBody.Text)), Reply: uniqueReplyTo}, tokenStr)
}

func handleMultipartRequest(ctx *fasthttp.RequestCtx) {
	file, err := ctx.FormFile("image")
	if err != nil {
		handleError(ctx, http.StatusBadRequest, "Failed to retrieve file")
		return
	}

	ff, err := file.Open()
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Error opening file")
		return
	}
	defer ff.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, ff); err != nil {
		handleError(ctx, http.StatusInternalServerError, "Error copying file")
		return
	}

	token := rand.Uint32()
	if err := CACHE.Set(fmt.Sprintf("%d", token), buf.Bytes(), 10*time.Second).Err(); err != nil {
		handleError(ctx, http.StatusInternalServerError, "Error saving image to Redis: "+err.Error())
		return
	}

	uniqueReplyTo := nats.NewInbox()
	subscribeAndPublish(ctx, uniqueReplyTo, "img.json.image", &Req{Token: token, Hextr: "", Reply: uniqueReplyTo}, fmt.Sprintf("%d", token))
}

func subscribeAndPublish(ctx *fasthttp.RequestCtx, uniqueReplyTo, subject string, req *Req, tokenStr string) {
	sub, err := NC.SubscribeSync(uniqueReplyTo)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Error subscribing to uniqueReplyTo: "+uniqueReplyTo)
		return
	}

	if err := EC.Publish(subject, req); err != nil {
		handleError(ctx, http.StatusInternalServerError, "Error publishing to subject: "+subject)
		return
	}

	sec, _ := time.ParseDuration(*Wait)
	msg, err := sub.NextMsg(sec)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Error receiving message from uniqueReplyTo: "+uniqueReplyTo)
		return
	}

	reply, err := CACHE.Get(string(msg.Data)).Result()

	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Error retrieving cache")
		return
	}

	if strings.Contains(subject, "ascii") {
		if decodedStr, err := hex.DecodeString(reply); err == nil {
			reply = string(decodedStr)
		} else {
			handleError(ctx, http.StatusInternalServerError, "Error decoding hex string")
			return
		}
	}

	ctx.Write([]byte(reply))
}
