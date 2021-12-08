package auth

import (
	"context"
	"explorer/models"
)

const CtxUserKey = "user"

// UserUsecase defines methods the handler layer expects
// any service it interacts with to implement
type UserUsecase interface {
	SignUp(ctx context.Context, u *models.User) error
	SignIn(ctx context.Context, u *models.User) error
	// ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producing JWTs as string
type TokenUsecase interface {
	NewPairTokens(ctx context.Context, u *models.User, prevTokenID string) (*models.TokenPair, error)
	ValidateIDToken(ctx context.Context, tokenString string) (*models.User, error)
}
