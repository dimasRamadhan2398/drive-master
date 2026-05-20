package repositories

import (
	"core-service/models"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) SaveProcessedEvent(eventType, payload string) error {
	return r.db.Create(&models.ProcessedEvent{
		EventType: eventType,
		Payload:   payload,
	}).Error
}
