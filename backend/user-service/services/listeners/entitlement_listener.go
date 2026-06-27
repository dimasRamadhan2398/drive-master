package listeners

import (
	"context"
	"fmt"
	"log"
	"time"

	"user-service/models"
	"user-service/models/dto"
	pkgKafka "user-service/pkg/kafka"

	"github.com/google/uuid"
)

// CertificationServiceInterface defines the certification service methods needed by the listener
type CertificationServiceInterface interface {
	IssueCertificate(ctx context.Context, input dto.IssueMemberCertificateInput) (*dto.IssueMemberCertificateResponse, error)
}

// CertificationRepositoryInterface defines the certification repository methods needed by the listener.
// FindByEntitlementID returns (nil, nil) when no certificate exists for the given entitlement yet.
type CertificationRepositoryInterface interface {
	FindByEntitlementID(ctx context.Context, entitlementID uuid.UUID) (*models.Certification, error)
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
	certRepo       CertificationRepositoryInterface
	eventPublisher EventPublisherInterface
}

// NewEntitlementCompletedListener creates a new listener instance.
// certRepo is used to check whether a certificate already exists for the entitlement
// before attempting issuance, making the handler fully idempotent.
func NewEntitlementCompletedListener(
	certService CertificationServiceInterface,
	certRepo CertificationRepositoryInterface,
	eventPublisher EventPublisherInterface,
) *EntitlementCompletedListener {
	return &EntitlementCompletedListener{
		certService:    certService,
		certRepo:       certRepo,
		eventPublisher: eventPublisher,
	}
}

// OnEntitlementCompleted is called when an entitlement is completed (remaining = 0).
// It is idempotent: if a certificate was already issued for this entitlement it skips
// creation and goes straight to publishing the Kafka event.
func (l *EntitlementCompletedListener) OnEntitlementCompleted(
	ctx context.Context,
	entitlement *models.Entitlement,
) error {
	// 1. Validate session integrity before doing anything
	if err := entitlement.ValidateSessionCount(); err != nil {
		return fmt.Errorf("entitlement integrity check failed: %w", err)
	}

	var certID uuid.UUID

	// 2. Check if a certificate already exists for this entitlement (via FK).
	//    FindByEntitlementID returns (nil, nil) when none is found.
	if l.certRepo != nil {
		existing, err := l.certRepo.FindByEntitlementID(ctx, entitlement.ID)
		if err != nil {
			// Log but continue — IssueCertificate is also idempotent; worst case it
			// returns the existing cert again.
			log.Printf("[WARN] failed to check existing certification for entitlement %s: %v",
				entitlement.ID, err)
		} else if existing != nil {
			log.Printf("[INFO] certificate %s already exists for entitlement %s, skipping issuance",
				existing.ID, entitlement.ID)
			certID = existing.ID
		}
	}

	// 3. Issue certification only if one does not exist yet.
	if certID == uuid.Nil {
		cert, err := l.certService.IssueCertificate(ctx, dto.IssueMemberCertificateInput{
			MemberID:      entitlement.MemberID,
			EntitlementID: entitlement.ID,
			PackageID:     entitlement.PackageID,
			PackageName:   entitlement.PackageName,
		})
		if err != nil {
			return fmt.Errorf("failed to issue certification: %w", err)
		}
		certID = cert.ID
	}

	// 4. Publish entitlement.completed event to Kafka for notification-service
	if l.eventPublisher != nil && l.eventPublisher.IsProducerEnabled() {
		err := l.eventPublisher.Publish(ctx, &pkgKafka.Event{
			Type:    "entitlement.completed",
			Success: true,
			Data: map[string]interface{}{
				"memberID":        entitlement.MemberID.String(),
				"entitlementID":   entitlement.ID.String(),
				"certificationID": certID.String(),
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
