package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type EntitlementRoute struct {
	controller     controllers.IControllerRegistry
	group          *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IEntitlementRoute interface {
	Run()
}

func NewEntitlementRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IEntitlementRoute {
	return &EntitlementRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *EntitlementRoute) Run() {
	group := r.group.Group("/entitlements")
	{
		// List entitlements for a member (with pagination)
		group.GET("/members/:memberId", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().ListEntitlements)

		// Get single entitlement by ID
		group.GET("/:id", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().GetEntitlement)

		// Create new entitlement
		group.POST("/", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().CreateEntitlement)

		// Update entitlement
		group.PUT("/:id", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().UpdateEntitlement)

		// Delete entitlement
		group.DELETE("/:id", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().DeleteEntitlement)

		// Use session (decrement remaining, increment used)
		group.POST("/:id/use-session", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().UseSession)
	}
}
