# --- Build Stage ---
# Drop the "1.22" and just use latest alpine to avoid version conflicts
FROM golang:alpine AS builder

WORKDIR /app

# Copy everything in the repository (including main.go and go.mod)
COPY . .

# Build the binary directly. 
# We can skip 'go mod download' because you have no external dependencies.
RUN CGO_ENABLED=0 GOOS=linux go build -o /smoothtime .

# --- Final Stage ---
FROM alpine:latest
COPY --from=builder /smoothtime /smoothtime

EXPOSE 123/udp 8080/tcp

ENTRYPOINT ["/smoothtime"]
