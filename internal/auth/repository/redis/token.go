package repository

import (
	"context"
	"log"
	"time"

	"explorer/internal/apperrors"
	"explorer/internal/auth"

	"github.com/go-redis/redis/v8"
)

// redisTokenRepository is data/repository implementation
// of service layer TokenRepository
type redisTokenRepository struct {
	Redis *redis.Client
}

// NewTokenRepository is a factory for initializing User Repositories
func NewTokenRepository(redisClient *redis.Client) auth.TokenRepo {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// SetRefreshToken stores token with an expiry time
func (r *redisTokenRepository) SetToken(ctx context.Context, key string, value string, expiresIn time.Duration) error {
	if err := r.Redis.Set(ctx, key, value, expiresIn).Err(); err != nil {
		log.Printf("Could not SET token to redis for %s/%s: %v\n", key, value, err)
		return apperrors.NewInternal()
	}
	return nil
}

// DeleteToken used to delete old  tokens
// Services my access this to revolve tokens
func (r *redisTokenRepository) DeleteToken(ctx context.Context, tokenId string) error {
	if err := r.Redis.Del(ctx, tokenId).Err(); err != nil {
		log.Printf("Could not delete token from redis for %s: %v\n", tokenId, err)
		return apperrors.NewInternal()
	}

	return nil
}

// FetchAuth use to fetch data from Redis by tokenId
func (r *redisTokenRepository) FetchToken(ctx context.Context, tokenId string) (string, error) {
	val, err := r.Redis.Get(ctx, tokenId).Result()
	if err != nil {
		log.Printf("Could not fetch token from redis for %s: %v\n", tokenId, err)
		return "", apperrors.NewInternal()
	}
	return val, nil
}
