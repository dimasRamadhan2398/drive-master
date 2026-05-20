package mocks

import (
	"context"
	"errors"
	"testing"

	"user-service/models/dto"
)

func TestEmailServiceMock_SendEmail_Success(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()
	request := dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: "muhammadrizqiko@gmail.com"}},
		Subject: "Zeta",
		Text:    "Yayaya",
	}

	err := mock.SendEmail(ctx, request)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if mock.SendEmailCallCount != 1 {
		t.Errorf("expected SendEmailCallCount to be 1, got %d", mock.SendEmailCallCount)
	}
	if len(mock.SentEmails) != 1 {
		t.Errorf("expected 1 email in SentEmails, got %d", len(mock.SentEmails))
	}
	if mock.SentEmails[0].Subject != "Zeta" {
		t.Errorf("expected subject 'Zeta', got '%s'", mock.SentEmails[0].Subject)
	}
}

func TestEmailServiceMock_SendEmail_Error(t *testing.T) {
	mock := NewEmailServiceMock()
	expectedErr := errors.New("email service unavailable")
	mock.SetError(expectedErr)

	ctx := context.Background()
	request := dto.SendEmailRequest{
		To:      []dto.EmailAddress{{Email: "test@example.com"}},
		Subject: "Test Subject",
	}

	err := mock.SendEmail(ctx, request)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if err != expectedErr {
		t.Errorf("expected error '%v', got '%v'", expectedErr, err)
	}
}

func TestEmailServiceMock_SendWelcomeEmail(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()

	err := mock.SendWelcomeEmail(ctx, "user@example.com", "John")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if mock.SendWelcomeEmailCallCount != 1 {
		t.Errorf("expected SendWelcomeEmailCallCount to be 1, got %d", mock.SendWelcomeEmailCallCount)
	}
	if len(mock.WelcomeEmailRecipients) != 1 || mock.WelcomeEmailRecipients[0] != "user@example.com" {
		t.Errorf("unexpected recipients: %v", mock.WelcomeEmailRecipients)
	}
	if len(mock.WelcomeEmailUsernames) != 1 || mock.WelcomeEmailUsernames[0] != "John" {
		t.Errorf("unexpected usernames: %v", mock.WelcomeEmailUsernames)
	}
}

func TestEmailServiceMock_SendPasswordResetEmail(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()
	token := "reset-token-123"

	err := mock.SendPasswordResetEmail(ctx, "user@example.com", token)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if mock.SendPasswordResetCallCount != 1 {
		t.Errorf("expected SendPasswordResetCallCount to be 1, got %d", mock.SendPasswordResetCallCount)
	}
	if len(mock.PasswordResetRecipients) != 1 || mock.PasswordResetRecipients[0] != "user@example.com" {
		t.Errorf("unexpected recipients: %v", mock.PasswordResetRecipients)
	}
	if len(mock.PasswordResetTokens) != 1 || mock.PasswordResetTokens[0] != token {
		t.Errorf("unexpected tokens: %v", mock.PasswordResetTokens)
	}
}

func TestEmailServiceMock_SendOTPEmail(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()
	otp := "123456"

	err := mock.SendOTPEmail(ctx, "user@example.com", otp)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if mock.SendOTPEmailCallCount != 1 {
		t.Errorf("expected SendOTPEmailCallCount to be 1, got %d", mock.SendOTPEmailCallCount)
	}
	if len(mock.OTPRecipients) != 1 || mock.OTPRecipients[0] != "user@example.com" {
		t.Errorf("unexpected recipients: %v", mock.OTPRecipients)
	}
	if len(mock.OTPCodes) != 1 || mock.OTPCodes[0] != otp {
		t.Errorf("unexpected OTPs: %v", mock.OTPCodes)
	}
}

func TestEmailServiceMock_SendBookingConfirmationEmail(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()

	err := mock.SendBookingConfirmationEmail(ctx, "student@example.com", "John", "Instructor", "2024-01-15 10:00", "Defensive Driving")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if mock.SendBookingConfirmCallCount != 1 {
		t.Errorf("expected SendBookingConfirmCallCount to be 1, got %d", mock.SendBookingConfirmCallCount)
	}
}

func TestEmailServiceMock_SendLessonReminderEmail(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()

	err := mock.SendLessonReminderEmail(ctx, "student@example.com", "John", "Instructor", "2024-01-15 10:00", "Defensive Driving")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if mock.SendLessonReminderCallCount != 1 {
		t.Errorf("expected SendLessonReminderCallCount to be 1, got %d", mock.SendLessonReminderCallCount)
	}
}

func TestEmailServiceMock_SendLessonCancellationEmail(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()

	err := mock.SendLessonCancellationEmail(ctx, "student@example.com", "John", "Instructor", "2024-01-15 10:00", "Instructor unavailable")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if mock.SendCancellationCallCount != 1 {
		t.Errorf("expected SendCancellationCallCount to be 1, got %d", mock.SendCancellationCallCount)
	}
}

func TestEmailServiceMock_Reset(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()

	// Make some calls
	mock.SendWelcomeEmail(ctx, "user1@example.com", "User1")
	mock.SendOTPEmail(ctx, "user2@example.com", "123456")
	mock.SetError(errors.New("test error"))

	// Reset
	mock.Reset()

	if mock.SendWelcomeEmailCallCount != 0 {
		t.Errorf("expected SendWelcomeEmailCallCount to be 0 after reset, got %d", mock.SendWelcomeEmailCallCount)
	}
	if mock.SendOTPEmailCallCount != 0 {
		t.Errorf("expected SendOTPEmailCallCount to be 0 after reset, got %d", mock.SendOTPEmailCallCount)
	}
	if mock.ShouldError {
		t.Error("expected ShouldError to be false after reset")
	}
	if len(mock.SentEmails) != 0 {
		t.Errorf("expected SentEmails to be empty after reset, got %d", len(mock.SentEmails))
	}
}

func TestEmailServiceMock_MultipleEmails(t *testing.T) {
	mock := NewEmailServiceMock()
	ctx := context.Background()

	// Send multiple emails
	mock.SendOTPEmail(ctx, "user1@example.com", "111111")
	mock.SendOTPEmail(ctx, "user2@example.com", "222222")
	mock.SendOTPEmail(ctx, "user3@example.com", "333333")

	if mock.SendOTPEmailCallCount != 3 {
		t.Errorf("expected SendOTPEmailCallCount to be 3, got %d", mock.SendOTPEmailCallCount)
	}
	if len(mock.OTPRecipients) != 3 {
		t.Errorf("expected 3 OTP recipients, got %d", len(mock.OTPRecipients))
	}
	if len(mock.OTPCodes) != 3 {
		t.Errorf("expected 3 OTP codes, got %d", len(mock.OTPCodes))
	}
}
