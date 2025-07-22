package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/henriquedessen/rate-limiter/internal/middleware"
	"github.com/stretchr/testify/assert"
)

type mockLimiter struct{}

func (m *mockLimiter) Allow(ctx context.Context, key string, limit int, blockTime time.Duration) (bool, error) {
	if key == "token:block" || key == "ip:1.2.3.4" {
		return false, nil
	}
	return true, nil
}

func (m *mockLimiter) GetTokenLimit(token string) (int, time.Duration) {
	if token == "valid" {
		return 5, time.Minute
	}
	return 0, 0
}

func TestMiddleware_AllowWithAPIKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "valid")
	w := httptest.NewRecorder()

	handler := middleware.RateLimiterMiddleware(&mockLimiter{}, 10, time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_BlockWithAPIKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "block")
	w := httptest.NewRecorder()

	handler := middleware.RateLimiterMiddleware(&mockLimiter{}, 10, time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestMiddleware_BlockWithIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:56789"
	w := httptest.NewRecorder()

	handler := middleware.RateLimiterMiddleware(&mockLimiter{}, 10, time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
