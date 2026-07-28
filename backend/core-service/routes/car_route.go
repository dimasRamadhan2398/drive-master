package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type CarRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type ICarRoute interface {
	Run()
}

func NewCarRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) ICarRoute {
	return &CarRoute{
		controller: controller,
		group:      group,
	}
}

func (r *CarRoute) Run() {
	cars := r.group.Group("/cars")
	{
		cars.GET("", r.controller.GetCarController().GetAllCars)
		cars.GET("/:id", r.controller.GetCarController().GetCarByID)
		cars.POST("", r.controller.GetCarController().CreateCar)
		cars.PUT("/:id", r.controller.GetCarController().UpdateCar)
		cars.DELETE("/:id", r.controller.GetCarController().DeleteCar)
	}
}