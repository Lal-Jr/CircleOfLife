package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis(redisURL string) {
	Client = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis at %s. Caching will gracefully bypass: %v", redisURL, err)
		return
	}

	log.Println("Connected to Redis successfully.")
}
