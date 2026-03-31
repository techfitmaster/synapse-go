package sequence

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Generator produces sequential IDs using Redis INCR.
// Format: "{bizPrefix}{YYYYMMDD}{seq}" e.g. "SH20260331001".
type Generator struct {
	rdb       *redis.Client
	keyPrefix string
}

// New creates a Generator with the given Redis client and key prefix.
// The key prefix is used to namespace Redis keys (e.g. "cargo" → "cargo:seq:SH:20260331").
func New(rdb *redis.Client, keyPrefix string) *Generator {
	return &Generator{rdb: rdb, keyPrefix: keyPrefix}
}

// Next generates the next sequential ID for the given business prefix.
// The sequence resets daily. Keys auto-expire after 48 hours.
func (g *Generator) Next(ctx context.Context, bizPrefix string) (string, error) {
	date := time.Now().Format("20060102")
	key := fmt.Sprintf("%s:seq:%s:%s", g.keyPrefix, bizPrefix, date)

	seq, err := g.rdb.Incr(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("sequence incr: %w", err)
	}

	g.rdb.Expire(ctx, key, 48*time.Hour)

	return fmt.Sprintf("%s%s%03d", bizPrefix, date, seq), nil
}
