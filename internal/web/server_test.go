package web_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/henriquedessen/rate-limiter/internal/web"
	"github.com/stretchr/testify/assert"
)

type mockLimiter struct{}

func (m *mockLimiter) Allow(_ context.Context, _ string, _ int, _ time.Duration) (bool, error) {
	return true, nil
}

func (m *mockLimiter) GetTokenLimit(_ string) (int, time.Duration) {
	return 0, 0
}

func TestNewServer(t *testing.T) {
	handler := web.NewServer(&mockLimiter{}, 10, time.Second)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Hello, World!")
}
