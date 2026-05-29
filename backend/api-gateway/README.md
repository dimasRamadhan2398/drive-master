# Kong API Gateway Configuration

This directory contains the configuration for Kong API Gateway in the Drive-Master application.

## Structure

```
api-gateway/
├── README.md                 # This file
├── setup-kong.sh            # Kong setup and management script
├── docker-compose.yml       # Docker compose configuration
├── kong/                    # Kong configuration directory
│   ├── kong.yml             # Main Kong configuration (DB-less mode)
│   └── kong.conf            # Kong instance configuration
└── pkg/
    └── config/              # Config files (legacy)
        └── config.yaml
```

## Quick Start

1. **Start Kong Gateway:**
   ```bash
   ./setup-kong.sh start
   ```

2. **Check Kong Health:**
   ```bash
   ./setup-kong.sh health
   ```

3. **View Kong Status:**
   ```bash
   ./setup-kong.sh status
   ```

4. **Open Kong Manager UI:**
   ```bash
   ./setup-kong.sh manager
   ```

## Configuration Details

### Services and Routes

The gateway is configured to route to three main microservices:

1. **User Service** (`/api/v1/users` or `/user-service`)
   - Upstream: `user-service:8001`
   - Features: Health checks, load balancing, retries

2. **Core Service** (`/api/v1/core` or `/core-service`)
   - Upstream: `core-service:8002`
   - Features: Health checks, load balancing, retries

3. **Booking Service** (`/api/v1/bookings` or `/booking-service`)
   - Upstream: `booking-service:8003`
   - Features: Health checks, load balancing, retries

### Features

- **Health Checks**: Active and passive health checking for all services
- **Load Balancing**: Round-robin algorithm with configurable weights
- **Retries**: Automatic retry mechanism for failed requests
- **CORS**: Configured for cross-origin requests
- **Rate Limiting**: IP-based rate limiting (1000/minute, 10000/hour, 100000/day)
- **Request ID**: Automatic request ID generation for tracing
- **Response Transforming**: Adds `X-Powered-By` header to responses

### Ports

- **Proxy (API)**: `8000`
- **Admin API**: `8001`
- **Manager UI**: `8002`

## API Endpoints

### Proxy Endpoints (for applications)

- User Service: `http://localhost:8000/api/v1/users`
- Core Service: `http://localhost:8000/api/v1/core`
- Booking Service: `http://localhost:8000/api/v1/bookings`

### Admin API (for management)

- Admin API: `http://localhost:8001`
- Manager UI: `http://localhost:8002`

## Management Script

The `setup-kong.sh` script provides convenient commands for managing Kong:

```bash
./setup-kong.sh [command]

Commands:
  health     - Check if Kong is running and healthy
  validate   - Validate Kong configuration file
  status     - Show Kong status and configuration
  test       - Test connectivity to microservices
  start      - Start Kong gateway
  stop       - Stop Kong gateway
  restart    - Restart Kong gateway
  logs       - Show Kong logs
  manager    - Open Kong Manager UI in browser
  reload     - Reload Kong configuration
```

## Configuration Validation

To validate the Kong configuration before deployment:

```bash
./setup-kong.sh validate
```

## Health Check Endpoints

Each microservice should implement a `/health` endpoint that returns HTTP 200 for Kong health checks to work properly.

## Monitoring

You can monitor Kong and its services through:

1. **Kong Manager UI**: `http://localhost:8002`
2. **Admin API**: `http://localhost:8001`
3. **Logs**: Use `./setup-kong.sh logs`
4. **Health Checks**: Use `./setup-kong.sh health`

## Troubleshooting

### Common Issues

1. **Services not reachable**: Ensure all microservices are running and accessible
2. **Health check failures**: Verify `/health` endpoints are implemented in services
3. **CORS errors**: Check CORS configuration in Kong settings
4. **Connection timeouts**: Adjust timeout values in Kong configuration

### Debug Commands

```bash
# Check Kong status
./setup-kong.sh status

# Test service connectivity
./setup-kong.sh test

# View logs
./setup-kong.sh logs

# Validate configuration
./setup-kong.sh validate
```

## Security Notes

- The Manager UI is exposed on port 8002 (consider restricting access in production)
- Rate limiting is configured but may need adjustment for production load
- CORS is set to allow all origins (`*`) - restrict to specific domains in production
- Authentication and authorization should be implemented at the service level