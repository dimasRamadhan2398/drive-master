# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with CGO disabled for alpine
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api-gateway ./cmd/main.go

# Runtime stage
FROM alpine:3.19

# Install ca-certificates for HTTPS and dumb-init for proper signal handling
RUN apk --no-cache add ca-certificates tzdata dumb-init

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api-gateway .

# Copy config file (adjust path based on your actual structure)
COPY --from=builder /app/pkg/config/config.yaml ./pkg/config/config.yaml

# Alternative if config is at root level:
# COPY --from=builder /app/config.yaml ./config.yaml

# Set environment variables
ENV CONFIG_PATH=./pkg/config/config.yaml
ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "if [ -n \"${API_GATEWAY_CONFIG}\" ]; then export CONFIG_PATH=\"${API_GATEWAY_CONFIG}\"; fi && exec ./api-gateway"]
