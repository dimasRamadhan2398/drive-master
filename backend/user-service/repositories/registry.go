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

// GetUser implements IRepositoryRegistry
func (r *Registry) GetUser() IUserRepository {
	return NewUserRepository(r.db)
}

// GetMember implements IRepositoryRegistry
func (r *Registry) GetMember() IMemberRepository {
	return NewMemberRepository(r.db)
}

// GetInstructor implements IRepositoryRegistry
func (r *Registry) GetInstructor() IInstructorRepository {
	return NewInstructorRepository(r.db)
}

// GetRole implements IRepositoryRegistry
func (r *Registry) GetRole() IRoleRepository {
	return NewRoleRepository(r.db)
}

// GetWorkExperience implements IRepositoryRegistry
func (r *Registry) GetWorkExperience() IWorkExperienceRepository {
	return NewWorkExperienceRepository(r.db)
}

// GetCoverageArea implements IRepositoryRegistry
func (r *Registry) GetCoverageArea() ICoverageAreaRepository {
	return NewCoverageArea(r.db)
}

// GetCertification implements IRepositoryRegistry
func (r *Registry) GetCertification() ICertificationRepository {
	return NewCertificationRepository(r.db)
}

// GetEntitlement implements IRepositoryRegistry
func (r *Registry) GetEntitlement() IEntitlementRepository {
	return NewEntitlementRepository(r.db)
}

// GetTestimonial implements IRepositoryRegistry
func (r *Registry) GetTestimonial() ITestimonialRepository {
	return NewTestimonialRepository(r.db)
}

// GetRecurringSchedule implements IRepositoryRegistry
func (r *Registry) GetRecurringSchedule() IRecurringScheduleRepository {
	return NewRecurringScheduleRepository(r.db)
}

type IRepositoryRegistry interface {
	GetUser() IUserRepository
	GetMember() IMemberRepository
	GetInstructor() IInstructorRepository
	GetRole() IRoleRepository
	GetWorkExperience() IWorkExperienceRepository
	GetCoverageArea() ICoverageAreaRepository
	GetCertification() ICertificationRepository
	GetEntitlement() IEntitlementRepository
	GetTestimonial() ITestimonialRepository
	GetRecurringSchedule() IRecurringScheduleRepository
}