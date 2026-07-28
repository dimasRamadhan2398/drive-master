# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install build dependencies (git only needed if you fetch dependencies, but vendor is used)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o payment-service ./cmd

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata dumb-init

WORKDIR /app

# Copy binary
COPY --from=builder /app/payment-service .

# Copy config file to the expected location
COPY --from=builder /app/pkg/config/config.yaml ./pkg/config/config.yaml

# Set config path environment variable
ENV CONFIG_PATH=./pkg/config/config.yaml

EXPOSE 8004

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
# If payment-service has its own migration/seed commands, adjust accordingly.
# For now, it directly serves.
CMD ["./payment-service", "serve"]