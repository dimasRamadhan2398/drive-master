package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apperrors "notification-service/pkg/errors"
)

// WhatsAppService handles WhatsApp message sending
type WhatsAppService struct {
	apiURL      string
	apiToken    string
	phoneNumber string
	client      *http.Client
}

// WhatsAppMessage represents a WhatsApp message
type WhatsAppMessage struct {
	To        string            `json:"to"`
	Type      string            `json:"type"`
	Text      *WhatsAppText     `json:"text,omitempty"`
	Template  *WhatsAppTemplate `json:"template,omitempty"`
	Image     *WhatsAppImage    `json:"image,omitempty"`
}

// WhatsAppText represents text content in a WhatsApp message
type WhatsAppText struct {
	Body string `json:"body"`
}

// WhatsAppTemplate represents a WhatsApp template message
type WhatsAppTemplate struct {
	Namespace string            `json:"namespace"`
	Name      string            `json:"name"`
	Language  WhatsAppLanguage  `json:"language"`
	Components []WhatsAppComponent `json:"components,omitempty"`
}

// WhatsAppLanguage represents the language for a template
type WhatsAppLanguage struct {
	Code string `json:"code"`
}

// WhatsAppComponent represents a component in a template
type WhatsAppComponent struct {
	Type       string   `json:"type"`
	Parameters []Parameter `json:"parameters,omitempty"`
}

// Parameter represents a parameter in a template component
type Parameter struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Currency *CurrencyParameter `json:"currency,omitempty"`
	DateTime *DateTimeParameter  `json:"date_time,omitempty"`
}

// CurrencyParameter represents a currency parameter
type CurrencyParameter struct {
	CurrencyCode string `json:"currency_code"`
	Amount       int64  `json:"amount_1000"`
}

// DateTimeParameter represents a date time parameter
type DateTimeParameter struct {
	DateTime string `json:"date_time"`
}

// WhatsAppImage represents an image in a WhatsApp message
type WhatsAppImage struct {
	Link    string `json:"link"`
	Caption string `json:"caption,omitempty"`
}

// NewWhatsAppService creates a new WhatsApp service
func NewWhatsAppService(apiURL, apiToken, phoneNumber string) *WhatsAppService {
	return &WhatsAppService{
		apiURL:      apiURL,
		apiToken:    apiToken,
		phoneNumber: phoneNumber,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendTextMessage sends a simple text WhatsApp message
func (s *WhatsAppService) SendTextMessage(ctx context.Context, to, message string) error {
	if to == "" {
		return apperrors.ErrInvalidPhoneNumber
	}

	// Format phone number (ensure it has country code)
	to = s.formatPhoneNumber(to)

	msg := WhatsAppMessage{
		To:   to,
		Type: "text",
		Text: &WhatsAppText{Body: message},
	}

	return s.sendMessage(ctx, msg)
}

// SendSessionReminderWhatsApp sends a session reminder via WhatsApp
func (s *WhatsAppService) SendSessionReminderWhatsApp(ctx context.Context, to, studentName, instructorName, dateTime, lessonType string) error {
	if to == "" {
		return apperrors.ErrInvalidPhoneNumber
	}

	to = s.formatPhoneNumber(to)

	message := fmt.Sprintf(`Hello %s!

⏰ *Session Reminder*

Your driving lesson is coming up tomorrow!

📚 Type: %s
👨‍🏫 Instructor: %s
🕐 Time: %s

Please arrive 5 minutes early with your learner's permit.

See you soon!
- Drive Master`, studentName, lessonType, instructorName, dateTime)

	msg := WhatsAppMessage{
		To:   to,
		Type: "text",
		Text: &WhatsAppText{Body: message},
	}

	return s.sendMessage(ctx, msg)
}

// SendBookingConfirmationWhatsApp sends a booking confirmation via WhatsApp
func (s *WhatsAppService) SendBookingConfirmationWhatsApp(ctx context.Context, to, studentName, instructorName, dateTime, lessonType string) error {
	if to == "" {
		return apperrors.ErrInvalidPhoneNumber
	}

	to = s.formatPhoneNumber(to)

	message := fmt.Sprintf(`Hello %s! ✅

*Booking Confirmed*

📚 Type: %s
👨‍🏫 Instructor: %s
🕐 Time: %s

Please arrive 5 minutes early.

Questions? Reply to this message!
- Drive Master`, studentName, lessonType, instructorName, dateTime)

	msg := WhatsAppMessage{
		To:   to,
		Type: "text",
		Text: &WhatsAppText{Body: message},
	}

	return s.sendMessage(ctx, msg)
}

// SendPromotionalWhatsApp sends a promotional message via WhatsApp
func (s *WhatsAppService) SendPromotionalWhatsApp(ctx context.Context, to string, title, content string) error {
	if to == "" {
		return apperrors.ErrInvalidPhoneNumber
	}

	to = s.formatPhoneNumber(to)

	message := fmt.Sprintf(`🎉 *%s*

%s

Don't miss out on this exclusive offer!
- Drive Master`, title, content)

	msg := WhatsAppMessage{
		To:   to,
		Type: "text",
		Text: &WhatsAppText{Body: message},
	}

	return s.sendMessage(ctx, msg)
}

// sendMessage sends a WhatsApp message via the API
func (s *WhatsAppService) sendMessage(ctx context.Context, msg WhatsAppMessage) error {
	if s.apiURL == "" {
		// Skip sending if no API URL configured (for development)
		return nil
	}

	jsonBody, err := json.Marshal(msg)
	if err != nil {
		return apperrors.ErrInternalServer
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL+"/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return apperrors.ErrInternalServer
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return apperrors.ErrWhatsAppSendFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return apperrors.ErrWhatsAppSendFailed
	}

	return nil
}

// formatPhoneNumber ensures the phone number is in the correct format
func (s *WhatsAppService) formatPhoneNumber(phone string) string {
	// Remove any existing + sign
	if len(phone) > 0 && phone[0] == '+' {
		phone = phone[1:]
	}

	// If it doesn't start with a country code, assume it's a local number
	// This should be configured based on your region
	if len(phone) > 0 && phone[0] != '1' && phone[0] != '7' && phone[0] != '2' && phone[0] != '3' && phone[0] != '4' && phone[0] != '5' && phone[0] != '6' && phone[0] != '8' && phone[0] != '9' {
		phone = "1" + phone // Default to US country code
	}

	return phone
}

// SendTemplateMessage sends a templated WhatsApp message
func (s *WhatsAppService) SendTemplateMessage(ctx context.Context, to, namespace, templateName, languageCode string, components []WhatsAppComponent) error {
	if to == "" {
		return apperrors.ErrInvalidPhoneNumber
	}

	to = s.formatPhoneNumber(to)

	msg := WhatsAppMessage{
		To:   to,
		Type: "template",
		Template: &WhatsAppTemplate{
			Namespace: namespace,
			Name:      templateName,
			Language:  WhatsAppLanguage{Code: languageCode},
			Components: components,
		},
	}

	return s.sendMessage(ctx, msg)
}