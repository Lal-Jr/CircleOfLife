// Command server runs the Circle of Life API. This file is the composition root: it builds
// each bounded context from its infrastructure adapters and mounts its HTTP routes.
//
//	Identity   members, sign-up, log-in and session tokens
//	Community  posts, comments and helpful votes, and the ranked nearby feed
//	Realtime   live updates to open browsers over SSE, fanned out through Redis
package main

import (
	"context"
	"log"
	"time"

	"circleoflife/internal/cache"
	community "circleoflife/internal/community/application"
	communityinfra "circleoflife/internal/community/infrastructure"
	communityhttp "circleoflife/internal/community/interfaces/httpapi"
	"circleoflife/internal/config"
	"circleoflife/internal/db"
	identity "circleoflife/internal/identity/application"
	identityinfra "circleoflife/internal/identity/infrastructure"
	identityhttp "circleoflife/internal/identity/interfaces/httpapi"
	"circleoflife/internal/middleware"
	"circleoflife/internal/realtime"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Shared infrastructure: Postgres (with PostGIS migrations and demo seed) and Redis.
	db.ConnectDB(cfg.DatabaseURL)
	cache.InitRedis(cfg.RedisURL)

	// Realtime context.
	bus := realtime.NewBus(cache.Client)
	bus.Listen(context.Background())

	// Identity context.
	identityApp := identity.NewService(
		identityinfra.NewUsers(db.Pool),
		identityinfra.Bcrypt{},
		identityinfra.NewJWTIssuer(cfg.JWTSecret),
	)
	identityAPI := identityhttp.NewHandler(identityApp)

	// Community context.
	communityApp := community.NewService(
		communityinfra.NewPosts(db.Pool),
		communityinfra.NewCachedFeed(communityinfra.NewFeed(db.Pool), cache.Client),
		communityinfra.NewLiveUpdates(bus),
		communityinfra.NewPlainText(),
	)
	communityAPI := communityhttp.NewHandler(communityApp)

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(cors)

	api := r.Group("/api")
	api.Use(middleware.RequestIDMiddleware())
	api.Use(middleware.LoggerMiddleware())
	api.GET("/health", health)

	auth := api.Group("/auth")
	auth.POST("/signup", identityAPI.Signup)
	auth.POST("/login", identityAPI.Login)

	// Everything below requires a signed-in member (JWT).
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	protected.GET("/events", bus.Stream)
	protected.GET("/users/me", identityAPI.Me)
	communityAPI.Register(protected.Group("/posts"), func(perMinute int) gin.HandlerFunc {
		return middleware.RateLimitMiddleware(perMinute, time.Minute)
	})

	log.Printf("Starting server on %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// cors allows the browser frontend, which is served from another origin, to call the API.
func cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}
