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
		return nil, fmt.Errorf("Error connecting to redis: %v", err)
	}

	log.Println("Connected to Redis.")
	return &RedisDataSource{
		Client: client,
	}, nil
}

// close to be used in graceful server shutdown
func (r *RedisDataSource) Close() error {
	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("error closing Redis: %v", err)
	}

	return nil
}
