# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o lazylms-macos ./cmd/lazylms

# Final stage
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /build/lazylms-macos /usr/local/bin/lazylms-macos

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/lazylms-macos"]

# Default command (can be overridden)
CMD ["--help"]
