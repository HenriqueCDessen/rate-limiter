package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/henriquedessen/rate-limiter/internal/limiter"
)

func RateLimiterMiddleware(limiter limiter.RateLimiter, defaultLimit int, defaultBlockTime time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ip := strings.Split(r.RemoteAddr, ":")[0]

			apiKey := r.Header.Get("API_KEY")

			var limit int
			var blockTime time.Duration

			if apiKey != "" {
				tokenLimit, tokenBlockTime := limiter.GetTokenLimit(apiKey)
				if tokenLimit > 0 {
					limit = tokenLimit
					blockTime = tokenBlockTime
				} else {
					limit = defaultLimit
					blockTime = defaultBlockTime
				}
				key := "token:" + apiKey
				allowed, err := limiter.Allow(ctx, key, limit, blockTime)
				if err != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
				if !allowed {
					w.Header().Set("Retry-After", blockTime.String())
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			} else {
				limit = defaultLimit
				blockTime = defaultBlockTime
				key := "ip:" + ip
				allowed, err := limiter.Allow(ctx, key, limit, blockTime)
				if err != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
				if !allowed {
					w.Header().Set("Retry-After", blockTime.String())
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
