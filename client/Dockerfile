FROM golang:1.20 as builder
WORKDIR /app
COPY . /app
RUN GO111MODULE=auto CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -o app cmd/rest-server-server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE ${CLIENT_PORT}
CMD ./app --port ${CLIENT_PORT} --host '0.0.0.0'
