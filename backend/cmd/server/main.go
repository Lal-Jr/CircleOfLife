package main

import (
	"context"
	"log"

	"circleoflife/internal/cache"
	"circleoflife/internal/config"
	"circleoflife/internal/db"
	"circleoflife/internal/events"
	"circleoflife/internal/handlers"
	"circleoflife/internal/middleware"
	"circleoflife/internal/repository"
	"circleoflife/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Setup Configuration
	cfg := config.LoadConfig()

	// 2. Setup Database Connection (includes running PostGIS migrations)
	db.ConnectDB(cfg.DatabaseURL)

	// 3. Setup Redis Cache and Pub/Sub Event Subscription
	cache.InitRedis(cfg.RedisURL)
	events.SubscribePostEvents(context.Background())

	// 4. Initialize Repositories
	userRepo := repository.NewUserRepository()
	postRepo := repository.NewPostRepository()
	commentRepo := repository.NewCommentRepository()

	// 5. Initialize Services
	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret)
	postSvc := services.NewPostService(postRepo)
	commentSvc := services.NewCommentService(commentRepo)

	// 6. Initialize Handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	postHandler := handlers.NewPostHandler(postSvc)
	commentHandler := handlers.NewCommentHandler(commentSvc)

	// 7. Setup Router
	r := gin.Default()

	// Global Middleware Setup
	r.Use(gin.Recovery())

	// Implement simple CORS assuming a local frontend running on 3000
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API Routes definition
	api := r.Group("/api")

	// Authentication API
	auth := api.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/login", authHandler.Login)
	}

	// Protected Routes (JWT required)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Real-time Event Stream (SSE)
	protected.GET("/events", handlers.HandleEvents)

	// Posts API
	posts := protected.Group("/posts")
	{
		posts.GET("", postHandler.GetNearbyPosts)
		posts.POST("", postHandler.CreatePost)
		posts.GET("/:id", postHandler.GetPostByID)

		// Nested Comments logic
		posts.GET("/:id/comments", commentHandler.GetComments)
		posts.POST("/:id/comments", commentHandler.CreateComment)
	}

	// 8. Start the server
	log.Printf("Starting Server natively running on %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server inherently failed to launch: %v", err)
	}
}
