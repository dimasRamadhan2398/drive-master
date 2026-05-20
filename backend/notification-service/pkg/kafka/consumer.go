package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
)

// Consumer handles Kafka message consumption
type Consumer struct {
	readers []*kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers, groupID string, topics []string) *Consumer {
	brokerList := strings.Split(brokers, ",")
	var readers []*kafka.Reader

	for _, topic := range topics {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokerList,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
		readers = append(readers, reader)
	}

	return &Consumer{readers: readers}
}

// MessageHandler is the interface for handling Kafka messages
type MessageHandler interface {
	ProcessSessionEvent(ctx context.Context, eventType string, payload []byte) error
	ProcessPromotionalEvent(ctx context.Context, payload []byte) error
}

// Start begins consuming messages from Kafka
func (c *Consumer) Start(handler MessageHandler) {
	log.Println("Starting Kafka consumer...")
	var wg sync.WaitGroup

	for _, reader := range c.readers {
		wg.Add(1)
		go func(r *kafka.Reader) {
			defer wg.Done()
			for {
				ctx := context.Background()
				msg, err := r.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message: %v", err)
					continue
				}

				log.Printf("Received message on topic %s: %s", msg.Topic, string(msg.Value))

				var errHandler error
				switch msg.Topic {
				case "session-upcoming":
					errHandler = handler.ProcessSessionEvent(ctx, "session-upcoming", msg.Value)
				case "booking-created":
					errHandler = handler.ProcessSessionEvent(ctx, "booking-created", msg.Value)
				case "promotional":
					errHandler = handler.ProcessPromotionalEvent(ctx, msg.Value)
				default:
					log.Printf("Unknown topic: %s", msg.Topic)
					continue
				}

				if errHandler != nil {
					log.Printf("Error processing message: %v", errHandler)
				}
			}
		}(reader)
	}

	wg.Wait()
}

// Close closes all Kafka readers
func (c *Consumer) Close() error {
	var lastErr error
	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// EventPayload represents common event payload structure
type EventPayload struct {
	Type   string          `json:"type"`
	UserID uint            `json:"userId"`
	Email  string          `json:"email"`
	Phone  string          `json:"phone"`
	Data   json.RawMessage `json:"data"`
}