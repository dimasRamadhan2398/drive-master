package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

// ArticleRoute handles article routes
type ArticleRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

// IArticleRoute defines the interface for article route
type IArticleRoute interface {
	Run()
}

func NewArticleRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IArticleRoute {
	return &ArticleRoute{
		controller: controller,
		group:      group,
	}
}

// Run registers all article routes
func (r *ArticleRoute) Run() {
	articles := r.group.Group("/articles")
	{
		// Admin CRUD endpoints
		articles.GET("", r.controller.GetArticleController().GetAllArticles)
		articles.GET("/:id", r.controller.GetArticleController().GetArticleByID)
		articles.POST("", r.controller.GetArticleController().CreateArticle)
		articles.PUT("/:id", r.controller.GetArticleController().UpdateArticle)
		articles.DELETE("/:id", r.controller.GetArticleController().DeleteArticle)

		// Public endpoints
		articles.GET("/slug/:slug", r.controller.GetArticleController().GetArticleBySlug)
		articles.GET("/search", r.controller.GetArticleController().SearchArticles)
		articles.GET("/featured", r.controller.GetArticleController().GetFeaturedArticles)
		articles.GET("/spotlight", r.controller.GetArticleController().GetSpotlightArticle)
		articles.GET("/tag/:tag", r.controller.GetArticleController().GetArticlesByTag)

		// Article-specific endpoints
		articles.GET("/:id/related", r.controller.GetArticleController().GetRelatedArticles)
		articles.POST("/:id/view", r.controller.GetArticleController().IncrementViewCount)
		articles.POST("/:id/publish", r.controller.GetArticleController().PublishArticle)
		articles.POST("/:id/archive", r.controller.GetArticleController().ArchiveArticle)

		// blog endpoints
		articles.GET("/blog", r.controller.GetArticleController().GetBlogArticles)
		articles.POST("/blog", r.controller.GetArticleController().CreateBlogPost)
		articles.PUT("/blog/:id", r.controller.GetArticleController().UpdateBlogPost)
		articles.DELETE("/blog/:id", r.controller.GetArticleController().DeleteBlogArticle)

		// FAQ endpoints
		articles.GET("/faq", r.controller.GetFAQController().GetAllFAQs)
		articles.GET("/faq/active", r.controller.GetFAQController().GetActiveFAQs)
		articles.POST("/faq", r.controller.GetFAQController().CreateFAQ)
		// Specific routes MUST come before parameterized routes
		articles.PUT("/faq/:id/reorder", r.controller.GetFAQController().ReorderFAQ)
		articles.GET("/faq/:id", r.controller.GetFAQController().GetFAQByID)
		articles.PUT("/faq/:id", r.controller.GetFAQController().UpdateFAQ)
		articles.DELETE("/faq/:id", r.controller.GetFAQController().DeleteFAQ)
	}
}