package services

import (
	"user-service/pkg/config"
	pkgKafka "user-service/pkg/kafka"
	"user-service/pkg/redis"
	"user-service/repositories"
)

type Registry struct {
	repoRegistry   *repositories.Registry
	eventPublisher pkgKafka.IEventPublisher
	redisClient    *redis.Client
}

type IServiceRegistry interface {
	GetUserService() IUserService
	GetAuthService() IAuthService
	GetMemberService() IMemberService
	GetInstructorService() IInstructorService
	GetRoleService() IRoleService
	GetEmailService() IMailtrapEmailService
	GetMediaService() IMediaService
	GetWorkExperienceService() IWorkExperienceService
	GetCoverageAreaService() ICoverageAreaService
}

func NewServiceRegistry(repoRegistry *repositories.Registry, eventPublisher pkgKafka.IEventPublisher, redisClient *redis.Client) IServiceRegistry {
	return &Registry{
		repoRegistry:   repoRegistry,
		eventPublisher: eventPublisher,
		redisClient:    redisClient,
	}
}

func (r *Registry) GetUserService() IUserService {
	return NewUserService(r.repoRegistry.GetUser())
}

func (r *Registry) GetAuthService() IAuthService {
	return NewAuthService(r.repoRegistry.GetUser(), r.redisClient, r.GetEmailService())
}

func (r *Registry) GetMemberService() IMemberService {
	return NewMemberService(r.repoRegistry.GetMember())
}

func (r *Registry) GetInstructorService() IInstructorService {
	return NewInstructorService(r.repoRegistry.GetInstructor())
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

func (r *Registry) GetWorkExperienceService() IWorkExperienceService {
	return NewWorkExperienceService(r.repoRegistry.GetWorkExperience(), r.repoRegistry.GetInstructor())
}

func (r *Registry) GetCoverageAreaService() ICoverageAreaService {
	return NewCoverageAreaService(r.repoRegistry.GetCoverageArea())
}