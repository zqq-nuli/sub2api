package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	rechargeProductCacheKey = "recharge:products:active"
)

type rechargeProductCache struct {
	rdb *redis.Client
}

// NewRechargeProductCache creates a new RechargeProductCache implementation.
func NewRechargeProductCache(rdb *redis.Client) service.RechargeProductCache {
	return &rechargeProductCache{rdb: rdb}
}

func (c *rechargeProductCache) GetActiveProducts(ctx context.Context) (string, error) {
	return c.rdb.Get(ctx, rechargeProductCacheKey).Result()
}

func (c *rechargeProductCache) SetActiveProducts(ctx context.Context, data string, ttl time.Duration) error {
	return c.rdb.Set(ctx, rechargeProductCacheKey, data, ttl).Err()
}

func (c *rechargeProductCache) InvalidateActiveProducts(ctx context.Context) error {
	return c.rdb.Del(ctx, rechargeProductCacheKey).Err()
}
