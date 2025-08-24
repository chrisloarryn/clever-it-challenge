# syntax=docker/dockerfile:1

##
## Build stage
##
FROM golang:1.17-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o beer-api \
    cmd/main.go

##
## Runtime stage
##
FROM scratch

# Copy CA certificates for HTTPS calls
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/beer-api /beer-api

# Expose port
EXPOSE 8080

# Add health check (removed since we can't use complex shell commands in scratch)
# Health check would be handled by the orchestrator (k8s, docker-compose, etc.)

# Run as non-root user (note: scratch doesn't have users, so we rely on the app)
ENTRYPOINT ["/beer-api"]