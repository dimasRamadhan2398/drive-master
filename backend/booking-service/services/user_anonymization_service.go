package services

import (
	"context"
	"time"

	"booking-service/repositories"

	"github.com/google/uuid"
)

// IUserAnonymizationService handles anonymization of user data when a user is deleted
type IUserAnonymizationService interface {
	AnonymizeUserData(ctx context.Context, userID string) error
}

type UserAnonymizationService struct {
	enrollmentRepo  repositories.IEnrollmentRepository
	sessionRepo     repositories.ISessionRepository
	entitlementRepo repositories.IEntitlementRepository
}

// NewUserAnonymizationService creates a new user anonymization service
func NewUserAnonymizationService(
	enrollmentRepo repositories.IEnrollmentRepository,
	sessionRepo repositories.ISessionRepository,
	entitlementRepo repositories.IEntitlementRepository,
) IUserAnonymizationService {
	return &UserAnonymizationService{
		enrollmentRepo:   enrollmentRepo,
		sessionRepo:      sessionRepo,
		entitlementRepo: entitlementRepo,
	}
}

// AnonymizeUserData marks all user-related records as anonymized when a user is deleted
func (s *UserAnonymizationService) AnonymizeUserData(ctx context.Context, userID string) error {
	now := time.Now()

	// Parse user ID from string (UUID format from user-service)
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		// If not a valid UUID, log and skip
		return nil
	}

	// Anonymize enrollments
	if err := s.enrollmentRepo.AnonymizeByUserID(ctx, parsedUserID, now); err != nil {
		return err
	}

	// Anonymize sessions
	if err := s.sessionRepo.AnonymizeByUserID(ctx, parsedUserID, now); err != nil {
		return err
	}

	// Anonymize entitlements
	if err := s.entitlementRepo.AnonymizeByUserID(ctx, parsedUserID, now); err != nil {
		return err
	}

	return nil
}