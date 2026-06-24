package services

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type ICertificationService interface {
	// Admin operations
	IssueCertificate(ctx context.Context, input dto.IssueMemberCertificateInput) (*dto.IssueMemberCertificateResponse, error)
	RevokeCertificate(ctx context.Context, certID uuid.UUID) error
	GetMemberCertificates(ctx context.Context, memberID uuid.UUID) ([]dto.MemberCertificateResponse, error)
	GetCertificateStats(ctx context.Context) (*dto.CertificateStatsResponse, error)

	// Member operations
	GetCertificate(ctx context.Context, certID uuid.UUID) (*dto.MemberCertificateDetail, error)
	DownloadCertificatePDF(ctx context.Context, certID uuid.UUID) ([]byte, string, error)
}

type CertificationService struct {
	repo         repositories.ICertificationRepository
	userRepo     repositories.IUserRepository
	emailService IMailtrapEmailService
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

// IssueCertificate creates a new certificate for a member (admin action)
func (s *CertificationService) IssueCertificate(ctx context.Context, input dto.IssueMemberCertificateInput) (*dto.IssueMemberCertificateResponse, error) {
	// Check if certificate already exists for this entitlement
	existingCerts, err := s.repo.FindAllByMemberID(ctx, input.MemberID)
	if err == nil {
		for _, cert := range existingCerts {
			if cert.Notes != "" && containsString(cert.Notes, input.EntitlementID.String()) {
				return nil, fmt.Errorf("certificate already issued for this entitlement")
			}
		}
	}

	// Get member info
	user, err := s.userRepo.FindByID(ctx, input.MemberID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	// Generate certificate number
	certNumber := fmt.Sprintf("CERT-%s-%s", input.PackageID.String()[:8], input.MemberID.String()[:8])
	now := time.Now()

	cert := &models.Certification{
		ID:         uuid.New(),
		MemberID:   input.MemberID,
		CertType:   "package_completion",
		CertNumber: certNumber,
		IssuedBy:   input.PackageName,
		IssuedDate: now,
		Status:     models.CertificationStatusVerified,
		Notes:      fmt.Sprintf("Package: %s | Entitlement: %s", input.PackageName, input.EntitlementID.String()),
		VerifiedAt: &now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(ctx, cert); err != nil {
		return nil, err
	}

	// Send email notification
	go func() {
		if s.emailService != nil {
			name := user.FirstName
			if name == "" {
				name = user.Username
			}
			_ = s.emailService.SendCertificationEmail(context.Background(), user.EmailAddress, name, certNumber, input.PackageName)
		}
	}()

	memberName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	return &dto.IssueMemberCertificateResponse{
		ID:         cert.ID,
		CertNumber: certNumber,
		MemberID:   cert.MemberID,
		IssuedDate: cert.IssuedDate,
		IssuedBy:   cert.IssuedBy,
		Message:    fmt.Sprintf("Certificate issued to %s for completing %s", memberName, input.PackageName),
	}, nil
}

// RevokeCertificate revokes a certificate (admin action)
func (s *CertificationService) RevokeCertificate(ctx context.Context, certID uuid.UUID) error {
	cert, err := s.repo.FindByID(ctx, certID)
	if err != nil {
		return fmt.Errorf("certificate not found: %w", err)
	}

	cert.Status = models.CertificationStatusRevoked
	cert.UpdatedAt = time.Now()

	return s.repo.Update(ctx, cert)
}

// GetMemberCertificates retrieves all certificates for a member
func (s *CertificationService) GetMemberCertificates(ctx context.Context, memberID uuid.UUID) ([]dto.MemberCertificateResponse, error) {
	user, err := s.userRepo.FindByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	certs, err := s.repo.FindAllByMemberID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificates: %w", err)
	}

	memberName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	responses := make([]dto.MemberCertificateResponse, 0, len(certs))
	for _, cert := range certs {
		if cert.Status == models.CertificationStatusRevoked {
			continue // Skip revoked certificates in member view
		}

		status := "issued"
		if cert.Status == models.CertificationStatusExpired {
			status = "expired"
		} else if cert.Status == models.CertificationStatusPending {
			status = "eligible"
		}

		completedAt := ""
		if cert.VerifiedAt != nil {
			completedAt = cert.VerifiedAt.Format("2006-01-02")
		}

		responses = append(responses, dto.MemberCertificateResponse{
			ID:          cert.ID,
			MemberID:    memberID,
			MemberName:  memberName,
			MemberEmail: user.EmailAddress,
			PackageName: cert.IssuedBy,
			CertNumber:  cert.CertNumber,
			IssuedDate:  cert.IssuedDate.Format("2006-01-02"),
			CompletedAt: completedAt,
			Status:      status,
		})
	}

	return responses, nil
}

// GetCertificate retrieves a specific certificate detail
func (s *CertificationService) GetCertificate(ctx context.Context, certID uuid.UUID) (*dto.MemberCertificateDetail, error) {
	cert, err := s.repo.FindByID(ctx, certID)
	if err != nil {
		return nil, fmt.Errorf("certificate not found: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, cert.MemberID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	memberName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	return &dto.MemberCertificateDetail{
		ID:            cert.ID,
		CertNumber:    cert.CertNumber,
		MemberName:    memberName,
		MemberEmail:   user.EmailAddress,
		PackageName:   cert.IssuedBy,
		IssuedDate:    cert.IssuedDate,
		IssuedBy:      "Drive Master",
		TrainingHours: 0,
		TotalSessions: 0,
		Status:        string(cert.Status),
	}, nil
}

// DownloadCertificatePDF generates a downloadable PDF certificate
func (s *CertificationService) DownloadCertificatePDF(ctx context.Context, certID uuid.UUID) ([]byte, string, error) {
	cert, err := s.GetCertificate(ctx, certID)
	if err != nil {
		return nil, "", err
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Set background color (light cream)
	pdf.SetFillColor(255, 251, 240)
	pdf.Rect(0, 0, 297, 210, "F")

	// Border
	pdf.SetDrawColor(139, 69, 19)
	pdf.SetLineWidth(2)
	pdf.Rect(10, 10, 277, 190, "D")

	// Inner border
	pdf.SetLineWidth(0.5)
	pdf.Rect(15, 15, 267, 180, "D")

	// Header - Logo area
	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetTextColor(139, 69, 19)
	pdf.SetXY(20, 25)
	pdf.Cell(257, 15, "DRIVE MASTER")

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(100, 100, 100)
	pdf.SetXY(20, 38)
	pdf.Cell(257, 8, "Professional Driving Academy")

	// Certificate Title
	pdf.SetFont("Helvetica", "B", 28)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(20, 60)
	pdf.Cell(257, 15, "CERTIFICATE OF COMPLETION")

	pdf.SetFont("Helvetica", "", 14)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(20, 80)
	pdf.Cell(257, 8, "This is to certify that")

	// Recipient Name
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetTextColor(139, 69, 19)
	pdf.SetXY(20, 95)
	pdf.Cell(257, 12, cert.MemberName)

	// Course completion text
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(20, 112)
	pdf.Cell(257, 8, "has successfully completed the")

	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(20, 125)
	pdf.Cell(257, 10, cert.PackageName)

	// Training details
	pdf.SetFont("Helvetica", "", 11)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(20, 140)
	if cert.TrainingHours > 0 {
		pdf.Cell(257, 8, fmt.Sprintf("Training Hours: %d hours", cert.TrainingHours))
	}

	// Certificate details section
	pdf.SetLineWidth(0.5)
	pdf.Line(25, 160, 272, 160)

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(100, 100, 100)

	// Certificate ID
	pdf.SetXY(25, 165)
	pdf.Cell(80, 6, fmt.Sprintf("Certificate ID: %s", cert.ID.String()))
	pdf.SetXY(25, 172)
	pdf.Cell(80, 6, fmt.Sprintf("Certificate Number: %s", cert.CertNumber))

	// Date
	pdf.SetXY(192, 165)
	pdf.Cell(80, 6, fmt.Sprintf("Issue Date: %s", cert.IssuedDate.Format("January 2, 2006")))

	// Issued by
	pdf.SetXY(192, 172)
	pdf.Cell(80, 6, fmt.Sprintf("Issued By: %s", cert.IssuedBy))

	// Footer signature line
	pdf.Line(60, 185, 120, 185)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetXY(60, 186)
	pdf.Cell(60, 5, "Authorized Signature")

	pdf.Line(177, 185, 237, 185)
	pdf.SetXY(177, 186)
	pdf.Cell(60, 5, "Date")

	// Generate PDF
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	filename := fmt.Sprintf("certificate-%s.pdf", cert.CertNumber)
	return buf.Bytes(), filename, nil
}

// GetCertificateStats retrieves certificate statistics
func (s *CertificationService) GetCertificateStats(ctx context.Context) (*dto.CertificateStatsResponse, error) {
	stats, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate growth rate
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	currentStart := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	currentEnd := currentStart.AddDate(0, 1, -1)

	lastMonth := currentMonth - 1
	lastYear := currentYear
	if lastMonth < 1 {
		lastMonth = 12
		lastYear--
	}
	lastMonthStart := time.Date(lastYear, lastMonth, 1, 0, 0, 0, 0, currentLocation)
	lastMonthEnd := lastMonthStart.AddDate(0, 1, -1)

	currentStats, err := s.repo.GetStatsByDateRange(ctx, currentStart, currentEnd)
	if err != nil {
		return nil, err
	}

	lastMonthStats, err := s.repo.GetStatsByDateRange(ctx, lastMonthStart, lastMonthEnd)
	if err != nil {
		return nil, err
	}

	var monthlyGrowth float64
	if lastMonthStats.Total > 0 {
		monthlyGrowth = float64(currentStats.Total-lastMonthStats.Total) / float64(lastMonthStats.Total) * 100
	}

	return &dto.CertificateStatsResponse{
		Total:            stats.Total,
		Verified:         stats.Verified,
		Pending:          stats.Pending,
		Expired:          stats.Expired,
		Revoked:          stats.Revoked,
		MonthlyTotal:     currentStats.Total,
		MonthlyGrowth:    monthlyGrowth,
		GrowthPercentage: monthlyGrowth,
	}, nil
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
