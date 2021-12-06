package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

type RedisDataSource struct {
	Client *redis.Client
}

func InitRedis() (*RedisDataSource, error) {
	// Initialize redis connection
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	log.Printf("Connecting to Redis\n")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	// verify redis connection
	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, fmt.Errorf("Error connecting to redis: %w", err)
	}

	return &RedisDataSource{
		Client: client,
	}, nil
}
