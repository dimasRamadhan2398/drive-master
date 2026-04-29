package routes

import (
	"user-service/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the user-service
func SetupRoutes(router *gin.Engine, userController *controllers.UserController) {
	userController.RegisterRoutes(router)
}