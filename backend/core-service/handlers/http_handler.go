package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct{}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.Health)
}

func (h *HTTPHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "core-service"})
}
