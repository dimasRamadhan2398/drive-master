package handlers

import (
	"net/http"
	"strconv"
	"user-service/models"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	userService *services.UserService
}

func NewHTTPHandler(userService *services.UserService) *HTTPHandler {
	return &HTTPHandler{userService: userService}
}

func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.Health)
	router.POST("/users", h.CreateUser)
	router.GET("/users/:id", h.GetUserByID)
}

func (h *HTTPHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "user-service"})
}

func (h *HTTPHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *HTTPHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
