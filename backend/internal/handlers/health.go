package handlers

import (
	"net/http"
	"circleoflife/internal/cache"
	"circleoflife/internal/db"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	status := map[string]string{
		"status": "ok",
		"database": "connected",
		"redis": "connected",
	}

	if err := db.Pool.Ping(c.Request.Context()); err != nil {
		status["database"] = "disconnected"
		status["status"] = "degraded"
	}

	if cache.Client != nil {
		if err := cache.Client.Ping(c.Request.Context()).Err(); err != nil {
			status["redis"] = "disconnected"
			status["status"] = "degraded"
		}
	} else {
		status["redis"] = "not_initialized"
	}

	code := http.StatusOK
	if status["status"] == "degraded" {
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, status)
}
