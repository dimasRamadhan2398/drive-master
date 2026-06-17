package services

import (
	"context"
	"fmt"
	"time"

	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"
	"user-service/services/listeners"

	"github.com/google/uuid"
)

type ICertificationService interface {
	CreateCertification(ctx context.Context, instructorID uuid.UUID, input dto.CreateCertificationInput) (*dto.CertificationResponse, error)
	UpdateCertification(ctx context.Context, instructorID, certID uuid.UUID, input dto.UpdateCertificationInput) (*dto.CertificationResponse, error)
	DeleteCertification(ctx context.Context, instructorID, certID uuid.UUID) error
	GetCertification(ctx context.Context, instructorID, certID uuid.UUID) (*dto.CertificationResponse, error)
	ListCertifications(ctx context.Context, instructorID uuid.UUID, page, limit int) (*dto.CertificationListResponse, error)
	VerifyCertification(ctx context.Context, instructorID, certID, verifiedBy uuid.UUID, input dto.VerifyCertificationInput) (*dto.CertificationResponse, error)
	IssueCertification(ctx context.Context, input dto.IssueCertificationInput) (*dto.CertificationResponse, error)
}

type CertificationService struct {
	repo     repositories.ICertificationRepository
	listener listeners.IEntitlementCompletedListener
}

func NewCertificationService(repo repositories.ICertificationRepository) ICertificationService {
	return &CertificationService{
		repo: repo,
	}
}

func (s *CertificationService) CreateCertification(ctx context.Context, instructorID uuid.UUID, input dto.CreateCertificationInput) (*dto.CertificationResponse, error) {
	issuedDate, err := time.Parse("2006-01-02", input.IssuedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid issued date format: %w", err)
	}

	cert := &models.Certification{
		ID:           uuid.New(),
		InstructorID: instructorID,
		CertType:     input.CertType,
		CertNumber:   input.CertNumber,
		IssuedBy:     input.IssuedBy,
		IssuedDate:   issuedDate,
		Status:       models.CertificationStatusPending,
		DocumentURL:  input.DocumentURL,
		Notes:        input.Notes,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if input.ExpiryDate != "" {
		expiryDate, err := time.Parse("2006-01-02", input.ExpiryDate)
		if err == nil {
			cert.ExpiryDate = &expiryDate
		}
	}

	if err := s.repo.Create(ctx, cert); err != nil {
		return nil, err
	}

	return toCertificationResponse(cert), nil
}

func (s *CertificationService) UpdateCertification(ctx context.Context, instructorID, certID uuid.UUID, input dto.UpdateCertificationInput) (*dto.CertificationResponse, error) {
	cert, err := s.repo.FindByInstructorAndID(ctx, instructorID, certID)
	if err != nil {
		return nil, err
	}

	if input.CertType != "" {
		cert.CertType = input.CertType
	}
	if input.CertNumber != "" {
		cert.CertNumber = input.CertNumber
	}
	if input.IssuedBy != "" {
		cert.IssuedBy = input.IssuedBy
	}
	if input.ExpiryDate != "" {
		expiryDate, err := time.Parse("2006-01-02", input.ExpiryDate)
		if err == nil {
			cert.ExpiryDate = &expiryDate
		}
	}
	if input.DocumentURL != "" {
		cert.DocumentURL = input.DocumentURL
	}
	if input.Notes != "" {
		cert.Notes = input.Notes
	}

	cert.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, cert); err != nil {
		return nil, err
	}

	return toCertificationResponse(cert), nil
}

func (s *CertificationService) DeleteCertification(ctx context.Context, instructorID, certID uuid.UUID) error {
	_, err := s.repo.FindByInstructorAndID(ctx, instructorID, certID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, certID)
}

func (s *CertificationService) GetCertification(ctx context.Context, instructorID, certID uuid.UUID) (*dto.CertificationResponse, error) {
	cert, err := s.repo.FindByInstructorAndID(ctx, instructorID, certID)
	if err != nil {
		return nil, err
	}
	return toCertificationResponse(cert), nil
}

func (s *CertificationService) ListCertifications(ctx context.Context, instructorID uuid.UUID, page, limit int) (*dto.CertificationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	certs, total, err := s.repo.FindByInstructorID(ctx, instructorID, page, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CertificationResponse, len(certs))
	for i, cert := range certs {
		responses[i] = *toCertificationResponse(&cert)
	}

	return &dto.CertificationListResponse{
		Data:       responses,
		Pagination: dto.NewPaginationMeta(total, page, limit),
	}, nil
}

func (s *CertificationService) VerifyCertification(ctx context.Context, instructorID, certID, verifiedBy uuid.UUID, input dto.VerifyCertificationInput) (*dto.CertificationResponse, error) {
	cert, err := s.repo.FindByInstructorAndID(ctx, instructorID, certID)
	if err != nil {
		return nil, err
	}

	cert.Status = models.CertificationStatusVerified
	now := time.Now()
	cert.VerifiedAt = &now
	cert.VerifiedBy = &verifiedBy
	cert.UpdatedAt = time.Now()

	if input.Notes != "" {
		cert.Notes = input.Notes
	}

	if err := s.repo.Update(ctx, cert); err != nil {
		return nil, err
	}

	return toCertificationResponse(cert), nil
}

// IssueCertification creates a certification for a member when their entitlement is completed
// This is called automatically by the entitlement service when all sessions are used
func (s *CertificationService) IssueCertification(ctx context.Context, input dto.IssueCertificationInput) (*dto.CertificationResponse, error) {
	// Generate a certificate number based on package and member
	certNumber := fmt.Sprintf("CERT-%s-%s", input.PackageID.String()[:8], input.MemberID.String()[:8])

	cert := &models.Certification{
		ID:           uuid.New(),
		InstructorID: input.MemberID, // Use member ID as instructor for member-issued certs
		CertType:     "package_completion",
		CertNumber:   certNumber,
		IssuedBy:     "System - " + input.PackageName,
		IssuedDate:   input.IssuedAt,
		Status:       models.CertificationStatusVerified, // Auto-verify for completed packages
		Notes:        fmt.Sprintf("Issued upon completion of package: %s (Entitlement: %s)", input.PackageName, input.EntitlementID.String()),
		VerifiedAt:   &input.IssuedAt, // Auto-verified by system
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, cert); err != nil {
		return nil, err
	}

	return toCertificationResponse(cert), nil
}

func toCertificationResponse(cert *models.Certification) *dto.CertificationResponse {
	resp := &dto.CertificationResponse{
		ID:           cert.ID,
		InstructorID: cert.InstructorID,
		CertType:     cert.CertType,
		CertNumber:   cert.CertNumber,
		IssuedBy:     cert.IssuedBy,
		IssuedDate:   cert.IssuedDate.Format("2006-01-02"),
		Status:       string(cert.Status),
		DocumentURL:  cert.DocumentURL,
		Notes:        cert.Notes,
		CreatedAt:    cert.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    cert.UpdatedAt.Format(time.RFC3339),
	}

	if cert.ExpiryDate != nil {
		expiryStr := cert.ExpiryDate.Format("2006-01-02")
		resp.ExpiryDate = &expiryStr
	}

	if cert.VerifiedAt != nil {
		verifiedStr := cert.VerifiedAt.Format(time.RFC3339)
		resp.VerifiedAt = &verifiedStr
	}

	return resp
}
