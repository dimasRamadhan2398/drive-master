package listeners

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"user-service/models/dto"
	pkgKafka "user-service/pkg/kafka"
)

// EntitlementServiceInterface defines the entitlement service methods needed by the enrollment paid handler
type EntitlementServiceInterface interface {
	SyncEntitlementFromBooking(ctx context.Context, memberID, bookingID uuid.UUID, packageID uuid.UUID, packageName string, totalSessions int, isNightSession, isWeekendSession bool) (*dto.EntitlementResponse, error)
}

// EnrollmentPaidHandler listens for enrollment.paid events from booking-service
// and creates entitlements in user-service when payment is confirmed.
type EnrollmentPaidHandler struct {
	entitlementService EntitlementServiceInterface
}

// NewEnrollmentPaidHandler creates a new enrollment paid handler
func NewEnrollmentPaidHandler(entitlementService EntitlementServiceInterface) *EnrollmentPaidHandler {
	return &EnrollmentPaidHandler{
		entitlementService: entitlementService,
	}
}

// GetEventTypes returns the event types this handler cares about
func (h *EnrollmentPaidHandler) GetEventTypes() []pkgKafka.EventType {
	return []pkgKafka.EventType{
		pkgKafka.EventType("enrollment.paid"),
	}
}

// HandleEvent processes the incoming enrollment.paid event
func (h *EnrollmentPaidHandler) HandleEvent(ctx context.Context, event *pkgKafka.Event) error {
	if event.Type != "enrollment.paid" {
		return nil
	}

	if h.entitlementService == nil {
		log.Printf("[WARN] EnrollmentPaidHandler: entitlement service not available")
		return nil
	}

	// Extract data from event
	var memberID, enrollmentID, packageID uuid.UUID
	var totalSessions int
	var packageName string

	// Helper to parse uuid strings
	parseUUID := func(v interface{}) (uuid.UUID, error) {
		if v == nil {
			return uuid.Nil, fmt.Errorf("empty")
		}
		if s, ok := v.(string); ok {
			if s == "" {
				return uuid.Nil, fmt.Errorf("empty string")
			}
			return uuid.Parse(s)
		}
		return uuid.Nil, fmt.Errorf("invalid uuid")
	}

	// Extract member_id (user_id in event or event.UserID)
	if event.UserID != "" {
		if id, err := uuid.Parse(event.UserID); err == nil {
			memberID = id
		}
	}
	if memberID == uuid.Nil {
		if v, ok := event.Data["user_id"]; ok {
			if id, err := parseUUID(v); err == nil {
				memberID = id
			}
		}
	}
	if memberID == uuid.Nil {
		if v, ok := event.Data["member_id"]; ok {
			if id, err := parseUUID(v); err == nil {
				memberID = id
			}
		}
	}

	// Extract enrollment_id
	if v, ok := event.Data["enrollment_id"]; ok {
		if id, err := parseUUID(v); err == nil {
			enrollmentID = id
		}
	}

	// Extract package_id
	if v, ok := event.Data["package_id"]; ok {
		if id, err := parseUUID(v); err == nil {
			packageID = id
		}
	}

	// Extract total_sessions (may come from core-service or event)
	if v, ok := event.Data["total_sessions"]; ok {
		switch t := v.(type) {
		case float64:
			totalSessions = int(t)
		case int:
			totalSessions = t
		case int64:
			totalSessions = int(t)
		case string:
			// Try to parse string
			_ = t
		}
	}

	// Extract package_name (may come from event)
	if v, ok := event.Data["package_name"]; ok {
		if s, ok := v.(string); ok {
			packageName = s
		}
	}

	// Validate required fields
	if memberID == uuid.Nil {
		return fmt.Errorf("missing member_id in enrollment.paid event")
	}
	if enrollmentID == uuid.Nil {
		return fmt.Errorf("missing enrollment_id in enrollment.paid event")
	}
	if packageID == uuid.Nil {
		return fmt.Errorf("missing package_id in enrollment.paid event")
	}

	// If total_sessions is not provided in event, use a default
	// In production, this should be fetched from core-service
	if totalSessions == 0 {
		// Default sessions based on common packages
		// This should be replaced with a core-service call
		totalSessions = 10 // Default fallback
		log.Printf("[WARN] EnrollmentPaidHandler: total_sessions not provided in event, using default: %d", totalSessions)
	}

	// If package_name is not provided, use a placeholder
	if packageName == "" {
		packageName = fmt.Sprintf("Package %s", packageID.String()[:8])
		log.Printf("[WARN] EnrollmentPaidHandler: package_name not provided in event, using placeholder: %s", packageName)
	}

	isNightSession := false
	if v, ok := event.Data["is_night_session"]; ok {
		if b, ok := v.(bool); ok {
			isNightSession = b
		}
	}
	isWeekendSession := false
	if v, ok := event.Data["is_weekend_session"]; ok {
		if b, ok := v.(bool); ok {
			isWeekendSession = b
		}
	}

	resp, err := h.entitlementService.SyncEntitlementFromBooking(
		ctx,
		memberID,
		enrollmentID,
		packageID,
		packageName,
		totalSessions,
		isNightSession,
		isWeekendSession,
	)

	if err != nil {
		return fmt.Errorf("failed to create entitlement: %w", err)
	}

	log.Printf("[INFO] Entitlement created successfully: %s (remaining: %d, night: %t, weekend: %t)",
		resp.ID, resp.Remaining, resp.IsNightSession, resp.IsWeekendSession)

	return nil
}
