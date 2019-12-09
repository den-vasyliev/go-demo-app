FROM golang:1.10 as builder
WORKDIR /go/src/app
COPY src/convert.diff src/convert.patch ./
COPY src/builder.dep ./main.go
RUN export GOPATH=/go
RUN go get -d -v ./...
RUN cd /go/src/github.com/qeesung/image2ascii/convert && git apply /go/src/app/convert.patch