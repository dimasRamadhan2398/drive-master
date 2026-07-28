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

// ProfileUpdateListener updates member profiles when enrollment/entitlement or session events occur
// Expected event types:
// - "entitlement.created" with Data: {"memberID":"<uuid>", "entitlementID":"<uuid>", "packageID":"<uuid>", "packageName":"...", "totalSessions":10}
// - "session.completed" with Data: {"memberID":"<uuid>", "entitlementID":"<uuid>", "duration_minutes":45, "rating":4.5}
// - "course.completed" behaves similar to session.completed
type ProfileUpdateListener struct {
	memberRepo repositories.IMemberRepository
	entRepo    repositories.IEntitlementRepository
}

func NewProfileUpdateListener(memberRepo repositories.IMemberRepository, entRepo repositories.IEntitlementRepository) *ProfileUpdateListener {
	return &ProfileUpdateListener{memberRepo: memberRepo, entRepo: entRepo}
}

func (h *ProfileUpdateListener) GetEventTypes() []pkgKafka.EventType {
	return []pkgKafka.EventType{
		pkgKafka.EventType("entitlement.created"),
		pkgKafka.EventType("session.completed"),
		pkgKafka.EventType("course.completed"),
	}
}

func (h *ProfileUpdateListener) HandleEvent(ctx context.Context, event *pkgKafka.Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	switch string(event.Type) {
	case "entitlement.created":
		return h.handleEntitlementCreated(ctx, event)
	case "session.completed", "course.completed":
		return h.handleSessionCompleted(ctx, event)
	default:
		// ignore other events
		return nil
	}
}

func (h *ProfileUpdateListener) handleEntitlementCreated(ctx context.Context, event *pkgKafka.Event) error {
	memberIDStr, _ := event.Data["memberID"].(string)
	if memberIDStr == "" {
		// fallback to user_id
		memberIDStr = event.UserID
	}
	if memberIDStr == "" {
		return fmt.Errorf("entitlement.created: missing memberID")
	}
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		return fmt.Errorf("invalid memberID: %w", err)
	}

	// Recalculate total available sessions for the member from entitlements
	total := 0
	ents, _, err := h.entRepo.FindByMemberID(ctx, memberID, 1, 1000)
	if err != nil {
		// If error, still return error so caller can log/handle
		return fmt.Errorf("failed to fetch entitlements for member %s: %w", memberID, err)
	}
	for _, e := range ents {
		if e.Status == models.EntitlementStatusActive {
			total += e.Remaining
		}
	}

	// Update member profile total available sessions
	if err := h.memberRepo.UpdateTotalAvailableSessions(ctx, memberID, total); err != nil {
		return fmt.Errorf("failed to update member total available sessions: %w", err)
	}

	return nil
}

func (h *ProfileUpdateListener) handleSessionCompleted(ctx context.Context, event *pkgKafka.Event) error {
	memberIDStr, _ := event.Data["memberID"].(string)
	if memberIDStr == "" {
		memberIDStr = event.UserID
	}
	if memberIDStr == "" {
		return fmt.Errorf("session.completed: missing memberID")
	}
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		return fmt.Errorf("invalid memberID: %w", err)
	}

	// duration
	durMinutes := 0
	if v, ok := event.Data["duration_minutes"]; ok {
		switch t := v.(type) {
		case float64:
			durMinutes = int(t)
		case int:
			durMinutes = t
		case int64:
			durMinutes = int(t)
		}
	}

	// rating (optional)
	rating := 0.0
	if v, ok := event.Data["rating"]; ok {
		switch t := v.(type) {
		case float64:
			rating = t
		case int:
			rating = float64(t)
		case string:
			// ignore string ratings
		}
	}

	// Load member profile
	profile, err := h.memberRepo.FindByUserID(ctx, memberID)
	if err != nil {
		return fmt.Errorf("failed to load member profile: %w", err)
	}

	// Update fields
	oldSessions := profile.SessionsCompleted
	profile.SessionsCompleted = profile.SessionsCompleted + 1
	profile.TrainingTime = profile.TrainingTime + durMinutes

	if rating > 0 {
		// Recompute average using sessions as count
		if profile.SessionsCompleted == 0 {
			profile.AverageRating = rating
		} else {
			// oldSessions is sessions before increment
			if oldSessions <= 0 {
				profile.AverageRating = rating
			} else {
				profile.AverageRating = ((profile.AverageRating * float64(oldSessions)) + rating) / float64(oldSessions+1)
			}
		}
	}

	profile.UpdatedAt = time.Now()

	if err := h.memberRepo.Update(ctx, profile); err != nil {
		return fmt.Errorf("failed to update member profile: %w", err)
	}

	// Recalculate total available sessions as some session consumptions may affect entitlements
	total := 0
	ents, _, err := h.entRepo.FindByMemberID(ctx, memberID, 1, 1000)
	if err == nil {
		for _, e := range ents {
			if e.Status == models.EntitlementStatusActive {
				total += e.Remaining
			}
		}
		_ = h.memberRepo.UpdateTotalAvailableSessions(ctx, memberID, total)
	}

	return nil
}
