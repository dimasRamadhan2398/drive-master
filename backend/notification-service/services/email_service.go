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

// EmailService handles email sending via Mailtrap
type EmailService struct {
	apiToken  string
	fromEmail string
	fromName  string
	client    *http.Client
}

// EmailAddress represents an email address
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
}

// SendEmailRequest represents the Mailtrap API request body
type SendEmailRequest struct {
	To          []EmailAddress     `json:"to"`
	From        EmailAddress       `json:"from"`
	Subject     string             `json:"subject"`
	Text        string             `json:"text,omitempty"`
	HTML        string             `json:"html,omitempty"`
	CC          []EmailAddress     `json:"cc,omitempty"`
	BCC         []EmailAddress     `json:"bcc,omitempty"`
	Attachments []EmailAttachment  `json:"attachments,omitempty"`
	CustomArgs  map[string]string  `json:"custom_args,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
}

// NewEmailService creates a new email service
func NewEmailService(apiToken, fromEmail, fromName string) *EmailService {
	return &EmailService{
		apiToken:  apiToken,
		fromEmail: fromEmail,
		fromName:  fromName,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendEmailInput contains the email data to be sent
type SendEmailInput struct {
	To          []string
	Subject     string
	Text        string
	HTML        string
	CC          []string
	BCC         []string
	Attachments []EmailAttachment
	CustomArgs  map[string]string
	Tags        []string
}

// SendEmail sends an email via Mailtrap
func (s *EmailService) SendEmail(ctx context.Context, input SendEmailInput) error {
	if len(input.To) == 0 {
		return apperrors.ErrBadRequest
	}

	// Build recipients
	toRecipients := make([]EmailAddress, len(input.To))
	for i, email := range input.To {
		toRecipients[i] = EmailAddress{Email: email}
	}

	// Build CC if provided
	var ccRecipients []EmailAddress
	for _, email := range input.CC {
		ccRecipients = append(ccRecipients, EmailAddress{Email: email})
	}

	// Build BCC if provided
	var bccRecipients []EmailAddress
	for _, email := range input.BCC {
		bccRecipients = append(bccRecipients, EmailAddress{Email: email})
	}

	// Build request body
	reqBody := SendEmailRequest{
		To:          toRecipients,
		From:        EmailAddress{Email: s.fromEmail, Name: s.fromName},
		Subject:     input.Subject,
		Text:        input.Text,
		HTML:        input.HTML,
		CC:          ccRecipients,
		BCC:         bccRecipients,
		Attachments: input.Attachments,
		CustomArgs:  input.CustomArgs,
		Tags:        input.Tags,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return apperrors.ErrInternalServer
	}

	url := "https://send.api.mailtrap.io/api/send"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return apperrors.ErrInternalServer
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return apperrors.ErrEmailSendFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return apperrors.ErrEmailSendFailed
	}

	return nil
}

// SendSessionReminderEmail sends a session reminder email
func (s *EmailService) SendSessionReminderEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType, location string) error {
	subject := fmt.Sprintf("Reminder: Your %s Lesson Tomorrow", lessonType)

	text := fmt.Sprintf(`Hello %s,

This is a reminder about your upcoming driving lesson:

- Instructor: %s
- Date & Time: %s
- Lesson Type: %s
- Location: %s

Please arrive 5 minutes early. Make sure you have your learner's permit ready.

Best regards,
Drive Master Team`, studentName, instructorName, dateTime, lessonType, location)

	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h1 style="color: #ffc107;">Lesson Reminder</h1>
	<p>Hello <strong>%s</strong>,</p>
	<p>This is a reminder about your upcoming driving lesson:</p>

	<div style="background-color: #fff3cd; padding: 20px; border-radius: 8px; margin: 20px 0;">
		<p><strong>Instructor:</strong> %s</p>
		<p><strong>Date & Time:</strong> %s</p>
		<p><strong>Lesson Type:</strong> %s</p>
		<p><strong>Location:</strong> %s</p>
	</div>

	<p>Please arrive 5 minutes early. Make sure you have your learner's permit ready.</p>

	<p style="color: #666; margin-top: 30px;">Best regards,<br>Drive Master Team</p>
</body>
</html>`, studentName, instructorName, dateTime, lessonType, location)

	return s.SendEmail(ctx, SendEmailInput{
		To:      []string{toEmail},
		Subject: subject,
		Text:    text,
		HTML:    html,
		Tags:    []string{"reminder", "session"},
	})
}

// SendBookingConfirmationEmail sends a booking confirmation email
func (s *EmailService) SendBookingConfirmationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error {
	subject := fmt.Sprintf("Booking Confirmed: %s Lesson", lessonType)

	text := fmt.Sprintf(`Hello %s,

Your driving lesson has been confirmed!

Details:
- Instructor: %s
- Date & Time: %s
- Lesson Type: %s

Please arrive 5 minutes early. If you need to reschedule, please contact us at least 24 hours in advance.

Best regards,
Drive Master Team`, studentName, instructorName, dateTime, lessonType)

	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h1 style="color: #28a745;">Booking Confirmed!</h1>
	<p>Hello <strong>%s</strong>,</p>
	<p>Your driving lesson has been confirmed!</p>

	<div style="background-color: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0;">
		<p><strong>Instructor:</strong> %s</p>
		<p><strong>Date & Time:</strong> %s</p>
		<p><strong>Lesson Type:</strong> %s</p>
	</div>

	<p>Please arrive 5 minutes early. If you need to reschedule, please contact us at least 24 hours in advance.</p>

	<p style="color: #666; margin-top: 30px;">Best regards,<br>Drive Master Team</p>
</body>
</html>`, studentName, instructorName, dateTime, lessonType)

	return s.SendEmail(ctx, SendEmailInput{
		To:      []string{toEmail},
		Subject: subject,
		Text:    text,
		HTML:    html,
		Tags:    []string{"booking", "confirmation"},
	})
}

// SendPromotionalEmail sends a promotional email
func (s *EmailService) SendPromotionalEmail(ctx context.Context, toEmails []string, subject, content string) error {
	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<div style="text-align: center;">
		<h1 style="color: #007bff;">Special Offer!</h1>
		<div style="padding: 20px;">%s</div>
		<p style="color: #666;">Don't miss out on this exclusive deal.</p>
		<p style="color: #666; margin-top: 30px;">Best regards,<br>Drive Master Team</p>
	</div>
</body>
</html>`, content)

	return s.SendEmail(ctx, SendEmailInput{
		To:      toEmails,
		Subject: subject,
		HTML:    html,
		Tags:    []string{"promotional", "offer"},
	})
}

// SendNewsletterEmail sends a newsletter email to subscribers
func (s *EmailService) SendNewsletterEmail(ctx context.Context, toEmails []string, subject, content string) error {
	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<div style="text-align: center;">
		<h1 style="color: #333;">%s</h1>
		<div style="padding: 20px; text-align: left;">%s</div>
		<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
		<p style="color: #999; font-size: 12px;">
			You're receiving this email because you subscribed to our newsletter.<br>
			<a href="{{unsubscribe_url}}" style="color: #999;">Unsubscribe</a>
		</p>
	</div>
</body>
</html>`, subject, content)

	return s.SendEmail(ctx, SendEmailInput{
		To:      toEmails,
		Subject: subject,
		HTML:    html,
		Tags:    []string{"newsletter"},
	})
}