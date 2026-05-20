package mocks

import (
	"context"
	"sync"

	"user-service/models/dto"
	"user-service/services"
)

// EmailServiceMock is a mock implementation of IMailtrapEmailService for testing
type EmailServiceMock struct {
	mu sync.RWMutex

	// Configuration
	ShouldError bool
	ErrorToReturn error

	// Call tracking
	SendEmailCallCount          int
	SendWelcomeEmailCallCount   int
	SendPasswordResetCallCount  int
	SendOTPEmailCallCount       int
	SendBookingConfirmCallCount int
	SendLessonReminderCallCount int
	SendCancellationCallCount   int

	// Captured data for assertions
	LastSendEmailRequest   dto.SendEmailRequest
	SentEmails             []dto.SendEmailRequest
	WelcomeEmailRecipients []string
	WelcomeEmailUsernames  []string
	PasswordResetRecipients []string
	PasswordResetTokens    []string
	OTPRecipients          []string
	OTPCodes               []string
}

// NewEmailServiceMock creates a new EmailServiceMock instance
func NewEmailServiceMock() *EmailServiceMock {
	return &EmailServiceMock{
		SentEmails:             make([]dto.SendEmailRequest, 0),
		WelcomeEmailRecipients: make([]string, 0),
		WelcomeEmailUsernames:  make([]string, 0),
		PasswordResetRecipients: make([]string, 0),
		PasswordResetTokens:    make([]string, 0),
		OTPRecipients:          make([]string, 0),
		OTPCodes:               make([]string, 0),
	}
}

// Ensure EmailServiceMock implements IMailtrapEmailService
var _ services.IMailtrapEmailService = (*EmailServiceMock)(nil)

// Reset clears all tracked data and call counts
func (m *EmailServiceMock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ShouldError = false
	m.ErrorToReturn = nil
	m.SendEmailCallCount = 0
	m.SendWelcomeEmailCallCount = 0
	m.SendPasswordResetCallCount = 0
	m.SendOTPEmailCallCount = 0
	m.SendBookingConfirmCallCount = 0
	m.SendLessonReminderCallCount = 0
	m.SendCancellationCallCount = 0
	m.LastSendEmailRequest = dto.SendEmailRequest{}
	m.SentEmails = make([]dto.SendEmailRequest, 0)
	m.WelcomeEmailRecipients = make([]string, 0)
	m.WelcomeEmailUsernames = make([]string, 0)
	m.PasswordResetRecipients = make([]string, 0)
	m.PasswordResetTokens = make([]string, 0)
	m.OTPRecipients = make([]string, 0)
	m.OTPCodes = make([]string, 0)
}

// SetError configures the mock to return an error
func (m *EmailServiceMock) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ShouldError = true
	m.ErrorToReturn = err
}

// SendEmail implements IMailtrapEmailService.SendEmail
func (m *EmailServiceMock) SendEmail(ctx context.Context, input dto.SendEmailRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendEmailCallCount++
	m.LastSendEmailRequest = input
	m.SentEmails = append(m.SentEmails, input)

	if m.ShouldError {
		return m.ErrorToReturn
	}

	return nil
}

// SendWelcomeEmail implements IMailtrapEmailService.SendWelcomeEmail
func (m *EmailServiceMock) SendWelcomeEmail(ctx context.Context, toEmail, username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendWelcomeEmailCallCount++
	m.WelcomeEmailRecipients = append(m.WelcomeEmailRecipients, toEmail)
	m.WelcomeEmailUsernames = append(m.WelcomeEmailUsernames, username)

	if m.ShouldError {
		return m.ErrorToReturn
	}

	return nil
}

// SendPasswordResetEmail implements IMailtrapEmailService.SendPasswordResetEmail
func (m *EmailServiceMock) SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendPasswordResetCallCount++
	m.PasswordResetRecipients = append(m.PasswordResetRecipients, toEmail)
	m.PasswordResetTokens = append(m.PasswordResetTokens, resetToken)

	if m.ShouldError {
		return m.ErrorToReturn
	}

	return nil
}

// SendOTPEmail implements IMailtrapEmailService.SendOTPEmail
func (m *EmailServiceMock) SendOTPEmail(ctx context.Context, toEmail, otp string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendOTPEmailCallCount++
	m.OTPRecipients = append(m.OTPRecipients, toEmail)
	m.OTPCodes = append(m.OTPCodes, otp)

	if m.ShouldError {
		return m.ErrorToReturn
	}

	return nil
}

// SendBookingConfirmationEmail implements IMailtrapEmailService.SendBookingConfirmationEmail
func (m *EmailServiceMock) SendBookingConfirmationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendBookingConfirmCallCount++

	if m.ShouldError {
		return m.ErrorToReturn
	}

	return nil
}

// SendLessonReminderEmail implements IMailtrapEmailService.SendLessonReminderEmail
func (m *EmailServiceMock) SendLessonReminderEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendLessonReminderCallCount++

	if m.ShouldError {
		return m.ErrorToReturn
	}

	return nil
}

// SendLessonCancellationEmail implements IMailtrapEmailService.SendLessonCancellationEmail
func (m *EmailServiceMock) SendLessonCancellationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendCancellationCallCount++

	if m.ShouldError {
		return m.ErrorToReturn
	}

	return nil
}
