FROM golang:1.10 as builder
ARG APP_BUILD_INFO
ARG APP_VERSION
ARG APP_ROLE
WORKDIR /go/src/app
COPY src/main.go app.go
RUN go get -d -v ./...
RUN go build -ldflags "-X main.AppRole=$APP_ROLE  -X main.BuildInfo=$APP_BUILD_INFO -X main.Version=$APP_VERSION" -v ./...
CMD ["pwd"]
CMD ["ls", "-l"]

FROM scratch
WORKDIR /
COPY --from=builder /go/src/app/app .
ENTRYPOINT ["/app"]
