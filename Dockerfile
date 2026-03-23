FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /smoothtime
FROM alpine
COPY --from=builder /smoothtime /smoothtime
EXPOSE 123/udp 8080
CMD ["/smoothtime", "-ntp-port", "123"]
LABEL org.opencontainers.image.source=https://github.com/tommodore/smoothtime
