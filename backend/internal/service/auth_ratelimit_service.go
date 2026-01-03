package service

import (
	"context"
	"time"
)

// RateLimitCache defines the interface for rate limiting storage
type RateLimitCache interface {
	// IncrementAndCheck increments counter and checks if limit is exceeded
	// Returns: allowed (bool), remaining requests (int), error
	IncrementAndCheck(ctx context.Context, key string, limit int, window time.Duration) (bool, int, error)
}

// AuthRateLimitConfig holds rate limit configuration for auth endpoints
type AuthRateLimitConfig struct {
	// LoginLimit is max login attempts per IP per window
	LoginLimit int
	// LoginWindow is the time window for login rate limiting
	LoginWindow time.Duration
	// RegisterLimit is max registration attempts per IP per window
	RegisterLimit int
	// RegisterWindow is the time window for registration rate limiting
	RegisterWindow time.Duration
	// VerifyCodeLimit is max verification code requests per IP per window
	VerifyCodeLimit int
	// VerifyCodeWindow is the time window for verification code rate limiting
	VerifyCodeWindow time.Duration
}

// DefaultAuthRateLimitConfig returns default rate limit configuration
func DefaultAuthRateLimitConfig() *AuthRateLimitConfig {
	return &AuthRateLimitConfig{
		LoginLimit:       10,              // 10 attempts per minute
		LoginWindow:      time.Minute,     // 1 minute window
		RegisterLimit:    5,               // 5 registrations per hour
		RegisterWindow:   time.Hour,       // 1 hour window
		VerifyCodeLimit:  3,               // 3 codes per 5 minutes
		VerifyCodeWindow: 5 * time.Minute, // 5 minute window
	}
}

// AuthRateLimitService provides rate limiting for authentication endpoints
type AuthRateLimitService struct {
	cache  RateLimitCache
	config *AuthRateLimitConfig
}

// NewAuthRateLimitService creates a new auth rate limit service
func NewAuthRateLimitService(cache RateLimitCache) *AuthRateLimitService {
	return &AuthRateLimitService{
		cache:  cache,
		config: DefaultAuthRateLimitConfig(),
	}
}

// CheckLogin checks if login attempt is allowed for the given IP
func (s *AuthRateLimitService) CheckLogin(ctx context.Context, ip string) (allowed bool, remaining int, err error) {
	if s.cache == nil {
		return true, s.config.LoginLimit, nil
	}
	return s.cache.IncrementAndCheck(ctx, "login:"+ip, s.config.LoginLimit, s.config.LoginWindow)
}

// CheckRegister checks if registration attempt is allowed for the given IP
func (s *AuthRateLimitService) CheckRegister(ctx context.Context, ip string) (allowed bool, remaining int, err error) {
	if s.cache == nil {
		return true, s.config.RegisterLimit, nil
	}
	return s.cache.IncrementAndCheck(ctx, "register:"+ip, s.config.RegisterLimit, s.config.RegisterWindow)
}

// CheckVerifyCode checks if verification code request is allowed for the given IP
func (s *AuthRateLimitService) CheckVerifyCode(ctx context.Context, ip string) (allowed bool, remaining int, err error) {
	if s.cache == nil {
		return true, s.config.VerifyCodeLimit, nil
	}
	return s.cache.IncrementAndCheck(ctx, "verify:"+ip, s.config.VerifyCodeLimit, s.config.VerifyCodeWindow)
}
