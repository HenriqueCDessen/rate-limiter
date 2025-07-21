package main

import (
	"log"
	"net/http"
	"time"

	"github.com/henriquedessen/rate-limiter/config"
	"github.com/henriquedessen/rate-limiter/internal/limiter"
	"github.com/henriquedessen/rate-limiter/internal/web"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	tokenConfigs := map[string]limiter.TokenConfig{
		"abc123": {
			Limit:     100,
			BlockTime: 5 * time.Minute,
		},
		"test123": {
			Limit:     50,
			BlockTime: 2 * time.Minute,
		},
	}

	rl := limiter.NewRedisRateLimiter(
		cfg.RedisAddr,
		cfg.RedisPassword,
		cfg.RedisDB,
		tokenConfigs,
	)

	server := web.NewServer(rl, cfg.DefaultIPLimit, cfg.DefaultIPBlockTime)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, server))
}
