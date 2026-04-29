package repositories

import (
	"notification-service/models"
	"notification-service/pkg/base"
	apperrors "notification-service/pkg/errors"

	"gorm.io/gorm"
)

// INotificationRepository defines the interface for notification repository operations
type INotificationRepository interface {
	// Notification CRUD
	CreateNotification(notification *models.Notification) error
	GetNotificationByID(id uint) (*models.Notification, error)
	GetNotificationsByUserID(userID uint, limit, offset int) ([]models.Notification, int64, error)
	GetNotificationsByStatus(status models.NotificationStatus, limit int) ([]models.Notification, error)
	GetNotificationsByCategory(category models.NotificationCategory, limit int) ([]models.Notification, error)
	UpdateNotificationStatus(id uint, status models.NotificationStatus) error
	MarkNotificationAsSent(id uint) error
	MarkNotificationAsFailed(id uint, errorMsg string) error
	IncrementRetryCount(id uint) error

	// Template operations
	CreateTemplate(template *models.NotificationTemplate) error
	GetTemplateByID(id uint) (*models.NotificationTemplate, error)
	GetTemplateByName(name string) (*models.NotificationTemplate, error)
	GetTemplatesByCategory(category models.NotificationCategory) ([]models.NotificationTemplate, error)
	UpdateTemplate(template *models.NotificationTemplate) error

	// Newsletter subscription
	Subscribe(email, name, source string) (*models.NewsletterSubscription, error)
	Unsubscribe(email string) error
	GetSubscriptionByEmail(email string) (*models.NewsletterSubscription, error)
	GetActiveSubscriptions(limit, offset int) ([]models.NewsletterSubscription, int64, error)

	// Preferences
	GetOrCreatePreference(userID uint) (*models.NotificationPreference, error)
	UpdatePreference(preference *models.NotificationPreference) error
	GetPreferenceByUserID(userID uint) (*models.NotificationPreference, error)
}

// NotificationRepository implements INotificationRepository
type NotificationRepository struct {
	*base.BaseRepository
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{BaseRepository: base.NewBaseRepository(db)}
}

// CreateNotification creates a new notification
func (r *NotificationRepository) CreateNotification(notification *models.Notification) error {
	return r.Create(notification)
}

// GetNotificationByID retrieves a notification by ID
func (r *NotificationRepository) GetNotificationByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := r.GetDB().First(&notification, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &notification, nil
}

// GetNotificationsByUserID retrieves notifications for a user with pagination
func (r *NotificationRepository) GetNotificationsByUserID(userID uint, limit, offset int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.GetDB().Model(&models.Notification{}).Where("user_id = ?", userID)
	query.Count(&total)

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}
	return notifications, total, nil
}

