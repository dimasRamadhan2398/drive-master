package repositories

import (
	"context"
	"errors"
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRecurringScheduleRepository interface {
	Create(ctx context.Context, schedule *models.InstructorRecurringSchedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.InstructorRecurringSchedule, error)
	GetByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorRecurringSchedule, error)
	GetActiveByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorRecurringSchedule, error)
	GetByDayOfWeek(ctx context.Context, dayOfWeek int) ([]models.InstructorRecurringSchedule, error)
	Update(ctx context.Context, schedule *models.InstructorRecurringSchedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByInstructorID(ctx context.Context, instructorID uuid.UUID) error
}

type RecurringScheduleRepository struct {
	*base.BaseRepository
}

func (r *RecurringScheduleRepository) Create(ctx context.Context, schedule *models.InstructorRecurringSchedule) error {
	return r.BaseRepository.Create(schedule)
}

func (r *RecurringScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.InstructorRecurringSchedule, error) {
	var schedule models.InstructorRecurringSchedule
	if err := r.BaseRepository.FindOne(&schedule, "id = ?", id); err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *RecurringScheduleRepository) GetByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorRecurringSchedule, error) {
	var schedules []models.InstructorRecurringSchedule
	if err := r.BaseRepository.DB.Where("instructor_id = ?", instructorID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *RecurringScheduleRepository) GetActiveByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorRecurringSchedule, error) {
	var schedules []models.InstructorRecurringSchedule
	if err := r.BaseRepository.DB.Where("instructor_id = ? AND is_active = ?", instructorID, true).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *RecurringScheduleRepository) GetByDayOfWeek(ctx context.Context, dayOfWeek int) ([]models.InstructorRecurringSchedule, error) {
	var schedules []models.InstructorRecurringSchedule
	if err := r.BaseRepository.DB.Where("day_of_week = ? AND is_active = ?", dayOfWeek, true).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *RecurringScheduleRepository) Update(ctx context.Context, schedule *models.InstructorRecurringSchedule) error {
	return r.BaseRepository.Update(schedule)
}

func (r *RecurringScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	schedule, err := r.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("recurring schedule not found")
		}
		return err
	}
	return r.BaseRepository.Delete(schedule)
}

func (r *RecurringScheduleRepository) DeleteByInstructorID(ctx context.Context, instructorID uuid.UUID) error {
	return r.BaseRepository.DB.Where("instructor_id = ?", instructorID).Delete(&models.InstructorRecurringSchedule{}).Error
}

func NewRecurringScheduleRepository(db *gorm.DB) IRecurringScheduleRepository {
	return &RecurringScheduleRepository{BaseRepository: base.NewBaseRepository(db)}
}
