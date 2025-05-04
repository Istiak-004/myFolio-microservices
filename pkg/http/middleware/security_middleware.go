package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; media-src 'self'; object-src 'none'; child-src 'none'; form-action 'self'; frame-ancestors 'none'; base-uri 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Feature-Policy", "geolocation 'none'; microphone 'none'; camera 'none'")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Remove server header
		c.Header("Server", "")

		c.Next()
	}
}
