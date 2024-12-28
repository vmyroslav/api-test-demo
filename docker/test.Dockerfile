FROM golang:1.23

WORKDIR /app

# Install curl for healthchecks
RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN go install gotest.tools/gotestsum@v1.12.0

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .