package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"user-service/models/dto"
	"user-service/pkg/base"
	apperrors "user-service/pkg/errors"
	"user-service/pkg/logger"
)

type IMailtrapEmailService interface {
	SendEmail(ctx context.Context, input dto.SendEmailRequest) error
	SendWelcomeEmail(ctx context.Context, toEmail, username string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error
	SendOTPEmail(ctx context.Context, toEmail, otp string) error
	SendBookingConfirmationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error
	SendLessonReminderEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error
	SendLessonCancellationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, reason string) error
}

// MailtrapEmailService sends emails via Mailtrap Sending API
type MailtrapEmailService struct {
	*base.BaseService
	fromEmail string
	fromName  string
	apiKey    string
	client    *http.Client
}

func NewMailtrapEmailService(fromEmail, fromName, apiKey string) IMailtrapEmailService {
	return &MailtrapEmailService{
		BaseService: base.NewBaseService(),
		fromEmail: fromEmail,
		fromName:  fromName,
		apiKey:    apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *MailtrapEmailService) SendEmail(ctx context.Context, input dto.SendEmailRequest) error {
	if len(input.To) == 0 {
		return apperrors.ErrBadRequest
	}

	// Build recipients
	toRecipients := make([]dto.EmailAddress, len(input.To))
	for i := range input.To {
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

	type mailtrapRequest struct {
		To          []dto.EmailAddress    `json:"to"`
		From        dto.EmailAddress      `json:"from"`
		Subject     string                `json:"subject"`
		Text        string                `json:"text,omitempty"`
		HTML        string                `json:"html,omitempty"`
		CC          []dto.EmailAddress    `json:"cc,omitempty"`
		BCC         []dto.EmailAddress    `json:"bcc,omitempty"`
		Attachments []dto.EmailAttachment `json:"attachments,omitempty"`
		CustomVariables map[string]string     `json:"custom_variables,omitempty"`
	}

	reqBody := mailtrapRequest{
		To:          toRecipients,
		From:        dto.EmailAddress{Email: s.fromEmail, Name: s.fromName},
		Subject:     input.Subject,
		Text:        input.Text,
		HTML:        input.HTML,
		CC:          ccRecipients,
		BCC:         bccRecipients,
		Attachments: input.Attachments,
		CustomVariables:  input.CustomVariables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		s.LogError("Failed to marshal email json", logger.LogField("error", err))
		return apperrors.ErrInternalServer
	}

	// Create HTTP request
	url := "https://send.api.mailtrap.io/api/send"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		s.LogError("Failed to create email request", logger.LogField("error", err))
		return apperrors.ErrInternalServer
	}

	// Set headers with Bearer token
	req.Header.Add("Authorization", "Bearer "+s.apiKey)
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		s.LogError("Failed to send email request", logger.LogField("error", err))
		return apperrors.ErrInternalServer
	}
	defer resp.Body.Close()

	// Read response body
	body, _ := ioutil.ReadAll(resp.Body)

	// Check response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		s.LogError("Email API error",
			logger.LogField("status", resp.StatusCode),
			logger.LogField("response", string(body)),
		)
		return apperrors.ErrInternalServer
	}

	s.LogInfo("Email sent successfully",
		logger.LogField("from", s.fromEmail),
		logger.LogField("subject", input.Subject),
	)

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

	err := s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
	})

	if err != nil {
		s.LogError("Failed to send welcome email", logger.LogField("error", err))
		return err
	}

	return nil
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

	err := s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
	})

	if err != nil {
		s.LogError("Failed to send password reset email", logger.LogField("error", err))
		return err
	}

	return nil
}

// SendOTPEmail sends an OTP verification email
func (s *MailtrapEmailService) SendOTPEmail(ctx context.Context, toEmail, otp string) error {
	subject := "Email Verification - Your OTP Code"

	text := fmt.Sprintf(`Hello,

Thank you for registering with us. To verify your email address, please use the following OTP code:

%s

This code is valid for 15 minutes. If you did not request this email, please ignore it.

Best regards,
The Team`, otp)

	html := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h1 style="color: #333;">Email Verification</h1>
	<p>Thank you for registering with us. To verify your email address, please use the following OTP code:</p>

	<div style="background-color: #f8f9fa; padding: 30px; border-radius: 8px; margin: 30px 0; text-align: center;">
		<p style="font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #007bff; margin: 0;">%s</p>
	</div>

	<p style="color: #666;">This code is valid for <strong>15 minutes</strong>.</p>
	<p style="color: #666;">If you did not request this email, please ignore it.</p>
	<p style="color: #666; margin-top: 30px;">Best regards,<br>The Team</p>
</body>
</html>`, otp)

	err := s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
	})

	if err != nil {
		s.LogError("Failed to send OTP email", logger.LogField("error", err))
		return err
	}

	return nil
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

	err := s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
	})

	if err != nil {
		s.LogError("Failed to send booking confirmation email", logger.LogField("error", err))
		return err
	}

	return nil
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

	err := s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
	})

	if err != nil {
		s.LogError("Failed to send lesson reminder email", logger.LogField("error", err))
		return err
	}

	return nil
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

	err := s.SendEmail(ctx, dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: toEmail}},
		Subject: subject,
		Text:    text,
		HTML:    html,
	})

	if err != nil {
		s.LogError("Failed to send lesson cancellation email", logger.LogField("error", err))
		return err
	}

	return nil
}
