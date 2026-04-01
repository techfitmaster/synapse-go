package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache provides a simple key-value caching interface.
type Cache interface {
	// Get retrieves a value by key. Returns redis.Nil error if not found.
	Get(ctx context.Context, key string) (string, error)
	// Set stores a value with the given TTL.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	// Del removes a value by key.
	Del(ctx context.Context, key string) error
}

// GetOrLoad retrieves a cached value, or calls loader to populate the cache on miss.
func GetOrLoad(ctx context.Context, c Cache, key string, ttl time.Duration, loader func() (string, error)) (string, error) {
	val, err := c.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	result, err := loader()
	if err != nil {
		return "", err
	}

	if setErr := c.Set(ctx, key, result, ttl); setErr != nil {
		// Cache set failure is non-fatal — log but return the loaded value
		return result, nil
	}

	return result, nil
}

// RedisCache implements Cache using Redis.
type RedisCache struct {
	rdb    *redis.Client
	prefix string
}

// NewRedis creates a Cache backed by Redis with the given key prefix.
func NewRedis(rdb *redis.Client, prefix string) *RedisCache {
	return &RedisCache{rdb: rdb, prefix: prefix}
}

// Get retrieves a value from Redis.
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, c.prefixedKey(key)).Result()
}

// Set stores a value in Redis with the given TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.rdb.Set(ctx, c.prefixedKey(key), value, ttl).Err()
}

// Del removes a value from Redis.
func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, c.prefixedKey(key)).Err()
}

func (c *RedisCache) prefixedKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}
