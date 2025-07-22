package main_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/henriquedessen/rate-limiter/config"
	"github.com/henriquedessen/rate-limiter/internal/web"
	"github.com/stretchr/testify/assert"
)

type mockLimiter struct{}

func (m *mockLimiter) Allow(ctx context.Context, key string, limit int, blockTime time.Duration) (bool, error) {
	return true, nil
}
func (m *mockLimiter) GetTokenLimit(token string) (int, time.Duration) {
	return 100, time.Minute
}

func TestRunServer_WithInvalidConfig(t *testing.T) {
	os.Clearenv()
	_, err := config.LoadConfig()
	assert.Error(t, err)
}

func TestRunServer_Integration(t *testing.T) {
	// usando porta fake para não subir servidor real
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")
	os.Setenv("DEFAULT_IP_LIMIT", "100")
	os.Setenv("DEFAULT_IP_BLOCK_TIME_SECONDS", "60")
	os.Setenv("SERVER_PORT", "9999")

	rl := &mockLimiter{}
	server := web.NewServer(rl, 100, time.Minute)

	go func() {
		_ = http.ListenAndServe(":9999", server)
	}()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:9999")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
