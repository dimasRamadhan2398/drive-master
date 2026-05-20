package services

import (
	"context"
	"errors"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/repositories"

	"gorm.io/gorm"
)

type EnrollmentService struct {
	enrollmentRepo  *repositories.EnrollmentRepository
	entitlementRepo *repositories.EntitlementRepository
}

func NewEnrollmentService(
	enrollmentRepo *repositories.EnrollmentRepository,
	entitlementRepo *repositories.EntitlementRepository,
) IEnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:  enrollmentRepo,
		entitlementRepo: entitlementRepo,
	}
}

func (s *EnrollmentService) CreateEnrollment(ctx context.Context, req dto.CreateEnrollmentRequest) (*dto.EnrollmentResponse, error) {
	// Create enrollment with pending_payment status
	enrollment := &models.Enrollment{
		UserID:    req.UserID,
		PackageID: req.PackageID,
		Status:    models.EnrollmentStatusPendingPayment,
		ExpiresAt: time.Now().AddDate(1, 0, 0), // Default 1 year validity
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToResponse(enrollment)
	return &resp, nil
}

func (s *EnrollmentService) GetEnrollment(ctx context.Context, id uint) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	resp := s.enrollmentRepo.ToResponse(enrollment)
	return &resp, nil
}

func (s *EnrollmentService) UpdateEnrollment(ctx context.Context, id uint, req dto.UpdateEnrollmentRequest) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	if req.ExpiresAt != nil {
		enrollment.ExpiresAt = *req.ExpiresAt
	}
	if req.Status != nil {
		enrollment.Status = models.EnrollmentStatus(*req.Status)
	}

	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToResponse(enrollment)
	return &resp, nil
}

func (s *EnrollmentService) CancelEnrollment(ctx context.Context, id uint) error {
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}
		return err
	}

	if enrollment.Status == models.EnrollmentStatusCancelled {
		return errors.New("enrollment already cancelled")
	}

	if enrollment.Status == models.EnrollmentStatusCompleted {
		return errors.New("cannot cancel a completed enrollment")
	}

	return s.enrollmentRepo.UpdateStatus(ctx, id, models.EnrollmentStatusCancelled)
}

func (s *EnrollmentService) MarkAsPaid(ctx context.Context, id uint, totalPrice float64) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	if enrollment.Status != models.EnrollmentStatusPendingPayment {
		return nil, errors.New("enrollment is not in pending_payment status")
	}

	now := time.Now()
	if err := s.enrollmentRepo.MarkAsPaid(ctx, id, now, totalPrice); err != nil {
		return nil, err
	}

	// Reload enrollment
	enrollment, err = s.enrollmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToResponse(enrollment)
	return &resp, nil
}

func (s *EnrollmentService) ListEnrollments(ctx context.Context, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, total, err := s.enrollmentRepo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}

func (s *EnrollmentService) ListUserEnrollments(ctx context.Context, userID uint, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, total, err := s.enrollmentRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}

func (s *EnrollmentService) ListEnrollmentsByStatus(ctx context.Context, status string, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, total, err := s.enrollmentRepo.GetByStatus(ctx, models.EnrollmentStatus(status), page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}

func (s *EnrollmentService) CreateEntitlementFromEnrollment(ctx context.Context, enrollmentID uint, sourceType, sourceID string, totalSessions int, expiresAt time.Time) (*dto.EntitlementResponse, error) {
	enrollment, err := s.enrollmentRepo.GetByID(ctx, enrollmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	entitlement := &models.UserEntitlement{
		EnrollmentID:  enrollmentID,
		UserID:        enrollment.UserID,
		SourceType:    sourceType,
		SourceID:      sourceID,
		TotalSessions: totalSessions,
		UsedSessions:  0,
		ExpiresAt:     expiresAt,
	}

	if err := s.entitlementRepo.Create(ctx, entitlement); err != nil {
		return nil, err
	}

	resp := s.entitlementRepo.ToResponse(entitlement)
	return &resp, nil
}