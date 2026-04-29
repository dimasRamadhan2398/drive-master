package routes

import (
	"user-service/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the user-service
// Returns the router group for documentation purposes
func SetupRoutes(router *gin.Engine, userController *controllers.UserController) {
	userController.RegisterRoutes(router)
}


// UserRoutes defines the user-related route groups
type UserRoutes struct {
	Users       string
	Auth        string
	Members     string
	Instructors string
	Media       string
	Roles       string
}

// GetRouteGroups returns all registered route groups for documentation
func GetRouteGroups() UserRoutes {
	return UserRoutes{
		Users:       "/users",
		Auth:        "/auth",
		Members:     "/members",
		Instructors: "/instructors",
		Media:       "/media",
		Roles:       "/roles",
	}
}

// RouteInfo represents information about a registered route
type RouteInfo struct {
	Method  string
	Path    string
	Handler string
}

// GetAllRoutes returns all registered routes with their handlers
func GetAllRoutes(router *gin.Engine) []RouteInfo {
	var routes []RouteInfo
	for _, route := range router.Routes() {
		routes = append(routes, RouteInfo{
			Method:  route.Method,
			Path:    route.Path,
			Handler: route.Handler,
		})
	}
	return routes
}
