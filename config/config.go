package config

import (
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

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	ipLimit, _ := strconv.Atoi(os.Getenv("DEFAULT_IP_LIMIT"))
	blockTime, _ := strconv.Atoi(os.Getenv("DEFAULT_IP_BLOCK_TIME_SECONDS"))

	return &Config{
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		RedisDB:            redisDB,
		DefaultIPLimit:     ipLimit,
		DefaultIPBlockTime: time.Duration(blockTime) * time.Second,
		ServerPort:         os.Getenv("SERVER_PORT"),
	}, nil
}
