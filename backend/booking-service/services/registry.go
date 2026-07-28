package services

// IServiceRegistry defines methods for getting services
type IServiceRegistry interface {
	GetSessionService() ISessionService
	GetEnrollmentService() IEnrollmentService
	GetScheduleService() IScheduleService
	GetPaymentService() IPaymentService
	GetRevenueService() IRevenueService
}