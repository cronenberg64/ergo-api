package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// RateLimiter interface allows switching between Redis and Memory implementations
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

// RedisLimiter implementation
type RedisLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRedisLimiter(addr string, limit int, window time.Duration) *RedisLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisLimiter{
		client: rdb,
		limit:  limit,
		window: window,
	}
}

func (l *RedisLimiter) Allow(ctx context.Context, key string) (bool, error) {
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

// MemoryLimiter implementation using golang.org/x/time/rate
type MemoryLimiter struct {
	mu       sync.Mutex
	visitors map[string]*rate.Limiter
	limit    rate.Limit
	burst    int
}

func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	// Convert req/window to req/second for rate.Limit
	r := rate.Limit(float64(limit) / window.Seconds())
	return &MemoryLimiter{
		visitors: make(map[string]*rate.Limiter),
		limit:    r,
		burst:    limit,
	}
}

func (l *MemoryLimiter) Allow(ctx context.Context, key string) (bool, error) {
	l.mu.Lock()
	limiter, exists := l.visitors[key]
	if !exists {
		limiter = rate.NewLimiter(l.limit, l.burst)
		l.visitors[key] = limiter
	}
	l.mu.Unlock()

	return limiter.Allow(), nil
}
