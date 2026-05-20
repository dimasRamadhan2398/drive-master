package handlers

import (
	"api-gateway/pkg/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
    userServiceURL        *url.URL
    coreServiceURL        *url.URL
    bookingServiceURL     *url.URL
    catalogServiceURL     *url.URL
    voucherServiceURL     *url.URL
    notificationServiceURL *url.URL
    contentServiceURL     *url.URL
}

func NewProxyHandler(cfg *config.Config) *ProxyHandler {
    parse := func(raw string) *url.URL {
        u, err := url.Parse(raw)
        if err != nil {
            panic("invalid service URL: " + raw)
        }
        return u
    }

    return &ProxyHandler{
        userServiceURL:         parse(cfg.Services.UserServiceURL),
        coreServiceURL:         parse(cfg.Services.CoreServiceURL),
        bookingServiceURL:      parse(cfg.Services.BookingServiceURL),
    }
}

func (h *ProxyHandler) ToUserService(c *gin.Context)        { h.proxy(c, h.userServiceURL, "/api/v1/users") }
func (h *ProxyHandler) ToCoreService(c *gin.Context)        { h.proxy(c, h.coreServiceURL, "/api/v1/core") }
func (h *ProxyHandler) ToBookingService(c *gin.Context)     { h.proxy(c, h.bookingServiceURL, "/api/v1/bookings") }
func (h *ProxyHandler) ToCatalogService(c *gin.Context)     { h.proxy(c, h.catalogServiceURL, "/api/v1/catalog") }
func (h *ProxyHandler) ToVoucherService(c *gin.Context)     { h.proxy(c, h.voucherServiceURL, "/api/v1/vouchers") }
func (h *ProxyHandler) ToNotificationService(c *gin.Context){ h.proxy(c, h.notificationServiceURL, "/api/v1/notifications") }
func (h *ProxyHandler) ToContentService(c *gin.Context)     { h.proxy(c, h.contentServiceURL, "/api/v1/content") }

func (h *ProxyHandler) proxy(c *gin.Context, target *url.URL, stripPrefix string) {
    proxy := httputil.NewSingleHostReverseProxy(target)

    // strip the gateway prefix before forwarding
    c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, stripPrefix)
    if c.Request.URL.Path == "" {
        c.Request.URL.Path = "/"
    }

    // forward user context injected by JWT middleware
    c.Request.Header.Set("X-User-ID",   c.GetString("userID"))
    c.Request.Header.Set("X-User-Role", c.GetString("userRole"))
    c.Request.Header.Set("X-Request-ID", c.GetString("requestID"))

    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
        c.JSON(http.StatusBadGateway, gin.H{
            "success": false,
            "message": "service unavailable",
            "error":   err.Error(),
        })
    }

    proxy.ServeHTTP(c.Writer, c.Request)
}
