package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisGetWithTTL fetches a value from Redis and returns it alongside its expiry time.
func redisGetWithTTL(ctx context.Context, client *redis.Client, key string) (string, time.Time, error) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", time.Time{}, ErrNotFound
		}
		return "", time.Time{}, err
	}

	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		return "", time.Time{}, err
	}
	if ttl <= 0 {
		return "", time.Time{}, ErrExpired
	}

	return val, time.Now().UTC().Add(ttl), nil
}
