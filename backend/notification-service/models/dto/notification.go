package dto

import (
	"time"

	"notification-service/models"
)

// SubscribeNewsletterRequest represents a newsletter subscription request
type SubscribeNewsletterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name"`
	Source string `json:"source"` // website, app, etc.
}

// UnsubscribeNewsletterRequest represents a newsletter unsubscription request
type UnsubscribeNewsletterRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// SendNotificationRequest represents a request to send a notification
type SendNotificationRequest struct {
	Type       string   `json:"type" binding:"required"` // email, whatsapp
	Recipients []string `json:"recipients" binding:"required"`
	Subject    string   `json:"subject"`
	Content    string   `json:"content" binding:"required"`
}

// ScheduleNotificationRequest represents a request to schedule a notification
type ScheduleNotificationRequest struct {
	UserID       uint   `json:"userId" binding:"required"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	Type         string `json:"type" binding:"required"` // email, whatsapp
	Category     string `json:"category"`
	Subject      string `json:"subject"`
	Content      string `json:"content" binding:"required"`
	ScheduledAt  string `json:"scheduledAt"` // RFC3339 format
	Metadata     string `json:"metadata"`
}

// ToNotification converts ScheduleNotificationRequest to a Notification model
func (r *ScheduleNotificationRequest) ToNotification() *models.Notification {
	notification := &models.Notification{
		UserID:   r.UserID,
		Email:    r.Email,
		PhoneNumber: r.PhoneNumber,
		Type:     models.NotificationType(r.Type),
		Category: models.NotificationCategory(r.Category),
		Subject:  r.Subject,
		Content:  r.Content,
		Metadata: r.Metadata,
	}

	if r.ScheduledAt != "" {
		t, err := time.Parse(time.RFC3339, r.ScheduledAt)
		if err == nil {
			notification.ScheduledAt = &t
		}
	}

	return notification
}

// UpdatePreferencesRequest represents a request to update notification preferences
type UpdatePreferencesRequest struct {
	UserID               uint `json:"userId" binding:"required"`
	EmailEnabled         bool `json:"emailEnabled"`
	WhatsAppEnabled      bool `json:"whatsAppEnabled"`
	SessionReminderHours int  `json:"sessionReminderHours"`
	PromotionEnabled     bool `json:"promotionEnabled"`
	NewsletterEnabled    bool `json:"newsletterEnabled"`
}

// ToPreference converts UpdatePreferencesRequest to a NotificationPreference model
func (r *UpdatePreferencesRequest) ToPreference() *models.NotificationPreference {
	return &models.NotificationPreference{
		UserID:               r.UserID,
		EmailEnabled:         r.EmailEnabled,
		WhatsAppEnabled:      r.WhatsAppEnabled,
		SessionReminderHours: r.SessionReminderHours,
		PromotionEnabled:     r.PromotionEnabled,
		NewsletterEnabled:    r.NewsletterEnabled,
	}
}

// PromotionalEmailRequest represents a promotional email request
type PromotionalEmailRequest struct {
	Recipients []string `json:"recipients" binding:"required"`
	Subject    string   `json:"subject" binding:"required"`
	Content    string   `json:"content" binding:"required"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
	Status    string    `json:"status"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"createdAt"`
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Data   interface{} `json:"data"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}