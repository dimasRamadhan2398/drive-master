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

type IRepositoryRegistry interface {
	GetUser() IUserRepository
	GetMember() IMemberRepository
	GetInstructor() IInstructorRepository
	GetRole() IRoleRepository
	GetWorkExperience() IWorkExperienceRepository
}