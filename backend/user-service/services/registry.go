package services

import (
	coreServices "core-service/services"
	"fmt"
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
	GetRecurringScheduleService() IRecurringScheduleService
	GetDashboardService() IDashboardService
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
	return NewUserService(r.repoRegistry.GetUser(), r.repoRegistry.GetRole(), r.GetInstructorService(), r.eventPublisher)
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
		r.GetMediaService(),
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

	// Debug logging for email configuration
	fmt.Printf("[EMAIL DEBUG] FromEmail: %s\n", cfg.Email.FromEmail)
	fmt.Printf("[EMAIL DEBUG] FromName: %s\n", cfg.Email.FromName)
	fmt.Printf("[EMAIL DEBUG] APIKey: %s\n", cfg.Email.APIKey)
	fmt.Printf("[EMAIL DEBUG] Host: %s\n", cfg.Email.Host)
	fmt.Printf("[EMAIL DEBUG] Port: %d\n", cfg.Email.Port)
	fmt.Printf("[EMAIL DEBUG] User: %s\n", cfg.Email.User)
	fmt.Printf("[EMAIL DEBUG] SMTPEnabled: %s\n", cfg.Email.Password)

	return NewMailtrapEmailService(cfg.Email.FromEmail, cfg.Email.FromName, cfg.Email.APIKey)
}

func (r *Registry) GetMediaService() IMediaService {
	cfg := config.Get()
	return coreServices.NewMediaService(cfg.ImageKit.PrivateKey, cfg.ImageKit.URLEndpoint)
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
	return NewCertificationService(r.repoRegistry.GetCertification(), r.repoRegistry.GetUser(), r.GetEmailService())
}

func (r *Registry) GetEntitlementService() IEntitlementService {
	certService := r.GetCertificationService()
	// Initialize the completion listener.
	// Passing certRepo allows the listener to check for an existing certificate
	// via the entitlement_id FK before attempting issuance (idempotent behaviour).
	listener := listeners.NewEntitlementCompletedListener(certService, r.repoRegistry.GetCertification(), r.eventPublisher)
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

func (r *Registry) GetRecurringScheduleService() IRecurringScheduleService {
	return NewRecurringScheduleService(r.repoRegistry.GetRecurringSchedule())
}

func (r *Registry) GetDashboardService() IDashboardService {
	return NewDashboardService(r.repoRegistry.GetUser(), r.repoRegistry.GetRole())
}
