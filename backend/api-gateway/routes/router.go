package routes

import (
	"api-gateway/handlers"
	"api-gateway/pkg/config"
	"api-gateway/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, cfg *config.Config) {
	auth := middlewares.NewAuthMiddleware(cfg.JWT.Secret)
    proxy := handlers.NewProxyHandler(cfg)
    dashboard := handlers.NewDashboardHandler(cfg)

    // ── PUBLIC routes — no JWT required ──────────────────
    public := r.Group("/api/v1")
    {
        // public content — articles
        public.GET("/content/articles/*path", proxy.ToContentService)

        // public catalog — browse cars and packages
        public.GET("/catalog/*path", proxy.ToCatalogService)

        // testimonials (public endpoints: published, featured)
        public.Any("/testimonials/*path", proxy.ToUserService)

        // packages (public browsing)
        public.Any("/packages/*path", proxy.ToCoreService)

        // articles (includes FAQ endpoints)
        public.Any("/articles/*path", proxy.ToCoreService)
    }

    // ── MIXED routes — conditional JWT ───────────────────────
    mixed := r.Group("/api/v1")
    {
        mixed.Any("/users/*path", func(c *gin.Context) {
            importStrings := "strings"
            _ = importStrings
            // If the path is an auth endpoint, skip JWT validation
            if len(c.Param("path")) >= 5 && c.Param("path")[0:6] == "/auth/" {
                proxy.ToUserService(c)
                return
            }
            if c.Param("path") == "/auth" {
                proxy.ToUserService(c)
                return
            }

            // Otherwise, validate JWT
            auth.Authenticate()(c)
            if !c.IsAborted() {
                proxy.ToUserService(c)
            }
        })

        // instructors routes
        mixed.Any("/instructors/*path", proxy.ToUserService)
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

        // testimonials admin operations (create, update, delete)
        protected.Any("/testimonials/*path", proxy.ToUserService)

        // testimonials admin operations (create, update, delete)
        protected.Any("/articles/*path", proxy.ToContentService)
    }

    // ── ADMIN routes — JWT + admin role required ──────────
    admin := r.Group("/api/v1/admin")
    admin.Use(auth.Authenticate())
    admin.Use(auth.RequireRole("admin"))
    {
        admin.Any("/content/*path",  proxy.ToContentService)
        admin.Any("/catalog/*path",  proxy.ToCatalogService)
        admin.Any("/users/*path",    proxy.ToUserService)
        admin.Any("/analytics/*path", proxy.ToCoreService)
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