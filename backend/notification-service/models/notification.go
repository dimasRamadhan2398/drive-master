package models

import (
	"time"
)

// NotificationType defines the type of notification
type NotificationType string

const (
	NotificationTypeEmail     NotificationType = "email"
	NotificationTypeWhatsApp NotificationType = "whatsapp"
	NotificationTypeSMS      NotificationType = "sms"
)

// NotificationStatus defines the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusFailed   NotificationStatus = "failed"
	NotificationStatusCancelled NotificationStatus = "cancelled"
)

// NotificationCategory defines the category of notification
type NotificationCategory string

const (
	CategorySessionReminder NotificationCategory = "session_reminder"
	CategoryBookingConfirm   NotificationCategory = "booking_confirmation"
	CategoryPromotion       NotificationCategory = "promotional"
	CategoryNewsletter      NotificationCategory = "newsletter"
	CategoryWhatsAppReminder NotificationCategory = "whatsapp_reminder"
)

// Notification represents a notification record
type Notification struct {
	ID           uint                `json:"id" gorm:"primaryKey"`
	UserID       uint                `json:"userId" gorm:"index;not null"`
	Email        string              `json:"email" gorm:"size:255"`
	PhoneNumber  string              `json:"phoneNumber" gorm:"size:20"`
	Type         NotificationType    `json:"type" gorm:"size:20;not null"`
	Category     NotificationCategory `json:"category" gorm:"size:50;not null"`
	Status       NotificationStatus  `json:"status" gorm:"size:20;default:'pending'"`
	Subject      string              `json:"subject" gorm:"size:255"`
	Content      string              `json:"content" gorm:"type:text"`
	ScheduledAt  *time.Time          `json:"scheduledAt" gorm:"index"`
	SentAt       *time.Time          `json:"sentAt"`
	FailedAt     *time.Time          `json:"failedAt"`
	ErrorMessage string              `json:"errorMessage" gorm:"type:text"`
	RetryCount   int                 `json:"retryCount" gorm:"default:0"`
	Metadata     string              `json:"metadata" gorm:"type:jsonb"` // JSON metadata for additional data
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
}

// NotificationTemplate represents a reusable notification template
type NotificationTemplate struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;uniqueIndex;not null"`
	Type      NotificationType `json:"type" gorm:"size:20;not null"`
	Category  NotificationCategory `json:"category" gorm:"size:50;not null"`
	Subject   string    `json:"subject" gorm:"size:255"`
	Content   string    `json:"content" gorm:"type:text"`
	Variables []string  `json:"variables" gorm:"type:text[]"` // Array of variable names
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewsletterSubscription represents a newsletter subscription
type NewsletterSubscription struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"size:100"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	SubscribedAt time.Time `json:"subscribedAt"`
	UnsubscribedAt *time.Time `json:"unsubscribedAt"`
	Source    string    `json:"source" gorm:"size:50"` // website, app, etc.
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NotificationPreference stores user notification preferences
type NotificationPreference struct {
	ID                    uint      `json:"id" gorm:"primaryKey"`
	UserID                uint      `json:"userId" gorm:"uniqueIndex;not null"`
	EmailEnabled          bool      `json:"emailEnabled" gorm:"default:true"`
	WhatsAppEnabled       bool      `json:"whatsAppEnabled" gorm:"default:true"`
	SessionReminderHours  int       `json:"sessionReminderHours" gorm:"default:24"` // Hours before session to send reminder
	PromotionEnabled      bool      `json:"promotionEnabled" gorm:"default:true"`
	NewsletterEnabled     bool      `json:"newsletterEnabled" gorm:"default:true"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// SessionInfo contains session details for reminders
type SessionInfo struct {
	SessionID     uint      `json:"sessionId"`
	UserID        uint      `json:"userId"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phoneNumber"`
	InstructorName string   `json:"instructorName"`
	StudentName   string    `json:"studentName"`
	SessionDate   time.Time `json:"sessionDate"`
	Duration      int       `json:"duration"` // in minutes
	LessonType    string    `json:"lessonType"`
	Location      string    `json:"location"`
}