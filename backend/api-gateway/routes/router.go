package routes

import (
	"api-gateway/handlers"
	"api-gateway/pkg/config"
	"api-gateway/pkg/middlewares"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, cfg *config.Config) {
	auth := middlewares.NewAuthMiddleware(cfg.JWT.Secret)
	proxy := handlers.NewProxyHandler(cfg)
	dashboard := handlers.NewDashboardHandler(cfg)

	// ── PUBLIC routes — no JWT required ──────────────────
	public := r.Group("/api/v1")
	{
		// public content — articles (GET only from content service)
		public.GET("/content/articles/*path", proxy.ToContentService)

		// public catalog — browse cars and packages (GET only)
		public.GET("/catalog/*path", proxy.ToCatalogService)

		// general settings (public — used by frontend on load)
		public.GET("/general-settings", proxy.ToCoreServiceDirect)
		public.GET("/general-settings/*path", proxy.ToCoreServiceDirect)

		// auth endpoints (login, register, forgot password, confirm reset, OTP, refresh)
		public.Any("/auth/*path", proxy.ToUserServiceDirect)

		// regions (GET only, public)
		public.GET("/regions", proxy.ToCoreServiceDirect)
		public.GET("/regions/*path", proxy.ToCoreServiceDirect)
	}

	// ── MIXED routes — conditional JWT ───────────────────────
	mixed := r.Group("/api/v1")
	{
		// users — auth routes skip JWT, everything else requires it
		mixed.Any("/users/*path", func(c *gin.Context) {
			path := c.Param("path")
			if path == "/auth" || strings.HasPrefix(path, "/auth/") {
				proxy.ToUserService(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserService(c)
			}
		})

		// instructors — public for GET, auth for write
		mixed.Any("/instructors/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToUserService(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserService(c)
			}
		})

		// cars — GET is public, write ops require auth
		mixed.Any("/cars", proxy.ToCoreServiceDirect)
		mixed.Any("/cars/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})

		// addons — GET is public, write ops require auth
		mixed.Any("/addons", proxy.ToCoreServiceDirect)
		mixed.Any("/addons/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})

		// testimonials — GET is public, mutating methods require auth
		mixed.Any("/testimonials/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToUserService(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserService(c)
			}
		})

		// articles — GET is public, mutating methods require auth
		mixed.Any("/articles/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreService(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreService(c)
			}
		})

		// packages — GET is public, mutating methods require auth
		mixed.Any("/packages/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreService(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreService(c)
			}
		})

		// payments — notification webhooks skip auth, everything else requires auth
		mixed.Any("/payments/*path", func(c *gin.Context) {
			path := c.Param("path")
			isWebhook := (path == "/notification" || path == "/doku/notification" ||
				strings.HasPrefix(path, "/doku/notify")) && c.Request.Method == http.MethodPost
			if isWebhook {
				proxy.ToPaymentService(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToPaymentService(c)
			}
		})

		// vehicles — GET is public (already registered), write ops require auth
		mixed.Any("/vehicles/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})

		// schedules — GET is public, write ops require auth
		mixed.Any("/schedules/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToBookingServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToBookingServiceDirect(c)
			}
		})
	}

	// ── PROTECTED routes — JWT required ──────────────────
	protected := r.Group("/api/v1")
	protected.Use(auth.Authenticate())
	{
		// booking
		protected.Any("/bookings/*path", proxy.ToBookingService)

		// vouchers
		protected.Any("/vouchers/*path", proxy.ToVoucherService)

		// notifications
		protected.Any("/notifications/*path", proxy.ToNotificationService)

		// sessions
		protected.Any("/sessions/*path", proxy.ToBookingServiceDirect)

		// members
		protected.Any("/members", proxy.ToUserServiceDirect)
		protected.Any("/members/*path", proxy.ToUserServiceDirect)

		// certificates
		protected.Any("/certificates", proxy.ToUserServiceDirect)
		protected.Any("/certificates/*path", proxy.ToUserServiceDirect)

		// entitlements
		protected.Any("/entitlements", proxy.ToUserServiceDirect)
		protected.Any("/entitlements/*path", proxy.ToUserServiceDirect)

		// transactions
		protected.Any("/transactions", proxy.ToPaymentService)
		protected.Any("/transactions/*path", proxy.ToPaymentService)
	}

	// ── ADMIN routes — JWT + admin role required ──────────
	admin := r.Group("/api/v1/admin")
	admin.Use(auth.Authenticate())
	admin.Use(auth.RequireRole("admin"))
	{
		admin.Any("/content/*path", proxy.ToContentService)
		admin.Any("/catalog/*path", proxy.ToCatalogService)
		admin.Any("/users/*path", proxy.ToUserService)
		admin.Any("/analytics/*path", proxy.ToCoreService)
		admin.Any("/sales", proxy.ToCoreService)
		admin.Any("/sales/*path", proxy.ToCoreService)
		admin.GET("/dashboard/stats", dashboard.GetDashboardStats)
	}

	// ── CORE ADMIN routes — routed under /api/v1/core/admin ──
	coreAdmin := r.Group("/api/v1/core/admin")
	coreAdmin.Use(auth.Authenticate())
	coreAdmin.Use(auth.RequireRole("admin"))
	{
		coreAdmin.Any("/analytics/*path", proxy.ToCoreService)
	}
}