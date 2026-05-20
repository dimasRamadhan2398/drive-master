package services

import (
	"context"
	"errors"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/repositories"

	"gorm.io/gorm"
)

type BookingService struct {
	bookingRepo     *repositories.BookingRepository
	entitlementRepo *repositories.EntitlementRepository
}

func NewBookingService(
	bookingRepo *repositories.BookingRepository,
	entitlementRepo *repositories.EntitlementRepository,
) IBookingService {
	return &BookingService{
		bookingRepo:     bookingRepo,
		entitlementRepo: entitlementRepo,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (*dto.EnrollmentResponse, error) {
	// Check if entitlement exists and has remaining sessions
	entitlement, err := s.entitlementRepo.GetByID(ctx, req.EntitlementID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entitlement not found")
		}
		return nil, err
	}

	// Calculate remaining sessions (TotalSessions - UsedSessions)
	remainingSessions := entitlement.TotalSessions - entitlement.UsedSessions
	if remainingSessions <= 0 {
		return nil, errors.New("no sessions remaining in entitlement")
	}

	if entitlement.ExpiresAt.Before(req.DateOfSession) {
		return nil, errors.New("entitlement has expired")
	}

	// Create enrollment (booking is now an enrollment)
	enrollment := &models.Enrollment{
		UserID:    req.UserID,
		PackageID: 0, // Not applicable for direct booking
		Status:    models.EnrollmentStatusInProgress,
		ExpiresAt: entitlement.ExpiresAt,
	}

	if err := s.bookingRepo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToResponse(enrollment)
	return &resp, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id uint) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	resp := s.bookingRepo.ToResponse(enrollment)
	return &resp, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, id uint, req dto.UpdateBookingRequest) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	// UpdateBookingRequest only has: DateOfSession, FromTime, ToTime, CarID, Area, Notes, Status
	// For enrollment updates, we use Status field
	if req.Status != nil {
		enrollment.Status = models.EnrollmentStatus(*req.Status)
	}

	if err := s.bookingRepo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToResponse(enrollment)
	return &resp, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, id uint) error {
	enrollment, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	if enrollment.Status == models.EnrollmentStatusCancelled {
		return errors.New("booking already cancelled")
	}

	if enrollment.Status == models.EnrollmentStatusCompleted {
		return errors.New("cannot cancel a completed booking")
	}

	return s.bookingRepo.UpdateStatus(ctx, id, models.EnrollmentStatusCancelled)
}

func (s *BookingService) ConfirmBooking(ctx context.Context, id uint) error {
	enrollment, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	if enrollment.Status != models.EnrollmentStatusPendingPayment {
		return errors.New("only pending bookings can be confirmed")
	}

	return s.bookingRepo.UpdateStatus(ctx, id, models.EnrollmentStatusInProgress)
}

func (s *BookingService) CompleteBooking(ctx context.Context, id uint) error {
	enrollment, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	if enrollment.Status != models.EnrollmentStatusInProgress {
		return errors.New("only in-progress bookings can be completed")
	}

	return s.bookingRepo.UpdateStatus(ctx, id, models.EnrollmentStatusCompleted)
}

func (s *BookingService) ListBookings(ctx context.Context, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, total, err := s.bookingRepo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}

func (s *BookingService) ListUserBookings(ctx context.Context, userID uint, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, total, err := s.bookingRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}

func (s *BookingService) ListInstructorBookings(ctx context.Context, instructorID uint, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, total, err := s.bookingRepo.GetByInstructorID(ctx, instructorID, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}

func (s *BookingService) ListBookingsByStatus(ctx context.Context, status string, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, total, err := s.bookingRepo.GetByStatus(ctx, models.EnrollmentStatus(status), page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}