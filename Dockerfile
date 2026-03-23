# Build Stage
FROM golang:1.22-alpine AS builder

# Install git and certs (often needed for fetching private modules)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 1. Copy and download dependencies first (improves build speed)
COPY go.mod ./
# If you don't have a go.sum yet, comment out the next line
# COPY go.sum ./ 
RUN go mod download

# 2. Copy the source code
COPY . .

# 3. Build the binary
# We use '.' to build the package in the current directory
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /smoothtime .

# Final Stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /smoothtime /smoothtime

EXPOSE 123/udp 8080/tcp

ENTRYPOINT ["/smoothtime"]
