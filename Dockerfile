FROM golang:1.10
ARG APP_BUILD_INFO
ARG APP_VERSION
ARG APP_NAME
ARG APP_DB
ARG NEW_FEATURE
WORKDIR /go/src/app
COPY main.go app.go
RUN go get -d -v ./...
RUN go build -ldflags "-X main.AppName=$APP_NAME  -X main.BuildInfo=$APP_BUILD_INFO -X main.Version=$APP_VERSION -X main.NewFeature=${NEW_FEATURE} -X main.AppDb=$APP_DB" -v ./...

CMD ["/go/src/app/app"]