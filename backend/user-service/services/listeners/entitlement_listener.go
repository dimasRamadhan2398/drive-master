package listeners

import (
	"context"
	"fmt"
	"log"
	"time"

	"user-service/models"
	"user-service/models/dto"
	pkgKafka "user-service/pkg/kafka"
)

// CertificationServiceInterface defines the certification service methods needed by the listener
type CertificationServiceInterface interface {
	IssueCertification(ctx context.Context, input dto.IssueCertificationInput) (*dto.CertificationResponse, error)
}

// EventPublisherInterface defines the event publisher methods needed by the listener
type EventPublisherInterface interface {
	IsProducerEnabled() bool
	Publish(ctx context.Context, event *pkgKafka.Event) error
}

type IEntitlementCompletedListener interface {
	OnEntitlementCompleted(ctx context.Context, entitlement *models.Entitlement) error
}

// EntitlementCompletedListener handles actions when an entitlement is completed
// (all sessions used). It issues a certification and publishes an event.
type EntitlementCompletedListener struct {
	certService    CertificationServiceInterface
	eventPublisher EventPublisherInterface
}

// NewEntitlementCompletedListener creates a new listener instance
func NewEntitlementCompletedListener(
	certService CertificationServiceInterface,
	eventPublisher EventPublisherInterface,
) *EntitlementCompletedListener {
	return &EntitlementCompletedListener{
		certService:    certService,
		eventPublisher: eventPublisher,
	}
}

// OnEntitlementCompleted is called when an entitlement is completed (remaining = 0)
func (l *EntitlementCompletedListener) OnEntitlementCompleted(
	ctx context.Context,
	entitlement *models.Entitlement,
) error {
	// 1. Validate session integrity before doing anything
	if err := entitlement.ValidateSessionCount(); err != nil {
		return fmt.Errorf("entitlement integrity check failed: %w", err)
	}

	// 2. Issue certification for the completed package
	cert, err := l.certService.IssueCertification(ctx, dto.IssueCertificationInput{
		MemberID:      entitlement.MemberID,
		EntitlementID: entitlement.ID,
		PackageID:     entitlement.PackageID,
		PackageName:   entitlement.PackageName,
		IssuedAt:      time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to issue certification: %w", err)
	}

	// 3. Publish entitlement.completed event to Kafka for notification-service
	if l.eventPublisher != nil && l.eventPublisher.IsProducerEnabled() {
		err = l.eventPublisher.Publish(ctx, &pkgKafka.Event{
			Type:    "entitlement.completed",
			Success: true,
			Data: map[string]interface{}{
				"memberID":        entitlement.MemberID.String(),
				"entitlementID":   entitlement.ID.String(),
				"certificationID": cert.ID.String(),
				"packageID":       entitlement.PackageID.String(),
				"packageName":     entitlement.PackageName,
				"completedAt":     time.Now().Format(time.RFC3339),
			},
			Timestamp: time.Now(),
		})
		if err != nil {
			// Non-fatal — certification was issued, notification can be retried
			log.Printf("[WARN] failed to publish entitlement.completed event: %v", err)
		}
	}

	return nil
}
