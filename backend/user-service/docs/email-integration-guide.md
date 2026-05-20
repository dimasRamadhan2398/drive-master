# Email Integration Guide

This guide explains how email sending is implemented in this codebase using **Mailtrap Sending API**.

---

## Overview

The service uses **Mailtrap's HTTP API** (not SMTP) for sending transactional emails. This approach:
- Requires only an API key (no SMTP credentials)
- Sends JSON requests via HTTP POST
- Supports HTML emails, attachments, and analytics via tags

---

## 1. Configuration

Add email settings to `pkg/config/config.yaml`:

```yaml
email:
  api_key: <your-mailtrap-api-key>
  from_email: "Admin Drive Master Indonesia"
  app_name: user_service
  enabled: true
```

---

## 2. Service Implementation

The email service is in `services/email_service.go`:

### Interface Definition

```go
type IMailtrapEmailService interface {
    SendEmail(ctx context.Context, input dto.SendEmailRequest) error
    SendWelcomeEmail(ctx context.Context, toEmail, username string) error
    SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error
    SendOTPEmail(ctx context.Context, toEmail, otp string) error
    SendBookingConfirmationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error
    SendLessonReminderEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error
    SendLessonCancellationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, reason string) error
}
```

### Key Implementation Details

```go
type MailtrapEmailService struct {
    *base.BaseService
    apiToken  string      // Mailtrap API token
    fromEmail string      // Sender email
    fromName  string      // Sender name
    client    *http.Client // HTTP client with 30s timeout
}
```

### Constructor

```go
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
```

---

## 3. Integration Steps

### Step 1: Create the Email Service

In `services/email_service.go`, implement the service:

```go
func (s *MailtrapEmailService) SendEmail(ctx context.Context, input dto.SendEmailRequest) error {
    // Use dedicated timeout context
    emailCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Build request body
    reqBody := dto.SendEmailRequest{
        To:      []dto.EmailAddress{{Email: input.To[0].Email}},
        From:    dto.EmailAddress{Email: s.fromEmail, Name: s.fromName},
        Subject: input.Subject,
        Text:    input.Text,
        HTML:    input.HTML,
        Tags:    input.Tags,
    }

    jsonBody, _ := json.Marshal(reqBody)

    // Create HTTP request
    url := "https://send.api.mailtrap.io/api/send"
    req, _ := http.NewRequestWithContext(emailCtx, http.MethodPost, url, bytes.NewBuffer(jsonBody))

    // Set headers
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
    req.Header.Set("Content-Type", "application/json")

    // Send request
    resp, err := s.client.Do(req)
    if err != nil {
        s.LogError("Failed to send email", logger.LogField("error", err))
        return apperrors.ErrInternalServer
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
        return apperrors.ErrInternalServer
    }

    return nil
}
```

### Step 2: Register the Service in Registry

In `services/registry.go`:

```go
type IServiceRegistry interface {
    // ... other services
    GetEmailService() IMailtrapEmailService
}

func (r *Registry) GetEmailService() IMailtrapEmailService {
    cfg := config.Get()
    return NewMailtrapEmailService(cfg.Email.APIKey, cfg.Email.FromEmail, cfg.Email.AppName)
}
```

### Step 3: Inject into Dependent Services

In `services/auth_service.go`:

```go
type AuthService struct {
    *base.BaseService
    userRepo      repositories.IUserRepository
    redisCli      *redis.Client
    emailService  IMailtrapEmailService  // Add email service
    memberService IMemberService
    instructorService IInstructorService
    roleService   IRoleService
}

func NewAuthService(
    userRepo repositories.IUserRepository,
    redisCli *redis.Client,
    emailService IMailtrapEmailService,  // Add parameter
    memberService IMemberService,
    instructorService IInstructorService,
    roleService IRoleService,
) IAuthService {
    return &AuthService{
        userRepo:        userRepo,
        redisCli:        redisCli,
        emailService:    emailService,  // Assign
        memberService:   memberService,
        instructorService: instructorService,
        roleService:     roleService,
    }
}
```

Update the registry to pass the email service:

```go
func (r *Registry) GetAuthService() IAuthService {
    return NewAuthService(
        r.repoRegistry.GetUser(),
        r.redisClient,
        r.GetEmailService(),  // Pass email service
        r.GetMemberService(),
        r.GetInstructorService(),
        r.GetRoleService(),
    )
}
```

### Step 4: Use in Controllers

In `controllers/auth_controller.go`:

```go
type AuthController struct {
    authService  services.IAuthService
    userService  services.IUserService
    emailService services.IMailtrapEmailService  // Add
    roleService  services.IRoleService
}

func NewAuthController(
    userService services.IUserService,
    authService services.IAuthService,
    emailService services.IMailtrapEmailService,  // Add parameter
    roleService services.IRoleService,
) IAuthController {
    return &AuthController{
        userService:  userService,
        authService:  authService,
        emailService: emailService,  // Assign
        roleService:  roleService,
    }
}
```

---

## 4. Sending Emails

### Send OTP Email

In `services/auth_service.go`:

