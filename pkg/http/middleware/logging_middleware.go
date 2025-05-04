package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	// Log the request details
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Log the request details
		c.Writer.Header().Set("X-Response-Time", duration.String())
	}
}
