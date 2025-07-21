package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client       *redis.Client
	tokenConfigs map[string]TokenConfig
}

type TokenConfig struct {
	Limit     int
	BlockTime time.Duration
}

func NewRedisRateLimiter(addr, password string, db int, tokenConfigs map[string]TokenConfig) *RedisRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisRateLimiter{
		client:       client,
		tokenConfigs: tokenConfigs,
	}
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit int, blockTime time.Duration) (bool, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.client.Expire(ctx, key, blockTime)
	}

	return count <= int64(limit), nil
}

func (r *RedisRateLimiter) GetTokenLimit(token string) (int, time.Duration) {
	if config, exists := r.tokenConfigs[token]; exists {
		return config.Limit, config.BlockTime
	}
	return 0, 0
}
