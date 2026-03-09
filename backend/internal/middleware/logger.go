package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware replaces the default Gin logger with structured Zerolog JSON traces
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		reqID, _ := c.Get("RequestID")
		userID, _ := c.Get("userID") // May be nil on unprotected routes e.g., Login

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		evt := log.Info()
		if len(c.Errors) > 0 {
			evt = log.Error().Err(c.Errors.Last())
		} else if status >= 400 {
			evt = log.Warn()
		}

		evt.
			Str("request_id", reqID.(string)).
			Str("method", c.Request.Method).
			Str("endpoint", c.Request.URL.Path).
			Int("status", status).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP())

		if userID != nil {
			evt.Str("user_id", userID.(string))
		}

		evt.Msg("Request processed")
	}
}
