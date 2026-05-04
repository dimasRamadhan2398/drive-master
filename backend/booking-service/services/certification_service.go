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

type CertificationService struct {
	certRepo     *repositories.CertificationRepository
	userCertRepo *repositories.UserCertificationRepository
}

func NewCertificationService(
	certRepo *repositories.CertificationRepository,
	userCertRepo *repositories.UserCertificationRepository,
) ICertificationService {
	return &CertificationService{
		certRepo:     certRepo,
		userCertRepo: userCertRepo,
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

func (s *CertificationService) GetCertification(ctx context.Context, id uint) (*dto.CertificationResponse, error) {
	certification, err := s.certRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("certification not found")
		}
		return nil, err
	}

	resp := s.certRepo.ToResponse(certification)
	return &resp, nil
}

func (s *CertificationService) UpdateCertificationStatus(ctx context.Context, id uint, status string) (*dto.CertificationResponse, error) {
	certification, err := s.certRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("certification not found")
		}
		return nil, err
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
	certification, err := s.certRepo.GetByID(ctx, req.CertificationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("certification not found")
		}
		return nil, err
	}

	// Check if already issued
	existing, err := s.userCertRepo.GetByUserAndCertification(ctx, req.UserID, req.CertificationID)
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
	return &resp, nil
}

func (s *CertificationService) RevokeCertification(ctx context.Context, userID, certificationID uint) error {
	certification, err := s.certRepo.GetByID(ctx, certificationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("certification not found")
		}
		return err
	}

	certification.Status = models.CertificationStatusRevoked
	return s.certRepo.Update(ctx, certification)
}

func (s *CertificationService) ListCertifications(ctx context.Context, page, limit int) (*dto.CertificationListResponse, error) {
	certifications, total, err := s.certRepo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.certRepo.ToListResponse(certifications, total, page, limit)
	return &resp, nil
}

func (s *CertificationService) GetUserCertifications(ctx context.Context, userID uint) ([]dto.UserCertificationResponse, error) {
	userCerts, err := s.userCertRepo.GetByUserID(ctx, userID)
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
	certifications, err := s.certRepo.GetByPackageID(ctx, packageID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CertificationResponse, len(certifications))
	for i, c := range certifications {
		responses[i] = s.certRepo.ToResponse(&c)
	}
	return responses, nil
}