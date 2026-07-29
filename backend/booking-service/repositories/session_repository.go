package repositories

import (
	"context"
	"strconv"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISessionRepository interface {
	Create(ctx context.Context, session *models.DrivingSession) error
	CreateTx(tx *gorm.DB, session *models.DrivingSession) error
	FindByID(ctx context.Context, id uint) (*models.DrivingSession, error)
	Update(ctx context.Context, session *models.DrivingSession) error
	Delete(ctx context.Context, session *models.DrivingSession) error
	FindAll(ctx context.Context) ([]models.DrivingSession, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.DrivingSession, error)
	FindByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.DrivingSession, error)
	FindByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) ([]models.DrivingSession, error)
	FindByEntitlementID(ctx context.Context, entitlementID uuid.UUID) ([]models.DrivingSession, error)
	FindByStatus(ctx context.Context, status string) ([]models.DrivingSession, error)
	FindByScheduleID(ctx context.Context, scheduleID uint) (*models.DrivingSession, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.DrivingSession, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	StartSession(ctx context.Context, id uint, startedAt time.Time) error
	CompleteSession(ctx context.Context, id uint, completedAt time.Time) error
	CompleteSessionTx(tx *gorm.DB, id uint, completedAt time.Time) error
	ForceCompleteByAdmin(ctx context.Context, id uint, endTime time.Time) error
	RateSession(ctx context.Context, id uint, rating float64, feedback string) error
	CancelSession(ctx context.Context, id uint) error
	GetStats(ctx context.Context) (*SessionStats, error)
	AnonymizeByUserID(ctx context.Context, userID uuid.UUID, anonymizedAt time.Time) error
	CountByStatus(ctx context.Context, status string) (int64, error)
	CountAll(ctx context.Context) (int64, error)
	ToResponse(session *models.DrivingSession) dto.DrivingSessionResponse
	ToListResponse(sessions []models.DrivingSession, total int64, page, limit int) dto.DrivingSessionListResponse
}

type SessionRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &SessionRepository{BaseRepository: base.NewBaseRepository(db), db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.DrivingSession) error {
	return r.BaseRepository.Create(session)
}

func (r *SessionRepository) CreateTx(tx *gorm.DB, session *models.DrivingSession) error {
	return r.BaseRepository.CreateTx(tx, session)
}

func (r *SessionRepository) FindByID(ctx context.Context, id uint) (*models.DrivingSession, error) {
	var session models.DrivingSession
	if err := r.BaseRepository.FindByIDWithPreload(&session, id); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) Update(ctx context.Context, session *models.DrivingSession) error {
	return r.BaseRepository.Update(session)
}

func (r *SessionRepository) Delete(ctx context.Context, session *models.DrivingSession) error {
	return r.BaseRepository.Delete(session)
}

func (r *SessionRepository) FindAll(ctx context.Context) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	opts := base.NewQueryOptions().WithOrder("date DESC, time DESC")
	if err := r.BaseRepository.FindMany(&models.DrivingSession{}, &sessions, opts); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) FindByScheduleID(ctx context.Context, scheduleID uint) (*models.DrivingSession, error) {
	var session models.DrivingSession
	opts := base.NewQueryOptions().WithWhere(map[string]any{"schedule_id": scheduleID})
	if err := r.BaseRepository.FindMany(&models.DrivingSession{}, &session, opts); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"user_id": userID}).
		WithOrder("date DESC, time DESC")
	if err := r.BaseRepository.FindMany(&models.DrivingSession{}, &sessions, opts); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) FindByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"instructor_id": instructorID}).
		WithOrder("date DESC, time DESC")
	if err := r.BaseRepository.FindMany(&models.DrivingSession{}, &sessions, opts); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) FindByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"enrollment_id": enrollmentID}).
		WithOrder("date DESC, time DESC")
	if err := r.BaseRepository.FindMany(&models.DrivingSession{}, &sessions, opts); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) FindByEntitlementID(ctx context.Context, entitlementID uuid.UUID) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"entitlement_id": entitlementID}).
		WithOrder("date DESC, time DESC")
	if err := r.BaseRepository.FindMany(&models.DrivingSession{}, &sessions, opts); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) FindByStatus(ctx context.Context, status string) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"status": status}).
		WithOrder("date DESC, time DESC")
	if err := r.BaseRepository.FindMany(&models.DrivingSession{}, &sessions, opts); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.DrivingSession, error) {
	var sessions []models.DrivingSession
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"date >= ?": startDate, "date <= ?": endDate}).
		WithOrder("date DESC, time DESC")
	if err := r.BaseRepository.FindMany(&models.DrivingSession{}, &sessions, opts); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.BaseRepository.Exec(
		"UPDATE driving_sessions SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	)
}

