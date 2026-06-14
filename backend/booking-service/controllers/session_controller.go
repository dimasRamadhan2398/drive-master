package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/pkg/base"
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

// CreateSession godoc
// @Summary Create a new session
// @Description Creates a new session with the provided details
// @Tags sessions
// @Accept json
// @Produce json
// @Param session body dto.CreateSessionRequest true "Session data"
// @Success 201 {object} dto.SessionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sessions [post]
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

// GetSession godoc
// @Summary Get a session by ID
// @Description Retrieves a session by its ID
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} dto.SessionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /sessions/{id} [get]
func (c *SessionController) GetSession(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
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

// ListSessions godoc
// @Summary List all sessions
// @Description Retrieves a paginated list of sessions
// @Tags sessions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.SessionListResponse
// @Failure 500 {object} map[string]string
// @Router /sessions [get]
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