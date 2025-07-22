package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/henriquedessen/rate-limiter/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "secret")
	os.Setenv("REDIS_DB", "0")
	os.Setenv("DEFAULT_IP_LIMIT", "100")
	os.Setenv("DEFAULT_IP_BLOCK_TIME_SECONDS", "60")
	os.Setenv("SERVER_PORT", "8080")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
	assert.Equal(t, "secret", cfg.RedisPassword)
	assert.Equal(t, 0, cfg.RedisDB)
	assert.Equal(t, 100, cfg.DefaultIPLimit)
	assert.Equal(t, 60*time.Second, cfg.DefaultIPBlockTime)
	assert.Equal(t, "8080", cfg.ServerPort)
}