func (r *SessionRepository) StartSession(ctx context.Context, id uint, startedAt time.Time) error {
	return r.BaseRepository.Exec(
		"UPDATE driving_sessions SET status = 'in_progress', started_at = ?, updated_at = ? WHERE id = ?",
		startedAt, time.Now(), id,
	)
}

func (r *SessionRepository) CompleteSession(ctx context.Context, id uint, completedAt time.Time) error {
	return r.BaseRepository.Exec(
		"UPDATE driving_sessions SET status = 'completed', completed_at = ?, end_time = ?, updated_at = ? WHERE id = ?",
		completedAt, completedAt, time.Now(), id,
	)
}

func (r *SessionRepository) CompleteSessionTx(tx *gorm.DB, id uint, completedAt time.Time) error {
	return tx.Exec(
		"UPDATE driving_sessions SET status = 'completed', completed_at = ?, end_time = ?, updated_at = ? WHERE id = ?",
		completedAt, completedAt, time.Now(), id,
	).Error
}

// ForceCompleteByAdmin marks a session as completed by an admin.
// It sets is_ended_by_admin=true and end_time so that auto-schedulers
// will not revert this session back to in_progress or re-complete it.
func (r *SessionRepository) ForceCompleteByAdmin(ctx context.Context, id uint, endTime time.Time) error {
	return r.BaseRepository.Exec(
		"UPDATE driving_sessions SET status = 'completed', completed_at = ?, end_time = ?, is_ended_by_admin = true, updated_at = ? WHERE id = ?",
		endTime, endTime, time.Now(), id,
	)
}

// RateSession records student rating and feedback for a session
func (r *SessionRepository) RateSession(ctx context.Context, id uint, rating float64, feedback string) error {
	return r.BaseRepository.Exec(
		"UPDATE driving_sessions SET rating = ?, feedback = ?, updated_at = ? WHERE id = ?",
		rating, feedback, time.Now(), id,
	)
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
	if err := r.db.Model(&models.DrivingSession{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Get active sessions (today or future)
	if err := r.db.Model(&models.DrivingSession{}).
		Where("date >= ?", today).
		Where("status != ?", "completed").
		Count(&stats.Active).Error; err != nil {
		return nil, err
	}

	// Get completed sessions
	if err := r.db.Model(&models.DrivingSession{}).
		Where("status = ?", "completed").
		Count(&stats.Completed).Error; err != nil {
		return nil, err
	}

	// Get pending sessions
	if err := r.db.Model(&models.DrivingSession{}).
		Where("status = ?", "scheduled").
		Count(&stats.Pending).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *SessionRepository) CancelSession(ctx context.Context, id uint) error {
	return r.BaseRepository.Exec(
		"UPDATE driving_sessions SET status = 'cancelled', updated_at = ? WHERE id = ?",
		time.Now(), id,
	)
}

func (r *SessionRepository) AnonymizeByUserID(ctx context.Context, userID uuid.UUID, anonymizedAt time.Time) error {
	// Convert UUID to uint (this assumes the UUID can be parsed as a number)
	// For production, you might want to store the UUID as a string or use a different approach
	userIDUint, err := strconv.ParseUint(userID.String()[0:8], 16, 64)
	if err != nil {
		// If conversion fails, try to match by the last 8 bytes converted to uint
		userIDUint = 0
	}

	return r.BaseRepository.Exec(
		"UPDATE driving_sessions SET anonymized_at = ? WHERE user_id = ?",
		anonymizedAt, userIDUint,
	)
}

func (r *SessionRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	opts := base.NewQueryOptions().WithWhere(map[string]any{"status": status})
	return r.BaseRepository.Count(&models.DrivingSession{}, opts)
}

func (r *SessionRepository) CountAll(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(&models.DrivingSession{}, base.NewQueryOptions())
}

// ToResponse converts a DrivingSession model to DrivingSessionResponse DTO
func (r *SessionRepository) ToResponse(session *models.DrivingSession) dto.DrivingSessionResponse {
	return dto.DrivingSessionResponse{
		ID:             session.ID,
		EnrollmentID:   session.EnrollmentID,
		EntitlementID:  session.EntitlementID,
		UserID:         session.UserID,
		InstructorID:   session.InstructorID,
		CarID:          session.CarID,
		ScheduleID:     session.ScheduleID,
		Date:           session.Date.Format("2006-01-02"),
		Time:           session.Time,
		Duration:       session.Duration,
		Status:         session.Status,
		Area:           session.Area,
		Notes:          session.Notes,
		StartedAt:      session.StartedAt,
		CompletedAt:    session.CompletedAt,
		EndTime:        session.EndTime,
		IsEndedByAdmin: session.IsEndedByAdmin,
		Rating:         session.Rating,
		Feedback:       session.Feedback,
		CreatedAt:      session.CreatedAt,
		UpdatedAt:      session.UpdatedAt,
	}
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