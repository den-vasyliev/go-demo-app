FROM golang:1.10 as builder
ARG APP_VERSION
WORKDIR /go/src/app
COPY src/ .
RUN export GOPATH=/go
RUN go get -d -v ./...
RUN cd /go/src/github.com/qeesung/image2ascii/convert && git apply /go/src/app/convert.patch
RUN CGO_ENABLED=0 GOOS=linux go build -o app -a -installsuffix cgo -ldflags "-X main.Version=$APP_VERSION" -v ./...

FROM scratch
WORKDIR /
COPY --from=builder /go/src/app/app .
COPY --from=builder /go/src/app/ml5 /ml5
ENTRYPOINT ["/app"]
