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
		public.GET("/content/articles", proxy.ToContentService)
		public.GET("/content/articles/*path", proxy.ToContentService)

		// public catalog — browse cars and packages (GET only)
		public.GET("/catalog", proxy.ToCatalogService)
		public.GET("/catalog/*path", proxy.ToCatalogService)



		// auth endpoints (login, register, forgot password, confirm reset, OTP, refresh)
		public.Any("/auth", proxy.ToUserServiceDirect)
		public.Any("/auth/*path", proxy.ToUserServiceDirect)

		// regions (GET only, public)
		public.GET("/regions", proxy.ToCoreServiceDirect)
		public.GET("/regions/*path", proxy.ToCoreServiceDirect)
	}

	// ── MIXED routes — conditional JWT ───────────────────────
	mixed := r.Group("/api/v1")
	{
		// contact — POST is public for inquiries submission, GET requires admin authentication
		mixed.Any("/contact", func(c *gin.Context) {
			if c.Request.Method == http.MethodPost {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			auth.RequireRole("admin")(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})
		// users — JWT required, proxies directly to user-service
		mixed.Any("/users", func(c *gin.Context) {
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserServiceDirect(c)
			}
		})
		mixed.Any("/users/*path", func(c *gin.Context) {
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserServiceDirect(c)
			}
		})

		// instructors — public for GET, auth for write
		mixed.Any("/instructors", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToUserServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserServiceDirect(c)
			}
		})
		mixed.Any("/instructors/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToUserServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserServiceDirect(c)
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
		mixed.Any("/testimonials", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToUserServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserServiceDirect(c)
			}
		})
		mixed.Any("/testimonials/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToUserServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToUserServiceDirect(c)
			}
		})

		// articles — GET is public, mutating methods require auth
		mixed.Any("/articles", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})
		mixed.Any("/articles/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})

		// packages — GET is public, mutating methods require auth
		mixed.Any("/packages", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})
		mixed.Any("/packages/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})

		// payments — notification webhooks skip auth, everything else requires auth
		mixed.Any("/payments", func(c *gin.Context) {
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToPaymentService(c)
			}
		})
		mixed.Any("/payments/*path", func(c *gin.Context) {
			path := c.Param("path")
			isBypassAuth := path == "/notification" || path == "/doku/notification" ||
				strings.HasPrefix(path, "/doku/notify") || strings.HasPrefix(path, "/pakasir") ||
				strings.HasSuffix(path, "/simulate") || path == "/callback"
			if isBypassAuth {
				proxy.ToPaymentService(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToPaymentService(c)
			}
		})

		mixed.Any("/entitlements/sync", func(c *gin.Context) {
			proxy.ToUserServiceDirect(c)
		})

		// vehicles — GET is public (already registered), write ops require auth
		mixed.Any("/vehicles", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})
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
		mixed.Any("/schedules", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToBookingServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToBookingServiceDirect(c)
			}
		})
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

		// pages — GET is public, write ops require auth
		mixed.Any("/pages", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})
		mixed.Any("/pages/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})

		// general settings — GET is public, write ops require admin role
		mixed.Any("/general-settings", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			auth.RequireRole("admin")(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})
		mixed.Any("/general-settings/*path", func(c *gin.Context) {
			if c.Request.Method == http.MethodGet {
				proxy.ToCoreServiceDirect(c)
				return
			}
			auth.Authenticate()(c)
			auth.RequireRole("admin")(c)
			if !c.IsAborted() {
				proxy.ToCoreServiceDirect(c)
			}
		})
	}

	// ── PROTECTED routes — JWT required ──────────────────
	protected := r.Group("/api/v1")
	protected.Use(auth.Authenticate())
	{
		// booking
		protected.Any("/bookings", proxy.ToBookingService)
		protected.Any("/bookings/*path", proxy.ToBookingService)

		// vouchers
		protected.Any("/vouchers", proxy.ToVoucherService)
		protected.Any("/vouchers/*path", proxy.ToVoucherService)

		// notifications
		protected.Any("/notifications", proxy.ToNotificationService)
		protected.Any("/notifications/*path", proxy.ToNotificationService)

		// enrollments
		protected.Any("/enrollments", proxy.ToBookingServiceDirect)
		protected.Any("/enrollments/*path", proxy.ToBookingServiceDirect)

		// sessions
		protected.Any("/sessions", proxy.ToBookingServiceDirect)
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

		// dashboard
		protected.Any("/dashboard/stats", proxy.ToUserServiceDirect)
		protected.Any("/dashboard/certifications/stats", proxy.ToUserServiceDirect)
		protected.Any("/dashboard/sessions/stats", proxy.ToBookingServiceDirect)
		protected.Any("/dashboard/revenue/stats", proxy.ToBookingServiceDirect)

		// transactions
		protected.Any("/transactions", proxy.ToPaymentService)
		protected.Any("/transactions/*path", proxy.ToPaymentService)
	}

	// ── ADMIN routes — JWT + admin role required ──────────
	admin := r.Group("/api/v1/admin")
	admin.Use(auth.Authenticate())
	admin.Use(auth.RequireRole("admin"))
	{
		admin.Any("/content", proxy.ToContentService)
		admin.Any("/content/*path", proxy.ToContentService)
		admin.Any("/catalog", proxy.ToCatalogService)
		admin.Any("/catalog/*path", proxy.ToCatalogService)
		admin.Any("/users", proxy.ToUserServiceDirect)
		admin.Any("/users/*path", proxy.ToUserServiceDirect)
		admin.Any("/analytics", proxy.ToCoreServiceDirect)
		admin.Any("/analytics/*path", proxy.ToCoreServiceDirect)
		admin.Any("/sales", proxy.ToCoreServiceDirect)
		admin.Any("/sales/*path", proxy.ToCoreServiceDirect)
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