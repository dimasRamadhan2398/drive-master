package services

import (
	"user-service/pkg/config"
	"user-service/repositories"
)

type Registry struct {
	repoRegistry *repositories.Registry
}

type IServiceRegistry interface {
	GetUserService() *UserService
	GetAuthService() *AuthService
	GetMemberService() *MemberService
	GetInstructorService() *InstructorService
	GetRoleService() *RoleService
	GetEmailService() IMailtrapEmailService
	GetMediaService() IMediaService
}

func NewServiceRegistry(repoRegistry *repositories.Registry) IServiceRegistry {
	return &Registry{repoRegistry: repoRegistry}
}

func (r *Registry) GetUserService() *UserService {
	return NewUserService(r.repoRegistry.GetUser())
}

func (r *Registry) GetAuthService() *AuthService {
	return NewAuthService(r.repoRegistry.GetUser())
}

func (r *Registry) GetMemberService() *MemberService {
	return NewMemberService(r.repoRegistry.GetMember())
}

func (r *Registry) GetInstructorService() *InstructorService {
	return NewInstructorService(r.repoRegistry.GetInstructor(), r.repoRegistry.GetWorkExperience())
}

func (r *Registry) GetRoleService() *RoleService {
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