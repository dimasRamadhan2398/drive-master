package services

import (
	"api-gateway/handlers"

	"github.com/gin-gonic/gin"
)

type RouterService struct {
	userServiceURL string
	coreServiceURL string
}

func NewRouterService(userServiceURL, coreServiceURL string) *RouterService {
	return &RouterService{
		userServiceURL: userServiceURL,
		coreServiceURL: coreServiceURL,
	}
}

func (s *RouterService) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "api-gateway"})
	})

	userProxy := handlers.NewReverseProxy(s.userServiceURL, "/api/users")
	coreProxy := handlers.NewReverseProxy(s.coreServiceURL, "/api/core")

	router.Any("/api/users", userProxy)
	router.Any("/api/users/*path", userProxy)
	router.Any("/api/core", coreProxy)
	router.Any("/api/core/*path", coreProxy)
}
