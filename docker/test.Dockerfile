FROM golang:1.23

WORKDIR /app

# Install curl for healthchecks
RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN go install gotest.tools/gotestsum@v1.12.0

COPY go.mod go.sum* ./

RUN go mod download

COPY . .