# Kong API Gateway Dockerfile
FROM kong:3.5.0

# Copy Kong declarative configuration
COPY kong/kong.yml /etc/kong/declarative/kong.yml

# Set environment variables for Kong
ENV KONG_DATABASE=off
ENV KONG_DECLARATIVE_CONFIG=/etc/kong/declarative/kong.yml
ENV KONG_PROXY_ACCESS_LOG=/dev/stdout
ENV KONG_ADMIN_ACCESS_LOG=/dev/stdout
ENV KONG_PROXY_ERROR_LOG=/dev/stderr
ENV KONG_ADMIN_ERROR_LOG=/dev/stderr
ENV KONG_ADMIN_LISTEN=0.0.0.0:8001
ENV KONG_ADMIN_GUI_URL=http://localhost:8002

# Create directory for config if it doesn't exist
RUN mkdir -p /etc/kong/declarative

# Expose ports
EXPOSE 8000 8001 8002

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD kong health

CMD ["kong", "start", "--vv"]
