# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Copy vendor and source
COPY vendor ./vendor
COPY . .

# Build the binary with vendor mode
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-w -s" -o core-service ./cmd

# Runtime stage
FROM alpine:3.19

# Install ca-certificates, tzdata, and dumb-init for signal handling
RUN apk --no-cache add ca-certificates tzdata dumb-init

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/core-service .

# Copy config file
COPY --from=builder /app/pkg/config/config.yaml ./pkg/config/config.yaml

# Create analytics directory for credentials file
RUN mkdir -p /app/analytics

# Copy database reference files (CSV for seeders)
COPY --from=builder /app/database/refs ./database/refs

# Copy analytics folder (including ga_credentials.json) to /app/analytics
COPY --from=builder /app/app/analytics ./analytics

# Environment variables
ENV CONFIG_PATH=./pkg/config/config.yaml
ENV GIN_MODE=release

EXPOSE 8002

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "if [ \"${RUN_MIGRATIONS}\" != \"false\" ]; then ./core-service migrate up || true; fi && if [ \"${RUN_SEEDERS}\" = \"true\" ]; then ./core-service migrate seed || true; fi && exec ./core-service serve"]