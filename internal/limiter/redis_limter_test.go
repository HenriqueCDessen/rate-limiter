package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/henriquedessen/rate-limiter/internal/limiter"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})
}

func TestAllowWithinLimit(t *testing.T) {
	ctx := context.Background()
	client := setupRedis()
	client.FlushDB(ctx)

	rl := limiter.NewRedisRateLimiter("localhost:6379", "", 1, nil)

	key := "test:key"
	limit := 3
	block := 5 * time.Second

	for i := 0; i < limit; i++ {
		allowed, err := rl.Allow(ctx, key, limit, block)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rl.Allow(ctx, key, limit, block)
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestGetTokenLimit(t *testing.T) {
	tokens := map[string]limiter.TokenConfig{
		"abc": {Limit: 10, BlockTime: time.Minute},
	}
	rl := limiter.NewRedisRateLimiter("localhost:6379", "", 1, tokens)

	limit, blockTime := rl.GetTokenLimit("abc")
	assert.Equal(t, 10, limit)
	assert.Equal(t, time.Minute, blockTime)

	limit, blockTime = rl.GetTokenLimit("invalid")
	assert.Equal(t, 0, limit)
	assert.Equal(t, time.Duration(0), blockTime)
}
