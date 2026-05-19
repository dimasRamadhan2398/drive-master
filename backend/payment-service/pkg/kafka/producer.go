package kafka

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Config struct {
	Brokers     []string
	Topic       string
	ServiceName string
	Enabled     bool
	UseAsync    bool
}

type Producer struct {
	producer    sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	topic       string
	serviceName string
	enabled     bool
	useAsync    bool
	logger      *zap.Logger
	wg          sync.WaitGroup
}

func NewProducer(cfg Config, logger *zap.Logger) (*Producer, error) {
	if !cfg.Enabled {
		logger.Info("Kafka producer is disabled")
		return &Producer{enabled: false, logger: logger}, nil
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Backoff = 100 * time.Millisecond

	var producer sarama.SyncProducer
	var asyncProducer sarama.AsyncProducer
	var err error

	if cfg.UseAsync {
		asyncProducer, err = sarama.NewAsyncProducer(cfg.Brokers, saramaConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create async producer: %w", err)
		}
	} else {
		producer, err = sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create sync producer: %w", err)
		}
	}

	logger.Info("Kafka producer initialized",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", cfg.Topic),
		zap.Bool("async", cfg.UseAsync),
	)

	return &Producer{
		producer:    producer,
		asyncProducer: asyncProducer,
		topic:       cfg.Topic,
		serviceName: cfg.ServiceName,
		enabled:     cfg.Enabled,
		useAsync:    cfg.UseAsync,
		logger:      logger,
	}, nil
}

func (p *Producer) SendMessage(key string, value interface{}) error {
	if !p.enabled {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
		Headers: []sarama.RecordHeader{
			{Key: []byte("service"), Value: []byte(p.serviceName)},
			{Key: []byte("timestamp"), Value: []byte(time.Now().Format(time.RFC3339))},
		},
	}

	if p.useAsync && p.asyncProducer != nil {
		p.asyncProducer.Input() <- msg
		return nil
	}

	if p.producer != nil {
		partition, offset, err := p.producer.SendMessage(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		p.logger.Debug("Message sent",
			zap.String("topic", p.topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
		)
	}

	return nil
}

func (p *Producer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	if p.asyncProducer != nil {
		return p.asyncProducer.Close()
	}
	return nil
}

type EventPublisherConfig struct {
	Producer *Producer
	Consumer *Consumer
	Topic    string
	Enabled  bool
}

type EventPublisher struct {
	producer *Producer
	consumer *Consumer
	topic    string
	enabled  bool
}

func NewEventPublisher(cfg EventPublisherConfig) *EventPublisher {
	return &EventPublisher{
		producer: cfg.Producer,
		consumer: cfg.Consumer,
		topic:    cfg.Topic,
		enabled:  cfg.Enabled,
	}
}

func (e *EventPublisher) Publish(key string, event interface{}) error {
	if !e.enabled {
		return nil
	}
	return e.producer.SendMessage(key, event)
}

func (e *EventPublisher) Subscribe(handler MessageHandler) {
	if e.consumer != nil {
		e.consumer.RegisterHandler(e.topic, handler)
	}
}

type IEventPublisher interface {
	Publish(key string, event interface{}) error
	Subscribe(handler MessageHandler)
}