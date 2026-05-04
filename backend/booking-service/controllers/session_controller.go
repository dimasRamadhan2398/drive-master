package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/services"

	"github.com/gin-gonic/gin"
)

type SessionController struct {
	sessionService services.ISessionService
}

func NewSessionController(sessionService services.ISessionService) ISessionController {
	return &SessionController{sessionService: sessionService}
}

type ISessionController interface {
	CreateSession(c *gin.Context)
	GetSession(c *gin.Context)
	ListSessions(c *gin.Context)
}

func (c *SessionController) CreateSession(ctx *gin.Context) {
	var req dto.CreateSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.sessionService.CreateSession(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *SessionController) GetSession(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	resp, err := c.sessionService.GetSession(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *SessionController) ListSessions(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.sessionService.ListSessions(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}