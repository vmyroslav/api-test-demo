FROM golang:1.23-alpine

RUN apk add --no-cache git

RUN go install github.com/go-swagger/go-swagger/cmd/swagger@latest

WORKDIR /app