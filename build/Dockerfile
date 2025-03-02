FROM golang:1.24-alpine

RUN go install github.com/kokardy/http-echo-server/cmd/http-echo-server@latest

ENTRYPOINT ["http-echo-server"]
