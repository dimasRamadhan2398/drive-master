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

    // ── PUBLIC routes — no JWT required ──────────────────
    public := r.Group("/api/v1")
    {
        // auth endpoints — anyone can hit these
        public.Any("/users/auth/*path",     proxy.ToUserService)
        public.Any("/users/auth/register",  proxy.ToUserService)

        // public content — articles, testimonials
        public.GET("/content/articles/*path",     proxy.ToContentService)
        public.GET("/content/testimonials/*path", proxy.ToContentService)

        // public catalog — browse cars and packages
        public.GET("/catalog/*path", proxy.ToCatalogService)
    }

    // ── PROTECTED routes — JWT required ──────────────────
    protected := r.Group("/api/v1")
    protected.Use(auth.Authenticate())
    {
        // user management
        protected.Any("/users/*path", proxy.ToUserService)

        // booking
        protected.Any("/bookings/*path", proxy.ToBookingService)

        // vouchers
        protected.Any("/vouchers/*path", proxy.ToVoucherService)

        // notifications
        protected.Any("/notifications/*path", proxy.ToNotificationService)
    }

    // ── ADMIN routes — JWT + admin role required ──────────
    admin := r.Group("/api/v1/admin")
    admin.Use(auth.Authenticate())
    admin.Use(auth.RequireRole("admin"))
    {
        admin.Any("/content/*path",  proxy.ToContentService)
        admin.Any("/catalog/*path",  proxy.ToCatalogService)
        admin.Any("/users/*path",    proxy.ToUserService)
    }
}