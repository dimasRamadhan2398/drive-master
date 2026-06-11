package services

import (
	"context"
	"time"

	"booking-service/models/dto"
)

type IBookingService interface {
	CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (*dto.EnrollmentResponse, error)
	GetBooking(ctx context.Context, id uint) (*dto.EnrollmentResponse, error)
	UpdateBooking(ctx context.Context, id uint, req dto.UpdateBookingRequest) (*dto.EnrollmentResponse, error)
	CancelBooking(ctx context.Context, id uint) error
	ConfirmBooking(ctx context.Context, id uint) error
	CompleteBooking(ctx context.Context, id uint) error
	ListBookings(ctx context.Context, page, limit int) (*dto.EnrollmentListResponse, error)
}

type ISessionService interface {
	CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error)
	GetSession(ctx context.Context, id uint) (*dto.SessionResponse, error)
	ListSessions(ctx context.Context, page, limit int) (*dto.SessionListResponse, error)
	GetStats(ctx context.Context) (*dto.SessionStatsResponse, error)
}

type IEntitlementService interface {
	CreateEntitlement(ctx context.Context, req dto.CreateEntitlementRequest) (*dto.EntitlementResponse, error)
	GetEntitlement(ctx context.Context, id uint) (*dto.EntitlementResponse, error)
	UpdateEntitlement(ctx context.Context, id uint, req dto.UpdateEntitlementRequest) (*dto.EntitlementResponse, error)
	DeleteEntitlement(ctx context.Context, id uint) error
	ListEntitlements(ctx context.Context, page, limit int) (*dto.EntitlementListResponse, error)
	GetUserEntitlements(ctx context.Context, userID uint) ([]dto.EntitlementResponse, error)
}

type ICertificationService interface {
	CreateCertification(ctx context.Context, req dto.CreateCertificationRequest) (*dto.CertificationResponse, error)
	GetCertification(ctx context.Context, id uint) (*dto.CertificationResponse, error)
	UpdateCertificationStatus(ctx context.Context, id uint, status string) (*dto.CertificationResponse, error)
	IssueCertification(ctx context.Context, req dto.IssueCertificationRequest) (*dto.UserCertificationResponse, error)
	RevokeCertification(ctx context.Context, userID, certificationID uint) error
	ListCertifications(ctx context.Context, page, limit int) (*dto.CertificationListResponse, error)
	GetUserCertifications(ctx context.Context, userID uint) ([]dto.UserCertificationResponse, error)
	GetCertificationsByPackage(ctx context.Context, packageID uint) ([]dto.CertificationResponse, error)
	GetStats(ctx context.Context) (*dto.CertificationStatsResponse, error)
}

type IEnrollmentService interface {
	CreateEnrollment(ctx context.Context, req dto.CreateEnrollmentRequest) (*dto.EnrollmentResponse, error)
	GetEnrollment(ctx context.Context, id uint) (*dto.EnrollmentResponse, error)
	UpdateEnrollment(ctx context.Context, id uint, req dto.UpdateEnrollmentRequest) (*dto.EnrollmentResponse, error)
	CancelEnrollment(ctx context.Context, id uint) error
	MarkAsPaid(ctx context.Context, id uint, totalPrice float64) (*dto.EnrollmentResponse, error)
	ListEnrollments(ctx context.Context, page, limit int) (*dto.EnrollmentListResponse, error)
	ListUserEnrollments(ctx context.Context, userID uint, page, limit int) (*dto.EnrollmentListResponse, error)
	ListEnrollmentsByStatus(ctx context.Context, status string, page, limit int) (*dto.EnrollmentListResponse, error)
	CreateEntitlementFromEnrollment(ctx context.Context, enrollmentID uint, sourceType, sourceID string, totalSessions int, expiresAt time.Time) (*dto.EntitlementResponse, error)
}

type IScheduleService interface {
	CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*dto.ScheduleResponse, error)
	GetSchedule(ctx context.Context, id uint) (*dto.ScheduleResponse, error)
	UpdateSchedule(ctx context.Context, id uint, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error)
	DeleteSchedule(ctx context.Context, id uint) error
	ListSchedules(ctx context.Context, page, limit int) (*dto.ScheduleListResponse, error)
	ListSchedulesFiltered(ctx context.Context, params dto.ScheduleFilterParams) (*dto.ScheduleListResponse, error)
	GetAvailableSchedules(ctx context.Context, startDate, endDate string) (*dto.ScheduleListResponse, error)
	BookSlot(ctx context.Context, slotID uint, req dto.BookSlotRequest) (*dto.ScheduleResponse, error)
	CancelBooking(ctx context.Context, slotID uint) error
}

// IServiceRegistry defines methods for getting services
type IServiceRegistry interface {
	GetBookingService() IBookingService
	GetSessionService() ISessionService
	GetEntitlementService() IEntitlementService
	GetCertificationService() ICertificationService
	GetEnrollmentService() IEnrollmentService
	GetScheduleService() IScheduleService
}