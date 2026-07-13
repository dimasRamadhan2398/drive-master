package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

// ArticleRoute handles article/blog routes
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

// Run registers all article and blog routes
func (r *ArticleRoute) Run() {
	articles := r.group.Group("/articles")
	{
		// Blog endpoints
		articles.GET("/blog", r.controller.GetArticleController().GetBlogArticles)
		articles.GET("/blog/:id", r.controller.GetArticleController().GetBlogPostByID)
		articles.POST("/blog", r.controller.GetArticleController().CreateBlogPost)
		articles.PUT("/blog/:id", r.controller.GetArticleController().UpdateBlogPost)
		articles.DELETE("/blog/:id", r.controller.GetArticleController().DeleteBlogPost)

		// Blog view metrics
		articles.POST("/:id/view", r.controller.GetArticleController().IncrementViewCount)

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