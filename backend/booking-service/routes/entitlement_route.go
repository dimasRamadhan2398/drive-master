package routes

import (
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type EntitlementRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IEntitlementRoute interface {
	Run()
}

func NewEntitlementRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IEntitlementRoute {
	return &BookingRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *EntitlementRoute) Run() {
	entitlements := r.group.Group("/entitlements")
	{
		entitlements.GET("", r.controller.GetEntitlementController().ListEntitlements)
		entitlements.POST("", r.controller.GetEntitlementController().CreateEntitlement)
		entitlements.GET("/:id", r.controller.GetEntitlementController().GetEntitlement)
		entitlements.PUT("/:id", r.controller.GetEntitlementController().UpdateEntitlement)
		entitlements.DELETE("/:id", r.controller.GetEntitlementController().DeleteEntitlement)
		entitlements.GET("/user/:userId", r.controller.GetEntitlementController().GetUserEntitlements)
	}
}