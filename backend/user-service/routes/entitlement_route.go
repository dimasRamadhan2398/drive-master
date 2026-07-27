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
		group.GET("/members/:id", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().ListEntitlements)

		// Get single entitlement by member ID and entitlement ID
		group.GET("/members/:id/entitlements/:entId", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().GetEntitlement)

		// Get single entitlement by ID only (for backward compatibility)
		group.GET("/:id", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().GetEntitlementByID)

		// Create new entitlement
		group.POST("/members/:id", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().CreateEntitlement)

		// Sync entitlement from booking (internal/direct HTTP sync)
		group.POST("/sync", r.controller.GetEntitlementController().SyncEntitlement)

		// Update entitlement
		group.PUT("/members/:id/entitlements/:entId", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().UpdateEntitlement)

		// Delete entitlement
		group.DELETE("/members/:id/entitlements/:entId", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().DeleteEntitlement)

		// Use session (decrement remaining, increment used)
		group.POST("/members/:id/entitlements/:entId/use-session", r.authMiddleware.Authenticate(), r.controller.GetEntitlementController().UseSession)
	}
}
