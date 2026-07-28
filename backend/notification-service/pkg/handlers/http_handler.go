package handlers

import (
	"net/http"
	"strconv"
	"notification-service/models/dto"
	"notification-service/services"

	"github.com/gin-gonic/gin"
)

// HTTPHandler handles HTTP requests for notification service
type HTTPHandler struct {
	notificationService *services.NotificationService
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(notificationService *services.NotificationService) *HTTPHandler {
	return &HTTPHandler{notificationService: notificationService}
}

// RegisterRoutes registers all HTTP routes
func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", h.Health)

	// Newsletter subscription routes
	router.POST("/newsletter/subscribe", h.SubscribeNewsletter)
	router.POST("/newsletter/unsubscribe", h.UnsubscribeNewsletter)
	router.GET("/newsletter/subscription/:email", h.GetSubscription)
	router.GET("/newsletter/subscribers", h.GetActiveSubscribers)

	// Notification routes
	router.POST("/notifications/send", h.SendNotification)
	router.POST("/notifications/schedule", h.ScheduleNotification)
	router.GET("/notifications/:id", h.GetNotification)
	router.GET("/notifications/user/:userId", h.GetUserNotifications)

	// Preference routes
	router.GET("/preferences/:userId", h.GetPreferences)
	router.PUT("/preferences", h.UpdatePreferences)

	// Admin routes
	router.POST("/admin/promotional/email", h.SendPromotionalEmail)
}

// Health returns service health status
func (h *HTTPHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "notification-service"})
}

// SubscribeNewsletter handles newsletter subscription
func (h *HTTPHandler) SubscribeNewsletter(c *gin.Context) {
	var req dto.SubscribeNewsletterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.notificationService.SubscribeNewsletter(c.Request.Context(), req.Email, req.Name, req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// UnsubscribeNewsletter handles newsletter unsubscription
func (h *HTTPHandler) UnsubscribeNewsletter(c *gin.Context) {
	var req dto.UnsubscribeNewsletterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationService.UnsubscribeNewsletter(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unsubscribed"})
}

// GetSubscription retrieves subscription by email
func (h *HTTPHandler) GetSubscription(c *gin.Context) {
	email := c.Param("email")

	subscription, err := h.notificationService.GetSubscriptionByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// GetActiveSubscribers retrieves all active subscribers
func (h *HTTPHandler) GetActiveSubscribers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	subscribers, total, err := h.notificationService.GetActiveSubscriptions(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscribers": subscribers,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// SendNotification sends a notification immediately
func (h *HTTPHandler) SendNotification(c *gin.Context) {
	var req dto.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process based on notification type
	var err error
	switch req.Type {
	case "email":
		err = h.notificationService.SendPromotionalEmail(c.Request.Context(), req.Recipients, req.Subject, req.Content)
	case "whatsapp":
		if len(req.Recipients) > 0 {
			err = h.notificationService.SendPromotionalWhatsApp(c.Request.Context(), req.Recipients[0], req.Subject, req.Content)
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}

// ScheduleNotification schedules a notification for future delivery
func (h *HTTPHandler) ScheduleNotification(c *gin.Context) {
	var req dto.ScheduleNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := req.ToNotification()
	if err := h.notificationService.ScheduleNotification(c.Request.Context(), notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notification scheduled"})
}

// GetNotification retrieves a notification by ID
func (h *HTTPHandler) GetNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	notification, err := h.notificationService.GetNotificationByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// GetUserNotifications retrieves notifications for a user
func (h *HTTPHandler) GetUserNotifications(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, total, err := h.notificationService.GetNotificationsByUserID(c.Request.Context(), uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// GetPreferences retrieves notification preferences for a user
func (h *HTTPHandler) GetPreferences(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	preferences, err := h.notificationService.GetNotificationPreferences(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

// UpdatePreferences updates notification preferences
func (h *HTTPHandler) UpdatePreferences(c *gin.Context) {
	var req dto.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	preference := req.ToPreference()
	if err := h.notificationService.UpdateNotificationPreferences(c.Request.Context(), preference); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, preference)
}

// SendPromotionalEmail sends a promotional email campaign
func (h *HTTPHandler) SendPromotionalEmail(c *gin.Context) {
	var req dto.PromotionalEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationService.SendPromotionalEmail(c.Request.Context(), req.Recipients, req.Subject, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotional email sent successfully"})
}