```go
func (s *AuthService) GenerateAndSendOTP(ctx context.Context, email string) error {
    // ... OTP generation logic ...

    // Send OTP via email
    if err := s.emailService.SendOTPEmail(ctx, email, otp); err != nil {
        s.LogError("Failed to send OTP email", logger.LogField("error", err))
        return errors.ErrInternalServer
    }

    return nil
}
```

### Send Password Reset Email

In controller or service:

```go
func (a *AuthController) ResetPassword(ctx *gin.Context) {
    // ... validate user exists ...

    // Send asynchronously (don't block response)
    go a.emailService.SendPasswordResetEmail(
        ctx.Request.Context(),
        user.Email,
        user.Username,
    )

    // Return success (don't reveal if email exists)
    responseRes.Success(ctx, http.StatusOK, "If the email exists, a reset link has been sent", nil)
}
```

---

## 5. Using Pre-built Templates

The service includes ready-to-use email templates:

### Welcome Email

```go
err := s.emailService.SendWelcomeEmail(ctx, user.Email, user.Username)
```

### OTP Email (with styled HTML)

```go
err := s.emailService.SendOTPEmail(ctx, email, otp)
```

### Password Reset Email

```go
err := s.emailService.SendPasswordResetEmail(ctx, email, resetToken)
```

### Booking-related Emails

```go
// Booking confirmation
err := s.emailService.SendBookingConfirmationEmail(
    ctx, studentEmail, studentName,
    instructorName, "2024-01-15 10:00", "Defensive Driving",
)

// Lesson reminder
err := s.emailService.SendLessonReminderEmail(
    ctx, studentEmail, studentName,
    instructorName, "2024-01-16 10:00", "Highway Driving",
)

// Cancellation
err := s.emailService.SendLessonCancellationEmail(
    ctx, studentEmail, studentName,
    instructorName, "2024-01-15 10:00", "Instructor unavailable",
)
```

---

## 6. Custom Emails

For custom email content, use `SendEmail`:

```go
err := s.emailService.SendEmail(ctx, dto.SendEmailRequest{
    To: []dto.EmailAddress{{Email: "user@example.com"}},
    Subject: "Custom Subject",
    Text:    "Plain text version",
    HTML:    "<html><body>HTML version</body></html>",
    Tags:    []string{"custom", "transactional"},
})
```

### With Attachments

```go
err := s.emailService.SendEmail(ctx, dto.SendEmailRequest{
    To:      []dto.EmailAddress{{Email: "user@example.com"}},
    Subject: "Your Invoice",
    Text:    "Please find your invoice attached.",
    HTML:    "<html><body><p>Please find your invoice attached.</p></body></html>",
    Attachments: []dto.Attachment{
        {
            Content:     base64.StdEncoding.EncodeToString(pdfBytes),
            Filename:    "invoice.pdf",
            ContentType: "application/pdf",
        },
    },
})
```

---

## 7. Best Practices

### Always Send Asynchronously

```go
// DON'T block the main request
user, err := s.userRepo.Create(ctx, req)
if err != nil {
    return nil, err
}
s.emailService.SendWelcomeEmail(ctx, user.Email, user.Username) // Bad!

// DO send in goroutine
go func() {
    s.emailService.SendWelcomeEmail(context.Background(), user.Email, user.Username)
}()

// Even better - handle errors
go func() {
    if err := s.emailService.SendWelcomeEmail(context.Background(), user.Email, user.Username); err != nil {
        s.LogError("Failed to send welcome email", logger.LogField("error", err))
    }
}()
```

### Don't Reveal Email Existence

For security, always return the same message whether email exists or not:

```go
// Bad - reveals if email exists
if userNotFound {
    return "Email not found"
}

// Good - always same response
user, _ := s.userService.GetUserByEmail(ctx, email)
if user != nil {
    go s.emailService.SendPasswordResetEmail(ctx, user.Email, token)
}
return "If the email exists, a reset link has been sent" // Always same message
```

### Use Tags for Analytics

All emails use consistent tags for Mailtrap tracking:
- Welcome: `["welcome", "onboarding"]`
- OTP: `["otp", "email-verification"]`
- Password Reset: `["password-reset", "security"]`
- Booking: `["booking", "confirmation"]`

---

## 8. Testing Locally

For development, use Mailtrap's fake API key. Emails will be captured in Mailtrap's inbox without being delivered to real users.

```yaml
email:
  api_key: test_api_key  # Works with fake key for development
  from_email: "Admin Test"
  app_name: user_service
  enabled: true
```

---

## 9. Error Handling

The service handles errors internally:

```go
func (s *MailtrapEmailService) SendEmail(...) error {
    // 1. Context timeout protection
    emailCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 2. JSON marshaling errors
    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return apperrors.ErrInternalServer
    }

    // 3. HTTP request errors
    resp, err := s.client.Do(req)
    if err != nil {
        s.LogError("Failed to send email client", logger.LogField("error", err))
        return apperrors.ErrInternalServer
    }

    // 4. API response errors
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
        s.LogError("Failed to check response", logger.LogField("status", resp.StatusCode))
        return apperrors.ErrInternalServer
    }

    return nil
}
```

Callers should log and handle gracefully:

```go
go func() {
    if err := s.emailService.SendOTPEmail(ctx, email, otp); err != nil {
        s.LogError("Failed to send OTP email", logger.LogField("error", err))
        // Don't fail the main operation
    }
}()
```