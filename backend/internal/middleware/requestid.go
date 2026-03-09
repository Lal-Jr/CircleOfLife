package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware injects a unique X-Request-ID onto every incoming API HTTP header
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set("RequestID", reqID)
		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next()
	}
}
