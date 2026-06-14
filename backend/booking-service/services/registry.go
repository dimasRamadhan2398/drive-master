package services

import (
	"context"

	"booking-service/models/dto"
)


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

// IServiceRegistry defines methods for getting services
type IServiceRegistry interface {
	GetSessionService() ISessionService
	GetEntitlementService() IEntitlementService
	GetCertificationService() ICertificationService
	GetEnrollmentService() IEnrollmentService
	GetScheduleService() IScheduleService
	GetPaymentService() IPaymentService
}
