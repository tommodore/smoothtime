FROM golang:1.22-alpine AS builder

WORKDIR /app

# 1. Copy the only file we KNOW exists
COPY go.mod ./

# 2. Try to copy go.sum, but don't crash if it's missing
COPY go.sum* ./

RUN go mod download

# 3. Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /smoothtime .

FROM alpine:latest
COPY --from=builder /smoothtime /smoothtime
EXPOSE 123/udp 8080/tcp
ENTRYPOINT ["/smoothtime"]
