package repositories

import (
	"context"

	"core-service/models"

	"gorm.io/gorm"
)

type IEventRepository interface {
	SaveProcessedEvent(ctx context.Context, eventType, payload string) error
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) SaveProcessedEvent(ctx context.Context, eventType, payload string) error {
	return r.db.WithContext(ctx).Create(&models.ProcessedEvent{
		EventType: eventType,
		Payload:   payload,
	}).Error
}