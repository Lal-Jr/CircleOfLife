package utils

import "github.com/gin-gonic/gin"

func JSONResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, gin.H{
		"message": message,
		"data":    data,
	})
}

func JSONError(c *gin.Context, status int, err string) {
	c.JSON(status, gin.H{
		"error": err,
	})
}
