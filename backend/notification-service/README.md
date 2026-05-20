# Notification Service

A microservice for handling email, WhatsApp, and newsletter notifications.

## Features

- **Email Notifications**: Mailtrap integration for sending emails
- **WhatsApp Reminders**: Send WhatsApp messages for session reminders 24 hours before
- **Promotional Updates**: Broadcast promotional content to users
- **Newsletter Subscription**: Manage newsletter subscriptions

## Structure

```
notification-service/
├── cmd/
│   └── main.go           # Entry point
├── models/
│   ├── notification.go   # Notification models
│   └── dto/
│       └── notification.go  # DTOs for API requests
├── repositories/
│   └── notification_repository.go  # Database operations
├── services/
│   ├── email_service.go       # Email sending logic
│   ├── whatsapp_service.go    # WhatsApp sending logic
│   └── notification_service.go  # Business logic
└── pkg/
    ├── config/
    │   └── config.go          # Configuration
    ├── errors/
    │   └── errors.go          # Error definitions
    ├── base/
    │   └── base_repository.go # Base repository
    ├── handlers/
    │   └── http_handler.go    # HTTP handlers
    └── kafka/
        └── consumer.go        # Kafka consumer
```

## API Endpoints

### Health
- `GET /health` - Service health check

### Newsletter
- `POST /newsletter/subscribe` - Subscribe to newsletter
- `POST /newsletter/unsubscribe` - Unsubscribe from newsletter
- `GET /newsletter/subscription/:email` - Get subscription by email
- `GET /newsletter/subscribers` - List active subscribers

### Notifications
- `POST /notifications/send` - Send immediate notification
- `POST /notifications/schedule` - Schedule a notification
- `GET /notifications/:id` - Get notification by ID
- `GET /notifications/user/:userId` - Get user notifications

### Preferences
- `GET /preferences/:userId` - Get user notification preferences
- `PUT /preferences` - Update preferences

### Admin
- `POST /admin/promotional/email` - Send promotional email campaign

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Service port | 8083 |
| POSTGRES_DSN | PostgreSQL connection string | localhost connection |
| MAILTRAP_API_TOKEN | Mailtrap API token | - |
| MAILTRAP_FROM_EMAIL | Sender email | noreply@example.com |
| MAILTRAP_FROM_NAME | Sender name | Drive Master |
| WHATSAPP_API_URL | WhatsApp API URL | - |
| WHATSAPP_API_TOKEN | WhatsApp API token | - |
| WHATSAPP_PHONE_NUMBER | WhatsApp phone number | - |
| KAFKA_BROKERS | Kafka broker addresses | localhost:9092 |

## Kafka Topics

- `session-upcoming` - Session reminder events (24 hours before)
- `booking-created` - New booking confirmation events
- `promotional` - Promotional campaign events

## Running the Service

```bash
# Install dependencies
go mod tidy

# Run the service
go run cmd/main.go
```