package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimitErrorHandler(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"error": "Too many requests, please try again later",
		"code":  "RATE_LIMIT_EXCEEDED",
	})
	c.Abort()
}
