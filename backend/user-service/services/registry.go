package services

import (
	"user-service/clients"
	"user-service/clients/region"
	"user-service/pkg/config"
	pkgKafka "user-service/pkg/kafka"
	"user-service/pkg/redis"
	"user-service/repositories"
	"user-service/services/listeners"
)

type Registry struct {
	repoRegistry   *repositories.Registry
	eventPublisher pkgKafka.IEventPublisher
	redisClient    *redis.Client
	regionService  IRegionService
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
	GetRegionService() IRegionService
	GetCertificationService() ICertificationService
	GetEntitlementService() IEntitlementService
	GetTestimonialService() ITestimonialService
}

func NewServiceRegistry(repoRegistry *repositories.Registry, eventPublisher pkgKafka.IEventPublisher, redisClient *redis.Client) IServiceRegistry {
	// Initialize region client and service using client config
	cfg := config.Get()
	clientConfig := clients.NewClientConfig(
		clients.WithBaseURL(cfg.CoreService.BaseURL),
	)
	regionClient := region.NewClient(clientConfig)
	regionService := NewRegionService(regionClient)

	return &Registry{
		repoRegistry:   repoRegistry,
		eventPublisher: eventPublisher,
		redisClient:    redisClient,
		regionService:  regionService,
	}
}

func (r *Registry) GetUserService() IUserService {
	return NewUserService(r.repoRegistry.GetUser(), r.repoRegistry.GetRole(), r.GetInstructorService())
}

func (r *Registry) GetMemberService() IMemberService {
	return NewMemberService(r.repoRegistry.GetMember(), r.repoRegistry.GetRole(), r.GetEntitlementService())
}

func (r *Registry) GetInstructorService() IInstructorService {
	return NewInstructorService(
		r.repoRegistry.GetInstructor(),
		r.repoRegistry.GetUser(),
		r.repoRegistry.GetRole(),
		r.GetEmailService(),
		r.redisClient,
	)
}

func (r *Registry) GetRoleService() IRoleService {
	return NewRoleService(r.repoRegistry.GetRole())
}

func (r *Registry) GetAuthService() IAuthService {
	return NewAuthService(r.repoRegistry.GetUser(), r.redisClient, r.GetEmailService(), r.GetMemberService(), r.GetInstructorService(), r.GetRoleService())
}

func (r *Registry) GetEmailService() IMailtrapEmailService {
	cfg := config.Get()
	return NewMailtrapEmailService(cfg.Email.Host, cfg.Email.Port, cfg.Email.User, cfg.Email.Password, cfg.Email.FromEmail, cfg.Email.FromName)
}

func (r *Registry) GetMediaService() IMediaService {
	cfg := config.Get()
	return NewMediaService(cfg.ImageKit.PrivateKey)
}

func (r *Registry) GetWorkExperienceService() IWorkExperienceService {
	return NewWorkExperienceService(r.repoRegistry.GetWorkExperience(), r.repoRegistry.GetInstructor())
}

func (r *Registry) GetCoverageAreaService() ICoverageAreaService {
	return NewCoverageAreaService(r.repoRegistry.GetCoverageArea(), r.regionService)
}

func (r *Registry) GetRegionService() IRegionService {
	return r.regionService
}

func (r *Registry) GetCertificationService() ICertificationService {
	return NewCertificationService(r.repoRegistry.GetCertification())
}

func (r *Registry) GetEntitlementService() IEntitlementService {
	certService := r.GetCertificationService()
	// Initialize the completion listener
	listener := listeners.NewEntitlementCompletedListener(certService, r.eventPublisher)
	return NewEntitlementService(
		r.repoRegistry.GetEntitlement(),
		r.repoRegistry.GetMember(),
		certService,
		r.eventPublisher,
		listener,
	)
}

func (r *Registry) GetTestimonialService() ITestimonialService {
	return NewTestimonialService(r.repoRegistry.GetTestimonial(), r.eventPublisher)
}