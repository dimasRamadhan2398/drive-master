package models

import "time"

type UserCreatedEvent struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type ProcessedEvent struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EventType string    `json:"eventType" gorm:"size:80;index;not null"`
	Payload   string    `json:"payload" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"createdAt"`
}
