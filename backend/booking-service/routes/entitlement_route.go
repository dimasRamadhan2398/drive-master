package routes

import (
	"booking-service/controllers"

	"github.com/gin-gonic/gin"
)

type EntitlementRoute struct {
	controller controllers.IControllerRegistry
}

type IEntitlementRoute interface {
	Run(group *gin.RouterGroup)
}

func (r *EntitlementRoute) Run(group *gin.RouterGroup) {
	entitlements := group.Group("/entitlements")
	{
		entitlements.GET("", r.controller.GetEntitlementController().ListEntitlements)
		entitlements.POST("", r.controller.GetEntitlementController().CreateEntitlement)
		entitlements.GET("/:id", r.controller.GetEntitlementController().GetEntitlement)
		entitlements.PUT("/:id", r.controller.GetEntitlementController().UpdateEntitlement)
		entitlements.DELETE("/:id", r.controller.GetEntitlementController().DeleteEntitlement)
		entitlements.GET("/user/:userId", r.controller.GetEntitlementController().GetUserEntitlements)
	}
}