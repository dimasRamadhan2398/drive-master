package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
)

type IMailtrapEmailService interface {
	SendEmail(ctx context.Context, input dto.SendEmailRequest) error
	SendWelcomeEmail(ctx context.Context, toEmail, username string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error
	SendBookingConfirmationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error
	SendLessonReminderEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error
	SendLessonCancellationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, reason string) error
}

// MailtrapEmailService sends emails via Mailtrap Sending API
type MailtrapEmailService struct {
	apiToken  string
	fromEmail string
	fromName  string
	client    *http.Client
}

func NewMailtrapEmailService(apiToken, fromEmail, fromName string) IMailtrapEmailService {
	return &MailtrapEmailService{
		apiToken:  apiToken,
		fromEmail: fromEmail,
		fromName:  fromName,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}



// SendEmail sends an email via Mailtrap
func (s *MailtrapEmailService) SendEmail(ctx context.Context, input dto.SendEmailRequest) error {
	if len(input.To) == 0 {
		return apperrors.ErrBadRequest
	}

	// Build recipients
	toRecipients := make([]dto.EmailAddress, len(input.To))
	for i, _ := range input.To {
		toRecipients[i] = dto.EmailAddress{Email: input.To[i].Email}
	}

	// Build CC if provided
	var ccRecipients []dto.EmailAddress
	for _, addr := range input.CC {
		ccRecipients = append(ccRecipients, dto.EmailAddress{Email: addr.Email})
	}

	// Build BCC if provided
	var bccRecipients []dto.EmailAddress
	for _, addr := range input.BCC {
		bccRecipients = append(bccRecipients, dto.EmailAddress{Email: addr.Email})
	}

	// Build request body
	reqBody := dto.SendEmailRequest{
		To:          toRecipients,
		From:        dto.EmailAddress{Email: s.fromEmail, Name: s.fromName},
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

	// Create HTTP request
	url := "https://send.api.mailtrap.io/api/send"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return apperrors.ErrInternalServer
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return apperrors.ErrInternalServer
	}

	return nil
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *MailtrapEmailService) SendWelcomeEmail(ctx context.Context, toEmail, username string) error {
	subject := "Welcome to Our Platform!"
	text := fmt.Sprintf(`Hello %s,

Welcome to our driving school platform! We're excited to have you on board.

With our platform, you can:
- Book driving lessons easily
- Track your progress
- Connect with certified instructors

Get started by logging in and exploring our features.

Best regards,
The Team`, username)

	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h1 style="color: #333;">Welcome, %s!</h1>
	<p>Welcome to our driving school platform! We're excited to have you on board.</p>
	<h2>What you can do:</h2>
	<ul>
		<li>Book driving lessons easily</li>
		<li>Track your progress</li>
		<li>Connect with certified instructors</li>
	</ul>
	<p>Get started by logging in and exploring our features.</p>
	<p style="color: #666; margin-top: 30px;">Best regards,<br>The Team</p>
</body>
</html>`, username)

	return s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
		Tags:    []string{"welcome", "onboarding"},
	})
}

// SendPasswordResetEmail sends a password reset email
func (s *MailtrapEmailService) SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error {
	subject := "Password Reset Request"
	resetLink := fmt.Sprintf("https://yourapp.com/reset-password?token=%s", resetToken)

	text := fmt.Sprintf(`Hello,

We received a request to reset your password. Click the link below to reset your password:

%s

If you didn't request this, please ignore this email.

Best regards,
The Team`, resetLink)

	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h1 style="color: #333;">Password Reset Request</h1>
	<p>We received a request to reset your password.</p>
	<p style="margin: 30px 0;">
		<a href="%s" style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px;">
			Reset Password
		</a>
	</p>
	<p>If you didn't request this, please ignore this email.</p>
	<p style="color: #666; margin-top: 30px;">Best regards,<br>The Team</p>
</body>
</html>`, resetLink)

	return s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
		Tags:    []string{"password-reset", "security"},
	})
}

// SendBookingConfirmationEmail sends a booking confirmation email
func (s *MailtrapEmailService) SendBookingConfirmationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error {
	subject := fmt.Sprintf("Booking Confirmed: %s Lesson", lessonType)

	text := fmt.Sprintf(`Hello %s,

Your driving lesson has been confirmed!

Details:
- Instructor: %s
- Date & Time: %s
- Lesson Type: %s

Please arrive 5 minutes early. If you need to reschedule, please contact us at least 24 hours in advance.

Best regards,
The Team`, studentName, instructorName, dateTime, lessonType)

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

	<p style="color: #666; margin-top: 30px;">Best regards,<br>The Team</p>
</body>
</html>`, studentName, instructorName, dateTime, lessonType)

	return s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
		Tags:    []string{"booking", "confirmation"},
	})
}

// SendLessonReminderEmail sends a reminder email before a lesson
func (s *MailtrapEmailService) SendLessonReminderEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error {
	subject := fmt.Sprintf("Reminder: Your %s Lesson Tomorrow", lessonType)

	text := fmt.Sprintf(`Hello %s,

This is a reminder about your upcoming driving lesson:

- Instructor: %s
- Date & Time: %s
- Lesson Type: %s

Please make sure you have your learner's permit ready. See you tomorrow!

Best regards,
The Team`, studentName, instructorName, dateTime, lessonType)

	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h1 style="color: #ffc107;">Lesson Reminder</h1>
	<p>Hello <strong>%s</strong>,</p>
	<p>This is a reminder about your upcoming driving lesson:</p>

	<div style="background-color: #fff3cd; padding: 20px; border-radius: 8px; margin: 20px 0;">
		<p><strong>Instructor:</strong> %s</p>
		<p><strong>Date & Time:</strong> %s</p>
		<p><strong>Lesson Type:</strong> %s</p>
	</div>

	<p>Please make sure you have your learner's permit ready. See you tomorrow!</p>

	<p style="color: #666; margin-top: 30px;">Best regards,<br>The Team</p>
</body>
</html>`, studentName, instructorName, dateTime, lessonType)

	return s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
		Tags:    []string{"reminder", "lesson"},
	})
}

// SendLessonCancellationEmail sends a lesson cancellation notification
func (s *MailtrapEmailService) SendLessonCancellationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, reason string) error {
	subject := "Lesson Cancelled"

	text := fmt.Sprintf(`Hello %s,

Your driving lesson has been cancelled.

Details:
- Instructor: %s
- Date & Time: %s
- Reason: %s

Please book a new lesson at your convenience.

Best regards,
The Team`, studentName, instructorName, dateTime, reason)

	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h1 style="color: #dc3545;">Lesson Cancelled</h1>
	<p>Hello <strong>%s</strong>,</p>
	<p>Your driving lesson has been cancelled.</p>

	<div style="background-color: #f8d7da; padding: 20px; border-radius: 8px; margin: 20px 0;">
		<p><strong>Instructor:</strong> %s</p>
		<p><strong>Date & Time:</strong> %s</p>
		<p><strong>Reason:</strong> %s</p>
	</div>

	<p>Please book a new lesson at your convenience.</p>

	<p style="color: #666; margin-top: 30px;">Best regards,<br>The Team</p>
</body>
</html>`, studentName, instructorName, dateTime, reason)

	return s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
		Tags:    []string{"cancellation", "lesson"},
	})
}
