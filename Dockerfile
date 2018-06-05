FROM golang:1.10
ARG APP_BUILD_INFO
ARG APP_VERSION
ARG APP_NAME
WORKDIR /go/src/app
COPY main.go app.go
RUN go get -d -v ./...
RUN go build -ldflags "-X main.AppName=$APP_NAME  -X main.BuildInfo=$APP_BUILD_INFO -X main.Version=$APP_VERSION" -v ./...

CMD ["/go/src/app/app"]