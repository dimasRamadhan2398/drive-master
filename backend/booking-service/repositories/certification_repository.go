package repositories

import (
	"context"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICertificationRepository interface {
	Create(ctx context.Context, certification *models.Certification) error
	CreateTx(tx *gorm.DB, certification *models.Certification) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Certification, error)
	Update(ctx context.Context, certification *models.Certification) error
	Delete(ctx context.Context, certification *models.Certification) error
	FindAll(ctx context.Context) ([]models.Certification, error)
	FindByPackageID(ctx context.Context, packageID uint) ([]models.Certification, error)
	FindByStatus(ctx context.Context, status models.CertificationStatus) ([]models.Certification, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.CertificationStatus) error
	CountAll(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status models.CertificationStatus) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	GetStats(ctx context.Context) (*CertificationStats, error)
	GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*CertificationStats, error)
	ToResponse(certification *models.Certification) dto.CertificationResponse
	ToListResponse(certifications []models.Certification, total int64, page, limit int) dto.CertificationListResponse
}

type CertificationRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewCertificationRepository(db *gorm.DB) ICertificationRepository {
	return &CertificationRepository{BaseRepository: base.NewBaseRepository(db), db: db}
}

func (r *CertificationRepository) Create(ctx context.Context, certification *models.Certification) error {
	return r.BaseRepository.Create(certification)
}

func (r *CertificationRepository) CreateTx(tx *gorm.DB, certification *models.Certification) error {
	return r.BaseRepository.CreateTx(tx, certification)
}

func (r *CertificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	var certification models.Certification
	if err := r.BaseRepository.FindByIDWithPreload(&certification, id); err != nil {
		return nil, err
	}
	return &certification, nil
}

func (r *CertificationRepository) Update(ctx context.Context, certification *models.Certification) error {
	return r.BaseRepository.Update(certification)
}

func (r *CertificationRepository) Delete(ctx context.Context, certification *models.Certification) error {
	return r.BaseRepository.Delete(certification)
}

func (r *CertificationRepository) FindAll(ctx context.Context) ([]models.Certification, error) {
	var certifications []models.Certification
	opts := base.NewQueryOptions().WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Certification{}, &certifications, opts); err != nil {
		return nil, err
	}
	return certifications, nil
}

func (r *CertificationRepository) FindByPackageID(ctx context.Context, packageID uint) ([]models.Certification, error) {
	var certifications []models.Certification
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"package_id": packageID}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Certification{}, &certifications, opts); err != nil {
		return nil, err
	}
	return certifications, nil
}

func (r *CertificationRepository) FindByStatus(ctx context.Context, status models.CertificationStatus) ([]models.Certification, error) {
	var certifications []models.Certification
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"status": status}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Certification{}, &certifications, opts); err != nil {
		return nil, err
	}
	return certifications, nil
}

func (r *CertificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.CertificationStatus) error {
	return r.BaseRepository.Exec(
		"UPDATE certifications SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	)
}

func (r *CertificationRepository) CountAll(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(&models.Certification{}, base.NewQueryOptions())
}

func (r *CertificationRepository) CountByStatus(ctx context.Context, status models.CertificationStatus) (int64, error) {
	opts := base.NewQueryOptions().WithWhere(map[string]any{"status": status})
	return r.BaseRepository.Count(&models.Certification{}, opts)
}

// CertificationStats holds certification statistics
type CertificationStats struct {
	Total   int64
	Issued  int64
	Active  int64
	Revoked int64
}

func (r *CertificationRepository) GetStats(ctx context.Context) (*CertificationStats, error) {
	stats := &CertificationStats{}

	// Get total certifications
	if err := r.db.Model(&models.Certification{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Get issued certifications
	if err := r.db.Model(&models.Certification{}).
		Where("status = ?", models.CertificationStatusIssued).
		Count(&stats.Issued).Error; err != nil {
		return nil, err
	}

	// Get active (issued and not revoked)
	if err := r.db.Model(&models.Certification{}).
		Where("status IN ?", []string{string(models.CertificationStatusIssued), string(models.CertificationStatusPending)}).
		Count(&stats.Active).Error; err != nil {
		return nil, err
	}

	// Get revoked certifications
	if err := r.db.Model(&models.Certification{}).
		Where("status = ?", models.CertificationStatusRevoked).
		Count(&stats.Revoked).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// ToResponse converts a Certification model to CertificationResponse DTO
func (r *CertificationRepository) ToResponse(certification *models.Certification) dto.CertificationResponse {
	return dto.CertificationResponse{
		ID:        certification.ID,
		Type:      certification.Type,
		Recipient: certification.Recipient,
		IssueDate: certification.IssueDate,
		PackageID: certification.PackageID,
		Status:    string(certification.Status),
		CreatedAt: certification.CreatedAt,
		UpdatedAt: certification.UpdatedAt,
	}
}

// ToListResponse converts a slice of Certifications to CertificationListResponse DTO
func (r *CertificationRepository) ToListResponse(certifications []models.Certification, total int64, page, limit int) dto.CertificationListResponse {
	items := make([]dto.CertificationResponse, len(certifications))
	for i, c := range certifications {
		items[i] = r.ToResponse(&c)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.CertificationListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}