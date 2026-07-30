package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"core-service/models"
	pkgkafka "core-service/pkg/kafka"
	"core-service/services"

	kafkago "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	readers      []*kafkago.Reader
	eventService *services.EventService
	stopCh       chan struct{}
}

func NewKafkaConsumer(brokers []string, topics []string, groupID string, eventService services.IEventService) *KafkaConsumer {
	readers := make([]*kafkago.Reader, 0, len(topics))
	for _, topic := range topics {
		reader := kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10e6, // 10MB
		})
		readers = append(readers, reader)
	}

	return &KafkaConsumer{
		readers:      readers,
		eventService: eventService.(*services.EventService),
		stopCh:       make(chan struct{}),
	}
}

func (k *KafkaConsumer) Consume(ctx context.Context) {
	for _, reader := range k.readers {
		go k.consumeReader(ctx, reader)
	}

	// Wait for stop signal
	<-k.stopCh
}

func (k *KafkaConsumer) consumeReader(ctx context.Context, reader *kafkago.Reader) {
	for {
		select {
		case <-k.stopCh:
			return
		case <-ctx.Done():
			return
		default:
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("[KafkaConsumer] read error from topic %s: %v", reader.Config().Topic, err)
			time.Sleep(2 * time.Second)
			continue
		}

		k.handleMessage(ctx, msg)
	}
}

func (k *KafkaConsumer) handleMessage(ctx context.Context, msg kafkago.Message) {
	// Extract event type from header first
	eventType := extractEventType(msg)

	// Fallback to key if header not present
	if eventType == "" {
		eventType = string(msg.Key)
	}

	if eventType == "" {
		log.Printf("[KafkaConsumer] missing event type, skipping message")
		return
	}

	log.Printf("[KafkaConsumer] Processing event type=%s from topic=%s", eventType, msg.Topic)

	var err error
	switch eventType {
	// User events
	case pkgkafka.EventUserCreated:
		var event models.UserCreatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleUserCreated(ctx, event)
		}
	case pkgkafka.EventUserUpdated:
		var event models.UserUpdatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleUserUpdated(ctx, event)
		}
	case pkgkafka.EventUserDeleted:
		var event models.UserDeletedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleUserDeleted(ctx, event)
		}

	// Car events
	case pkgkafka.EventCarCreated:
		var event models.CarCreatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleCarCreated(ctx, event)
		}
	case pkgkafka.EventCarUpdated:
		var event models.CarUpdatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleCarUpdated(ctx, event)
		}
	case pkgkafka.EventCarDeleted:
		var event models.CarDeletedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleCarDeleted(ctx, event)
		}

	// Package events
	case pkgkafka.EventPackageCreated:
		var event models.PackageCreatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandlePackageCreated(ctx, event)
		}
	case pkgkafka.EventPackageUpdated:
		var event models.PackageUpdatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandlePackageUpdated(ctx, event)
		}
	case pkgkafka.EventPackageDeleted:
		var event models.PackageDeletedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandlePackageDeleted(ctx, event)
		}
	case pkgkafka.EventEnrollmentPaid:
		var event models.EnrollmentPaidEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleEnrollmentPaid(ctx, event)
		}

	// Article events
	case pkgkafka.EventArticleCreated:
		var event models.ArticleCreatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleArticleCreated(ctx, event)
		}
	case pkgkafka.EventArticleUpdated:
		var event models.ArticleUpdatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleArticleUpdated(ctx, event)
		}
	case pkgkafka.EventArticleDeleted:
		var event models.ArticleDeletedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleArticleDeleted(ctx, event)
		}
	case pkgkafka.EventArticlePublished:
		var event models.ArticlePublishedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleArticlePublished(ctx, event)
		}
	case pkgkafka.EventArticleArchived:
		var event models.ArticleArchivedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleArticleArchived(ctx, event)
		}

	// Region events
	case pkgkafka.EventRegionProvinceUpdated:
		var event models.RegionProvinceUpdatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleRegionProvinceUpdated(ctx, event)
		}
	case pkgkafka.EventRegionRegencyUpdated:
		var event models.RegionRegencyUpdatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleRegionRegencyUpdated(ctx, event)
		}
	case pkgkafka.EventRegionDistrictUpdated:
		var event models.RegionDistrictUpdatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleRegionDistrictUpdated(ctx, event)
		}

	// System events
	case pkgkafka.EventCacheInvalidated:
		var event models.CacheInvalidatedEvent
		if err = json.Unmarshal(msg.Value, &event); err == nil {
			err = k.eventService.HandleCacheInvalidated(ctx, event)
		}

	default:
		log.Printf("[KafkaConsumer] unknown event type: %s", eventType)
		return
	}

	if err != nil {
		log.Printf("[KafkaConsumer] failed to process %s: %v", eventType, err)
	} else {
		log.Printf("[KafkaConsumer] successfully processed %s", eventType)
	}
}

// extractEventType extracts the event type from Kafka message headers
func extractEventType(msg kafkago.Message) string {
	for _, header := range msg.Headers {
		if strings.ToLower(header.Key) == "event_type" {
			return string(header.Value)
		}
	}
	return ""
}

// Stop gracefully stops the consumer
func (k *KafkaConsumer) Stop() {
	close(k.stopCh)
	for _, reader := range k.readers {
		if err := reader.Close(); err != nil {
			log.Printf("[KafkaConsumer] error closing reader: %v", err)
		}
	}
}

// Legacy single-topic consumer for backwards compatibility
type LegacyKafkaConsumer struct {
	reader       *kafkago.Reader
	eventService *services.EventService
}

func NewLegacyKafkaConsumer(broker, topic, groupID string, eventService services.IEventService) *LegacyKafkaConsumer {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
	return &LegacyKafkaConsumer{
		reader:       reader,
		eventService: eventService.(*services.EventService),
	}
}

func (k *LegacyKafkaConsumer) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
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

func (k *LegacyKafkaConsumer) Close() error {
	return k.reader.Close()
}

// Config holds Kafka consumer configuration
type ConsumerConfig struct {
	Brokers []string
	Topics  []string
	GroupID string
}

// Validate checks if the configuration is valid
func (c *ConsumerConfig) Validate() error {
	if len(c.Brokers) == 0 {
		return fmt.Errorf("at least one broker is required")
	}
	if c.GroupID == "" {
		return fmt.Errorf("group ID is required")
	}
	if len(c.Topics) == 0 {
		return fmt.Errorf("at least one topic is required")
	}
	return nil
}
