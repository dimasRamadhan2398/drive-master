package repositories

import (
	"time"
	"user-service/models"
	"user-service/pkg/base"
	apperrors "user-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserSessionRepository interface {
	Create(userSession *models.UserSession) error
	FindByID(id uuid.UUID) (*models.UserSession, error)
	FindByRefreshToken(token string) (*models.UserSession, error)
	FindActiveByUserID(userID uuid.UUID) ([]models.UserSession, error)
	Update(userSession *models.UserSession) error
	Delete(userSession *models.UserSession) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpired() error
}

type UserSessionRepository struct {
	*base.BaseRepository // ← no separate db field, use r.DB from BaseRepository
}

func NewUserSessionRepository(db *gorm.DB) *UserSessionRepository {
	return &UserSessionRepository{
		BaseRepository: base.NewBaseRepository(db),
	}
}

// Create creates a new user session
func (r *UserSessionRepository) Create(userSession *models.UserSession) error {
	return r.BaseRepository.Create(userSession)
}

// FindByID finds a session by UUID
func (r *UserSessionRepository) FindByID(id uuid.UUID) (*models.UserSession, error) {
	var userSession models.UserSession // ← fixed typo: UserSEssion → UserSession
	if err := r.BaseRepository.FindByID(&userSession, id); err != nil {
		return nil, err
	}
	return &userSession, nil // ← fixed: was returning &user which doesn't exist
}

// FindByRefreshToken finds a session by its refresh token
func (r *UserSessionRepository) FindByRefreshToken(token string) (*models.UserSession, error) {
	var userSession models.UserSession
	if err := r.BaseRepository.FindOne(&userSession, "refresh_token = ?", token); err != nil {
		return nil, err
	}
	return &userSession, nil
}

// FindActiveByUserID finds all non-expired sessions for a user
func (r *UserSessionRepository) FindActiveByUserID(userID uuid.UUID) ([]models.UserSession, error) {
	var sessions []models.UserSession
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"user_id": userID}).
		WithOrder("created_at DESC")

	if err := r.BaseRepository.FindMany(&models.UserSession{}, &sessions, opts); err != nil {
		return nil, err
	}

	// Filter out expired sessions in memory
	active := make([]models.UserSession, 0)
	for _, s := range sessions {
		if !s.IsExpired() {
			active = append(active, s)
		}
	}
	return active, nil
}

// Update saves changes to an existing session
func (r *UserSessionRepository) Update(userSession *models.UserSession) error {
	return r.BaseRepository.Update(userSession) // ← fixed: was passing undefined `user`
}

// Delete soft-deletes a specific session
func (r *UserSessionRepository) Delete(userSession *models.UserSession) error {
	return r.BaseRepository.Delete(userSession) // ← fixed: BaseRepository.Delete needs the entity
}

// DeleteByUserID invalidates all sessions for a user (e.g. on logout all devices)
func (r *UserSessionRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.BaseRepository.Exec(
		"DELETE FROM user_sessions WHERE user_id = ?", userID,
	)
}

// DeleteExpired removes all expired sessions (for cleanup jobs)
func (r *UserSessionRepository) DeleteExpired() error {
	return r.BaseRepository.Exec(
		"DELETE FROM user_sessions WHERE expires_at < ?", time.Now(),
	)
}

// InvalidateByRefreshToken marks a specific refresh token session as expired
func (r *UserSessionRepository) InvalidateByRefreshToken(token string) error {
	session, err := r.FindByRefreshToken(token)
	if err != nil {
		return apperrors.ErrNotFound
	}
	now := time.Now()
	session.ExpiresAt = now
	return r.BaseRepository.Update(session)
}