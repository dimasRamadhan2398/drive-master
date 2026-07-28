package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// EventType represents the type of auth event
type EventType string

const (
	// Authentication Events
	EventUserLogin       EventType = "user.login"
	EventUserLogout      EventType = "user.logout"
	EventUserLogoutAll   EventType = "user.logout_all"
	EventUserRegister    EventType = "user.register"
	EventUserLoginFailed EventType = "user.login_failed"
	EventTokenRefresh    EventType = "token.refresh"

	// Password Events
	EventPasswordChanged       EventType = "password.changed"
	EventPasswordReset         EventType = "password.reset_requested"
	EventPasswordResetComplete EventType = "password.reset_completed"

	// MFA Events
	EventMFAEnabled  EventType = "mfa.enabled"
	EventMFADisabled EventType = "mfa.disabled"
	EventMFAVerified EventType = "mfa.verified"
	EventMFAFailed   EventType = "mfa.failed"

	// Account Events
	EventAccountLocked      EventType = "account.locked"
	EventAccountUnlocked    EventType = "account.unlocked"
	EventAccountDeactivated EventType = "account.deactivated"
	EventAccountReactivated EventType = "account.reactivated"

	// Permission Events
	EventRoleAssigned     EventType = "role.assigned"
	EventRoleRevoked      EventType = "role.revoked"
	EventPermissionDenied EventType = "permission.denied"

	// Security Events
	EventSuspiciousActivity   EventType = "security.suspicious"
	EventRateLimitExceeded    EventType = "security.rate_limit"
	EventNewDeviceLogin       EventType = "security.new_device"
	EventTrustedDeviceAdded   EventType = "device.trusted_added"
	EventTrustedDeviceRemoved EventType = "device.trusted_removed"

	// Testimonial Events
	EventTestimonialCreated   EventType = "testimonial.created"
	EventTestimonialUpdated   EventType = "testimonial.updated"
	EventTestimonialDeleted   EventType = "testimonial.deleted"
	EventTestimonialPublished EventType = "testimonial.published"
	EventTestimonialArchived  EventType = "testimonial.archived"

	// Course Events
	EventCourseCompleted EventType = "course.completed"
	EventCourseUpdated   EventType = "course.updated"
	EventCourseDeleted   EventType = "course.deleted"
	EventCoursePublished EventType = "course.published"
	EventCourseArchived  EventType = "course.archived"

	// Session Events
	EventSessionCompleted EventType = "session.completed"

	// User Lifecycle Events
	EventUserDeleted EventType = "user.deleted"
)

