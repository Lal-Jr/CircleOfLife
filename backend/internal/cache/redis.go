package cache

import (
	"context"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis(redisURL string) {
	opts := &redis.Options{Addr: redisURL}
	// Accept either a bare host:port or a full redis:// / rediss:// URL (with password/TLS).
	if strings.Contains(redisURL, "://") {
		parsed, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Printf("Warning: invalid REDIS_URL: %v", err)
			return
		}
		opts = parsed
	}
	Client = redis.NewClient(opts)

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis at %s. Caching will gracefully bypass: %v", redisURL, err)
		return
	}

	log.Println("Connected to Redis successfully.")
}
