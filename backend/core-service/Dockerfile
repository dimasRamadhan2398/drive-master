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
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o core-service ./cmd

# Runtime stage
FROM alpine:3.19

# Install ca-certificates for HTTPS and dumb-init for proper signal handling
RUN apk --no-cache add ca-certificates tzdata dumb-init

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/core-service .

# Copy config file
COPY --from=builder /app/pkg/config/config.yaml ./pkg/config/config.yaml

# Copy database refs (CSV files for seeders)
COPY --from=builder /app/database/refs ./database/refs

# Set environment variables
ENV CONFIG_PATH=./pkg/config/config.yaml
ENV GIN_MODE=release

EXPOSE 8002

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "if [ \"${RUN_MIGRATIONS}\" != \"false\" ]; then ./core-service migrate || true; fi && if [ \"${RUN_SEEDERS}\" = \"true\" ]; then ./core-service seed || true; fi && exec ./core-service serve"]
