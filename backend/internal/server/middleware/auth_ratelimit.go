package middleware

import (
	"net/http"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthRateLimitMiddleware provides rate limiting for authentication endpoints
type AuthRateLimitMiddleware struct {
	service *service.AuthRateLimitService
}

// NewAuthRateLimitMiddleware creates a new auth rate limit middleware
func NewAuthRateLimitMiddleware(service *service.AuthRateLimitService) *AuthRateLimitMiddleware {
	return &AuthRateLimitMiddleware{service: service}
}

// getClientIP extracts the real client IP from the request
func (m *AuthRateLimitMiddleware) getClientIP(c *gin.Context) string {
	// Try X-Forwarded-For first (common for proxied requests)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		return xff
	}
	// Try X-Real-IP
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to remote address
	return c.ClientIP()
}

// setRateLimitHeaders sets standard rate limit response headers
func (m *AuthRateLimitMiddleware) setRateLimitHeaders(c *gin.Context, remaining int) {
	c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
}

// rateLimitExceeded sends a 429 response
func (m *AuthRateLimitMiddleware) rateLimitExceeded(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"code":    "RATE_LIMIT_EXCEEDED",
		"message": "Too many requests, please try again later",
	})
	c.Abort()
}

// LoginRateLimit returns middleware for login endpoint rate limiting
func (m *AuthRateLimitMiddleware) LoginRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := m.getClientIP(c)
		allowed, remaining, _ := m.service.CheckLogin(c.Request.Context(), ip)
		m.setRateLimitHeaders(c, remaining)
		if !allowed {
			m.rateLimitExceeded(c)
			return
		}
		c.Next()
	}
}

// RegisterRateLimit returns middleware for registration endpoint rate limiting
func (m *AuthRateLimitMiddleware) RegisterRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := m.getClientIP(c)
		allowed, remaining, _ := m.service.CheckRegister(c.Request.Context(), ip)
		m.setRateLimitHeaders(c, remaining)
		if !allowed {
			m.rateLimitExceeded(c)
			return
		}
		c.Next()
	}
}

// VerifyCodeRateLimit returns middleware for verification code endpoint rate limiting
func (m *AuthRateLimitMiddleware) VerifyCodeRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := m.getClientIP(c)
		allowed, remaining, _ := m.service.CheckVerifyCode(c.Request.Context(), ip)
		m.setRateLimitHeaders(c, remaining)
		if !allowed {
			m.rateLimitExceeded(c)
			return
		}
		c.Next()
	}
}
