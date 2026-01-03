package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const ssoSessionKeyPrefix = "sso:session:"

// ssoSessionKey generates the Redis key for SSO session.
func ssoSessionKey(sessionID string) string {
	return ssoSessionKeyPrefix + sessionID
}

type ssoSessionCache struct {
	rdb *redis.Client
}

// NewSSOSessionCache creates a new SSOSessionCache implementation.
func NewSSOSessionCache(rdb *redis.Client) service.OIDCSessionCache {
	return &ssoSessionCache{rdb: rdb}
}

func (c *ssoSessionCache) Set(ctx context.Context, sessionID string, session *service.OIDCSession, ttl time.Duration) error {
	key := ssoSessionKey(sessionID)
	val, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, val, ttl).Err()
}

func (c *ssoSessionCache) Get(ctx context.Context, sessionID string) (*service.OIDCSession, error) {
	key := ssoSessionKey(sessionID)
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err // includes redis.Nil for key not found
	}
	var session service.OIDCSession
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (c *ssoSessionCache) Delete(ctx context.Context, sessionID string) error {
	key := ssoSessionKey(sessionID)
	return c.rdb.Del(ctx, key).Err()
}
