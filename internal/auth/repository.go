package auth

import (
	"context"
	"explorer/models"
	"time"

	"github.com/google/uuid"
)

// UserRepo defines methods the service layer expects
// any repository it interacts with to implement
type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, uid uuid.UUID) (*models.User, error)
}

// TokenRepository defines methods it expects a repository
// it interacts with to implement
type TokenRepo interface {
	FetchToken(ctx context.Context, tokenId string) (string, error)
	SetToken(ctx context.Context, tokenId, token string, expiresIn time.Duration) error
	DeleteToken(ctx context.Context, tokenId string) error
}
