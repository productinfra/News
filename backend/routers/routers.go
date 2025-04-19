package routers

import (
	"backend/controller"
	_ "backend/docs" // Don't forget to import the docs generated in the previous step
	"backend/logger"
	"backend/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-contrib/pprof"
)

// SetupRouter Set up the router
func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // Set to release mode
	}
	// Initialize gin Engine with no default middlewares
	r := gin.New()
	// Set up middleware
	r.Use(logger.GinLogger(),
		logger.GinRecovery(true),                           // Recovery middleware will recover from panics and log relevant logs using zap
		middlewares.RateLimitMiddleware(2*time.Second, 40), // Global rate limiting: 40 tokens every 2 seconds
	)

	r.LoadHTMLFiles("templates/index.html") // Load HTML file
	r.Static("/static", "./static")         // Serve static files
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// Register swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	// Register login business
	v1.POST("/login", controller.LoginHandler)               // Login business
	v1.POST("/signup", controller.SignUpHandler)             // Registration business
	v1.GET("/refresh_token", controller.RefreshTokenHandler) // Refresh accessToken
	// Post business
	v1.GET("/posts", controller.PostListHandler)      // Paginated post list
	v1.GET("/posts2", controller.PostList2Handler)    // Paginated post list sorted by community ID, time, or score
	v1.GET("/post/:id", controller.PostDetailHandler) // Query post details
	v1.GET("/search", controller.PostSearchHandler)   // Search business - search posts
	// Community business
	v1.GET("/community", controller.CommunityHandler)           // Get community list
	v1.GET("/community/:id", controller.CommunityDetailHandler) // Get community details by ID
	// Github Trending
	v1.GET("/github_trending", controller.GithubTrendingHandler) // Github Trending
	v1.GET("/news", controller.News)
	v1.POST("/gemini", controller.Gemini)

	// Middleware
	v1.Use(middlewares.JWTAuthMiddleware()) // Apply JWT authentication middleware
	{
		v1.POST("/post", controller.CreatePostHandler) // Create a post

		v1.POST("/vote", controller.VoteHandler) // Voting

		v1.POST("/comment", controller.CommentHandler)    // Post comment
		v1.GET("/comment", controller.CommentListHandler) // Comment list

		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}

	pprof.Register(r) // Register pprof related routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
