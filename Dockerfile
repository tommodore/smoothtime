# 1. Tell Docker to ALWAYS use the native speed of the GitHub runner (AMD64) for the builder
FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine3.24 AS builder

# 2. These variables are automatically injected by Buildx (e.g., linux and arm64)
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY . .

# 3. Cross-compile instantly using Go's native engine!
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /smoothtime .

# 4. Put the finished binary into the correct final container
FROM alpine:3.24.3
COPY --from=builder /smoothtime /smoothtime

EXPOSE 123/udp 8080/tcp

ENTRYPOINT ["/smoothtime"]
