package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"circleoflife/pkg/auth"
	"circleoflife/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		// Browsers do not support Authorization headers on generic EventSource objects.
		// Fallback to checking a token query param.
		if authHeader == "" {
			tokenString = c.Query("token")
			if tokenString == "" {
				utils.JSONError(c, http.StatusUnauthorized, "Authorization required")
				c.Abort()
				return
			}
		} else {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				utils.JSONError(c, http.StatusUnauthorized, "Authorization format must be Bearer {token}")
				c.Abort()
				return
			}
			tokenString = parts[1]
		}
		userID, err := auth.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			utils.JSONError(c, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("userID", userID)
		c.Next()
	}
}
