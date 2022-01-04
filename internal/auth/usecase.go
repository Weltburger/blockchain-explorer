package auth

import (
	"context"
	"explorer/models"
	"net/http"

	"github.com/google/uuid"
)

const CtxUserKey = "user"

// UserUsecase defines methods the handler layer expects
// any service it interacts with to implement
type UserUsecase interface {
	SignUp(ctx context.Context, u *models.User) error
	SignIn(ctx context.Context, u *models.User) error
	Get(ctx context.Context, uid uuid.UUID) (*models.User, error)
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producing JWTs as string
type TokenUsecase interface {
	GenerateTokens(ctx context.Context, user string) (*models.TokenPair, error)
	DeleteTokens(ctx context.Context, r *http.Request) error
	ValidateToken(ctx context.Context, token string, refresh bool) (*models.ValidationDetails, error)
	DeleteRefreshToken(ctx context.Context, tokenStr string) error
}
