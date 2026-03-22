package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRequestCache struct {
	client *redis.Client
	keyNS  string
}

// NewRedisRequestCache creates a Redis-backed RequestCache using an existing redis.Client.
func NewRedisRequestCache(client *redis.Client) RequestCache {
	return &redisRequestCache{client: client, keyNS: "req_cache"}
}

func (c *redisRequestCache) buildKey(key string) string {
	return c.keyNS + ":" + key
}

func (c *redisRequestCache) Set(ctx context.Context, key string, value string, expiresIn time.Duration) error {
	return c.client.Set(ctx, c.buildKey(key), value, expiresIn).Err()
}

func (c *redisRequestCache) Get(ctx context.Context, key string) (string, time.Time, error) {
	return redisGetWithTTL(ctx, c.client, c.buildKey(key))
}

func (c *redisRequestCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.buildKey(key)).Err()
}

func (c *redisRequestCache) Close() error {
	return c.client.Close()
}

func (c *redisRequestCache) HealthCheck(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
