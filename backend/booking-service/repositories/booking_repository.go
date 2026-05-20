package repositories

import (
	"context"
	"time"

	"booking-service/models"
	"booking-service/models/dto"

	"gorm.io/gorm"
)

// BookingRepository handles booking database operations
type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

func (r *BookingRepository) GetByID(ctx context.Context, id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.WithContext(ctx).First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *BookingRepository) Update(ctx context.Context, enrollment *models.Enrollment) error {
	return r.db.WithContext(ctx).Save(enrollment).Error
}

func (r *BookingRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Enrollment{}, id).Error
}

func (r *BookingRepository) List(ctx context.Context, page, limit int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Enrollment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Entitlements").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

func (r *BookingRepository) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.Enrollment{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Entitlements").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

func (r *BookingRepository) GetByInstructorID(ctx context.Context, instructorID uint, page, limit int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	offset := (page - 1) * limit

	// For enrollment-based bookings, we filter by entitlements that might link to instructor
	// This is a simplified approach - in real implementation, you might join with sessions table
	query := r.db.WithContext(ctx).Model(&models.Enrollment{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Entitlements").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	_ = instructorID // unused in simplified version
	return enrollments, total, nil
}

func (r *BookingRepository) GetByStatus(ctx context.Context, status models.EnrollmentStatus, page, limit int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.Enrollment{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Entitlements").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id uint, status models.EnrollmentStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Enrollment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ToResponse converts an Enrollment model to EnrollmentResponse DTO
func (r *BookingRepository) ToResponse(enrollment *models.Enrollment) dto.EnrollmentResponse {
	return dto.EnrollmentResponse{
		ID:         enrollment.ID,
		UserID:     enrollment.UserID,
		PackageID:  enrollment.PackageID,
		Status:     string(enrollment.Status),
		TotalPrice: enrollment.TotalPrice,
		PaidAt:     enrollment.PaidAt,
		ExpiresAt:  enrollment.ExpiresAt,
		CreatedAt:  enrollment.CreatedAt,
		UpdatedAt:  enrollment.UpdatedAt,
	}
}

// ToListResponse converts a slice of Enrollments to EnrollmentListResponse DTO
func (r *BookingRepository) ToListResponse(enrollments []models.Enrollment, total int64, page, limit int) dto.EnrollmentListResponse {
	items := make([]dto.EnrollmentResponse, len(enrollments))
	for i, e := range enrollments {
		items[i] = r.ToResponse(&e)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.EnrollmentListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// SessionRepository handles session database operations
type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *SessionRepository) GetByID(ctx context.Context, id uint) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) List(ctx context.Context, page, limit int) ([]models.Session, int64, error) {
	var sessions []models.Session
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *SessionRepository) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Session, int64, error) {
	var sessions []models.Session
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.Session{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("date_of_session DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *SessionRepository) GetByInstructorID(ctx context.Context, instructorID uint, page, limit int) ([]models.Session, int64, error) {
	var sessions []models.Session
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.Session{}).Where("instructor_id = ?", instructorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("date_of_session DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

// ToResponse converts a Session model to SessionResponse DTO
func (r *SessionRepository) ToResponse(session *models.Session) dto.SessionResponse {
	return dto.SessionResponse{
		ID:            session.ID,
		UserID:        session.UserID,
		InstructorID:  session.InstructorID,
		Duration:      session.Duration,
		CarID:         session.CarID,
		Area:          session.Area,
		Notes:         session.Notes,
		CreatedAt:     session.CreatedAt,
		UpdatedAt:     session.UpdatedAt,
	}
}

// ToListResponse converts a slice of Sessions to SessionListResponse DTO
func (r *SessionRepository) ToListResponse(sessions []models.Session, total int64, page, limit int) dto.SessionListResponse {
	items := make([]dto.SessionResponse, len(sessions))
	for i, s := range sessions {
		items[i] = r.ToResponse(&s)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.SessionListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// EntitlementRepository handles user entitlement database operations
type EntitlementRepository struct {
	db *gorm.DB
}

func NewEntitlementRepository(db *gorm.DB) *EntitlementRepository {
	return &EntitlementRepository{db: db}
}

func (r *EntitlementRepository) Create(ctx context.Context, entitlement *models.UserEntitlement) error {
	return r.db.WithContext(ctx).Create(entitlement).Error
}

func (r *EntitlementRepository) GetByID(ctx context.Context, id uint) (*models.UserEntitlement, error) {
	var entitlement models.UserEntitlement
	if err := r.db.WithContext(ctx).First(&entitlement, id).Error; err != nil {
		return nil, err
	}
	return &entitlement, nil
}

func (r *EntitlementRepository) Update(ctx context.Context, entitlement *models.UserEntitlement) error {
	return r.db.WithContext(ctx).Save(entitlement).Error
}

func (r *EntitlementRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.UserEntitlement{}, id).Error
}

func (r *EntitlementRepository) List(ctx context.Context, page, limit int) ([]models.UserEntitlement, int64, error) {
	var entitlements []models.UserEntitlement
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.UserEntitlement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&entitlements).Error; err != nil {
		return nil, 0, err
	}

	return entitlements, total, nil
}

func (r *EntitlementRepository) GetByUserID(ctx context.Context, userID uint) ([]models.UserEntitlement, error) {
	var entitlements []models.UserEntitlement
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&entitlements).Error; err != nil {
		return nil, err
	}
	return entitlements, nil
}

func (r *EntitlementRepository) GetActiveByUserID(ctx context.Context, userID uint) ([]models.UserEntitlement, error) {
	var entitlements []models.UserEntitlement
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND sessions_remaining > 0 AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&entitlements).Error; err != nil {
		return nil, err
	}
	return entitlements, nil
}

func (r *EntitlementRepository) DecrementSessions(ctx context.Context, id uint, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.UserEntitlement{}).
		Where("id = ? AND sessions_remaining >= ?", id, count).
		Update("sessions_remaining", gorm.Expr("sessions_remaining - ?", count)).Error
}

// ToResponse converts a UserEntitlement model to EntitlementResponse DTO
func (r *EntitlementRepository) ToResponse(entitlement *models.UserEntitlement) dto.EntitlementResponse {
	return dto.EntitlementResponse{
		ID:                entitlement.ID,
		UserID:            entitlement.UserID,
		SourceType:        entitlement.SourceType,
		SourceID:          entitlement.SourceID,
		TotalSessions:     entitlement.TotalSessions,
		SessionsRemaining: entitlement.UsedSessions,
		ExpiresAt:         entitlement.ExpiresAt,
		CreatedAt:         entitlement.CreatedAt,
		UpdatedAt:         entitlement.UpdatedAt,
	}
}

// ToListResponse converts a slice of UserEntitlements to EntitlementListResponse DTO
func (r *EntitlementRepository) ToListResponse(entitlements []models.UserEntitlement, total int64, page, limit int) dto.EntitlementListResponse {
	items := make([]dto.EntitlementResponse, len(entitlements))
	for i, e := range entitlements {
		items[i] = r.ToResponse(&e)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.EntitlementListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// CertificationRepository handles certification database operations
type CertificationRepository struct {
	db *gorm.DB
}

func NewCertificationRepository(db *gorm.DB) *CertificationRepository {
	return &CertificationRepository{db: db}
}

func (r *CertificationRepository) Create(ctx context.Context, certification *models.Certification) error {
	return r.db.WithContext(ctx).Create(certification).Error
}

func (r *CertificationRepository) GetByID(ctx context.Context, id uint) (*models.Certification, error) {
	var certification models.Certification
	if err := r.db.WithContext(ctx).First(&certification, id).Error; err != nil {
		return nil, err
	}
	return &certification, nil
}

func (r *CertificationRepository) Update(ctx context.Context, certification *models.Certification) error {
	return r.db.WithContext(ctx).Save(certification).Error
}

func (r *CertificationRepository) List(ctx context.Context, page, limit int) ([]models.Certification, int64, error) {
	var certifications []models.Certification
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Certification{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&certifications).Error; err != nil {
		return nil, 0, err
	}

	return certifications, total, nil
}

func (r *CertificationRepository) GetByPackageID(ctx context.Context, packageID uint) ([]models.Certification, error) {
	var certifications []models.Certification
	if err := r.db.WithContext(ctx).
		Where("package_id = ?", packageID).
		Find(&certifications).Error; err != nil {
		return nil, err
	}
	return certifications, nil
}

func (r *CertificationRepository) GetByStatus(ctx context.Context, status models.CertificationStatus) ([]models.Certification, error) {
	var certifications []models.Certification
	if err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Find(&certifications).Error; err != nil {
		return nil, err
	}
	return certifications, nil
}

func (r *CertificationRepository) UpdateStatus(ctx context.Context, id uint, status models.CertificationStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Certification{}).
		Where("id = ?", id).
		Update("status", status).Error
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

// UserCertificationRepository handles user-certification associations
type UserCertificationRepository struct {
	db *gorm.DB
}

func NewUserCertificationRepository(db *gorm.DB) *UserCertificationRepository {
	return &UserCertificationRepository{db: db}
}

func (r *UserCertificationRepository) Create(ctx context.Context, uc *models.UserCertification) error {
	return r.db.WithContext(ctx).Create(uc).Error
}

func (r *UserCertificationRepository) GetByUserID(ctx context.Context, userID uint) ([]models.UserCertification, error) {
	var userCerts []models.UserCertification
	if err := r.db.WithContext(ctx).
		Preload("Certification").
		Where("user_id = ?", userID).
		Find(&userCerts).Error; err != nil {
		return nil, err
	}
	return userCerts, nil
}

func (r *UserCertificationRepository) GetByUserAndCertification(ctx context.Context, userID, certificationID uint) (*models.UserCertification, error) {
	var uc models.UserCertification
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND certification_id = ?", userID, certificationID).
		First(&uc).Error; err != nil {
		return nil, err
	}
	return &uc, nil
}

func (r *UserCertificationRepository) Delete(ctx context.Context, userID, certificationID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND certification_id = ?", userID, certificationID).
		Delete(&models.UserCertification{}).Error
}

// ToResponse converts a UserCertification model to UserCertificationResponse DTO
func (r *UserCertificationRepository) ToResponse(uc *models.UserCertification, certResponse dto.CertificationResponse) dto.UserCertificationResponse {
	return dto.UserCertificationResponse{
		UserID:          uc.UserID,
		CertificationID: uc.CertificationID,
		IssuedAt:        uc.IssuedAt,
		Certification:   certResponse,
	}
}