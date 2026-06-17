package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"booking-service/clients/user"
	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/kafka"
	"booking-service/repositories"

	"github.com/google/uuid"
)

type CertificationService struct {
	certRepo     repositories.ICertificationRepository
	userCertRepo repositories.IUserCertificationRepository
	userClient   user.IUserClient
	eventPublisher kafka.IEventPublisher
}

func NewCertificationService(
	certRepo repositories.ICertificationRepository,
	userCertRepo repositories.IUserCertificationRepository,
	userClient user.IUserClient,
	eventPublisher kafka.IEventPublisher,
) ICertificationService {
	return &CertificationService{
		certRepo:     certRepo,
		userCertRepo: userCertRepo,
		userClient:   userClient,
		eventPublisher: eventPublisher,
	}
}

func (s *CertificationService) CreateCertification(ctx context.Context, req dto.CreateCertificationRequest) (*dto.CertificationResponse, error) {
	certification := &models.Certification{
		Type:       req.Type,
		Recipient:  req.Recipient,
		IssueDate:  req.IssueDate,
		PackageID:  req.PackageID,
		Status:     models.CertificationStatusPending,
	}

	if err := s.certRepo.Create(ctx, certification); err != nil {
		return nil, err
	}

	resp := s.certRepo.ToResponse(certification)
	return &resp, nil
}

func (s *CertificationService) GetCertification(ctx context.Context, id uuid.UUID) (*dto.CertificationResponse, error) {
	certification, err := s.certRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("certification not found")
	}

	resp := s.certRepo.ToResponse(certification)
	return &resp, nil
}

func (s *CertificationService) UpdateCertificationStatus(ctx context.Context, id uuid.UUID, status string) (*dto.CertificationResponse, error) {
	certification, err := s.certRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("certification not found")
	}

	certification.Status = models.CertificationStatus(status)
	if err := s.certRepo.Update(ctx, certification); err != nil {
		return nil, err
	}

	resp := s.certRepo.ToResponse(certification)
	return &resp, nil
}

func (s *CertificationService) IssueCertification(ctx context.Context, req dto.IssueCertificationRequest) (*dto.UserCertificationResponse, error) {
	// Check if certification exists
	certification, err := s.certRepo.FindByID(ctx, req.CertificationID)
	if err != nil {
		return nil, errors.New("certification not found")
	}

	// Check if already issued
	existing, err := s.userCertRepo.FindByUserAndCertification(ctx, req.UserID, req.CertificationID)
	if err == nil && existing != nil {
		return nil, errors.New("certification already issued to this user")
	}

	userCert := &models.UserCertification{
		UserID:          req.UserID,
		CertificationID: req.CertificationID,
		IssuedAt:        time.Now(),
	}

	if err := s.userCertRepo.Create(ctx, userCert); err != nil {
		return nil, err
	}

	// Update certification status to issued
	certification.Status = models.CertificationStatusIssued
	if err := s.certRepo.Update(ctx, certification); err != nil {
		return nil, err
	}

	resp := s.userCertRepo.ToResponse(userCert, s.certRepo.ToResponse(certification))

	// Send email notification asynchronously
	go s.sendCertificationEmail(context.Background(), req.UserID, certification.Type, userCert.IssuedAt)

	return &resp, nil
}

// sendCertificationEmail fetches user info and sends the certification email
func (s *CertificationService) sendCertificationEmail(ctx context.Context, userID uuid.UUID, certType string, issueDate time.Time) {
	if s.userClient == nil {
		return
	}

	userInfo, err := s.userClient.GetUserByID(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to get user info for certification email: %v\n", err)
		return
	}

	if s.eventPublisher != nil {
		userName := userInfo.FirstName
		if userName == "" {
			userName = userInfo.Username
		}
		if err := s.eventPublisher.PublishCertificationIssued(
			ctx,
			fmt.Sprintf("%d", userID),
			userInfo.Email,
			userName,
			certType,
			issueDate,
		); err != nil {
			fmt.Printf("Failed to publish certification issued event: %v\n", err)
		}
	}
}

func (s *CertificationService) RevokeCertification(ctx context.Context, userID, certificationID uuid.UUID) error {
	certification, err := s.certRepo.FindByID(ctx, certificationID)
	if err != nil {
		return errors.New("certification not found")
	}

	certification.Status = models.CertificationStatusRevoked
	return s.certRepo.Update(ctx, certification)
}

func (s *CertificationService) ListCertifications(ctx context.Context, page, limit int) (*dto.CertificationListResponse, error) {
	certifications, err := s.certRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	total, err := s.certRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := s.certRepo.ToListResponse(certifications, total, page, limit)
	return &resp, nil
}

func (s *CertificationService) GetUserCertifications(ctx context.Context, userID uuid.UUID) ([]dto.UserCertificationResponse, error) {
	userCerts, err := s.userCertRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserCertificationResponse, len(userCerts))
	for i, uc := range userCerts {
		responses[i] = s.userCertRepo.ToResponse(&uc, s.certRepo.ToResponse(&uc.Certification))
	}
	return responses, nil
}

func (s *CertificationService) GetCertificationsByPackage(ctx context.Context, packageID uint) ([]dto.CertificationResponse, error) {
	certifications, err := s.certRepo.FindByPackageID(ctx, packageID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CertificationResponse, len(certifications))
	for i, c := range certifications {
		responses[i] = s.certRepo.ToResponse(&c)
	}
	return responses, nil
}

func (s *CertificationService) GetStats(ctx context.Context) (*dto.CertificationStatsResponse, error) {
	stats, err := s.certRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.CertificationStatsResponse{
		TotalCertifications:   stats.Total,
		IssuedCertifications:  stats.Issued,
		ActiveCertifications:  stats.Active,
		RevokedCertifications: stats.Revoked,
	}, nil
}