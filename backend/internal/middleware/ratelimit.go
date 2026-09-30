package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"circleoflife/internal/cache"
	"circleoflife/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware blocks abuse with a fixed-window counter in Redis, per user (or IP when
// signed out): each request INCRs a key that expires after the window, and requests past the limit
// get 429 until the key expires.
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cache.Client == nil {
			c.Next() // Bypass silently if no Redis infrastructure exists natively yet
			return
		}

		identifier := c.ClientIP()
		if userID, exists := c.Get("userID"); exists {
			identifier = userID.(string) // Prefer hard user locks if logged in
		}

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		// c.FullPath() returns the same route pattern regardless of which
		// middleware in the chain reads it (it doesn't vary by HTTP method
		// either), so a broad group-level limiter and a stricter per-route
		// limiter stacked on the same route would otherwise collide on one
		// shared key: every request then increments both counters, and the
		// stricter one trips long before its real limit. Namespacing the key
		// by this middleware's own (method, limit, window) keeps every
		// configured limiter's counter independent.
		key := fmt.Sprintf("ratelimit:%s:%s:%s:%d:%s", identifier, c.Request.Method, endpoint, limit, window)

		count, err := cache.Client.Incr(context.Background(), key).Result()
		if err != nil {
			c.Next() // Fail open to maintain availability
			return
		}

		if count == 1 {
			// First request in bucket window, start the TTL
			cache.Client.Expire(context.Background(), key, window)
		}

		if count > int64(limit) {
			utils.JSONError(c, http.StatusTooManyRequests, "Rate limit exceeded. Try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}
