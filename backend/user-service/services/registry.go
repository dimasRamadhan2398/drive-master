package services

import (
	"user-service/pkg/config"
	"user-service/repositories"
)

type Registry struct {
	repoRegistry *repositories.Registry
}

type IServiceRegistry interface {
	GetUserService() IUserService
	GetAuthService() IAuthService
	GetMemberService() IMemberService
	GetInstructorService() IInstructorService
	GetRoleService() IRoleService
	GetEmailService() IMailtrapEmailService
	GetMediaService() IMediaService
}

func NewServiceRegistry(repoRegistry *repositories.Registry) IServiceRegistry {
	return &Registry{repoRegistry: repoRegistry}
}

func (r *Registry) GetUserService() IUserService {
	return NewUserService(r.repoRegistry.GetUser())
}

func (r *Registry) GetAuthService() IAuthService {
	return NewAuthService(r.repoRegistry.GetUser())
}

func (r *Registry) GetMemberService() IMemberService {
	return NewMemberService(r.repoRegistry.GetMember())
}

func (r *Registry) GetInstructorService() IInstructorService {
	return NewInstructorService(r.repoRegistry.GetInstructor(), r.repoRegistry.GetWorkExperience())
}

func (r *Registry) GetRoleService() IRoleService {
	return NewRoleService(r.repoRegistry.GetRole())
}

func (r *Registry) GetEmailService() IMailtrapEmailService {
	cfg := config.Get()
	return NewMailtrapEmailService(cfg.Email.Token, cfg.Email.FromEmail, cfg.Email.AppName)
}

func (r *Registry) GetMediaService() IMediaService {
	cfg := config.Get()
	return NewMediaService(cfg.ImageKit.PrivateKey)
}