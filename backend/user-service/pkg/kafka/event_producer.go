package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// EventProducer is a thin wrapper around Producer to publish Event objects
// as structured log messages to Kafka. It uses Producer.SendLogSync for
// reliable delivery.
type EventProducer struct {
	producer *Producer
	service  string
}

// NewEventProducer creates a new EventProducer
func NewEventProducer(p *Producer, service string) *EventProducer {
	return &EventProducer{producer: p, service: service}
}

// Publish marshals the Event and sends it via the underlying Producer
func (e *EventProducer) Publish(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	// Ensure event metadata
	if event.ID == "" {
		event.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	fields := map[string]interface{}{
		"event":      event,
		"event_type": event.Type,
		"event_data": string(data),
	}

	// Use the global getLogLevel helper from events.go
	level := getLogLevel(event.Type)

	// Use the producer's synchronous send to ensure delivery before returning
	return e.producer.SendLogSync(ctx, level, string(event.Type), fields)
}

// PublishUserDeleted is a convenience method to publish user.deleted events
func (e *EventProducer) PublishUserDeleted(ctx context.Context, userID, username, email string) error {
	event := &Event{
		Type:      EventUserDeleted,
		UserID:    userID,
		Username:  username,
		Success:   true,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"email":      email,
			"deleted_at": time.Now().Format(time.RFC3339),
		},
	}
	return e.Publish(ctx, event)
}

func (e *EventProducer) PublishCourseCompleted(ctx context.Context, userID, username, email string) error {
	event := &Event{
		Type:      EventCourseCompleted,
		UserID:    userID,
		Username:  username,
		Success:   true,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"email":      email,
			"deleted_at": time.Now().Format(time.RFC3339),
		},
	}
	return e.Publish(ctx, event)
}
