package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimitConfig struct {
	Enabled  bool
	Requests int           // e.g., 100
	Interval time.Duration // e.g., 1 * time.Minute
}

var (
	mu              sync.Mutex
	clientLimiters  = make(map[string]*clientLimiter)
	cleanupInterval = time.Minute * 5
)

// Middleware
func RateLimiterMiddleware(config RateLimitConfig) gin.HandlerFunc {
	if !config.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	go cleanupStaleClients()

	return func(c *gin.Context) {
		ip := getClientIP(c)

		limiter := getRateLimiter(ip, config.Requests, config.Interval)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		c.Next()
	}
}

// Get IP from common headers or remote addr
func getClientIP(c *gin.Context) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"} {
		if ip := c.GetHeader(h); ip != "" {
			return strings.Split(ip, ",")[0]
		}
	}
	ip := c.ClientIP()
	if ip == "" {
		ip = c.Request.RemoteAddr
	}
	return ip
}

// Create or retrieve a rate limiter for a given IP
func getRateLimiter(ip string, rps int, interval time.Duration) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := clientLimiters[ip]
	if !exists {
		limiter = &clientLimiter{
			limiter:  rate.NewLimiter(rate.Every(interval/time.Duration(rps)), rps),
			lastSeen: time.Now(),
		}
		clientLimiters[ip] = limiter
	} else {
		limiter.lastSeen = time.Now()
	}

	return limiter.limiter
}

// Clean up old entries to prevent memory leaks
func cleanupStaleClients() {
	for {
		time.Sleep(cleanupInterval)
		mu.Lock()
		for ip, client := range clientLimiters {
			if time.Since(client.lastSeen) > cleanupInterval {
				delete(clientLimiters, ip)
			}
		}
		mu.Unlock()
	}
}
