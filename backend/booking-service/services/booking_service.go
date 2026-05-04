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

func (s *BookingService) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	// Check if entitlement exists and has remaining sessions
	entitlement, err := s.entitlementRepo.GetByID(ctx, req.EntitlementID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entitlement not found")
		}
		return nil, err
	}

	if entitlement.SessionsRemaining <= 0 {
		return nil, errors.New("no sessions remaining in entitlement")
	}

	if entitlement.ExpiresAt.Before(req.DateOfSession) {
		return nil, errors.New("entitlement has expired")
	}

	booking := &models.Booking{
		UserID:        req.UserID,
		InstructorID:  req.InstructorID,
		EntitlementID: req.EntitlementID,
		DateOfSession: req.DateOfSession,
		FromTime:      req.FromTime,
		ToTime:        req.ToTime,
		CarID:         req.CarID,
		Area:          req.Area,
		Notes:         req.Notes,
		Status:        models.BookingStatusPending,
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToResponse(booking)
	return &resp, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id uint) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	resp := s.bookingRepo.ToResponse(booking)
	return &resp, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, id uint, req dto.UpdateBookingRequest) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	if req.DateOfSession != nil {
		booking.DateOfSession = *req.DateOfSession
	}
	if req.FromTime != nil {
		booking.FromTime = *req.FromTime
	}
	if req.ToTime != nil {
		booking.ToTime = *req.ToTime
	}
	if req.CarID != nil {
		booking.CarID = *req.CarID
	}
	if req.Area != nil {
		booking.Area = *req.Area
	}
	if req.Notes != nil {
		booking.Notes = *req.Notes
	}
	if req.Status != nil {
		booking.Status = models.BookingStatus(*req.Status)
	}

	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToResponse(booking)
	return &resp, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, id uint) error {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	if booking.Status == models.BookingStatusCancelled {
		return errors.New("booking already cancelled")
	}

	if booking.Status == models.BookingStatusCompleted {
		return errors.New("cannot cancel a completed booking")
	}

	return s.bookingRepo.UpdateStatus(ctx, id, models.BookingStatusCancelled)
}

func (s *BookingService) ConfirmBooking(ctx context.Context, id uint) error {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	if booking.Status != models.BookingStatusPending {
		return errors.New("only pending bookings can be confirmed")
	}

	return s.bookingRepo.UpdateStatus(ctx, id, models.BookingStatusConfirmed)
}

func (s *BookingService) CompleteBooking(ctx context.Context, id uint) error {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	if booking.Status != models.BookingStatusConfirmed {
		return errors.New("only confirmed bookings can be completed")
	}

	// Decrement entitlement sessions
	if err := s.entitlementRepo.DecrementSessions(ctx, booking.EntitlementID, 1); err != nil {
		return errors.New("failed to decrement entitlement sessions")
	}

	return s.bookingRepo.UpdateStatus(ctx, id, models.BookingStatusCompleted)
}

func (s *BookingService) ListBookings(ctx context.Context, page, limit int) (*dto.BookingListResponse, error) {
	bookings, total, err := s.bookingRepo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToListResponse(bookings, total, page, limit)
	return &resp, nil
}

func (s *BookingService) ListUserBookings(ctx context.Context, userID uint, page, limit int) (*dto.BookingListResponse, error) {
	bookings, total, err := s.bookingRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToListResponse(bookings, total, page, limit)
	return &resp, nil
}

func (s *BookingService) ListInstructorBookings(ctx context.Context, instructorID uint, page, limit int) (*dto.BookingListResponse, error) {
	bookings, total, err := s.bookingRepo.GetByInstructorID(ctx, instructorID, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToListResponse(bookings, total, page, limit)
	return &resp, nil
}

func (s *BookingService) ListBookingsByStatus(ctx context.Context, status string, page, limit int) (*dto.BookingListResponse, error) {
	bookings, total, err := s.bookingRepo.GetByStatus(ctx, models.BookingStatus(status), page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.bookingRepo.ToListResponse(bookings, total, page, limit)
	return &resp, nil
}