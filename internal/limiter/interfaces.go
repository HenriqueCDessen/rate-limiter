package limiter

import (
	"context"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, blockTime time.Duration) (bool, error)
	GetTokenLimit(token string) (int, time.Duration)
}
