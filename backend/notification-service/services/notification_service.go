package services

import (
	"context"
	"encoding/json"
	"time"

	"notification-service/models"
	"notification-service/repositories"
	apperrors "notification-service/pkg/errors"
)

// NotificationService handles notification business logic
type NotificationService struct {
	repo           repositories.INotificationRepository
	emailService   *EmailService
	whatsAppService *WhatsAppService
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	repo repositories.INotificationRepository,
	emailService *EmailService,
	whatsAppService *WhatsAppService,
) *NotificationService {
	return &NotificationService{
		repo:           repo,
		emailService:   emailService,
		whatsAppService: whatsAppService,
	}
}

// SendSessionReminder sends session reminder notifications (email and WhatsApp)
func (s *NotificationService) SendSessionReminder(ctx context.Context, session models.SessionInfo) error {
	var lastErr error

	// Send email notification
	if session.Email != "" {
		dateTimeStr := session.SessionDate.Format("Monday, January 2, 2006 at 3:04 PM")
		if err := s.emailService.SendSessionReminderEmail(
			ctx,
			session.Email,
			session.StudentName,
			session.InstructorName,
			dateTimeStr,
			session.LessonType,
			session.Location,
		); err != nil {
			lastErr = err
		}
	}

	// Send WhatsApp notification
	if session.PhoneNumber != "" {
		dateTimeStr := session.SessionDate.Format("Monday 3:04 PM")
		if err := s.whatsAppService.SendSessionReminderWhatsApp(
			ctx,
			session.PhoneNumber,
			session.StudentName,
			session.InstructorName,
			dateTimeStr,
			session.LessonType,
		); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// SendBookingConfirmation sends booking confirmation notifications
func (s *NotificationService) SendBookingConfirmation(ctx context.Context, session models.SessionInfo) error {
	var lastErr error

	// Send email notification
	if session.Email != "" {
		dateTimeStr := session.SessionDate.Format("Monday, January 2, 2006 at 3:04 PM")
		if err := s.emailService.SendBookingConfirmationEmail(
			ctx,
			session.Email,
			session.StudentName,
			session.InstructorName,
			dateTimeStr,
			session.LessonType,
		); err != nil {
			lastErr = err
		}
	}

	// Send WhatsApp notification
	if session.PhoneNumber != "" {
		dateTimeStr := session.SessionDate.Format("Monday 3:04 PM")
		if err := s.whatsAppService.SendBookingConfirmationWhatsApp(
			ctx,
			session.PhoneNumber,
			session.StudentName,
			session.InstructorName,
			dateTimeStr,
			session.LessonType,
		); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// SendPromotionalEmail sends promotional email to multiple recipients
func (s *NotificationService) SendPromotionalEmail(ctx context.Context, toEmails []string, subject, content string) error {
	return s.emailService.SendPromotionalEmail(ctx, toEmails, subject, content)
}

// SendPromotionalWhatsApp sends promotional WhatsApp message
func (s *NotificationService) SendPromotionalWhatsApp(ctx context.Context, to, title, content string) error {
	return s.whatsAppService.SendPromotionalWhatsApp(ctx, to, title, content)
}

// SubscribeNewsletter subscribes an email to the newsletter
func (s *NotificationService) SubscribeNewsletter(ctx context.Context, email, name, source string) (*models.NewsletterSubscription, error) {
	return s.repo.Subscribe(email, name, source)
}

// UnsubscribeNewsletter unsubscribes an email from the newsletter
func (s *NotificationService) UnsubscribeNewsletter(ctx context.Context, email string) error {
	return s.repo.Unsubscribe(email)
}

// GetSubscriptionByEmail retrieves subscription by email
func (s *NotificationService) GetSubscriptionByEmail(ctx context.Context, email string) (*models.NewsletterSubscription, error) {
	return s.repo.GetSubscriptionByEmail(email)
}

// GetActiveSubscriptions retrieves all active newsletter subscriptions
func (s *NotificationService) GetActiveSubscriptions(ctx context.Context, limit, offset int) ([]models.NewsletterSubscription, int64, error) {
	return s.repo.GetActiveSubscriptions(limit, offset)
}

// CreateNotification creates a new notification record
func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	return s.repo.CreateNotification(notification)
}

// GetNotificationByID retrieves a notification by ID
func (s *NotificationService) GetNotificationByID(ctx context.Context, id uint) (*models.Notification, error) {
	return s.repo.GetNotificationByID(id)
}

// GetNotificationsByUserID retrieves notifications for a user
func (s *NotificationService) GetNotificationsByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.Notification, int64, error) {
	return s.repo.GetNotificationsByUserID(userID, limit, offset)
}

// GetNotificationPreferences gets notification preferences for a user
func (s *NotificationService) GetNotificationPreferences(ctx context.Context, userID uint) (*models.NotificationPreference, error) {
	return s.repo.GetOrCreatePreference(userID)
}

// UpdateNotificationPreferences updates notification preferences
func (s *NotificationService) UpdateNotificationPreferences(ctx context.Context, preference *models.NotificationPreference) error {
	return s.repo.UpdatePreference(preference)
}

// ScheduleNotification schedules a notification for future delivery
func (s *NotificationService) ScheduleNotification(ctx context.Context, notification *models.Notification) error {
	notification.Status = models.NotificationStatusPending
	return s.repo.CreateNotification(notification)
}

// ProcessScheduledNotifications processes pending scheduled notifications
func (s *NotificationService) ProcessScheduledNotifications(ctx context.Context) error {
	notifications, err := s.repo.GetNotificationsByStatus(models.NotificationStatusPending, 100)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		// Check if scheduled time has passed
		if notification.ScheduledAt != nil && notification.ScheduledAt.After(time.Now()) {
			continue
		}

		// Process based on notification type
		var sendErr error
		switch notification.Type {
		case models.NotificationTypeEmail:
			if notification.Email != "" {
				sendErr = s.emailService.SendEmail(ctx, SendEmailInput{
					To:      []string{notification.Email},
					Subject: notification.Subject,
					Text:    notification.Content,
					HTML:    notification.Content,
				})
			}
		case models.NotificationTypeWhatsApp:
			if notification.PhoneNumber != "" {
				sendErr = s.whatsAppService.SendTextMessage(ctx, notification.PhoneNumber, notification.Content)
			}
		}

		if sendErr != nil {
			s.repo.MarkNotificationAsFailed(notification.ID, sendErr.Error())
		} else {
			s.repo.MarkNotificationAsSent(notification.ID)
		}
	}

	return nil
}

// ProcessSessionEvent processes session-related events from Kafka
func (s *NotificationService) ProcessSessionEvent(ctx context.Context, eventType string, payload []byte) error {
	switch eventType {
	case "session-upcoming":
		var session models.SessionInfo
		if err := json.Unmarshal(payload, &session); err != nil {
			return err
		}
		return s.SendSessionReminder(ctx, session)
	case "booking-created":
		var session models.SessionInfo
		if err := json.Unmarshal(payload, &session); err != nil {
			return err
		}
		return s.SendBookingConfirmation(ctx, session)
	default:
		return apperrors.ErrBadRequest
	}
}

// ProcessPromotionalEvent processes promotional events from Kafka
func (s *NotificationService) ProcessPromotionalEvent(ctx context.Context, payload []byte) error {
	var data struct {
		EmailRecipients []string `json:"email_recipients"`
		Subject         string   `json:"subject"`
		Content         string   `json:"content"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}

	if len(data.EmailRecipients) > 0 {
		return s.SendPromotionalEmail(ctx, data.EmailRecipients, data.Subject, data.Content)
	}

	return nil
}