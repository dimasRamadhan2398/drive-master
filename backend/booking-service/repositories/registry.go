package repositories

import (
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) *Registry {
	return &Registry{db: db}
}

// GetSchedule returns a new ScheduleRepository
func (r *Registry) GetSchedule() IScheduleRepository {
	return NewScheduleRepository(r.db)
}

// GetEnrollment returns a new EnrollmentRepository
func (r *Registry) GetEnrollment() IEnrollmentRepository {
	return NewEnrollmentRepository(r.db)
}

// GetSession returns a new SessionRepository
func (r *Registry) GetSession() ISessionRepository {
	return NewSessionRepository(r.db)
}

type IRepositoryRegistry interface {
	GetSchedule() IScheduleRepository
	GetEnrollment() IEnrollmentRepository
	GetSession() ISessionRepository
}