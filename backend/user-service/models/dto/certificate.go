package dto

import (
	"time"

	"github.com/google/uuid"
)

// MemberCertificateResponse represents a member's certificate
type MemberCertificateResponse struct {
	ID           uuid.UUID `json:"id"`
	MemberID     uuid.UUID `json:"memberId"`
	MemberName   string    `json:"memberName"`
	MemberEmail  string    `json:"memberEmail"`
	PackageName  string    `json:"packageName"`
	CertNumber   string    `json:"certNumber"`
	IssuedDate   string    `json:"issuedDate"`
	CompletedAt  string    `json:"completedAt"`
	Status       string    `json:"status"` // "eligible", "issued", "expired"
}

// MemberCertificateDetail represents detailed certificate info for PDF generation
type MemberCertificateDetail struct {
	ID           uuid.UUID `json:"id"`
	CertNumber   string    `json:"certNumber"`
	MemberName   string    `json:"memberName"`
	MemberEmail  string    `json:"memberEmail"`
	PackageName  string    `json:"packageName"`
	IssuedDate   time.Time `json:"issuedDate"`
	IssuedBy     string    `json:"issuedBy"`
	TrainingHours int       `json:"trainingHours"`
	TotalSessions int       `json:"totalSessions"`
	Status       string    `json:"status"`
}

// IssueMemberCertificateInput represents input for issuing a member certificate
type IssueMemberCertificateInput struct {
	MemberID      uuid.UUID `json:"memberId" binding:"required"`
	PackageID     uuid.UUID `json:"packageId" binding:"required"`
	PackageName   string    `json:"packageName" binding:"required"`
	EntitlementID uuid.UUID `json:"entitlementId" binding:"required"`
}

// IssueMemberCertificateResponse represents response after issuing a certificate
type IssueMemberCertificateResponse struct {
	ID          uuid.UUID `json:"id"`
	CertNumber  string    `json:"certNumber"`
	MemberID    uuid.UUID `json:"memberId"`
	IssuedDate  time.Time `json:"issuedDate"`
	IssuedBy    string    `json:"issuedBy"`
	Message     string    `json:"message"`
}