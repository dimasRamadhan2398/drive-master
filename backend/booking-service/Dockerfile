# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./

COPY vendor ./vendor
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-w -s" -o booking-service ./cmd

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata dumb-init

WORKDIR /app

# Create non-root user for security
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Copy binary
COPY --from=builder /app/booking-service .
COPY --from=builder /app/pkg/config/config.yaml ./pkg/config/config.yaml

# Set config path environment variable
ENV CONFIG_PATH=./pkg/config/config.yaml

# Switch to non-root user
USER appuser

EXPOSE 8003

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "if [ \"${RUN_MIGRATIONS}\" != \"false\" ]; then ./booking-service migrate || true; fi && if [ \"${RUN_SEEDERS}\" = \"true\" ]; then ./booking-service seed || true; fi && exec ./booking-service serve"]