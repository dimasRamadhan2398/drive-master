package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:120;not null"`
	Email     string    `json:"email" gorm:"size:190;uniqueIndex;not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateUserInput struct {
	Name  string `json:"name" binding:"required,min=2"`
	Email string `json:"email" binding:"required,email"`
}

type UserCreatedEvent struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}
