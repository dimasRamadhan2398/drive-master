# Backend Microservices

This backend follows a clean-layered architecture for each service:

- `cmd/` for application entry points.
- `handlers/` for HTTP delivery and Kafka consumers.
- `services/` for use cases and business logic.
- `repositories/` for PostgreSQL and Redis access.
- `models/` for domain structs.

## Services

- `api-gateway` (`:8080`) routes Nuxt requests to internal services.
- `user-service` (`:8081`) creates users in PostgreSQL and publishes `user.created` events.
- `core-service` (`:8082`) consumes `user.created` events, persists processed event logs, and caches user payloads in Redis.

## Event Flow

1. `POST /api/users` hits API Gateway.
2. Gateway forwards to `user-service`.
3. `user-service` writes user into PostgreSQL and publishes to Kafka topic `user.created`.
4. `core-service` consumes `user.created`, stores an audit row in PostgreSQL, and caches user data in Redis.

## Run

From repository root:

```bash
docker compose up --build
```

Useful smoke tests:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/users/health
curl http://localhost:8080/api/core/health
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```
