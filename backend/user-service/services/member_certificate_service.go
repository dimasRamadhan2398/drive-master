package services

import (
	"bytes"
	"context"
	"fmt"

	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type IMemberCertificateService interface {
	GetCertificate(ctx context.Context, memberID uuid.UUID, certID uuid.UUID) (*dto.MemberCertificateDetail, error)
	GetCertificatesByMember(ctx context.Context, memberID uuid.UUID) ([]dto.MemberCertificateResponse, error)
	GenerateCertificatePDF(ctx context.Context, memberID uuid.UUID, certID uuid.UUID) ([]byte, string, error)
}

type MemberCertificateService struct {
	userRepo        repositories.IUserRepository
	entitlementRepo repositories.IEntitlementRepository
	certRepo        repositories.ICertificationRepository
}

func NewMemberCertificateService(
	userRepo repositories.IUserRepository,
	entitlementRepo repositories.IEntitlementRepository,
	certRepo repositories.ICertificationRepository,
) IMemberCertificateService {
	return &MemberCertificateService{
		userRepo:        userRepo,
		entitlementRepo: entitlementRepo,
		certRepo:        certRepo,
	}
}

// GetCertificate retrieves a specific certificate for a member
func (s *MemberCertificateService) GetCertificate(ctx context.Context, memberID uuid.UUID, certID uuid.UUID) (*dto.MemberCertificateDetail, error) {
	// Get user info
	user, err := s.userRepo.FindByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get certification
	cert, err := s.certRepo.FindByMemberIDAndCertID(ctx, memberID, certID)
	if err != nil {
		return nil, fmt.Errorf("certificate not found: %w", err)
	}

	memberName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	return &dto.MemberCertificateDetail{
		ID:            cert.ID,
		CertNumber:    cert.CertNumber,
		MemberName:    memberName,
		MemberEmail:   user.EmailAddress,
		PackageName:   cert.IssuedBy, // Using IssuedBy as package name
		IssuedDate:    cert.IssuedDate,
		IssuedBy:      "Drive Master",
		TrainingHours: 0,
		TotalSessions: 0,
		Status:        string(cert.Status),
	}, nil
}

// GetCertificatesByMember retrieves all certificates for a member
func (s *MemberCertificateService) GetCertificatesByMember(ctx context.Context, memberID uuid.UUID) ([]dto.MemberCertificateResponse, error) {
	// Get user info
	user, err := s.userRepo.FindByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get all certifications for the member
	certs, err := s.certRepo.FindAllByMemberID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificates: %w", err)
	}

	memberName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	responses := make([]dto.MemberCertificateResponse, 0, len(certs))
	for _, cert := range certs {
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

// GenerateCertificatePDF generates a PDF certificate for download
func (s *MemberCertificateService) GenerateCertificatePDF(ctx context.Context, memberID uuid.UUID, certID uuid.UUID) ([]byte, string, error) {
	cert, err := s.GetCertificate(ctx, memberID, certID)
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
