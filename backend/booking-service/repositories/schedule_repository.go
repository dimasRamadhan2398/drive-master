package repositories

import (
	"context"
	"time"

	"booking-service/models"
	"booking-service/models/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduleRepository handles schedule database operations
type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id uint) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.WithContext(ctx).First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) Update(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

func (r *ScheduleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Schedule{}, id).Error
}

func (r *ScheduleRepository) List(ctx context.Context, page, limit int) ([]models.Schedule, int64, error) {
	var schedules []models.Schedule
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Schedule{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("date ASC, time ASC").
		Offset(offset).
		Limit(limit).
		Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *ScheduleRepository) ListFiltered(ctx context.Context, params dto.ScheduleFilterParams) ([]models.Schedule, int64, error) {
	var schedules []models.Schedule
	var total int64

	offset := (params.Page - 1) * params.Limit
	if params.Limit <= 0 {
		params.Limit = 10
	}

	query := r.db.WithContext(ctx).Model(&models.Schedule{})

	// Apply filters
	if params.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", params.Date)
		if err == nil {
			query = query.Where("date = ?", parsedDate)
		}
	}

	if params.StartDate != "" {
		parsedDate, err := time.Parse("2006-01-02", params.StartDate)
		if err == nil {
			query = query.Where("date >= ?", parsedDate)
		}
	}

	if params.EndDate != "" {
		parsedDate, err := time.Parse("2006-01-02", params.EndDate)
		if err == nil {
			query = query.Where("date <= ?", parsedDate)
		}
	}

	if params.InstructorID != "" {
		query = query.Where("instructor_id = ?", params.InstructorID)
	}

	if params.CarID != 0 {
		query = query.Where("car_id = ?", params.CarID)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("date ASC, time ASC").
		Offset(offset).
		Limit(params.Limit).
		Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *ScheduleRepository) GetByDateAndInstructor(ctx context.Context, date time.Time, instructorID uuid.UUID) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := r.db.WithContext(ctx).
		Where("date = ? AND instructor_id = ?", date, instructorID).
		Order("time ASC").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) GetByDateAndTime(ctx context.Context, date time.Time, time string, instructorID uuid.UUID, carID uint) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.WithContext(ctx).
		Where("date = ? AND time = ? AND instructor_id = ? AND car_id = ?", date, time, instructorID, carID).
		First(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) GetAvailableByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := r.db.WithContext(ctx).
		Where("date >= ? AND date <= ? AND status = ?", startDate, endDate, models.ScheduleStatusAvailable).
		Order("date ASC, time ASC").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) UpdateStatus(ctx context.Context, id uint, status models.ScheduleStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Schedule{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *ScheduleRepository) BookSlot(ctx context.Context, id uint, userID, enrollmentID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Schedule{}).
		Where("id = ? AND status = ?", id, models.ScheduleStatusAvailable).
		Updates(map[string]interface{}{
			"user_id":       userID,
			"enrollment_id": enrollmentID,
			"status":        models.ScheduleStatusBooked,
		}).Error
}

func (r *ScheduleRepository) ReleaseSlot(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Schedule{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"user_id":       nil,
			"enrollment_id": nil,
			"status":        models.ScheduleStatusAvailable,
		}).Error
}

// ToResponse converts a Schedule model to ScheduleResponse DTO
func (r *ScheduleRepository) ToResponse(schedule *models.Schedule) dto.ScheduleResponse {
	dateStr := ""
	if !schedule.Date.IsZero() {
		dateStr = schedule.Date.Format("2006-01-02")
	}

	return dto.ScheduleResponse{
		ID:           schedule.ID,
		Date:         dateStr,
		Time:         schedule.Time,
		Duration:     schedule.Duration,
		InstructorID: schedule.InstructorID,
		CarID:        schedule.CarID,
		UserID:       schedule.UserID,
		Status:       string(schedule.Status),
		Notes:        schedule.Notes,
		CreatedAt:    schedule.CreatedAt,
		UpdatedAt:    schedule.UpdatedAt,
	}
}

// ToListResponse converts a slice of Schedules to ScheduleListResponse DTO
func (r *ScheduleRepository) ToListResponse(schedules []models.Schedule, total int64, page, limit int) dto.ScheduleListResponse {
	items := make([]dto.ScheduleResponse, len(schedules))
	for i, s := range schedules {
		items[i] = r.ToResponse(&s)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.ScheduleListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// ExistsForInstructorAndDateTime checks if a schedule slot already exists
func (r *ScheduleRepository) ExistsForInstructorAndDateTime(ctx context.Context, instructorID uuid.UUID, date time.Time, timeStr string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Schedule{}).
		Where("instructor_id = ? AND date = ? AND time = ?", instructorID, date, timeStr).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}