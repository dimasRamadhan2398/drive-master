package clients

import (
	"context"
	"encoding/json"
	"time"

	"user-service/models/dto"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *KafkaProducer) PublishUserCreated(event dto.CreateUserResponse) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Email),
		Value: payload,
	})
}
