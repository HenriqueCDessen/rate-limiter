package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	DefaultIPLimit     int
	DefaultIPBlockTime time.Duration
	ServerPort         string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	redisDBStr := os.Getenv("REDIS_DB")
	ipLimitStr := os.Getenv("DEFAULT_IP_LIMIT")
	blockTimeStr := os.Getenv("DEFAULT_IP_BLOCK_TIME_SECONDS")

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	ipLimit, err := strconv.Atoi(ipLimitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DEFAULT_IP_LIMIT: %w", err)
	}

	blockTime, err := strconv.Atoi(blockTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DEFAULT_IP_BLOCK_TIME_SECONDS: %w", err)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return nil, fmt.Errorf("SERVER_PORT is required")
	}

	return &Config{
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		RedisDB:            redisDB,
		DefaultIPLimit:     ipLimit,
		DefaultIPBlockTime: time.Duration(blockTime) * time.Second,
		ServerPort:         port,
	}, nil
}
