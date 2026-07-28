package services

import (
	"context"
	"testing"
	"time"

	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockCertRepo struct {
	repositories.ICertificationRepository
	cert *models.Certification
	err  error
}

func (m *mockCertRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	return m.cert, m.err
}

type mockUserRepo struct {
	repositories.IUserRepository
	user *models.User
	err  error
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return m.user, m.err
}

func TestCertificationService_DownloadCertificatePDF(t *testing.T) {
	certID := uuid.New()
	memberID := uuid.New()

	mockCert := &models.Certification{
		ID:         certID,
		MemberID:   memberID,
		CertType:   "package_completion",
		CertNumber: "CERT-TEST-12345",
		IssuedBy:   "Gold Pack completion",
		IssuedDate: time.Now(),
		Status:     models.CertificationStatusVerified,
	}

	mockUser := &models.User{
		ID:           memberID,
		Username:     "teststudent",
		EmailAddress: "student@example.com",
		FirstName:    "John",
		LastName:     "Doe",
	}

	certRepo := &mockCertRepo{cert: mockCert}
	userRepo := &mockUserRepo{user: mockUser}

	service := NewCertificationService(certRepo, userRepo, nil)

	// Test DownloadCertificatePDF
	pdfBytes, filename, err := service.DownloadCertificatePDF(context.Background(), certID)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	assert.Equal(t, "certificate-CERT-TEST-12345.pdf", filename)

	// Check PDF header (%PDF)
	assert.True(t, len(pdfBytes) > 4)
	assert.Equal(t, "%PDF", string(pdfBytes[:4]))
}
