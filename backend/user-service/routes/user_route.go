package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IUserRoute interface {
	Run()
}

func NewUserRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IUserRoute {
	return &UserRoute{controller: controller, group: group}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/users")
	group.GET("/", middlewares.Authenticate(), u.controller.GetUserController().GetAllUsers)
	group.GET("/:id", middlewares.Authenticate(), u.controller.GetUserController().GetUserByID)
	group.PUT("/:id", middlewares.Authenticate(), u.controller.GetUserController().UpdateUser)
	group.DELETE("/:id", middlewares.Authenticate(), u.controller.GetUserController().DeleteUser)
	group.POST("/", middlewares.Authenticate(), u.controller.GetUserController().CreateUser)
}
