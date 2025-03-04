FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o testserver ./tests/testserver

FROM alpine:3.18

RUN apk add --no-cache curl

COPY --from=builder /app/testserver /usr/local/bin/testserver

EXPOSE 8080

CMD ["testserver"]