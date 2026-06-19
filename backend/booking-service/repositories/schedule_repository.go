package repositories

import (
	"context"
	"time"

	"booking-service/models/dto"
	"booking-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IScheduleRepository interface {
	Create(ctx context.Context, schedule *dto.Schedule) error
	CreateTx(tx *gorm.DB, schedule *dto.Schedule) error
	FindByID(ctx context.Context, id uint) (*dto.Schedule, error)
	Update(ctx context.Context, schedule *dto.Schedule) error
	Delete(ctx context.Context, schedule *dto.Schedule) error
	FindAll(ctx context.Context) ([]dto.Schedule, error)
	FindByDateAndInstructor(ctx context.Context, date time.Time, instructorID uuid.UUID) ([]dto.Schedule, error)
	FindByDateAndTime(ctx context.Context, date time.Time, time string, instructorID uuid.UUID, carID uuid.UUID) (*dto.Schedule, error)
	FindAvailableByDateRange(ctx context.Context, startDate, endDate time.Time) ([]dto.Schedule, error)
	UpdateStatus(ctx context.Context, id uint, status dto.ScheduleStatus) error
	BookSlot(ctx context.Context, id uint, userID, enrollmentID uuid.UUID) error
	ReleaseSlot(ctx context.Context, id uint) error
	ExistsForInstructorAndDateTime(ctx context.Context, instructorID uuid.UUID, date time.Time, timeStr string) (bool, error)
	CountAll(ctx context.Context) (int64, error)
	GetStats(ctx context.Context) (*dto.ScheduleStatsResponse, error)
	ToResponse(schedule *dto.Schedule) dto.ScheduleResponse
	ToListResponse(schedules []dto.Schedule, total int64, page, limit int) dto.ScheduleListResponse
}

type ScheduleRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) IScheduleRepository {
	return &ScheduleRepository{BaseRepository: base.NewBaseRepository(db), db: db}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule *dto.Schedule) error {
	return r.BaseRepository.Create(schedule)
}

func (r *ScheduleRepository) CreateTx(tx *gorm.DB, schedule *dto.Schedule) error {
	return r.BaseRepository.CreateTx(tx, schedule)
}

func (r *ScheduleRepository) FindByID(ctx context.Context, id uint) (*dto.Schedule, error) {
	var schedule dto.Schedule
	if err := r.BaseRepository.FindByIDWithPreload(&schedule, id); err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) Update(ctx context.Context, schedule *dto.Schedule) error {
	return r.BaseRepository.Update(schedule)
}

func (r *ScheduleRepository) Delete(ctx context.Context, schedule *dto.Schedule) error {
	return r.BaseRepository.Delete(schedule)
}

func (r *ScheduleRepository) FindAll(ctx context.Context) ([]dto.Schedule, error) {
	var schedules []dto.Schedule
	opts := base.NewQueryOptions().WithOrder("date ASC, time ASC")
	if err := r.BaseRepository.FindMany(&dto.Schedule{}, &schedules, opts); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) FindByDateAndInstructor(ctx context.Context, date time.Time, instructorID uuid.UUID) ([]dto.Schedule, error) {
	var schedules []dto.Schedule
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"date": date, "instructor_id": instructorID}).
		WithOrder("time ASC")
	if err := r.BaseRepository.FindMany(&dto.Schedule{}, &schedules, opts); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) FindByDateAndTime(ctx context.Context, date time.Time, time string, instructorID uuid.UUID, carID uuid.UUID) (*dto.Schedule, error) {
	var schedule dto.Schedule
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"date": date, "time": time, "instructor_id": instructorID, "car_id": carID})
	if err := r.BaseRepository.FindMany(&dto.Schedule{}, &schedule, opts); err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) FindAvailableByDateRange(ctx context.Context, startDate, endDate time.Time) ([]dto.Schedule, error) {
	var schedules []dto.Schedule
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"date >= ?": startDate, "date <= ?": endDate, "status": dto.ScheduleStatusAvailable}).
		WithOrder("date ASC, time ASC")
	if err := r.BaseRepository.FindMany(&dto.Schedule{}, &schedules, opts); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) UpdateStatus(ctx context.Context, id uint, status dto.ScheduleStatus) error {
	return r.BaseRepository.Exec(
		"UPDATE schedules SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	)
}

func (r *ScheduleRepository) BookSlot(ctx context.Context, id uint, userID, enrollmentID uuid.UUID) error {
	return r.BaseRepository.Exec(
		"UPDATE schedules SET user_id = ?, enrollment_id = ?, status = ?, updated_at = ? WHERE id = ? AND status = ?",
		userID, enrollmentID, dto.ScheduleStatusBooked, time.Now(), id, dto.ScheduleStatusAvailable,
	)
}

func (r *ScheduleRepository) ReleaseSlot(ctx context.Context, id uint) error {
	return r.BaseRepository.Exec(
		"UPDATE schedules SET user_id = NULL, enrollment_id = NULL, status = ?, updated_at = ? WHERE id = ?",
		dto.ScheduleStatusAvailable, time.Now(), id,
	)
}

func (r *ScheduleRepository) ExistsForInstructorAndDateTime(ctx context.Context, instructorID uuid.UUID, date time.Time, timeStr string) (bool, error) {
	return r.BaseRepository.Exists(&dto.Schedule{}, "instructor_id = ? AND date = ? AND time = ?", instructorID, date, timeStr)
}

func (r *ScheduleRepository) CountAll(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(&dto.Schedule{}, base.NewQueryOptions())
}

func (r *ScheduleRepository) GetStats(ctx context.Context) (*dto.ScheduleStatsResponse, error) {
	stats := &dto.ScheduleStatsResponse{}

	// Get available schedules
	if err := r.db.Model(&dto.Schedule{}).
		Where("status = ?", dto.ScheduleStatusAvailable).
		Count(&stats.AvailableSchedule).Error; err != nil {
		return nil, err
	}

	// Get booked schedules
	if err := r.db.Model(&dto.Schedule{}).
		Where("status = ?", dto.ScheduleStatusBooked).
		Count(&stats.BookedSchedule).Error; err != nil {
		return nil, err
	}

	// Get in-progress schedules
	if err := r.db.Model(&dto.Schedule{}).
		Where("status = ?", dto.ScheduleStatusInProgress).
		Count(&stats.InProgressSchedule).Error; err != nil {
		return nil, err
	}

	// Get completed schedules
	if err := r.db.Model(&dto.Schedule{}).
		Where("status = ?", dto.ScheduleStatusCompleted).
		Count(&stats.CompletedSchedule).Error; err != nil {
		return nil, err
	}

	// Get blocked schedules
	if err := r.db.Model(&dto.Schedule{}).
		Where("status = ?", dto.ScheduleStatusBlocked).
		Count(&stats.BlockedSchedule).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// ToResponse converts a Schedule model to ScheduleResponse DTO
func (r *ScheduleRepository) ToResponse(schedule *dto.Schedule) dto.ScheduleResponse {
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
func (r *ScheduleRepository) ToListResponse(schedules []dto.Schedule, total int64, page, limit int) dto.ScheduleListResponse {
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
		Pagination: dto.PaginationMeta{
			Page: page,
			Total: total,
			Limit: limit,
			TotalPages: totalPages,
		},
	}
}