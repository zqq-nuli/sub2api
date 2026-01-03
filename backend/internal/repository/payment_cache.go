package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	paymentLockKeyPrefix = "payment:lock:"
)

// paymentLockKey generates the Redis key for payment lock.
func paymentLockKey(orderNo string) string {
	return paymentLockKeyPrefix + orderNo
}

type paymentCache struct {
	rdb *redis.Client
}

// NewPaymentCache creates a new PaymentCache implementation.
func NewPaymentCache(rdb *redis.Client) service.PaymentCache {
	return &paymentCache{rdb: rdb}
}

func (c *paymentCache) AcquirePaymentLock(ctx context.Context, orderNo string, ttl time.Duration) (bool, error) {
	key := paymentLockKey(orderNo)
	return c.rdb.SetNX(ctx, key, 1, ttl).Result()
}

func (c *paymentCache) ReleasePaymentLock(ctx context.Context, orderNo string) error {
	key := paymentLockKey(orderNo)
	return c.rdb.Del(ctx, key).Err()
}
