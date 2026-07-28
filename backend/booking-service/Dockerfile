# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o booking-service ./cmd

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata dumb-init

WORKDIR /app

# Create non-root user
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Copy binary
COPY --from=builder /app/booking-service .

# Copy config file
COPY --from=builder /app/pkg/config/config.yaml ./pkg/config/config.yaml

# Set config path
ENV CONFIG_PATH=./pkg/config/config.yaml

# Switch to non-root user
USER appuser

EXPOSE 8003

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "if [ \"${RUN_MIGRATIONS}\" != \"false\" ]; then ./booking-service migrate || true; fi && if [ \"${RUN_SEEDERS}\" = \"true\" ]; then ./booking-service seed || true; fi && exec ./booking-service serve"]