// GetNotificationsByStatus retrieves notifications by status
func (r *NotificationRepository) GetNotificationsByStatus(status models.NotificationStatus, limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := r.GetDB().Where("status = ?", status).Order("scheduled_at ASC").Limit(limit).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

// GetNotificationsByCategory retrieves notifications by category
func (r *NotificationRepository) GetNotificationsByCategory(category models.NotificationCategory, limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := r.GetDB().Where("category = ?", category).Order("created_at DESC").Limit(limit).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

// UpdateNotificationStatus updates the status of a notification
func (r *NotificationRepository) UpdateNotificationStatus(id uint, status models.NotificationStatus) error {
	return r.GetDB().Model(&models.Notification{}).Where("id = ?", id).Update("status", status).Error
}

// MarkNotificationAsSent marks a notification as sent
func (r *NotificationRepository) MarkNotificationAsSent(id uint) error {
	sql := "UPDATE notifications SET status = 'sent', sent_at = NOW() WHERE id = ?"
	return r.Exec(sql, id)
}

// MarkNotificationAsFailed marks a notification as failed
func (r *NotificationRepository) MarkNotificationAsFailed(id uint, errorMsg string) error {
	sql := "UPDATE notifications SET status = 'failed', failed_at = NOW(), error_message = ? WHERE id = ?"
	return r.Exec(sql, errorMsg, id)
}

// IncrementRetryCount increments the retry count for a notification
func (r *NotificationRepository) IncrementRetryCount(id uint) error {
	return r.GetDB().Model(&models.Notification{}).Where("id = ?").Update("retry_count", gorm.Expr("retry_count + 1")).Error
}

// CreateTemplate creates a new notification template
func (r *NotificationRepository) CreateTemplate(template *models.NotificationTemplate) error {
	return r.Create(template)
}

// GetTemplateByID retrieves a template by ID
func (r *NotificationRepository) GetTemplateByID(id uint) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	if err := r.GetDB().First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &template, nil
}

// GetTemplateByName retrieves a template by name
func (r *NotificationRepository) GetTemplateByName(name string) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	if err := r.GetDB().Where("name = ? AND is_active = true", name).First(&template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &template, nil
}

// GetTemplatesByCategory retrieves templates by category
func (r *NotificationRepository) GetTemplatesByCategory(category models.NotificationCategory) ([]models.NotificationTemplate, error) {
	var templates []models.NotificationTemplate
	if err := r.GetDB().Where("category = ? AND is_active = true", category).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// UpdateTemplate updates a notification template
func (r *NotificationRepository) UpdateTemplate(template *models.NotificationTemplate) error {
	return r.Update(template)
}

// Subscribe adds a new newsletter subscription
func (r *NotificationRepository) Subscribe(email, name, source string) (*models.NewsletterSubscription, error) {
	// Check if already subscribed
	var existing models.NewsletterSubscription
	err := r.GetDB().Where("email = ?", email).First(&existing).Error
	if err == nil {
		if existing.IsActive {
			return nil, apperrors.ErrSubscriptionExists
		}
		// Reactivate subscription
		existing.IsActive = true
		existing.UnsubscribedAt = nil
		existing.SubscribedAt = existing.CreatedAt
		if err := r.Update(&existing); err != nil {
			return nil, err
		}
		return &existing, nil
	}

	subscription := &models.NewsletterSubscription{
		Email: email,
		Name: name,
		IsActive: true,
		Source: source,
	}
	if err := r.Create(subscription); err != nil {
		return nil, err
	}
	return subscription, nil
}

// Unsubscribe unsubscribes an email from the newsletter
func (r *NotificationRepository) Unsubscribe(email string) error {
	now := models.Notification{}.CreatedAt
	return r.GetDB().Model(&models.NewsletterSubscription{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"is_active": false,
			"unsubscribed_at": now,
		}).Error
}

// GetSubscriptionByEmail retrieves a subscription by email
func (r *NotificationRepository) GetSubscriptionByEmail(email string) (*models.NewsletterSubscription, error) {
	var subscription models.NewsletterSubscription
	if err := r.GetDB().Where("email = ?", email).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

// GetActiveSubscriptions retrieves all active subscriptions with pagination
func (r *NotificationRepository) GetActiveSubscriptions(limit, offset int) ([]models.NewsletterSubscription, int64, error) {
	var subscriptions []models.NewsletterSubscription
	var total int64

	query := r.GetDB().Model(&models.NewsletterSubscription{}).Where("is_active = true")
	query.Count(&total)

	if err := query.Order("subscribed_at DESC").Limit(limit).Offset(offset).Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}
	return subscriptions, total, nil
}

// GetOrCreatePreference gets or creates notification preferences for a user
func (r *NotificationRepository) GetOrCreatePreference(userID uint) (*models.NotificationPreference, error) {
	var preference models.NotificationPreference
	if err := r.GetDB().Where("user_id = ?", userID).First(&preference).Error; err == nil {
		return &preference, nil
	}

	// Create new preference
	preference = models.NotificationPreference{
		UserID: userID,
	}
	if err := r.Create(&preference); err != nil {
		return nil, err
	}
	return &preference, nil
}

// UpdatePreference updates notification preferences
func (r *NotificationRepository) UpdatePreference(preference *models.NotificationPreference) error {
	return r.Update(preference)
}

// GetPreferenceByUserID retrieves notification preferences by user ID
func (r *NotificationRepository) GetPreferenceByUserID(userID uint) (*models.NotificationPreference, error) {
	var preference models.NotificationPreference
	if err := r.GetDB().Where("user_id = ?", userID).First(&preference).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &preference, nil
}