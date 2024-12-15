FROM golang:1.23-alpine

RUN apk add --no-cache git

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

WORKDIR /app
ENTRYPOINT ["oapi-codegen"]