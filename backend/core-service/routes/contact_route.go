package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type ContactRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IContactRoute interface {
	Run()
}

func NewContactRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IContactRoute {
	return &ContactRoute{
		controller: controller,
		group:      group,
	}
}

func (r *ContactRoute) Run() {
	contact := r.group.Group("/contact")
	{
		contact.POST("", r.controller.GetContactController().CreateInquiry)
		contact.GET("", r.controller.GetContactController().GetAllInquiries)
	}
}
