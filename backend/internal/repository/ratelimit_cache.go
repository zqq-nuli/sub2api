package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const authRateLimitKeyPrefix = "auth_ratelimit:"

type rateLimitCache struct {
	rdb *redis.Client
}

// NewRateLimitCache creates a new rate limit cache implementation
func NewRateLimitCache(rdb *redis.Client) service.RateLimitCache {
	return &rateLimitCache{rdb: rdb}
}

// IncrementAndCheck increments the counter for the given key and returns
// whether the limit has been exceeded. Uses sliding window algorithm.
func (c *rateLimitCache) IncrementAndCheck(ctx context.Context, key string, limit int, window time.Duration) (allowed bool, remaining int, err error) {
	if c.rdb == nil {
		// Fail-open if Redis is unavailable
		return true, limit, nil
	}

	fullKey := authRateLimitKeyPrefix + key
	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	// Use Lua script for atomic sliding window rate limiting
	script := redis.NewScript(`
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local window_start = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		local ttl = tonumber(ARGV[4])

		-- Remove old entries outside the window
		redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start)

		-- Count current entries
		local count = redis.call('ZCARD', key)

		if count < limit then
			-- Add new entry with current timestamp as score
			redis.call('ZADD', key, now, now .. ':' .. math.random(1000000))
			redis.call('EXPIRE', key, ttl)
			return {1, limit - count - 1}
		else
			-- Rate limited
			redis.call('EXPIRE', key, ttl)
			return {0, 0}
		end
	`)

	result, err := script.Run(ctx, c.rdb, []string{fullKey}, now, windowStart, limit, int64(window.Seconds())+1).Slice()
	if err != nil {
		// Fail-open on error
		return true, limit, err
	}

	allowedInt := result[0].(int64)
	remainingInt := result[1].(int64)

	return allowedInt == 1, int(remainingInt), nil
}
