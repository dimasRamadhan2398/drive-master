package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FAQ represents a frequently asked question
type FAQ struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Question  string         `json:"question" gorm:"size:500;not null"`
	Answer    string         `json:"answer" gorm:"type:text;not null"`
	Order     int            `json:"order" gorm:"default:0"`
	Category  string         `json:"category" gorm:"size:100;default:'general'"`
	IsActive  bool           `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}