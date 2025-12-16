package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
	limit  int           // Requests per window
	window time.Duration // Window size
}

func NewLimiter(addr string, limit int, window time.Duration) *Limiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Limiter{
		client: rdb,
		limit:  limit,
		window: window,
	}
}

func (l *Limiter) Allow(ctx context.Context, key string) (bool, error) {
	// Simple sliding window using Redis expiration
	// Key: rate_limit:<key>
	// We'll use INCR and EXPIRE
	
	rKey := fmt.Sprintf("rate_limit:%s", key)
	
	count, err := l.client.Incr(ctx, rKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		l.client.Expire(ctx, rKey, l.window)
	}

	if count > int64(l.limit) {
		return false, nil
	}

	return true, nil
}