// Event represents a generic auth event
type Event struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	Timestamp   time.Time              `json:"timestamp"`
	UserID      string                 `json:"user_id,omitempty"`
	Username    string                 `json:"username,omitempty"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Success     bool                   `json:"success"`
	ErrorReason string                 `json:"error_reason,omitempty"`
	Metadata    map[string]string      `json:"metadata,omitempty"`
}

// EventHandler defines the interface for handling events
type EventHandler interface {
	HandleEvent(ctx context.Context, event *Event) error
	GetEventTypes() []EventType
}

// IEventPublisher defines the interface for event publishing
type IEventPublisher interface {
	// Publishing
	Publish(ctx context.Context, event *Event) error
	PublishLogin(ctx context.Context, userID, username, ip, userAgent string, success bool) error
	PublishLogout(ctx context.Context, userID, username, ip, userAgent string) error
	PublishLogoutAll(ctx context.Context, userID, username string) error
	PublishLoginFailed(ctx context.Context, username, ip, reason string) error
	PublishPasswordChanged(ctx context.Context, userID, username, ip string) error
	PublishPasswordReset(ctx context.Context, userID, username, email string) error
	PublishMFAEvent(ctx context.Context, userID, username string, eventType EventType, success bool, ip string) error
	PublishAccountEvent(ctx context.Context, userID, username string, eventType EventType) error
	PublishSecurityEvent(ctx context.Context, userID, username string, eventType EventType, ip string, data map[string]interface{}) error
	PublishRoleEvent(ctx context.Context, userID, username, role string, eventType EventType) error

	// Testimonial publishing
	PublishTestimonialCreated(ctx context.Context, testimonialID string, userID, userName string, rating float64, status string, isFeatured bool) error
	PublishTestimonialUpdated(ctx context.Context, testimonialID string, userName string, content string, rating float64, status string) error
	PublishTestimonialDeleted(ctx context.Context, testimonialID string) error
	PublishTestimonialPublished(ctx context.Context, testimonialID, publishedBy string) error
	PublishTestimonialArchived(ctx context.Context, testimonialID, archivedBy string) error

	// User lifecycle publishing
	PublishUserDeleted(ctx context.Context, userID string, username, email string) error

	// Consumer management
	StartConsumer(ctx context.Context) error
	StopConsumer() error
	RegisterHandler(handler EventHandler) error

	// Utilities
	IsConsumerEnabled() bool
	IsProducerEnabled() bool
}

// EventPublisher implements IEventPublisher
type EventPublisher struct {
	producer *Producer
	consumer *Consumer
	topic    string
	handlers []EventHandler
	enabled  bool
}

// EventPublisherConfig holds configuration for the event Publisher
type EventPublisherConfig struct {
	Producer *Producer
	Consumer *Consumer
	Topic    string
	Enabled  bool
}

// NewEventPublisher creates a new event Publisher
func NewEventPublisher(cfg EventPublisherConfig) IEventPublisher {
	return &EventPublisher{
		producer: cfg.Producer,
		consumer: cfg.Consumer,
		topic:    cfg.Topic,
		enabled:  cfg.Enabled,
		handlers: make([]EventHandler, 0),
	}
}

// Publish publishes an event to Kafka
func (r *EventPublisher) Publish(ctx context.Context, event *Event) error {
	if !r.enabled || r.producer == nil {
		return nil
	}

	if event.ID == "" {
		event.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Use the producer's SendLogSync for reliable delivery
	return r.producer.SendLogSync(ctx, "info", string(event.Type), map[string]interface{}{
		"event":      event,
		"event_type": event.Type,
		"event_data": string(data),
		"log_level":  getLogLevel(event.Type),
	})
}

// PublishLogin publishes a login event
func (r *EventPublisher) PublishLogin(ctx context.Context, userID, username, ip, userAgent string, success bool) error {
	event := &Event{
		Type:      EventUserLogin,
		UserID:    userID,
		Username:  username,
		IPAddress: ip,
		UserAgent: userAgent,
		Success:   success,
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishLogout publishes a logout event
func (r *EventPublisher) PublishLogout(ctx context.Context, userID, username, ip, userAgent string) error {
	event := &Event{
		Type:      EventUserLogout,
		UserID:    userID,
		Username:  username,
		IPAddress: ip,
		UserAgent: userAgent,
		Success:   true,
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishLogoutAll publishes a logout all devices event
func (r *EventPublisher) PublishLogoutAll(ctx context.Context, userID, username string) error {
	event := &Event{
		Type:      EventUserLogoutAll,
		UserID:    userID,
		Username:  username,
		Success:   true,
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishLoginFailed publishes a failed login attempt event
func (r *EventPublisher) PublishLoginFailed(ctx context.Context, username, ip, reason string) error {
	event := &Event{
		Type:        EventUserLoginFailed,
		Username:    username,
		IPAddress:   ip,
		Success:     false,
		ErrorReason: reason,
		Timestamp:   time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishPasswordChanged publishes a password changed event
func (r *EventPublisher) PublishPasswordChanged(ctx context.Context, userID, username, ip string) error {
	event := &Event{
		Type:      EventPasswordChanged,
		UserID:    userID,
		Username:  username,
		IPAddress: ip,
		Success:   true,
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishPasswordReset publishes a password reset request event
func (r *EventPublisher) PublishPasswordReset(ctx context.Context, userID, username, email string) error {
	event := &Event{
		Type:     EventPasswordReset,
		UserID:   userID,
		Username: username,
		Success:  true,
		Data: map[string]interface{}{
			"email": email,
		},
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishMFAEvent publishes an MFA-related event
func (r *EventPublisher) PublishMFAEvent(ctx context.Context, userID, username string, eventType EventType, success bool, ip string) error {
	event := &Event{
		Type:      eventType,
		UserID:    userID,
		Username:  username,
		IPAddress: ip,
		Success:   success,
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishAccountEvent publishes an account-related event
func (r *EventPublisher) PublishAccountEvent(ctx context.Context, userID, username string, eventType EventType) error {
	event := &Event{
		Type:      eventType,
		UserID:    userID,
		Username:  username,
		Success:   true,
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishSecurityEvent publishes a security-related event
func (r *EventPublisher) PublishSecurityEvent(ctx context.Context, userID, username string, eventType EventType, ip string, data map[string]interface{}) error {
	event := &Event{
		Type:      eventType,
		UserID:    userID,
		Username:  username,
		IPAddress: ip,
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// PublishRoleEvent publishes a role assignment/revocation event
func (r *EventPublisher) PublishRoleEvent(ctx context.Context, userID, username, role string, eventType EventType) error {
	event := &Event{
		Type:     eventType,
		UserID:   userID,
		Username: username,
		Success:  true,
		Data: map[string]interface{}{
			"role": role,
		},
		Timestamp: time.Now().UTC(),
	}
	return r.Publish(ctx, event)
}

// ========== TESTIMONIAL EVENTS ==========

// PublishTestimonialCreated publishes a testimonial created event
func (r *EventPublisher) PublishTestimonialCreated(ctx context.Context, testimonialID string, userID, userName string, rating float64, status string, isFeatured bool) error {
	event := &Event{
		Type:      EventTestimonialCreated,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"testimonial_id": testimonialID,
			"user_id":        userID,
			"user_name":      userName,
			"rating":         rating,
			"status":         status,
			"is_featured":    isFeatured,
		},
		Success: true,
	}
	return r.Publish(ctx, event)
}

// PublishTestimonialUpdated publishes a testimonial updated event
func (r *EventPublisher) PublishTestimonialUpdated(ctx context.Context, testimonialID string, userName string, content string, rating float64, status string) error {
	event := &Event{
		Type:      EventTestimonialUpdated,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"testimonial_id": testimonialID,
			"user_name":      userName,
			"content":        content,
			"rating":         rating,
			"status":         status,
		},
		Success: true,
	}
	return r.Publish(ctx, event)
}

// PublishTestimonialDeleted publishes a testimonial deleted event
func (r *EventPublisher) PublishTestimonialDeleted(ctx context.Context, testimonialID string) error {
	event := &Event{
		Type:      EventTestimonialDeleted,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"testimonial_id": testimonialID,
			"deleted_at":     time.Now().Format(time.RFC3339),
		},
		Success: true,
	}
	return r.Publish(ctx, event)
}

// PublishTestimonialPublished publishes a testimonial published event
func (r *EventPublisher) PublishTestimonialPublished(ctx context.Context, testimonialID, publishedBy string) error {
	event := &Event{
		Type:      EventTestimonialPublished,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"testimonial_id": testimonialID,
			"published_by":   publishedBy,
			"published_at":   time.Now().Format(time.RFC3339),
		},
		Success: true,
	}
	return r.Publish(ctx, event)
}

// PublishTestimonialArchived publishes a testimonial archived event
func (r *EventPublisher) PublishTestimonialArchived(ctx context.Context, testimonialID, archivedBy string) error {
	event := &Event{
		Type:      EventTestimonialArchived,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"testimonial_id": testimonialID,
			"archived_by":    archivedBy,
			"archived_at":    time.Now().Format(time.RFC3339),
		},
		Success: true,
	}
	return r.Publish(ctx, event)
}

// ========== USER LIFECYCLE EVENTS ==========

// PublishUserDeleted publishes a user deleted event
// This event is consumed by other services (booking, payment, notification, certificate)
// to anonymize or clean up user-related data while preserving transactional records
func (r *EventPublisher) PublishUserDeleted(ctx context.Context, userID string, username, email string) error {
	event := &Event{
		Type:      EventUserDeleted,
		UserID:    userID,
		Username:  username,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"email":      email,
			"deleted_at": time.Now().Format(time.RFC3339),
		},
		Success: true,
	}
	return r.Publish(ctx, event)
}

// StartConsumer starts the event consumer
func (r *EventPublisher) StartConsumer(ctx context.Context) error {
	if r.consumer == nil || !r.consumer.IsEnabled() {
		return nil
	}

	handler := NewJSONMessageHandler(func(ctx context.Context, msg LogMessage) error {
		if eventData, ok := msg.Fields["event"]; ok {
			if eventBytes, err := json.Marshal(eventData); err == nil {
				var event Event
				if err := json.Unmarshal(eventBytes, &event); err == nil {
					for _, h := range r.handlers {
						if err := h.HandleEvent(ctx, &event); err != nil {
							// Log error but don't fail the consumer
						}
					}
				}
			}
		}
		return nil
	})

	for _, topic := range r.consumer.topics {
		r.consumer.RegisterHandler(topic, handler)
	}

	return r.consumer.Start()
}

// StopConsumer stops the event consumer
func (r *EventPublisher) StopConsumer() error {
	if r.consumer == nil {
		return nil
	}
	return r.consumer.Stop()
}

// RegisterHandler registers an event handler
func (r *EventPublisher) RegisterHandler(handler EventHandler) error {
	r.handlers = append(r.handlers, handler)
	return nil
}

// IsConsumerEnabled returns whether the consumer is enabled
func (r *EventPublisher) IsConsumerEnabled() bool {
	return r.consumer != nil && r.consumer.IsEnabled()
}

// IsProducerEnabled returns whether the producer is enabled
func (r *EventPublisher) IsProducerEnabled() bool {
	return r.producer != nil && r.producer.IsEnabled()
}

// wrappedEventHandler wraps EventHandler for Kafka consumer
type wrappedEventHandler struct {
	handlers []EventHandler
}

func (h *wrappedEventHandler) Handle(ctx context.Context, message []byte) error {
	var event Event
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	for _, handler := range h.handlers {
		if err := handler.HandleEvent(ctx, &event); err != nil {
			return err
		}
	}
	return nil
}

// ============================================
// Pre-built Event Handlers
// ============================================

// UserActivityHandler logs user activity events
type UserActivityHandler struct {
	logger func(msg string, fields map[string]interface{})
}

// NewUserActivityHandler creates a new user activity handler
func NewUserActivityHandler(logger func(msg string, fields map[string]interface{})) *UserActivityHandler {
	return &UserActivityHandler{logger: logger}
}

// HandleEvent handles user activity events
func (h *UserActivityHandler) HandleEvent(ctx context.Context, event *Event) error {
	if h.logger == nil {
		return nil
	}

	msg := fmt.Sprintf("User activity: %s", event.Type)
	fields := map[string]interface{}{
		"user_id":   event.UserID,
		"username":  event.Username,
		"ip":        event.IPAddress,
		"success":   event.Success,
		"timestamp": event.Timestamp,
	}
	h.logger(msg, fields)
	return nil
}

// GetEventTypes returns the event types this handler handles
func (h *UserActivityHandler) GetEventTypes() []EventType {
	return []EventType{
		EventUserLogin,
		EventUserLogout,
		EventUserLogoutAll,
		EventUserRegister,
	}
}

// SecurityAlertHandler handles security-related events
type SecurityAlertHandler struct {
	logger func(msg string, fields map[string]interface{})
}

// NewSecurityAlertHandler creates a new security alert handler
func NewSecurityAlertHandler(logger func(msg string, fields map[string]interface{})) *SecurityAlertHandler {
	return &SecurityAlertHandler{logger: logger}
}

// HandleEvent handles security events
func (h *SecurityAlertHandler) HandleEvent(ctx context.Context, event *Event) error {
	if h.logger == nil {
		return nil
	}

	msg := fmt.Sprintf("Security event: %s", event.Type)
	fields := map[string]interface{}{
		"user_id":      event.UserID,
		"username":     event.Username,
		"ip":           event.IPAddress,
		"error_reason": event.ErrorReason,
		"data":         event.Data,
		"timestamp":    event.Timestamp,
	}
	h.logger(msg, fields)
	return nil
}

// GetEventTypes returns the event types this handler handles
func (h *SecurityAlertHandler) GetEventTypes() []EventType {
	return []EventType{
		EventUserLoginFailed,
		EventAccountLocked,
		EventSuspiciousActivity,
		EventRateLimitExceeded,
		EventMFAFailed,
	}
}

// AuditLogHandler exports events as audit logs
type AuditLogHandler struct {
	exporter func(event *Event) error
}

// NewAuditLogHandler creates a new audit log handler
func NewAuditLogHandler(exporter func(event *Event) error) *AuditLogHandler {
	return &AuditLogHandler{exporter: exporter}
}

// HandleEvent handles all events for audit logging
func (h *AuditLogHandler) HandleEvent(ctx context.Context, event *Event) error {
	if h.exporter == nil {
		return nil
	}
	return h.exporter(event)
}

// GetEventTypes returns all event types for audit logging
func (h *AuditLogHandler) GetEventTypes() []EventType {
	return []EventType{
		EventUserLogin,
		EventUserLogout,
		EventUserLogoutAll,
		EventUserRegister,
		EventUserLoginFailed,
		EventPasswordChanged,
		EventPasswordReset,
		EventPasswordResetComplete,
		EventMFAEnabled,
		EventMFADisabled,
		EventMFAVerified,
		EventMFAFailed,
		EventAccountLocked,
		EventAccountUnlocked,
		EventAccountDeactivated,
		EventAccountReactivated,
		EventRoleAssigned,
		EventRoleRevoked,
		EventPermissionDenied,
		EventSuspiciousActivity,
		EventRateLimitExceeded,
		EventNewDeviceLogin,
		EventTrustedDeviceAdded,
		EventTrustedDeviceRemoved,
	}
}

// getLogLevel returns the appropriate log level for an event type
func getLogLevel(eventType EventType) string {
	switch eventType {
	case EventUserLoginFailed, EventAccountLocked, EventSuspiciousActivity,
		EventRateLimitExceeded, EventMFAFailed, EventPermissionDenied:
		return "warn"
	default:
		return "info"
	}
}
