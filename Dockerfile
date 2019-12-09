FROM denvasyliev/k8sdiy:builder as builder
ARG APP_VERSION
WORKDIR /go/src/app
COPY src/ .
RUN export GOPATH=/go
RUN CGO_ENABLED=0 GOOS=linux go build -o app -a -installsuffix cgo -ldflags "-X main.Version=$APP_VERSION" -v ./...

FROM scratch
WORKDIR /
COPY --from=builder /go/src/app/app .
ENTRYPOINT ["/app"]
