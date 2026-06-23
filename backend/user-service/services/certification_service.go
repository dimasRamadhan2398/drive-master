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
	CreateCertification(ctx context.Context, memberID uuid.UUID, input dto.CreateCertificationInput) (*dto.CertificationResponse, error)
	UpdateCertification(ctx context.Context, memberID, certID uuid.UUID, input dto.UpdateCertificationInput) (*dto.CertificationResponse, error)
	DeleteCertification(ctx context.Context, memberID, certID uuid.UUID) error
	GetCertification(ctx context.Context, memberID, certID uuid.UUID) (*dto.CertificationResponse, error)
	ListCertifications(ctx context.Context, memberID uuid.UUID, page, limit int) (*dto.CertificationListResponse, error)
	VerifyCertification(ctx context.Context, memberID, certID, verifiedBy uuid.UUID, input dto.VerifyCertificationInput) (*dto.CertificationResponse, error)
	IssueCertification(ctx context.Context, input dto.IssueCertificationInput) (*dto.CertificationResponse, error)
	RevokeCertification(ctx context.Context, certID uuid.UUID) error
	GetStats(ctx context.Context) (*dto.CertificateStatsResponse, error)
}

type CertificationService struct {
	repo         repositories.ICertificationRepository
	userRepo     repositories.IUserRepository
	emailService IMailtrapEmailService
	listener     listeners.IEntitlementCompletedListener
}

func NewCertificationService(
	repo repositories.ICertificationRepository,
	userRepo repositories.IUserRepository,
	emailService IMailtrapEmailService,
) ICertificationService {
	return &CertificationService{
		repo:         repo,
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *CertificationService) CreateCertification(ctx context.Context, memberID uuid.UUID, input dto.CreateCertificationInput) (*dto.CertificationResponse, error) {
	issuedDate, err := time.Parse("2006-01-02", input.IssuedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid issued date format: %w", err)
	}

	cert := &models.Certification{
		ID:         uuid.New(),
		MemberID:   memberID,
		CertType:   input.CertType,
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

func (s *CertificationService) UpdateCertification(ctx context.Context, memberID, certID uuid.UUID, input dto.UpdateCertificationInput) (*dto.CertificationResponse, error) {
	cert, err := s.repo.FindByMemberIDAndCertID(ctx, memberID, certID)
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

func (s *CertificationService) DeleteCertification(ctx context.Context, memberID, certID uuid.UUID) error {
	_, err := s.repo.FindByMemberIDAndCertID(ctx, memberID, certID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, certID)
}

func (s *CertificationService) GetCertification(ctx context.Context, memberID, certID uuid.UUID) (*dto.CertificationResponse, error) {
	cert, err := s.repo.FindByMemberIDAndCertID(ctx, memberID, certID)
	if err != nil {
		return nil, err
	}
	return toCertificationResponse(cert), nil
}

func (s *CertificationService) ListCertifications(ctx context.Context, memberID uuid.UUID, page, limit int) (*dto.CertificationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	certs, total, err := s.repo.FindByMemberID(ctx, memberID, page, limit)
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

func (s *CertificationService) VerifyCertification(ctx context.Context, memberID, certID, verifiedBy uuid.UUID, input dto.VerifyCertificationInput) (*dto.CertificationResponse, error) {
	cert, err := s.repo.FindByMemberIDAndCertID(ctx, memberID, certID)
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
		ID:         uuid.New(),
		MemberID:   input.MemberID,
		CertType:   "package_completion",
		CertNumber: certNumber,
		IssuedBy:   "System - " + input.PackageName,
		IssuedDate: input.IssuedAt,
		Status:     models.CertificationStatusVerified, // Auto-verify for completed packages
		Notes:      fmt.Sprintf("Issued upon completion of package: %s (Entitlement: %s)", input.PackageName, input.EntitlementID.String()),
		VerifiedAt: &input.IssuedAt, // Auto-verified by system
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, cert); err != nil {
		return nil, err
	}

	// Send email notification
	go func() {
		user, err := s.userRepo.FindByID(context.Background(), input.MemberID)
		if err == nil && s.emailService != nil {
			name := user.FirstName
			if name == "" {
				name = user.Username
			}
			_ = s.emailService.SendCertificationEmail(context.Background(), user.EmailAddress, name, certNumber, input.PackageName)
		}
	}()

	return toCertificationResponse(cert), nil
}

func (s *CertificationService) RevokeCertification(ctx context.Context, certID uuid.UUID) error {
	cert, err := s.repo.FindByID(ctx, certID)
	if err != nil {
		return err
	}

	cert.Status = models.CertificationStatusRevoked
	cert.UpdatedAt = time.Now()

	return s.repo.Update(ctx, cert)
}

func toCertificationResponse(cert *models.Certification) *dto.CertificationResponse {
	instructorID := ""
	if cert.InstructorID != nil {
		instructorID = cert.InstructorID.String()
	}

	resp := &dto.CertificationResponse{
		ID:           cert.ID,
		MemberID:     cert.MemberID,
		InstructorID: instructorID,
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

// GetStats retrieves certificate statistics with growth compared to previous month
func (s *CertificationService) GetStats(ctx context.Context) (*dto.CertificateStatsResponse, error) {
	stats, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate growth rate compared to previous month
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	// Current month date range
	currentStart := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	currentEnd := currentStart.AddDate(0, 1, -1)

	// Last month date range
	lastMonth := currentMonth - 1
	lastYear := currentYear
	if lastMonth < 1 {
		lastMonth = 12
		lastYear--
	}
	lastMonthStart := time.Date(lastYear, lastMonth, 1, 0, 0, 0, 0, currentLocation)
	lastMonthEnd := lastMonthStart.AddDate(0, 1, -1)

	// Get current month stats
	currentStats, err := s.repo.GetStatsByDateRange(ctx, currentStart, currentEnd)
	if err != nil {
		return nil, err
	}

	// Get last month stats
	lastMonthStats, err := s.repo.GetStatsByDateRange(ctx, lastMonthStart, lastMonthEnd)
	if err != nil {
		return nil, err
	}

	// Calculate growth rate
	var monthlyGrowth float64
	if lastMonthStats.Total > 0 {
		monthlyGrowth = float64(currentStats.Total-lastMonthStats.Total) / float64(lastMonthStats.Total) * 100
	}

	return &dto.CertificateStatsResponse{
		Total:              stats.Total,
		Verified:           stats.Verified,
		Pending:            stats.Pending,
		Expired:            stats.Expired,
		Revoked:            stats.Revoked,
		MonthlyTotal:       currentStats.Total,
		MonthlyGrowth:      monthlyGrowth,
		GrowthPercentage:   monthlyGrowth,
	}, nil
}
