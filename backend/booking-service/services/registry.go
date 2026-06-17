package services

import (
	"context"

	"booking-service/models/dto"

	"github.com/google/uuid"
)


type IEntitlementService interface {
	CreateEntitlement(ctx context.Context, req dto.CreateEntitlementRequest) (*dto.EntitlementResponse, error)
	GetEntitlement(ctx context.Context, id uuid.UUID) (*dto.EntitlementResponse, error)
	UpdateEntitlement(ctx context.Context, id uuid.UUID, req dto.UpdateEntitlementRequest) (*dto.EntitlementResponse, error)
	DeleteEntitlement(ctx context.Context, id uuid.UUID) error
	ListEntitlements(ctx context.Context, page, limit int) (*dto.EntitlementListResponse, error)
	GetUserEntitlements(ctx context.Context, userID uint) ([]dto.EntitlementResponse, error)
}

type ICertificationService interface {
	CreateCertification(ctx context.Context, req dto.CreateCertificationRequest) (*dto.CertificationResponse, error)
	GetCertification(ctx context.Context, id uuid.UUID) (*dto.CertificationResponse, error)
	UpdateCertificationStatus(ctx context.Context, id uuid.UUID, status string) (*dto.CertificationResponse, error)
	IssueCertification(ctx context.Context, req dto.IssueCertificationRequest) (*dto.UserCertificationResponse, error)
	RevokeCertification(ctx context.Context, userID, certificationID uuid.UUID) error
	ListCertifications(ctx context.Context, page, limit int) (*dto.CertificationListResponse, error)
	GetUserCertifications(ctx context.Context, userID uuid.UUID) ([]dto.UserCertificationResponse, error)
	GetCertificationsByPackage(ctx context.Context, packageID uint) ([]dto.CertificationResponse, error)
	GetStats(ctx context.Context) (*dto.CertificationStatsResponse, error)
}

// IServiceRegistry defines methods for getting services
type IServiceRegistry interface {
	GetSessionService() ISessionService
	GetEntitlementService() IEntitlementService
	GetCertificationService() ICertificationService
	GetEnrollmentService() IEnrollmentService
	GetScheduleService() IScheduleService
	GetPaymentService() IPaymentService
}
