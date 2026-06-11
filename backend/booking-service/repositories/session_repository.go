package repositories

import (
	"context"
	"time"

	"booking-service/models"
	"booking-service/models/dto"

	"gorm.io/gorm"
)

// SessionRepository handles session database operations
type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.DrivingSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *SessionRepository) GetByID(ctx context.Context, id uint) (*models.DrivingSession, error) {
	var session models.DrivingSession
	if err := r.db.WithContext(ctx).First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) Update(ctx context.Context, session *models.DrivingSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *SessionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.DrivingSession{}, id).Error
}

func (r *SessionRepository) List(ctx context.Context, page, limit int) ([]models.DrivingSession, int64, error) {
	var sessions []models.DrivingSession
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.DrivingSession{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("date DESC, time DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *SessionRepository) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.DrivingSession, int64, error) {
	var sessions []models.DrivingSession
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.DrivingSession{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("date DESC, time DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *SessionRepository) GetByInstructorID(ctx context.Context, instructorID uint, page, limit int) ([]models.DrivingSession, int64, error) {
	var sessions []models.DrivingSession
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.DrivingSession{}).Where("instructor_id = ?", instructorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("date DESC, time DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *SessionRepository) GetByEnrollmentID(ctx context.Context, enrollmentID uint) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	if err := r.db.WithContext(ctx).
		Where("enrollment_id = ?", enrollmentID).
		Order("date DESC, time DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) GetByEntitlementID(ctx context.Context, entitlementID uint) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	if err := r.db.WithContext(ctx).
		Where("entitlement_id = ?", entitlementID).
		Order("date DESC, time DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) GetByStatus(ctx context.Context, status string, page, limit int) ([]models.DrivingSession, int64, error) {
	var sessions []models.DrivingSession
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.DrivingSession{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("date DESC, time DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *SessionRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]models.DrivingSession, int64, error) {
	var sessions []models.DrivingSession
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.DrivingSession{}).
		Where("date >= ? AND date <= ?", startDate, endDate)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("date DESC, time DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.DrivingSession{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *SessionRepository) StartSession(ctx context.Context, id uint, startedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.DrivingSession{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "in_progress",
			"started_at": startedAt,
		}).Error
}

func (r *SessionRepository) CompleteSession(ctx context.Context, id uint, completedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.DrivingSession{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": completedAt,
		}).Error
}

// ToResponse converts a DrivingSession model to DrivingSessionResponse DTO
func (r *SessionRepository) ToResponse(session *models.DrivingSession) dto.DrivingSessionResponse {
	return dto.DrivingSessionResponse{
		ID:            session.ID,
		EnrollmentID:  session.EnrollmentID,
		EntitlementID: session.EntitlementID,
		UserID:        session.UserID,
		InstructorID:  session.InstructorID,
		CarID:         session.CarID,
		ScheduleID:    session.ScheduleID,
		Date:          session.Date.Format("2006-01-02"),
		Time:          session.Time,
		Duration:      session.Duration,
		Status:        session.Status,
		Area:          session.Area,
		Notes:         session.Notes,
		StartedAt:     session.StartedAt,
		CompletedAt:   session.CompletedAt,
		CreatedAt:     session.CreatedAt,
		UpdatedAt:     session.UpdatedAt,
	}
}

// SessionStats holds session statistics
type SessionStats struct {
	Total     int64
	Active    int64
	Completed int64
	Pending   int64
}

func (r *SessionRepository) GetStats(ctx context.Context) (*SessionStats, error) {
	stats := &SessionStats{}
	today := time.Now().Truncate(24 * time.Hour)

	// Get total sessions
	if err := r.db.WithContext(ctx).Model(&models.DrivingSession{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Get active sessions (today or future)
	if err := r.db.WithContext(ctx).Model(&models.DrivingSession{}).
		Where("date >= ?", today).
		Where("status != ?", "completed").
		Count(&stats.Active).Error; err != nil {
		return nil, err
	}

	// Get completed sessions
	if err := r.db.WithContext(ctx).Model(&models.DrivingSession{}).
		Where("status = ?", "completed").
		Count(&stats.Completed).Error; err != nil {
		return nil, err
	}

	// Get pending sessions
	if err := r.db.WithContext(ctx).Model(&models.DrivingSession{}).
		Where("status = ?", "scheduled").
		Count(&stats.Pending).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// ToListResponse converts a slice of DrivingSessions to DrivingSessionListResponse DTO
func (r *SessionRepository) ToListResponse(sessions []models.DrivingSession, total int64, page, limit int) dto.DrivingSessionListResponse {
	items := make([]dto.DrivingSessionResponse, len(sessions))
	for i, s := range sessions {
		items[i] = r.ToResponse(&s)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.DrivingSessionListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
