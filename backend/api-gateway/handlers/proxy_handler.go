package handlers

import (
	"api-gateway/pkg/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	userServiceURL         *url.URL
	coreServiceURL         *url.URL
	bookingServiceURL      *url.URL
	paymentServiceURL      *url.URL
	catalogServiceURL      *url.URL
	voucherServiceURL      *url.URL
	notificationServiceURL *url.URL
	contentServiceURL      *url.URL
}

func NewProxyHandler(cfg *config.Config) *ProxyHandler {
	parse := func(raw string) *url.URL {
		u, err := url.Parse(raw)
		if err != nil {
			panic("invalid service URL: " + raw)
		}
		return u
	}

	paymentURL := parse(cfg.Services.PaymentServiceURL)

	return &ProxyHandler{
		userServiceURL:    parse(cfg.Services.UserServiceURL),
		coreServiceURL:    parse(cfg.Services.CoreServiceURL),
		bookingServiceURL: parse(cfg.Services.BookingServiceURL),
		paymentServiceURL: paymentURL,
	}
}

func (h *ProxyHandler) ToUserService(c *gin.Context) { h.proxy(c, h.userServiceURL, "/api/v1/users") }
func (h *ProxyHandler) ToUserServiceDirect(c *gin.Context) { h.proxy(c, h.userServiceURL, "") }
func (h *ProxyHandler) ToCoreService(c *gin.Context) { h.proxy(c, h.coreServiceURL, "/api/v1/core") }

// ToCoreServiceDirect proxies to core-service WITHOUT stripping any prefix.
// Use this for routes registered directly under /api/v1 (e.g. /api/v1/general-settings).
func (h *ProxyHandler) ToCoreServiceDirect(c *gin.Context) { h.proxy(c, h.coreServiceURL, "") }
func (h *ProxyHandler) ToBookingService(c *gin.Context) {
	h.proxy(c, h.bookingServiceURL, "/api/v1/bookings")
}
func (h *ProxyHandler) ToBookingServiceDirect(c *gin.Context) {
	h.proxy(c, h.bookingServiceURL, "")
}
func (h *ProxyHandler) ToPaymentService(c *gin.Context) {
	h.proxy(c, h.paymentServiceURL, "")
}
func (h *ProxyHandler) ToCatalogService(c *gin.Context) {
	h.proxy(c, h.catalogServiceURL, "/api/v1/catalog")
}
func (h *ProxyHandler) ToVoucherService(c *gin.Context) {
	h.proxy(c, h.voucherServiceURL, "/api/v1/vouchers")
}
func (h *ProxyHandler) ToNotificationService(c *gin.Context) {
	h.proxy(c, h.notificationServiceURL, "/api/v1/notifications")
}
func (h *ProxyHandler) ToContentService(c *gin.Context) {
	h.proxy(c, h.contentServiceURL, "/api/v1/content")
}

func (h *ProxyHandler) proxy(c *gin.Context, target *url.URL, stripPrefix string) {
	proxy := httputil.NewSingleHostReverseProxy(target)

	// strip the gateway prefix before forwarding
	c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, stripPrefix)
	if len(c.Request.URL.Path) > 1 && strings.HasSuffix(c.Request.URL.Path, "/") {
		c.Request.URL.Path = strings.TrimSuffix(c.Request.URL.Path, "/")
	}
	if c.Request.URL.Path == "" {
		c.Request.URL.Path = "/"
	}

	// forward user context injected by JWT middleware
	c.Request.Header.Set("X-User-ID", c.GetString("userID"))
	c.Request.Header.Set("X-User-Role", c.GetString("userRole"))
	c.Request.Header.Set("X-Request-ID", c.GetString("requestID"))

	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("[PROXY] ModifyResponse called for path: %s", c.Request.URL.Path)
		log.Printf("[PROXY] Headers before delete: %v", resp.Header)
		// Strip CORS headers from the downstream service to avoid duplication
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Credentials")
		resp.Header.Del("Access-Control-Allow-Methods")
		resp.Header.Del("Access-Control-Allow-Headers")
		resp.Header.Del("Access-Control-Expose-Headers")
		log.Printf("[PROXY] Headers after delete: %v", resp.Header)
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"message": "service unavailable",
			"error":   err.Error(),
		})
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
