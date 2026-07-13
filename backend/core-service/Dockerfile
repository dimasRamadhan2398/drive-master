# Runtime stage
FROM alpine:3.19

# Install ca-certificates, tzdata, and dumb-init for signal handling
RUN apk --no-cache add ca-certificates tzdata dumb-init

WORKDIR /app

# Copy pre-compiled binary
COPY core-service .
RUN chmod +x core-service

# Copy config file
COPY pkg/config/config.yaml ./pkg/config/config.yaml

# Create analytics directory for credentials file
RUN mkdir -p /app/analytics

# Copy database reference files (CSV for seeders)
COPY database/refs ./database/refs

# Copy analytics folder (including ga_credentials.json) to /app/analytics
COPY app/analytics ./analytics

# Environment variables
ENV CONFIG_PATH=./pkg/config/config.yaml
ENV GIN_MODE=release

EXPOSE 8002

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "if [ \"${RUN_MIGRATIONS}\" != \"false\" ]; then ./core-service migrate up || true; fi && if [ \"${RUN_SEEDERS}\" = \"true\" ]; then ./core-service migrate seed || true; fi && exec ./core-service serve"]