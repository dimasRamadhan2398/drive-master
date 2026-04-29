package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Scheduler handles notification scheduling
type Scheduler struct {
	redis       *redis.Client
	emailService EmailServiceInterface
	whatsAppService WhatsAppServiceInterface
}

// EmailServiceInterface defines email service methods
type EmailServiceInterface interface {
	SendSessionReminderEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType, location string) error
}

// WhatsAppServiceInterface defines WhatsApp service methods
type WhatsAppServiceInterface interface {
	SendSessionReminderWhatsApp(ctx context.Context, to, studentName, instructorName, dateTime, lessonType string) error
}

// NewScheduler creates a new scheduler
func NewScheduler(redisClient *redis.Client, emailSvc EmailServiceInterface, whatsappSvc WhatsAppServiceInterface) *Scheduler {
	return &Scheduler{
		redis:          redisClient,
		emailService:   emailSvc,
		whatsAppService: whatsappSvc,
	}
}

// ScheduledReminder represents a scheduled reminder
type ScheduledReminder struct {
	SessionID      uint      `json:"sessionId"`
	UserID         uint      `json:"userId"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phoneNumber"`
	InstructorName string    `json:"instructorName"`
	StudentName    string    `json:"studentName"`
	SessionDate    time.Time `json:"sessionDate"`
	LessonType     string    `json:"lessonType"`
	Location       string    `json:"location"`
	EmailSent      bool      `json:"emailSent"`
	WhatsAppSent   bool      `json:"whatsAppSent"`
}

// ScheduleSessionReminder schedules a session reminder
func (s *Scheduler) ScheduleSessionReminder(ctx context.Context, reminder *ScheduledReminder) error {
	data, err := json.Marshal(reminder)
	if err != nil {
		return err
	}

	// Schedule email reminder (24 hours before)
	emailKey := sessionReminderKeyPrefix + "email:" + string(rune(reminder.SessionID))
	emailTime := reminder.SessionDate.Add(-24 * time.Hour)

	if emailTime.After(time.Now()) {
		s.redis.Set(ctx, emailKey, data, time.Until(emailTime))
	}

	// Schedule WhatsApp reminder (24 hours before)
	whatsAppKey := whatsAppReminderKeyPrefix + reminder.PhoneNumber + ":" + string(rune(reminder.SessionID))
	whatsAppTime := reminder.SessionDate.Add(-24 * time.Hour)

	if whatsAppTime.After(time.Now()) {
		s.redis.Set(ctx, whatsAppKey, data, time.Until(whatsAppTime))
	}

	return nil
}

// Start starts the scheduler
func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(reminderCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.processReminders(ctx)
		}
	}
}

// processReminders processes scheduled reminders
func (s *Scheduler) processReminders(ctx context.Context) {
	// Get all pending email reminders
	emailKeys, err := s.redis.Keys(ctx, sessionReminderKeyPrefix+"*").Result()
	if err != nil {
		log.Printf("Error fetching email reminder keys: %v", err)
		return
	}

	for _, key := range emailKeys {
		data, err := s.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var reminder ScheduledReminder
		if err := json.Unmarshal([]byte(data), &reminder); err != nil {
			continue
		}

		if !reminder.EmailSent && reminder.Email != "" {
			dateTimeStr := reminder.SessionDate.Format("Monday, January 2, 2006 at 3:04 PM")
			if err := s.emailService.SendSessionReminderEmail(
				ctx,
				reminder.Email,
				reminder.StudentName,
				reminder.InstructorName,
				dateTimeStr,
				reminder.LessonType,
				reminder.Location,
			); err != nil {
				log.Printf("Failed to send email reminder: %v", err)
			} else {
				reminder.EmailSent = true
				s.redis.Del(ctx, key)
			}
		}
	}

	// Get all pending WhatsApp reminders
	whatsAppKeys, err := s.redis.Keys(ctx, whatsAppReminderKeyPrefix+"*").Result()
	if err != nil {
		log.Printf("Error fetching WhatsApp reminder keys: %v", err)
		return
	}

	for _, key := range whatsAppKeys {
		data, err := s.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var reminder ScheduledReminder
		if err := json.Unmarshal([]byte(data), &reminder); err != nil {
			continue
		}

		if !reminder.WhatsAppSent && reminder.PhoneNumber != "" {
			dateTimeStr := reminder.SessionDate.Format("Monday 3:04 PM")
			if err := s.whatsAppService.SendSessionReminderWhatsApp(
				ctx,
				reminder.PhoneNumber,
				reminder.StudentName,
				reminder.InstructorName,
				dateTimeStr,
				reminder.LessonType,
			); err != nil {
				log.Printf("Failed to send WhatsApp reminder: %v", err)
			} else {
				reminder.WhatsAppSent = true
				s.redis.Del(ctx, key)
			}
		}
	}
}

// CancelReminder cancels a scheduled reminder
func (s *Scheduler) CancelReminder(ctx context.Context, sessionID uint, phoneNumber string) error {
	emailKey := sessionReminderKeyPrefix + "email:" + string(rune(sessionID))
	whatsAppKey := whatsAppReminderKeyPrefix + phoneNumber + ":" + string(rune(sessionID))

	s.redis.Del(ctx, emailKey)
	s.redis.Del(ctx, whatsAppKey)

	return nil
}

const (
	sessionReminderKeyPrefix = "session:reminder:"
	whatsAppReminderKeyPrefix = "whatsapp:reminder:"
	reminderCheckInterval = 1 * time.Hour
)