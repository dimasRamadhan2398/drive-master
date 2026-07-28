package routes

import (
	"user-service/controllers"

	"github.com/gin-gonic/gin"
)

type TestimonialRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type ITestimonialRoute interface {
	Run()
}

func NewTestimonialRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) ITestimonialRoute {
	return &TestimonialRoute{
		controller: controller,
		group:      group,
	}
}

func (r *TestimonialRoute) Run() {
	testimonials := r.group.Group("/testimonials")
	{
		// Admin endpoints
		testimonials.GET("", r.controller.GetTestimonialController().GetAllTestimonials)
		testimonials.GET("/all", r.controller.GetTestimonialController().GetAllTestimonials)
		testimonials.GET("/:id", r.controller.GetTestimonialController().GetTestimonialByID)
		testimonials.POST("", r.controller.GetTestimonialController().CreateTestimonial)
		testimonials.PUT("/:id", r.controller.GetTestimonialController().UpdateTestimonial)
		testimonials.DELETE("/:id", r.controller.GetTestimonialController().DeleteTestimonial)
		testimonials.PUT("/:id/featured", r.controller.GetTestimonialController().ToggleFeatured)
		testimonials.PATCH("/:id/status", r.controller.GetTestimonialController().UpdateStatus)

		// Public endpoints
		testimonials.GET("/published", r.controller.GetTestimonialController().GetPublishedTestimonials)
		testimonials.GET("/featured", r.controller.GetTestimonialController().GetFeaturedTestimonials)
		testimonials.GET("/user/:userId", r.controller.GetTestimonialController().GetTestimonialsByUserID)
	}
}
