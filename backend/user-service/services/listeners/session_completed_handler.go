package listeners

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"user-service/models"
	pkgKafka "user-service/pkg/kafka"
	"user-service/repositories"
)

// SessionCompletedHandler listens for session/course completion events and
// updates entitlements accordingly. When an entitlement reaches 0 remaining
// sessions, it triggers the EntitlementCompletedListener.
type SessionCompletedHandler struct {
	entRepo  repositories.IEntitlementRepository
	listener IEntitlementCompletedListener
}

func NewSessionCompletedHandler(entRepo repositories.IEntitlementRepository, listener IEntitlementCompletedListener) *SessionCompletedHandler {
	return &SessionCompletedHandler{entRepo: entRepo, listener: listener}
}

// GetEventTypes returns the event types this handler cares about
func (h *SessionCompletedHandler) GetEventTypes() []pkgKafka.EventType {
	return []pkgKafka.EventType{
		pkgKafka.EventType("session.completed"),
		pkgKafka.EventType("course.completed"),
	}
}

// HandleEvent processes the incoming event
func (h *SessionCompletedHandler) HandleEvent(ctx context.Context, event *pkgKafka.Event) error {
	// Extract useful fields from event.Data
	var entitlementID uuid.UUID
	var bookingID uuid.UUID
	var memberID uuid.UUID
	var packageID uuid.UUID
	var sessionsUsed int = 1

	// helper to parse uuid strings
	parseUUID := func(v interface{}) (uuid.UUID, error) {
		if v == nil {
			return uuid.Nil, fmt.Errorf("empty")
		}
		if s, ok := v.(string); ok {
			return uuid.Parse(s)
		}
		return uuid.Nil, fmt.Errorf("invalid uuid")
	}

	if v, ok := event.Data["entitlementID"]; ok {
		if id, err := parseUUID(v); err == nil {
			entitlementID = id
		}
	}
	if v, ok := event.Data["bookingID"]; ok {
		if id, err := parseUUID(v); err == nil {
			bookingID = id
		}
	}
	if v, ok := event.Data["memberID"]; ok {
		if id, err := parseUUID(v); err == nil {
			memberID = id
		}
	}
	if v, ok := event.Data["packageID"]; ok {
		if id, err := parseUUID(v); err == nil {
			packageID = id
		}
	}
	if v, ok := event.Data["sessions_used"]; ok {
		switch t := v.(type) {
		case float64:
			sessionsUsed = int(t)
		case int:
			sessionsUsed = t
		case int64:
			sessionsUsed = int(t)
		case string:
			// try parse
			// ignore error, default to 1
			// fallthrough
			_ = t
		}
	}

	// Find entitlement
	var ent *models.Entitlement
	var err error
	if entitlementID != uuid.Nil {
		ent, err = h.entRepo.FindByID(ctx, entitlementID)
		if err != nil {
			return fmt.Errorf("entitlement not found: %w", err)
		}
	} else if bookingID != uuid.Nil {
		ent, err = h.entRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			return fmt.Errorf("entitlement not found by booking: %w", err)
		}
	} else if memberID != uuid.Nil && packageID != uuid.Nil {
		// Search by member entitlements and package
		ents, _, err := h.entRepo.FindByMemberID(ctx, memberID, 1, 50)
		if err != nil {
			return fmt.Errorf("failed to query entitlements: %w", err)
		}
		for _, e := range ents {
			if e.PackageID == packageID {
				ent = &e
				break
			}
		}
		if ent == nil {
			return fmt.Errorf("entitlement not found for member+package")
		}
	} else {
		return fmt.Errorf("insufficient data to locate entitlement")
	}

	// If entitlement already used, nothing to do
	if ent.Remaining <= 0 {
		return nil
	}

	// Decrement remaining sessions (atomic decrement in repo)
	for i := 0; i < sessionsUsed; i++ {
		if err := h.entRepo.DecrementRemaining(ctx, ent.ID); err != nil {
			// If decrement fails, continue to next to avoid partial failures causing inconsistencies
			return fmt.Errorf("failed to decrement entitlement: %w", err)
		}
	}

	// Refresh
	ent, err = h.entRepo.FindByID(ctx, ent.ID)
	if err != nil {
		return fmt.Errorf("failed to refresh entitlement: %w", err)
	}

	// If no remaining sessions, mark used and trigger completion listener
	if ent.Remaining == 0 && ent.Status != models.EntitlementStatusUsed {
		ent.Status = models.EntitlementStatusUsed
		ent.UpdatedAt = time.Now()
		if err := h.entRepo.Update(ctx, ent); err != nil {
			return fmt.Errorf("failed to update entitlement status: %w", err)
		}

		if h.listener != nil {
			if err := h.listener.OnEntitlementCompleted(ctx, ent); err != nil {
				// log and continue
				return fmt.Errorf("failed to trigger entitlement completion: %w", err)
			}
		}
	}

	return nil
}
