# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata dumb-init

WORKDIR /app

# Create non-root user for security
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Copy pre-compiled binary and config
COPY user-service .
RUN chmod +x user-service
COPY pkg/config/config.yaml ./pkg/config/config.yaml

# Set config path environment variable
ENV CONFIG_PATH=./pkg/config/config.yaml

# Switch to non-root user
USER appuser

EXPOSE 8001

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "if [ \"${RUN_MIGRATIONS}\" != \"false\" ]; then ./user-service migrate || true; fi && if [ \"${RUN_SEEDERS}\" = \"true\" ]; then ./user-service seed || true; fi && exec ./user-service serve"]
