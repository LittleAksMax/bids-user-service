package health

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// CacheHealthChecker wraps *redis.Client to implement the HealthChecker interface.
type CacheHealthChecker struct {
	cache *redis.Client
}

// NewCacheHealthChecker creates a new database health checker.
func NewCacheHealthChecker(redis *redis.Client) HealthChecker {
	return &CacheHealthChecker{cache: redis}
}

// HealthCheck pings the database to verify connectivity.
func (h *CacheHealthChecker) HealthCheck(ctx context.Context) error {
	return h.cache.Ping(ctx).Err()
}
