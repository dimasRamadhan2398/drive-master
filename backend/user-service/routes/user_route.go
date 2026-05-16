package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IUserRoute interface {
	Run()
}

func NewUserRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IUserRoute {
	return &UserRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/users")
	group.GET("/", u.authMiddleware.Authenticate(), u.controller.GetUserController().GetAllUsers)
	group.GET("/:id", u.authMiddleware.Authenticate(), u.controller.GetUserController().GetUserByID)
	group.PUT("/:id", u.authMiddleware.Authenticate(), u.controller.GetUserController().UpdateUser)
	group.DELETE("/:id", u.authMiddleware.Authenticate(), u.controller.GetUserController().DeleteUser)
	group.POST("/", u.authMiddleware.Authenticate(), u.controller.GetUserController().CreateUser)
}
