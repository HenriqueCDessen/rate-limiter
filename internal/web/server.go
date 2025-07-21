package web

import (
	"net/http"
	"time"

	"github.com/henriquedessen/rate-limiter/internal/limiter"
	"github.com/henriquedessen/rate-limiter/internal/middleware"
)

func NewServer(limiter limiter.RateLimiter, defaultLimit int, defaultBlockTime time.Duration) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	handler := middleware.RateLimiterMiddleware(limiter, defaultLimit, defaultBlockTime)(mux)

	return handler
}
