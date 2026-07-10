package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type PageRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IPageRoute interface {
	Run()
}

func NewPageRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IPageRoute {
	return &PageRoute{
		controller: controller,
		group:      group,
	}
}

func (r *PageRoute) Run() {
	pages := r.group.Group("/pages")
	{
		pages.GET("", r.controller.GetPageController().GetAllPages)
		pages.GET("/:id", r.controller.GetPageController().GetPageByID)
		pages.GET("/slug/*slug", r.controller.GetPageController().GetPageBySlug) // Use *slug wildcard to catch path slugs like / or /about
		pages.POST("", r.controller.GetPageController().CreatePage)
		pages.PUT("/:id", r.controller.GetPageController().UpdatePage)
		pages.DELETE("/:id", r.controller.GetPageController().DeletePage)
	}
}
