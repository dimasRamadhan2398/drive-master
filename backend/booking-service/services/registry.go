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

// IServiceRegistry defines methods for getting services
type IServiceRegistry interface {
	GetSessionService() ISessionService
	GetEntitlementService() IEntitlementService
	GetEnrollmentService() IEnrollmentService
	GetScheduleService() IScheduleService
	GetPaymentService() IPaymentService
	GetRevenueService() IRevenueService
}