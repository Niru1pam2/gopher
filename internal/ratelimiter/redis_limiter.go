package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRateLimiter struct {
	rdb    *redis.Client
	limit  int64
	window time.Duration
}

func NewRedisRateLimiter(rdb *redis.Client, limit int64, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		rdb:    rdb,
		limit:  limit,
		window: window,
	}
}

func (rl *RedisRateLimiter) Allow(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", ip)

	count, err := rl.rdb.Incr(ctx, key).Result()

	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit: %w", err)
	}

	if count == 1 {
		err := rl.rdb.Expire(ctx, key, rl.window).Err()
		if err != nil {
			return false, fmt.Errorf("failed to set expiration: %w", err)
		}
	}

	if count > rl.limit {
		return false, nil
	}

	return true, nil
}
