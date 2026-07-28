package seeders

import (
	"time"

	"user-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CertificationSeeder struct {
	db *gorm.DB
}

func NewCertificationSeeder(db *gorm.DB) *CertificationSeeder {
	return &CertificationSeeder{db: db}
}

func (s *CertificationSeeder) Seed() error {
	entitlementID := uuid.MustParse("88888888-8888-8888-8888-888888888888")
	memberID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	packageName := "10x Session"

	now := time.Now()
	certNumber := "CERT-11111111-99999999"

	cert := models.Certification{
		ID:            uuid.MustParse("77777777-7777-7777-7777-777777777777"),
		MemberID:      memberID,
		EntitlementID: &entitlementID,
		CertType:      "package_completion",
		CertNumber:    certNumber,
		IssuedBy:      packageName,
		IssuedDate:    now.AddDate(0, 0, -30),
		Status:        models.CertificationStatusVerified,
		Notes:         "Package: " + packageName + " | Entitlement: " + entitlementID.String(),
		VerifiedAt:    &now,
		CreatedAt:     now.AddDate(0, 0, -30),
		UpdatedAt:     now.AddDate(0, 0, -30),
	}

	var existing models.Certification
	if err := s.db.Where("id = ?", cert.ID).First(&existing).Error; err != nil {
		if err := s.db.Create(&cert).Error; err != nil {
			return err
		}
	}

	return nil
}
