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

// SetRefreshToken stores a refresh token with an expiry time
func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, tokenID string, userID string, expiresIn time.Duration) error {
	// key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, tokenID, userID, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for tokenID/userID: %s/%s: %v\n", tokenID, userID, err)
		return apperrors.NewInternal()
	}
	return nil
}

// SetAccessToken stores a access token with an expiry time
func (r *redisTokenRepository) SetAccessToken(ctx context.Context, tokenID string, userID string, expiresIn time.Duration) error {
	if err := r.Redis.Set(ctx, tokenID, userID, expiresIn).Err(); err != nil {
		log.Printf("Could not SET access token to redis for tokenID/userID: %s/%s: %v\n", tokenID, userID, err)
		return apperrors.NewInternal()
	}
	return nil
}

// DeleteRefreshToken used to delete old  refresh tokens
// Services my access this to revolve tokens
func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, tokenID string) error {
	// key := fmt.Sprintf("%s:%s", tokenID, userID)
	if err := r.Redis.Del(ctx, tokenID).Err(); err != nil {
		log.Printf("Could not delete refresh token from redis for tokenID: %s: %v\n", tokenID, err)
		return apperrors.NewInternal()
	}

	return nil
}

// DeleteAccessToken used to delete old access tokens
// Services my access this to revolve tokens
func (r *redisTokenRepository) DeleteAccessToken(ctx context.Context, tokenID string) error {
	if err := r.Redis.Del(ctx, tokenID).Err(); err != nil {
		log.Printf("Could not delete access token from redis for tokenID: %s: %v\n", tokenID, err)
		return apperrors.NewInternal()
	}

	return nil
}

// FetchAuth use to fetch data from Redis by tokenId
func (r *redisTokenRepository) FetchAuth(ctx context.Context, tokenID string) (string, error) {
	userid, err := r.Redis.Get(ctx, tokenID).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

// DeleteTokens used to delete tokens when logout
// func (r *redisTokenRepository) DeleteTokens(ctx context.Context, tokenID string, userID string) error {
// 	//get the refresh uuid
// 	refreshID := fmt.Sprintf("%s++%s", tokenID, userID)
// 	//delete access token
// 	if err := r.Redis.Del(ctx, tokenID).Err(); err != nil {
// 		if err != nil {
// 			log.Printf("Could not delete access token userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
// 			return apperrors.NewInternal()
// 		}
// 	}
// 	//delete refresh token
// 	if err := r.Redis.Del(ctx, refreshID).Err(); err != nil {
// 		if err != nil {
// 			log.Printf("Could not delete refresh token userID/tokenID: %s/%s: %v\n", userID, refreshID, err)
// 			return apperrors.NewInternal()
// 		}
// 	}

// 	return nil
// }
