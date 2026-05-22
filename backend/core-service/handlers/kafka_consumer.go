package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"core-service/models"
	"core-service/services"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader       *kafka.Reader
	eventService *services.EventService
}

func NewKafkaConsumer(broker, topic, groupID string, eventService *services.EventService) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaConsumer{
		reader:       reader,
		eventService: eventService,
	}
}

func (k *KafkaConsumer) Consume(ctx context.Context) {
	for {
		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("kafka read error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var event models.UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("invalid event payload: %v", err)
			continue
		}

		if err := k.eventService.HandleUserCreated(ctx, event); err != nil {
			log.Printf("failed to process user.created: %v", err)
			continue
		}

		log.Printf("processed user.created event for userId=%d", event.UserID)
	}
}
