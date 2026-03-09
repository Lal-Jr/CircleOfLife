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

// RateLimitMiddleware blocks abuse by applying a Redis Token Bucket counter per IP/User over the given limit and window.
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

		key := fmt.Sprintf("ratelimit:%s:%s", identifier, endpoint)

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
