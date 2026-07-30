package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type PackageRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IPackageRoute interface {
	Run()
}

func NewPackageRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IPackageRoute {
	return &PackageRoute{
		controller: controller,
		group:      group,
	}
}

func (r *PackageRoute) Run() {
	packages := r.group.Group("/packages")
	{
		packages.GET("/all", r.controller.GetPackageController().GetAllPackages)
		packages.GET("/:id", r.controller.GetPackageController().GetPackageByID)
		packages.POST("/create", r.controller.GetPackageController().CreatePackage)
		packages.PUT("/:id", r.controller.GetPackageController().UpdatePackage)
		packages.DELETE("/:id", r.controller.GetPackageController().DeletePackage)
		packages.PUT("/toggle-status/:id", r.controller.GetPackageController().ToggleStatusPackage)
		packages.POST("/:id/increment-count", r.controller.GetPackageController().IncrementCount)
	}
}