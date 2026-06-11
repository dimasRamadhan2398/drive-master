package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"api-gateway/pkg/config"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	userServiceURL     *url.URL
	bookingServiceURL  *url.URL
	coreServiceURL     *url.URL
	httpClient         *http.Client
}

func NewDashboardHandler(cfg *config.Config) *DashboardHandler {
	parse := func(raw string) *url.URL {
		u, err := url.Parse(raw)
		if err != nil {
			panic("invalid service URL: " + raw)
		}
		return u
	}

	return &DashboardHandler{
		userServiceURL:    parse(cfg.Services.UserServiceURL),
		bookingServiceURL: parse(cfg.Services.BookingServiceURL),
		coreServiceURL:    parse(cfg.Services.CoreServiceURL),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type DashboardStats struct {
	TotalUsers           int64 `json:"totalUsers"`
	TotalMembers         int64 `json:"totalMembers"`
	TotalInstructors     int64 `json:"totalInstructors"`
	RecentRegistrations  int64 `json:"recentRegistrations"`
	ActiveSessions       int64 `json:"activeSessions"`
	TotalSessions        int64 `json:"totalSessions"`
	RevenueMTD           int64 `json:"revenueMTD"`
	RevenueCurrency      string `json:"revenueCurrency"`
	CertificatesIssued    int64 `json:"certificatesIssued"`
	TotalCertifications   int64 `json:"totalCertifications"`
}

type ServiceResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// GetDashboardStats handles GET /api/v1/admin/dashboard/stats
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Create channels for concurrent calls
	type statsResult struct {
		data interface{}
		err  error
	}

	userChan := make(chan statsResult, 1)
	sessionChan := make(chan statsResult, 1)
	certChan := make(chan statsResult, 1)
	salesChan := make(chan statsResult, 1)

	// Fetch from user-service
	go func() {
		data, err := h.callService(c, ctx, h.userServiceURL, "/api/v1/dashboard/stats", "GET", nil)
		userChan <- statsResult{data: data, err: err}
	}()

	// Fetch from booking-service (sessions)
	go func() {
		data, err := h.callService(c, ctx, h.bookingServiceURL, "/api/v1/dashboard/sessions/stats", "GET", nil)
		sessionChan <- statsResult{data: data, err: err}
	}()

	// Fetch from booking-service (certifications)
	go func() {
		data, err := h.callService(c, ctx, h.bookingServiceURL, "/api/v1/dashboard/certifications/stats", "GET", nil)
		certChan <- statsResult{data: data, err: err}
	}()

	// Fetch from core-service (sales)
	go func() {
		data, err := h.callService(c, ctx, h.coreServiceURL, "/api/v1/admin/sales/analytics/overview", "GET", nil)
		salesChan <- statsResult{data: data, err: err}
	}()

	// Wait for all results
	userResult := <-userChan
	sessionResult := <-sessionChan
	certResult := <-certChan
	salesResult := <-salesChan

	// Initialize stats with defaults
	stats := DashboardStats{
		RevenueCurrency: "IDR",
	}

	// Parse user service stats
	if userResult.err == nil && userResult.data != nil {
		if resp, ok := userResult.data.(map[string]interface{}); ok {
			if data, ok := resp["data"].(map[string]interface{}); ok {
				if v, ok := data["totalUsers"].(float64); ok {
					stats.TotalUsers = int64(v)
				}
				if v, ok := data["totalMembers"].(float64); ok {
					stats.TotalMembers = int64(v)
				}
				if v, ok := data["totalInstructors"].(float64); ok {
					stats.TotalInstructors = int64(v)
				}
				if v, ok := data["recentRegistrations"].(float64); ok {
					stats.RecentRegistrations = int64(v)
				}
			}
		}
	}

	// Parse session stats
	if sessionResult.err == nil && sessionResult.data != nil {
		if resp, ok := sessionResult.data.(map[string]interface{}); ok {
			if data, ok := resp["data"].(map[string]interface{}); ok {
				if v, ok := data["activeSessions"].(float64); ok {
					stats.ActiveSessions = int64(v)
				}
				if v, ok := data["totalSessions"].(float64); ok {
					stats.TotalSessions = int64(v)
				}
			}
		}
	}

	// Parse certification stats
	if certResult.err == nil && certResult.data != nil {
		if resp, ok := certResult.data.(map[string]interface{}); ok {
			if data, ok := resp["data"].(map[string]interface{}); ok {
				if v, ok := data["issuedCertifications"].(float64); ok {
					stats.CertificatesIssued = int64(v)
				}
				if v, ok := data["totalCertifications"].(float64); ok {
					stats.TotalCertifications = int64(v)
				}
			}
		}
	}

	// Parse sales stats
	if salesResult.err == nil && salesResult.data != nil {
		if resp, ok := salesResult.data.(map[string]interface{}); ok {
			if data, ok := resp["data"].(map[string]interface{}); ok {
				if v, ok := data["totalRevenue"].(float64); ok {
					stats.RevenueMTD = int64(v)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard stats retrieved successfully",
		"data":    stats,
	})
}

func (h *DashboardHandler) callService(c *gin.Context, ctx context.Context, baseURL *url.URL, path, method string, body interface{}) (interface{}, error) {
	reqURL := baseURL.String() + path

	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = nil
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Forward authorization token if present
	if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("service returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}